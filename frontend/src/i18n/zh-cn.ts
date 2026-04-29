export default {
  route: {
    title: {
      overview: '概览',
      routes: '路由管理',
      upstreams: '上游服务',
      circuitBreakers: '熔断器',
      rateLimiting: '限流规则',
      metrics: '监控指标',
      settings: '系统设置',
      login: '登录'
    },
    defaultTitle: 'Aegis',
    titleTemplate: '{title} - Aegis 网关'
  },
  header: {
    uptime: '运行时间:',
    port: '端口:',
    status: { running: '运行中' },
    logout: '退出登录'
  },
  sidebar: {
    brand: 'Aegis',
    version: 'API 网关 v1.0.0',
    section: {
      monitor: '监控中心',
      management: '服务管理',
      security: '安全策略',
      system: '系统配置'
    },
    nav: {
      overview: '概览',
      metrics: '监控指标',
      routes: '路由管理',
      upstreams: '上游服务',
      circuitBreakers: '熔断器',
      rateLimiting: '限流规则',
      settings: '系统设置'
    }
  },
  common: {
    cancel: '取消',
    confirm: '确认',
    loading: '加载中...',
    noData: '暂无数据',
    edit: '编辑',
    delete: '删除',
    reset: '重置',
    register: '注册',
    deregister: '注销',
    add: '添加',
    actions: '操作',
    address: '地址',
    service: '服务',
    status: '状态',
    latency: '延迟',
    weight: '权重',
    route: '路由',
    healthy: '健康',
    down: '宕机'
  },
  login: {
    brand: 'AEGIS',
    subtitle: 'API 网关管理平台',
    username: '用户名',
    usernamePlaceholder: '请输入用户名',
    password: '密码',
    passwordPlaceholder: '请输入密码',
    signIn: '登 录',
    signingIn: '登录中...',
    error: {
      missingFields: '请输入用户名和密码',
      failed: '登录失败'
    }
  },
  dashboard: {
    metric: {
      requestsPerSec: '每秒请求数',
      avgLatency: '平均延迟',
      totalRequests: '总请求数',
      inFlight: '进行中请求'
    },
    upstreamHealth: '上游健康状态',
    empty: '暂无注册的上游服务'
  },
  routes: {
    header: '路由配置',
    addRoute: '+ 添加路由',
    col: {
      path: '路径',
      method: '方法',
      matchType: '匹配类型',
      balancer: '负载均衡',
      rateLimit: '限流',
      upstreams: '上游数'
    },
    empty: '暂无配置的路由',
    confirm: {
      title: '删除路由',
      message: '确定要删除路由 "{path}" 吗？'
    },
    success: { deleted: '路由已删除' },
    error: { deleteFailed: '删除路由失败' }
  },
  upstreams: {
    header: '上游服务器',
    register: '+ 注册服务',
    col: { connections: '连接数' },
    empty: '暂无注册的上游服务',
    dialog: { title: '注册上游服务' },
    form: {
      serviceName: '服务名称',
      servicePlaceholder: '例如: user-service',
      address: '服务地址',
      addressPlaceholder: '例如: http://localhost:9001',
      weight: '权重'
    },
    confirm: {
      title: '注销上游服务',
      message: '确定要从 {service} 中移除 {address} 吗？'
    },
    success: {
      deregistered: '上游服务已注销',
      registered: '上游服务已注册'
    },
    error: {
      deregisterFailed: '注销失败',
      registerFailed: '注册失败',
      required: '服务名称和地址不能为空'
    }
  },
  circuit: {
    header: '熔断器状态',
    col: {
      state: '状态',
      errorRate: '错误率',
      minRequests: '最小请求数',
      window: '时间窗口',
      recovery: '恢复时间'
    },
    empty: '暂无配置的熔断器',
    stateMachine: '状态机',
    state: {
      closed: '关闭',
      closedDesc: '正常运行',
      open: '打开',
      openDesc: '拒绝请求',
      halfOpen: '半开',
      halfOpenDesc: '试探恢复'
    },
    success: { reset: '熔断器已重置' },
    error: { resetFailed: '重置失败' }
  },
  ratelimit: {
    header: '限流规则',
    addRule: '+ 添加规则',
    col: {
      strategy: '策略',
      rate: '速率',
      burst: '突发'
    },
    empty: '暂无配置的限流规则',
    dialog: { title: '添加限流规则' },
    form: {
      route: '路由',
      routePlaceholder: '例如: /api/users',
      strategy: '限流策略',
      strategyOptions: {
        tokenBucket: '令牌桶',
        redis: 'Redis 滑动窗口'
      },
      rate: '速率 (请求/秒)',
      burst: '突发容量'
    },
    confirm: {
      title: '删除规则',
      message: '确定要删除 "{route}" 的限流规则吗？'
    },
    success: {
      deleted: '规则已删除',
      added: '规则已添加'
    },
    error: {
      deleteFailed: '删除规则失败',
      addFailed: '添加规则失败',
      routeRequired: '路由不能为空'
    }
  },
  metrics: {
    requestsPerSec: '每秒请求数',
    latency: '延迟',
    statusDistribution: '状态码分布'
  },
  settings: {
    serverConfig: '服务器配置',
    form: {
      port: '端口',
      host: '监听地址',
      readTimeout: '读取超时',
      writeTimeout: '写入超时',
      etcdEndpoints: 'etcd 地址',
      etcdPrefix: 'etcd 前缀',
      redisAddress: 'Redis 地址',
      redisDb: 'Redis DB',
      tokenExpiry: 'Token 有效期',
      logLevel: '日志级别',
      logFormat: '日志格式'
    },
    externalServices: '外部服务',
    authLogging: '认证与日志',
    error: { loadFailed: '加载配置失败' }
  },
  chart: {
    legend: { avg: '平均', p99: 'P99' },
    axis: { ms: '毫秒' },
    status: {
      '2xx': '2xx 成功',
      '4xx': '4xx 客户端错误',
      '5xx': '5xx 服务端错误'
    }
  }
}
