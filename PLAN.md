📌 项目愿景 (Project Vision)
打造一个基于云原生架构、高吞吐、低延迟的轻量级 API 网关。实现流量的动态路由、负载均衡、高可用防护（限流/熔断），并提供易用的终端（TUI/CLI）管理能力。

🛠️ 核心技术栈 (Tech Stack)
核心语言： Go (Golang)

网络组件： 原生 net/http / 反向代理 httputil.ReverseProxy (后期可按需演进至 fasthttp)

配置与服务发现： etcd (或 Consul)

RPC 通信： gRPC (用于控制面与数据面通信)

可观测性： Prometheus + Jaeger (链路追踪)

日志系统： Uber zap + lumberjack

🗺️ 研发路线图 (Roadmap)
Phase 0: 架构设计与基础设施 (Week 1)
目标：搭建项目骨架，确立控制面 (Control Plane) 与数据面 (Data Plane) 分离的设计思想。

[x] 初始化 Go Modules 并在 GitHub/GitCode 创建仓库。

[x] 设计配置结构体（路由规则、Upstream、限流配置等）。

[x] 引入 viper 实现本地 YAML 配置文件的加载与热重载。

[x] 搭建标准化日志组件（基于 zap），支持级别动态调整。

Phase 1: 核心路由与反向代理 (Week 2)
目标：让请求能够正确、高效地转发到后端服务。

[x] 基于 httputil.ReverseProxy 实现基础的 HTTP/HTTPS 反向代理。

[x] 实现路由树算法（Trie Tree/基数树），支持精准匹配、前缀匹配和正则匹配。

[x] 封装 Context，为后续的中间件（Plugin）生态打下基础。

[x] 性能基准测试： 编写 wrk/hey 压测脚本，记录 baseline 数据。

Phase 2: 服务发现与负载均衡 (Week 3)
目标：打破单点僵局，支持集群内部服务的动态伸缩。

[x] 接入 etcd，实现网关节点对 etcd 中服务注册信息的 Watch 机制。

[x] 抽象 LoadBalancer 接口。

[x] 实现基础负载均衡算法：Round-Robin（轮询）、Random（随机）。

[x] 实现高级负载均衡算法：Weight-Round-Robin（加权轮询）、Consistent Hash（一致性哈希，确保同一用户的请求落到同一节点）。

Phase 3: 流量治理与安全防护 (Week 4)
目标：这是简历上最硬核的部分，保护后端不被打挂。

[x] 限流 (Rate Limiting)： 基于单机 x/time/rate 实现令牌桶算法；基于 Redis + Lua 脚本实现分布式限流。

[x] 熔断降级 (Circuit Breaking)： 借鉴 Hystrix 的状态机模型（Closed, Open, Half-Open），实现滑动窗口的错误率统计与熔断切换。

[x] 鉴权中心 (Auth)： 开发 JWT 鉴权中间件，实现统一的 Token 拦截与校验。

Phase 4: 极客管理面板与可观测性 (Week 5)
目标：提升项目品味，实现对网关的监控与动态干预。

[ ] 引入 prometheus/client_golang，暴露网关 QPS、延迟统计、HTTP 状态码等 Metrics 接口。

[ ] 开发一个基于 charmbracelet/bubbletea 的 TUI (终端用户界面) Dashboard，在终端中实时监控网关的吞吐量和节点存活状态。

[ ] 支持通过 TUI 或 CLI 动态修改路由规则，实时同步到 etcd，实现网关规则零停机（Zero-Downtime）热更新。

Phase 5: 云原生部署与文档沉淀 (Week 6)
目标：准备面试展示素材，达到 Production Ready 级别。

[ ] 编写 Dockerfile 容器化打包。

[ ] 编写 Kubernetes 部署清单 (Deployment, Service, ConfigMap)，实现网关自身的水平扩容。

[ ] 补充完善 README.md，包含架构图、压测报告（Before vs After 优化）、以及 Quick Start 指南。
