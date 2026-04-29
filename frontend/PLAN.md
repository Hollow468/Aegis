# Aegis Frontend — Development Plan

## Current Status

| Layer | Status | Detail |
|-------|--------|--------|
| API client | ✅ Done | Axios + interceptors, 5 endpoint files |
| Stores | ✅ Done | 5 stores wired to API, loading states |
| Layout | ✅ Done | Sidebar + Header, router-aware |
| Views | ❌ Static mockups | 8 views, all hardcoded ref() data, zero store integration |
| Common components | ❌ Missing | No shared table/modal/form/button |
| Composables | ❌ Missing | No reuse layer |
| Charts | ❌ Missing | No ECharts integration |

**Core problem:** The data layer (API + stores) is complete, but views are wireframes that don't use it.

## Backend API Gap

Current backend only exposes `/metrics` (Prometheus) and `/health`. The frontend API layer targets admin endpoints (`/api/routes`, `/api/upstreams`, etc.) that **do not exist yet**. Strategy:

- Phase 1-2: Build frontend with **mock API fallback** — views call stores, stores call API, but if API fails, show graceful empty state
- Phase 3+: Backend admin API can be added independently; frontend is ready to consume it

## Phase 1: Shared Components & Composables

Extract reusable pieces from duplicated view code.

### 1.1 Common Components

```
src/components/common/
├── DataTable.vue       # Generic table with header slot, row slot, loading, empty state
├── StatusDot.vue       # Animated green/red/yellow dot
├── StateTag.vue        # Colored tag (CLOSED/OPEN/HALF-OPEN, GET/POST, etc.)
├── MetricCard.vue      # Label + value + trend
├── ConfirmDialog.vue   # Simple confirm/cancel modal
└── PageHeader.vue      # Title + action button slot
```

**DataTable.vue** — the most impactful extraction. All 6 table views share identical `<table>` markup with `<thead>`/`<tbody>`, border styles, and loading/empty patterns. One component with `<slot name="row">` eliminates ~200 lines of duplication.

**Props:**
```
DataTable: { loading: boolean, emptyText: string }
  slot #headers — <th> elements
  slot #row="{ item }" — row template
  slot #empty — empty state
```

### 1.2 Composables

```
src/composables/
├── usePolling.ts       # setInterval with auto-cleanup, configurable interval
└── useConfirm.ts       # Confirm dialog state management
```

**usePolling(fn, interval)** — wraps `onMounted`/`onUnmounted` with `setInterval`. Every dashboard-style view needs periodic refresh. Returns `{ start, stop }`.

### 1.3 Commit
`feat(frontend): extract shared components and composables`

---

## Phase 2: Wire Views to Stores (Core Pages)

Replace hardcoded `ref()` data with store calls. Each view becomes: `onMounted → store.fetchX()` + template reads from store.

### 2.1 LoginView → authStore

- Call `authStore.login(username, password)`
- On success: `router.push('/')`
- On error: show error message
- Wire auth guard in router to check `authStore.isAuthenticated`

### 2.2 DashboardView → metricsStore + upstreamsStore

- `onMounted`: call `metricsStore.fetchSummary()` + `upstreamsStore.fetchUpstreams()`
- Metric cards read from `metricsStore.summary` (qps, avg_latency, total_requests, in_flight)
- Upstream health table reads from `upstreamsStore.upstreams`
- Add `usePolling(5000)` for auto-refresh every 5s

### 2.3 RoutesView → routesStore

- `onMounted`: call `routesStore.fetchRoutes()`
- Table reads from `routesStore.routes`
- Delete button calls `routesStore.removeRoute(path)` with confirm dialog
- Add Route button opens a form (can be a simple dialog or inline form initially)

### 2.4 UpstreamsView → upstreamsStore

- `onMounted`: call `upstreamsStore.fetchUpstreams()`
- Table reads from `upstreamsStore.upstreams`
- Deregister button calls `upstreamsStore.deregister(service, address)` with confirm
- Register button opens dialog: service name, address, weight fields

### 2.5 Commit
`feat(frontend): wire core views to stores — login, dashboard, routes, upstreams`

---

## Phase 3: Wire Views to Stores (Advanced Pages)

### 3.1 CircuitView — needs new store + API

Current: no circuit store or API file. Create:

```
src/stores/circuit.ts    # useCircuitStore — fetchBreakers(), resetBreaker(route)
src/api/circuit.ts       # GET /api/circuit-breakers, PUT /api/circuit-breakers/:route/reset
```

- Table reads from store
- Reset button calls `store.resetBreaker(route)`
- State machine diagram stays as-is (static visualization, no API needed)

### 3.2 RateLimitView — needs new store + API

Current: no ratelimit store or API file. Create:

```
src/stores/ratelimit.ts  # useRateLimitStore — fetchLimits(), addLimit(), removeLimit()
src/api/ratelimit.ts     # GET/POST/DELETE /api/rate-limits
```

- Table reads from store
- Add/Delete buttons wired to store actions

### 3.3 MetricsView — metricsStore + ECharts

- Replace hardcoded top routes with real data from `metricsStore`
- Status code distribution: parse from `/metrics` raw Prometheus text, or add a backend endpoint
- Install `echarts` + `vue-echarts` (already in package.json)

Create chart components:
```
src/components/charts/
├── QpsChart.vue          # Line chart — requests/sec over time
├── LatencyChart.vue      # Line chart — avg/p99 latency
└── StatusCodeChart.vue   # Pie/bar chart — 2xx/4xx/5xx distribution
```

- QpsChart and LatencyChart need **time-series data** — currently the backend only exposes instantaneous metrics. Two options:
  - **Option A (simple):** Frontend polls every 2s, accumulates last N data points in memory, renders line chart
  - **Option B (proper):** Backend exposes `/api/metrics/history` with time-bucketed data
  - **Recommendation:** Start with Option A (pure frontend accumulation), upgrade to Option B later

### 3.4 SettingsView — needs new API

```
src/api/settings.ts      # GET /api/settings, PUT /api/settings
src/stores/settings.ts   # useSettingsStore
```

- Currently read-only display; wire to `settingsStore.fetchSettings()`
- Edit mode toggle with save/cancel (future enhancement)

### 3.5 Commit
`feat(frontend): wire advanced views — circuit, ratelimit, metrics charts, settings`

---

## Phase 4: Polish & UX

### 4.1 Loading & Error States

- Every view shows skeleton/spinner during `store.loading`
- API errors show toast notification (use Element Plus `ElMessage`)
- Empty states with icon + message when data is empty

### 4.2 Responsive Design

- Sidebar collapses to icons on < 768px
- Tables scroll horizontally on mobile
- Metric cards stack vertically on small screens

### 4.3 Real-time Updates

- Dashboard: polling every 5s via `usePolling`
- Metrics charts: polling every 2s, accumulate data points
- Circuit/Routes: manual refresh button (no polling — these change infrequently)

### 4.4 Header Enhancements

- Fetch real server status from `/health` endpoint
- Show actual port from config (if settings API exists)
- Logout button that calls `authStore.logout()` + redirects to `/login`

### 4.5 Commit
`refactor(frontend): loading states, responsive design, real-time polling`

---

## Execution Order

```
Phase 1 → Phase 2 → Phase 3 → Phase 4
  │           │           │           │
  │           │           │           └─ responsive, polish, UX
  │           │           └─ circuit/ratelimit stores, ECharts, settings
  │           └─ wire 4 core views to existing stores
  └─ extract DataTable, StatusDot, composables
```

Each phase is one git commit. Estimated scope:

| Phase | Files changed | New files | Key deliverable |
|-------|--------------|-----------|-----------------|
| 1 | 0 | 8 | DataTable, StatusDot, StateTag, MetricCard, ConfirmDialog, PageHeader, usePolling, useConfirm |
| 2 | 4 views + 1 store | 0 | All core views show real data |
| 3 | 2 views | 6 | circuit/ratelimit stores+API, ECharts, settings |
| 4 | 8+ views | 0 | Loading, responsive, polling, polish |

## Notes

- Element Plus is in `package.json` but not yet imported in `main.ts`. Add `import ElementPlus from 'element-plus'` + `app.use(ElementPlus)` in Phase 1.
- The `src/components/common/` and `src/components/charts/` directories need to be created.
- Sidebar badges (route count, upstream count) should read from stores, not hardcoded.
