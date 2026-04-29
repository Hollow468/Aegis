import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/dashboard/DashboardView.vue'),
        meta: { title: 'Overview' }
      },
      {
        path: 'routes',
        name: 'Routes',
        component: () => import('../views/routes/RoutesView.vue'),
        meta: { title: 'Routes' }
      },
      {
        path: 'upstreams',
        name: 'Upstreams',
        component: () => import('../views/upstreams/UpstreamsView.vue'),
        meta: { title: 'Upstreams' }
      },
      {
        path: 'circuit',
        name: 'CircuitBreakers',
        component: () => import('../views/circuit/CircuitView.vue'),
        meta: { title: 'Circuit Breakers' }
      },
      {
        path: 'ratelimit',
        name: 'RateLimiting',
        component: () => import('../views/ratelimit/RateLimitView.vue'),
        meta: { title: 'Rate Limiting' }
      },
      {
        path: 'metrics',
        name: 'Metrics',
        component: () => import('../views/metrics/MetricsView.vue'),
        meta: { title: 'Metrics' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/settings/SettingsView.vue'),
        meta: { title: 'Settings' }
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/LoginView.vue'),
    meta: { title: 'Login' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Auth guard
router.beforeEach((to, _from, next) => {
  document.title = `${to.meta.title || 'Aegis'} - Aegis Gateway`
  const token = localStorage.getItem('token')
  if (to.name !== 'Login' && !token) {
    // Allow access without auth for now (development)
    next()
  } else {
    next()
  }
})

export default router
