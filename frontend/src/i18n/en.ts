export default {
  route: {
    title: {
      overview: 'Overview',
      routes: 'Routes',
      upstreams: 'Upstreams',
      circuitBreakers: 'Circuit Breakers',
      rateLimiting: 'Rate Limiting',
      metrics: 'Metrics',
      settings: 'Settings',
      login: 'Login'
    },
    defaultTitle: 'Aegis',
    titleTemplate: '{title} - Aegis Gateway'
  },
  header: {
    uptime: 'Uptime:',
    port: 'Port:',
    status: { running: 'Running' },
    logout: 'Logout'
  },
  sidebar: {
    brand: 'Aegis',
    version: 'API Gateway v1.0.0',
    section: {
      monitor: 'Monitor',
      management: 'Management',
      security: 'Security',
      system: 'System'
    },
    nav: {
      overview: 'Overview',
      metrics: 'Metrics',
      routes: 'Routes',
      upstreams: 'Upstreams',
      circuitBreakers: 'Circuit Breakers',
      rateLimiting: 'Rate Limiting',
      settings: 'Settings'
    }
  },
  common: {
    cancel: 'Cancel',
    confirm: 'Confirm',
    loading: 'Loading...',
    noData: 'No data',
    edit: 'Edit',
    delete: 'Delete',
    reset: 'Reset',
    register: 'Register',
    deregister: 'Deregister',
    add: 'Add',
    actions: 'Actions',
    address: 'Address',
    service: 'Service',
    status: 'Status',
    latency: 'Latency',
    weight: 'Weight',
    route: 'Route',
    healthy: 'Healthy',
    down: 'Down'
  },
  login: {
    brand: 'AEGIS',
    subtitle: 'API Gateway Management',
    username: 'Username',
    usernamePlaceholder: 'admin',
    password: 'Password',
    passwordPlaceholder: 'password',
    signIn: 'Sign In',
    signingIn: 'Signing in...',
    error: {
      missingFields: 'Please enter username and password',
      failed: 'Login failed'
    }
  },
  dashboard: {
    metric: {
      requestsPerSec: 'Requests / sec',
      avgLatency: 'Avg Latency',
      totalRequests: 'Total Requests',
      inFlight: 'In Flight'
    },
    upstreamHealth: 'Upstream Health',
    empty: 'No upstreams registered'
  },
  routes: {
    header: 'Route Configuration',
    addRoute: '+ Add Route',
    col: {
      path: 'Path',
      method: 'Method',
      matchType: 'Match Type',
      balancer: 'Balancer',
      rateLimit: 'Rate Limit',
      upstreams: 'Upstreams'
    },
    empty: 'No routes configured',
    confirm: {
      title: 'Delete Route',
      message: 'Are you sure you want to delete route "{path}"?'
    },
    success: { deleted: 'Route deleted' },
    error: { deleteFailed: 'Failed to delete route' }
  },
  upstreams: {
    header: 'Upstream Servers',
    register: '+ Register',
    col: { connections: 'Connections' },
    empty: 'No upstreams registered',
    dialog: { title: 'Register Upstream' },
    form: {
      serviceName: 'Service Name',
      servicePlaceholder: 'user-service',
      address: 'Address',
      addressPlaceholder: 'http://localhost:9001',
      weight: 'Weight'
    },
    confirm: {
      title: 'Deregister Upstream',
      message: 'Remove {address} from {service}?'
    },
    success: {
      deregistered: 'Upstream deregistered',
      registered: 'Upstream registered'
    },
    error: {
      deregisterFailed: 'Failed to deregister',
      registerFailed: 'Failed to register',
      required: 'Service and address are required'
    }
  },
  circuit: {
    header: 'Circuit Breaker States',
    col: {
      state: 'State',
      errorRate: 'Error Rate',
      minRequests: 'Min Requests',
      window: 'Window',
      recovery: 'Recovery'
    },
    empty: 'No circuit breakers configured',
    stateMachine: 'State Machine',
    state: {
      closed: 'CLOSED',
      closedDesc: 'Normal operation',
      open: 'OPEN',
      openDesc: 'Rejecting requests',
      halfOpen: 'HALF-OPEN',
      halfOpenDesc: 'Testing recovery'
    },
    success: { reset: 'Circuit breaker reset' },
    error: { resetFailed: 'Failed to reset' }
  },
  ratelimit: {
    header: 'Rate Limiting Rules',
    addRule: '+ Add Rule',
    col: {
      strategy: 'Strategy',
      rate: 'Rate',
      burst: 'Burst'
    },
    empty: 'No rate limit rules configured',
    dialog: { title: 'Add Rate Limit Rule' },
    form: {
      route: 'Route',
      routePlaceholder: '/api/users',
      strategy: 'Strategy',
      strategyOptions: {
        tokenBucket: 'Token Bucket',
        redis: 'Redis Sliding Window'
      },
      rate: 'Rate (req/s)',
      burst: 'Burst'
    },
    confirm: {
      title: 'Delete Rule',
      message: 'Remove rate limit for "{route}"?'
    },
    success: {
      deleted: 'Rule deleted',
      added: 'Rule added'
    },
    error: {
      deleteFailed: 'Failed to delete rule',
      addFailed: 'Failed to add rule',
      routeRequired: 'Route is required'
    }
  },
  metrics: {
    requestsPerSec: 'Requests / sec',
    latency: 'Latency',
    statusDistribution: 'Status Code Distribution'
  },
  settings: {
    serverConfig: 'Server Configuration',
    form: {
      port: 'Port',
      host: 'Host',
      readTimeout: 'Read Timeout',
      writeTimeout: 'Write Timeout',
      etcdEndpoints: 'etcd Endpoints',
      etcdPrefix: 'etcd Prefix',
      redisAddress: 'Redis Address',
      redisDb: 'Redis DB',
      tokenExpiry: 'Token Expiry',
      logLevel: 'Log Level',
      logFormat: 'Log Format'
    },
    externalServices: 'External Services',
    authLogging: 'Authentication & Logging',
    error: { loadFailed: 'Failed to load settings' }
  },
  chart: {
    legend: { avg: 'Avg', p99: 'P99' },
    axis: { ms: 'ms' },
    status: {
      '2xx': '2xx',
      '4xx': '4xx',
      '5xx': '5xx'
    }
  }
}
