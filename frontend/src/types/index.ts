// Route types
export interface Route {
  path: string
  method: string
  match_type: 'exact' | 'prefix' | 'regex'
  balancer: 'round_robin' | 'random' | 'weighted_round_robin' | 'consistent_hash'
  hash_key?: string
  upstreams: Upstream[]
  rate_limit?: RateLimit
  circuit_breaker?: CircuitBreaker
  service_discovery?: string
}

export interface Upstream {
  address: string
  weight: number
  healthy?: boolean
  latency?: number
}

export interface RateLimit {
  type: 'token_bucket' | 'redis'
  limit: number
  burst: number
}

export interface CircuitBreaker {
  error_threshold: number
  min_requests: number
  window_seconds: number
  recovery_time: number
}

// API response types
export interface ApiResponse<T> {
  code: number
  data: T
  message?: string
}

export interface RouteListResponse {
  routes: Route[]
  total: number
}

// Metrics types
export interface MetricsSummary {
  qps: number
  avg_latency: number
  p99_latency: number
  total_requests: number
  in_flight: number
  status_2xx: number
  status_4xx: number
  status_5xx: number
}

export interface UpstreamStatus {
  address: string
  service: string
  healthy: boolean
  latency: number
  weight: number
  connections: number
}

// Circuit breaker state
export type CircuitState = 'closed' | 'open' | 'half-open'

export interface CircuitBreakerStatus {
  route: string
  state: CircuitState
  error_rate: number
  min_requests: number
  window_seconds: number
  recovery_time: number
}

// Auth types
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  expires_in: number
}

// App types
export interface NavItem {
  key: string
  label: string
  icon: string
  badge?: number
  section?: string
}
