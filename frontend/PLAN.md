# Aegis Frontend - Development Plan

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Framework | Vue 3 | 3.5+ |
| Build Tool | Vite | 6.x |
| Language | TypeScript | 5.x |
| State Management | Pinia | 3.x |
| Router | Vue Router | 4.x |
| HTTP Client | Axios | 1.x |
| UI Library | Element Plus | 2.x |
| Charts | ECharts + vue-echarts | 5.x |
| CSS | CSS Variables + Scoped Styles |

## Architecture

### Component Architecture (Atomic Design)

```
atoms:     ElButton, ElTag, ElBadge, StatusDot, MetricValue
molecules: MetricCard, RouteTableRow, UpstreamCard, CircuitStateBadge
organisms: RouteTable, UpstreamList, CircuitBreakerPanel, MetricsChart
pages:     DashboardView, RoutesView, UpstreamsView, ...
```

### State Management (Pinia Stores)

```
stores/
├── routes.ts       # Route CRUD, filtering, sorting
├── upstreams.ts    # Upstream health, registration
├── metrics.ts      # Prometheus metrics, QPS, latency
├── circuit.ts      # Circuit breaker states
├── ratelimit.ts    # Rate limit rules
├── auth.ts         # JWT token, user info
└── app.ts          # Global app state (sidebar, theme)
```

### API Layer

```
api/
├── index.ts        # Axios instance, interceptors, error handling
├── routes.ts       # GET/POST/PUT/DELETE /api/routes
├── upstreams.ts    # GET/POST/DELETE /api/upstreams
├── metrics.ts      # GET /api/metrics, /api/metrics/summary
├── circuit.ts      # GET/PUT /api/circuit-breakers
├── ratelimit.ts    # GET/PUT /api/rate-limits
└── auth.ts         # POST /api/auth/login, /api/auth/refresh
```

### Routing Strategy

- Lazy loading for all views
- Auth guard: redirect to /login if no token
- Route meta: title, requiresAuth

## Project Structure

```
frontend/
├── src/
│   ├── api/              # API service layer
│   ├── assets/           # Static assets, global styles
│   │   └── styles/       # CSS variables, theme
│   ├── components/       # Reusable components
│   │   ├── common/       # StatusDot, MetricCard, etc.
│   │   ├── layout/       # Sidebar, Header
│   │   └── charts/       # ECharts wrappers
│   ├── composables/      # Vue composables
│   ├── layouts/          # MainLayout
│   ├── router/           # Vue Router config
│   ├── stores/           # Pinia stores
│   ├── types/            # TypeScript interfaces
│   └── views/            # Page components
│       ├── dashboard/    # Overview dashboard
│       ├── routes/       # Route management
│       ├── upstreams/    # Upstream management
│       ├── circuit/      # Circuit breaker monitoring
│       ├── ratelimit/    # Rate limit config
│       ├── auth/         # Login page
│       ├── metrics/      # Metrics visualization
│       └── settings/     # System settings
├── index.html
├── vite.config.ts
├── tsconfig.json
└── package.json
```

## Development Phases

### Phase 1: Foundation (Current)
- [x] Project scaffold (Vue 3 + Vite + TS)
- [x] Install dependencies
- [x] Router with lazy loading
- [x] Pinia stores skeleton
- [x] Axios API layer
- [x] TypeScript types
- [x] MainLayout (Sidebar + Header)
- [x] CSS variables dark theme

### Phase 2: Core Pages
- [ ] Dashboard overview (metric cards, charts)
- [ ] Routes page (table, CRUD modal)
- [ ] Upstreams page (health status, register)
- [ ] Circuit Breakers page (state visualization)

### Phase 3: Advanced Features
- [ ] Metrics page (ECharts graphs)
- [ ] Rate Limiting page
- [ ] Auth page (JWT login)
- [ ] Settings page

### Phase 4: Polish
- [ ] Responsive design
- [ ] Error handling & loading states
- [ ] Real-time data polling
- [ ] Unit tests

## Design Tokens

```css
--bg-primary: #0f1117
--bg-secondary: #1a1d27
--bg-tertiary: #242736
--border: #2e3148
--text-primary: #e4e6f0
--text-secondary: #8b8fa3
--accent: #3b82f6
--success: #22c55e
--warning: #f59e0b
--danger: #ef4444
```

## Git Workflow

Each phase is committed separately:
1. Foundation → `feat(frontend): project scaffold and core infrastructure`
2. Core Pages → `feat(frontend): dashboard, routes, upstreams pages`
3. Advanced → `feat(frontend): metrics, ratelimit, auth, settings`
4. Polish → `refactor(frontend): responsive design and error handling`
