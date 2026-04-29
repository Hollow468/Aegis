package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"apigateway/internal/discovery"
	"apigateway/internal/model"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// RouteCLI provides CLI commands for managing routes dynamically.
type RouteCLI struct {
	etcdClient *clientv3.Client
	prefix     string
}

// NewRouteCLI creates a new route CLI.
func NewRouteCLI(cfg model.EtcdConfig) (*RouteCLI, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: cfg.Endpoints,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "/services"
	}

	return &RouteCLI{etcdClient: client, prefix: prefix}, nil
}

// Run executes the CLI with the given arguments.
func (c *RouteCLI) Run(args []string) error {
	if len(args) < 1 {
		return c.printUsage()
	}

	switch args[0] {
	case "list":
		return c.listServices()
	case "register":
		return c.registerService(args[1:])
	case "deregister":
		return c.deregisterService(args[1:])
	case "status":
		return c.showStatus()
	default:
		return c.printUsage()
	}
}

func (c *RouteCLI) printUsage() error {
	fmt.Println("Aegis API Gateway - Route Management CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  aegis-cli list                    List all registered services")
	fmt.Println("  aegis-cli register <name> <addr>  Register a service instance")
	fmt.Println("  aegis-cli deregister <name> <addr> Remove a service instance")
	fmt.Println("  aegis-cli status                  Show gateway status")
	return nil
}

func (c *RouteCLI) listServices() error {
	ctx := context.Background()
	resp, err := c.etcdClient.Get(ctx, c.prefix+"/", clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	if len(resp.Kvs) == 0 {
		fmt.Println("No services registered.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SERVICE\tADDRESS\tWEIGHT")
	fmt.Fprintln(w, "-------\t-------\t------")

	for _, kv := range resp.Kvs {
		var inst discovery.ServiceInstance
		if err := json.Unmarshal(kv.Value, &inst); err != nil {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%d\n", inst.Name, inst.Address, inst.Weight)
	}

	return w.Flush()
}

func (c *RouteCLI) registerService(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: aegis-cli register <name> <address> [weight]")
	}

	name := args[0]
	addr := args[1]
	weight := 1
	if len(args) > 2 {
		fmt.Sscanf(args[2], "%d", &weight)
	}

	instance := discovery.ServiceInstance{
		Name:    name,
		Address: addr,
		Weight:  weight,
	}

	ctx := context.Background()
	key := fmt.Sprintf("%s/%s/%s", c.prefix, name, addr)
	value, err := json.Marshal(instance)
	if err != nil {
		return err
	}

	_, err = c.etcdClient.Put(ctx, key, string(value))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	fmt.Printf("Registered: %s -> %s (weight=%d)\n", name, addr, weight)
	return nil
}

func (c *RouteCLI) deregisterService(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: aegis-cli deregister <name> <address>")
	}

	name := args[0]
	addr := args[1]

	ctx := context.Background()
	key := fmt.Sprintf("%s/%s/%s", c.prefix, name, addr)

	_, err := c.etcdClient.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	fmt.Printf("Deregistered: %s -> %s\n", name, addr)
	return nil
}

func (c *RouteCLI) showStatus() error {
	ctx := context.Background()

	// Count services
	resp, err := c.etcdClient.Get(ctx, c.prefix+"/", clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to query etcd: %w", err)
	}

	services := make(map[string]int)
	for _, kv := range resp.Kvs {
		var inst discovery.ServiceInstance
		if err := json.Unmarshal(kv.Value, &inst); err != nil {
			continue
		}
		services[inst.Name]++
	}

	fmt.Println("Aegis API Gateway Status")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Registered Services: %d\n", len(services))
	fmt.Printf("Total Instances: %d\n", len(resp.Kvs))
	fmt.Println()

	for name, count := range services {
		fmt.Printf("  %s: %d instance(s)\n", name, count)
	}

	return nil
}

// Close closes the CLI's etcd connection.
func (c *RouteCLI) Close() error {
	return c.etcdClient.Close()
}
