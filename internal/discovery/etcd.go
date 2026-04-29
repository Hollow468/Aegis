package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"apigateway/internal/logger"
	"apigateway/internal/model"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// ServiceInstance represents a registered service instance in etcd.
type ServiceInstance struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Weight  int    `json:"weight"`
}

// ServiceDiscovery watches etcd for service registration changes.
type ServiceDiscovery struct {
	client  *clientv3.Client
	prefix  string
	mu      sync.RWMutex
	// serviceName -> list of upstreams
	services map[string][]model.Upstream
	// callbacks for upstream changes
	watchers []func(serviceName string, upstreams []model.Upstream)
}

// NewServiceDiscovery creates a new etcd-based service discovery.
func NewServiceDiscovery(cfg model.EtcdConfig) (*ServiceDiscovery, error) {
	if len(cfg.Endpoints) == 0 {
		return nil, fmt.Errorf("etcd endpoints cannot be empty")
	}

	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "/services"
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	sd := &ServiceDiscovery{
		client:   client,
		prefix:   prefix,
		services: make(map[string][]model.Upstream),
	}

	// Load existing services
	if err := sd.loadServices(); err != nil {
		logger.Log.Warn("failed to load existing services from etcd", zap.Error(err))
	}

	// Start watching for changes
	go sd.watch()

	return sd, nil
}

// loadServices loads all existing services from etcd.
func (sd *ServiceDiscovery) loadServices() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := sd.client.Get(ctx, sd.prefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to get services: %w", err)
	}

	for _, kv := range resp.Kvs {
		sd.processKeyValue(string(kv.Key), string(kv.Value))
	}

	logger.Log.Info("loaded services from etcd",
		zap.Int("total_keys", len(resp.Kvs)),
		zap.Int("services", len(sd.services)),
	)
	return nil
}

// watch watches etcd for service registration changes.
func (sd *ServiceDiscovery) watch() {
	watchCh := sd.client.Watch(context.Background(), sd.prefix, clientv3.WithPrefix())

	for resp := range watchCh {
		for _, event := range resp.Events {
			key := string(event.Kv.Key)
			switch event.Type {
			case clientv3.EventTypePut:
				sd.processKeyValue(key, string(event.Kv.Value))
				logger.Log.Info("service updated",
					zap.String("key", key),
				)
			case clientv3.EventTypeDelete:
				sd.removeKey(key)
				logger.Log.Info("service removed",
					zap.String("key", key),
				)
			}
			sd.notifyWatchers()
		}
	}
}

// processKeyValue parses a key-value pair and updates the services map.
// Key format: /prefix/serviceName/instanceID
// Value: JSON ServiceInstance
func (sd *ServiceDiscovery) processKeyValue(key, value string) {
	var instance ServiceInstance
	if err := json.Unmarshal([]byte(value), &instance); err != nil {
		logger.Log.Warn("failed to parse service instance",
			zap.String("key", key),
			zap.Error(err),
		)
		return
	}

	serviceName := instance.Name
	if serviceName == "" {
		// Extract from key
		parts := strings.Split(strings.TrimPrefix(key, sd.prefix), "/")
		if len(parts) > 0 {
			serviceName = parts[0]
		}
	}
	if serviceName == "" {
		return
	}

	weight := instance.Weight
	if weight <= 0 {
		weight = 1
	}

	sd.mu.Lock()
	defer sd.mu.Unlock()

	// Check if this instance already exists (update case)
	upstreams := sd.services[serviceName]
	found := false
	for i, u := range upstreams {
		if u.Address == instance.Address {
			upstreams[i].Weight = weight
			found = true
			break
		}
	}
	if !found {
		upstreams = append(upstreams, model.Upstream{
			Address: instance.Address,
			Weight:  weight,
		})
	}
	sd.services[serviceName] = upstreams
}

// removeKey removes a service instance by key.
func (sd *ServiceDiscovery) removeKey(key string) {
	parts := strings.Split(strings.TrimPrefix(key, sd.prefix), "/")
	if len(parts) < 2 {
		return
	}
	serviceName := parts[0]
	instanceID := parts[1]

	sd.mu.Lock()
	defer sd.mu.Unlock()

	upstreams := sd.services[serviceName]
	for i, u := range upstreams {
		// Match by address containing instanceID or by address
		if strings.Contains(u.Address, instanceID) {
			sd.services[serviceName] = append(upstreams[:i], upstreams[i+1:]...)
			return
		}
	}
}

// GetUpstreams returns the current upstreams for a service name.
func (sd *ServiceDiscovery) GetUpstreams(serviceName string) []model.Upstream {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	upstreams := sd.services[serviceName]
	result := make([]model.Upstream, len(upstreams))
	copy(result, upstreams)
	return result
}

// OnChange registers a callback for upstream changes.
func (sd *ServiceDiscovery) OnChange(callback func(serviceName string, upstreams []model.Upstream)) {
	sd.mu.Lock()
	defer sd.mu.Unlock()
	sd.watchers = append(sd.watchers, callback)
}

// notifyWatchers notifies all registered callbacks about changes.
func (sd *ServiceDiscovery) notifyWatchers() {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	for serviceName, upstreams := range sd.services {
		for _, watcher := range sd.watchers {
			watcher(serviceName, upstreams)
		}
	}
}

// Register registers a service instance in etcd.
func (sd *ServiceDiscovery) Register(ctx context.Context, instance ServiceInstance) error {
	key := fmt.Sprintf("%s%s/%s", sd.prefix, instance.Name, instance.Address)
	value, err := json.Marshal(instance)
	if err != nil {
		return fmt.Errorf("failed to marshal instance: %w", err)
	}

	// Use lease for automatic deregistration
	lease, err := sd.client.Grant(ctx, 10) // 10 second TTL
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}

	_, err = sd.client.Put(ctx, key, string(value), clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	// Keep alive
	ch, err := sd.client.KeepAlive(ctx, lease.ID)
	if err != nil {
		return fmt.Errorf("failed to keep alive: %w", err)
	}

	// Drain keepalive responses
	go func() {
		for range ch {
		}
	}()

	return nil
}

// Close closes the etcd client connection.
func (sd *ServiceDiscovery) Close() error {
	if sd.client != nil {
		return sd.client.Close()
	}
	return nil
}
