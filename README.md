# Aegis API Gateway

高性能、云原生的轻量级 API 网关，基于 Go 构建。支持动态路由、负载均衡、限流熔断、JWT 鉴权，并内置 Web 管理面板。

## Features

**路由引擎**
- Trie 树路由算法，支持精确匹配、前缀匹配、正则匹配
- 按 HTTP 方法路由，支持通配符路径

**负载均衡**
- Round Robin（轮询）
- Random（随机）
- Weighted Round Robin（加权平滑轮询）
- Consistent Hash（一致性哈希，基于 CRC32）

**流量治理**
- 令牌桶限流（单机，`golang.org/x/time/rate`）
- Redis 滑动窗口限流（分布式，Lua 脚本）
- Hystrix 风格熔断器（Closed → Open → Half-Open 状态机）

**安全**
- JWT Bearer Token 鉴权中间件
- 路径白名单（支持通配符）

**可观测性**
- Prometheus 指标暴露（QPS、延迟、状态码、在途请求数）
- 结构化日志（zap + lumberjack，支持日志轮转）

**服务发现**
- etcd Watch 机制，动态感知上游节点变化
- Lease + KeepAlive 自动注销

**管理面板**
- Vue 3 Web Dashboard（暗色主题，8 个页面）
- 实时 ECharts 图表（QPS、延迟、状态码分布）
- 中英文国际化（i18n）
- 响应式布局，支持移动端

**终端工具**
- Bubbletea TUI 实时监控仪表盘
- CLI 命令行管理工具

## Architecture

```
┌─────────────────────────────────────────────────┐
│                   Client Request                │
└───────────────────────┬─────────────────────────┘
                        │
                   ┌────▼────┐
                   │ Gateway │
                   │  :8080  │
                   └────┬────┘
                        │
         ┌──────────────┼──────────────┐
         │              │              │
    ┌────▼────┐   ┌─────▼─────┐  ┌────▼────┐
    │JWT Auth │   │Route Match│  │/metrics │
    └────┬────┘   │  (Trie)   │  │/health  │
         │        └─────┬─────┘  └─────────┘
         │              │
    ┌────▼────┐   ┌─────▼─────┐
    │Rate     │   │Circuit    │
    │Limit    │   │Breaker    │
    └────┬────┘   └─────┬─────┘
         │              │
         └──────┬───────┘
                │
         ┌──────▼──────┐
         │Load Balancer│
         │(RR/WRR/Hash)│
         └──────┬──────┘
                │
    ┌───────────┼───────────┐
    │           │           │
┌───▼──┐  ┌────▼──┐  ┌────▼──┐
│:9001 │  │:9002  │  │:9003  │
│User  │  │File   │  │Order  │
│Svc   │  │Svc    │  │Svc    │
└──────┘  └───────┘  └───────┘
```

## Quick Start

### 前置依赖

- Go 1.22+
- Node.js 20+（构建前端）
- etcd（可选，服务发现）
- Redis（可选，分布式限流）

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/Hollow468/Aegis.git
cd Aegis

# 一键构建（前端 + Go 二进制）
make all

# 运行
./aegis
```

### 分步构建

```bash
# 仅构建前端
make frontend

# 仅构建 Go 二进制（需要先构建前端）
go build -o aegis ./cmd/apigateway

# 前端开发模式（热更新，代理到 :8080）
make frontend-dev
```

### Docker

```bash
docker build -t aegis .
docker run -p 8080:8080 aegis
```

### 配置

编辑 `config.yaml`：

```yaml
server:
  port: 8080

routes:
  - path: "/api/users"
    method: "GET"
    match_type: "exact"
    balancer: "round_robin"
    rate_limit:
      type: "token_bucket"
      limit: 100
      burst: 200
    circuit_breaker:
      error_threshold: 0.5
      min_requests: 10
      window_seconds: 60
      recovery_time: 30
    upstreams:
      - address: "http://localhost:9001"
        weight: 100
```

## API Endpoints

| Path | Method | Description |
|------|--------|-------------|
| `/health` | GET | 健康检查 |
| `/metrics` | GET | Prometheus 指标（无需鉴权） |
| `/*` | * | 代理到匹配的上游服务 |

## Tech Stack

**后端**

| Component | Technology |
|-----------|-----------|
| Language | Go 1.22+ |
| HTTP | net/http + httputil.ReverseProxy |
| Config | Viper (YAML) |
| Logging | zap + lumberjack |
| Service Discovery | etcd client v3 |
| Rate Limiting | golang.org/x/time/rate, Redis + Lua |
| Auth | golang-jwt/jwt v5 |
| Metrics | prometheus/client_golang |
| TUI | charmbracelet/bubbletea |
| CLI | charmbracelet/bubbletea |

**前端**

| Component | Technology |
|-----------|-----------|
| Framework | Vue 3 (Composition API) |
| Build Tool | Vite |
| Language | TypeScript |
| State | Pinia |
| Router | Vue Router 4 |
| UI | Element Plus |
| Charts | ECharts + vue-echarts |
| HTTP | Axios |
| i18n | vue-i18n (zh-cn / en) |

## Project Structure

```
.
├── cmd/
│   ├── apigateway/       # 主服务入口
│   ├── aegis-cli/        # CLI 工具
│   └── aegis-tui/        # TUI 仪表盘
├── internal/
│   ├── server/           # HTTP 服务器与中间件编排
│   ├── router/           # Trie 路由引擎
│   ├── proxy/            # 反向代理与负载均衡
│   ├── context/          # Gateway Context 封装
│   ├── middleware/        # JWT、限流、熔断、Metrics 中间件
│   ├── circuit/          # 熔断器状态机
│   ├── limiter/          # 令牌桶 + Redis 限流
│   ├── discovery/        # etcd 服务发现
│   ├── metrics/          # Prometheus 指标定义
│   ├── config/           # 配置加载
│   ├── logger/           # 日志组件
│   ├── model/            # 数据模型
│   ├── web/              # go:embed 前端资源
│   └── tui/              # TUI Dashboard
├── frontend/             # Vue 3 Web 管理面板
│   ├── src/
│   │   ├── api/          # Axios API 层
│   │   ├── stores/       # Pinia 状态管理
│   │   ├── views/        # 页面视图
│   │   ├── components/   # 通用组件 + 图表
│   │   ├── composables/  # 组合式函数
│   │   ├── i18n/         # 国际化 (zh-cn / en)
│   │   ├── router/       # 路由配置
│   │   └── layouts/      # 布局组件
│   └── vite.config.ts
├── web/                  # HTML 原型
├── Makefile              # 构建脚本
├── Dockerfile            # 多阶段容器构建
├── .github/workflows/    # CI/CD 流水线
└── config.yaml           # 示例配置
```

## Development

```bash
# 运行 Go 测试
make test

# 跨平台构建
make build-all

# Go 基准测试
go test -bench=. ./benchmark/
```

## License

MIT
