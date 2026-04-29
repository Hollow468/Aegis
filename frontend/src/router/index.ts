import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import i18n from '../i18n'

const t = i18n.global.t

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/dashboard/DashboardView.vue'),
        meta: { titleKey: 'route.title.overview' }
      },
      {
        path: 'routes',
        name: 'Routes',
        component: () => import('../views/routes/RoutesView.vue'),
        meta: { titleKey: 'route.title.routes' }
      },
      {
        path: 'upstreams',
        name: 'Upstreams',
        component: () => import('../views/upstreams/UpstreamsView.vue'),
        meta: { titleKey: 'route.title.upstreams' }
      },
      {
        path: 'circuit',
        name: 'CircuitBreakers',
        component: () => import('../views/circuit/CircuitView.vue'),
        meta: { titleKey: 'route.title.circuitBreakers' }
      },
      {
        path: 'ratelimit',
        name: 'RateLimiting',
        component: () => import('../views/ratelimit/RateLimitView.vue'),
        meta: { titleKey: 'route.title.rateLimiting' }
      },
      {
        path: 'metrics',
        name: 'Metrics',
        component: () => import('../views/metrics/MetricsView.vue'),
        meta: { titleKey: 'route.title.metrics' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/settings/SettingsView.vue'),
        meta: { titleKey: 'route.title.settings' }
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/LoginView.vue'),
    meta: { titleKey: 'route.title.login' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  const titleKey = to.meta.titleKey as string
  const title = titleKey ? t(titleKey) : t('route.defaultTitle')
  document.title = t('route.titleTemplate', { title })
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router
