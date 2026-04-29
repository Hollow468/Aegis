package model

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
	Proxy  ProxyConfig  `mapstructure:"proxy"`
	Etcd   EtcdConfig   `mapstructure:"etcd"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Redis  RedisConfig  `mapstructure:"redis"`
	Routes []Route      `mapstructure:"routes"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type EtcdConfig struct {
	Endpoints []string `mapstructure:"endpoints"`
	Timeout   int      `mapstructure:"timeout"`
	Prefix    string   `mapstructure:"prefix"`
}

type JWTConfig struct {
	Secret     string   `mapstructure:"secret"`
	Expiration int      `mapstructure:"expiration"` // seconds
	WhiteList  []string `mapstructure:"white_list"`  // paths that skip auth
	AdminUser  string   `mapstructure:"admin_user"`
	AdminPass  string   `mapstructure:"admin_pass"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Route struct {
	Path             string           `mapstructure:"path"`
	Method           string           `mapstructure:"method"`
	MatchType        string           `mapstructure:"match_type"`
	Balancer         string           `mapstructure:"balancer"`
	HashKey          string           `mapstructure:"hash_key"`
	Upstreams        []Upstream       `mapstructure:"upstreams"`
	RateLimit        *RateLimit       `mapstructure:"rate_limit"`
	CircuitBreaker   *CircuitBreaker  `mapstructure:"circuit_breaker"`
	ServiceDiscovery string           `mapstructure:"service_discovery"`
}

type ProxyConfig struct {
	MaxIdleConns        int `mapstructure:"max_idle_conns"`
	MaxIdleConnsPerHost int `mapstructure:"max_idle_conns_per_host"`
	IdleConnTimeout     int `mapstructure:"idle_conn_timeout"`
	ResponseTimeout     int `mapstructure:"response_timeout"`
}

type Upstream struct {
	Address string `mapstructure:"address"`
	Weight  int    `mapstructure:"weight"`
}

type RateLimit struct {
	Type  string  `mapstructure:"type"`  // "token_bucket" or "redis"
	Limit float64 `mapstructure:"limit"` // Rate per second
	Burst int     `mapstructure:"burst"` // Max burst
}

type CircuitBreaker struct {
	ErrorThreshold float64 `mapstructure:"error_threshold"` // error rate threshold (0.0-1.0)
	MinRequests    int     `mapstructure:"min_requests"`    // min requests before tripping
	WindowSeconds  int     `mapstructure:"window_seconds"`  // sliding window size
	RecoveryTime   int     `mapstructure:"recovery_time"`   // seconds in Open before Half-Open
}
