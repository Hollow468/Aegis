<script setup lang="ts">
import { useMetricsStore } from '../../stores/metrics'
import { useUpstreamsStore } from '../../stores/upstreams'
import { usePolling } from '../../composables/usePolling'
import { useI18n } from 'vue-i18n'
import MetricCard from '../../components/common/MetricCard.vue'
import StatusDot from '../../components/common/StatusDot.vue'

const { t } = useI18n()
const metricsStore = useMetricsStore()
const upstreamsStore = useUpstreamsStore()

async function refresh() {
  await Promise.allSettled([
    metricsStore.fetchSummary(),
    upstreamsStore.fetchUpstreams()
  ])
}

usePolling(refresh, 5000)

const metricCards = [
  { key: 'qps', labelKey: 'dashboard.metric.requestsPerSec', format: (v: number) => v.toLocaleString() },
  { key: 'avg_latency', labelKey: 'dashboard.metric.avgLatency', format: (v: number) => `${v.toFixed(1)}ms` },
  { key: 'total_requests', labelKey: 'dashboard.metric.totalRequests', format: (v: number) => v >= 1e6 ? `${(v / 1e6).toFixed(1)}M` : v.toLocaleString() },
  { key: 'in_flight', labelKey: 'dashboard.metric.inFlight', format: (v: number) => String(v) }
]
</script>

<template>
  <div class="dashboard">
    <div class="metrics-grid">
      <MetricCard
        v-for="m in metricCards"
        :key="m.key"
        :label="t(m.labelKey)"
        :value="m.format((metricsStore.summary as any)[m.key] ?? 0)"
      />
    </div>

    <div class="panel">
      <div class="panel-header"><h3>{{ t('dashboard.upstreamHealth') }}</h3></div>
      <div class="panel-body">
        <div v-if="upstreamsStore.loading" class="loading">{{ t('common.loading') }}</div>
        <table v-else>
          <thead>
            <tr>
              <th>{{ t('common.address') }}</th>
              <th>{{ t('common.service') }}</th>
              <th>{{ t('common.status') }}</th>
              <th>{{ t('common.latency') }}</th>
              <th>{{ t('common.weight') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in upstreamsStore.upstreams" :key="u.address">
              <td><code>{{ u.address }}</code></td>
              <td>{{ u.service }}</td>
              <td><StatusDot :status="u.healthy ? 'healthy' : 'down'" :label="u.healthy ? t('common.healthy') : t('common.down')" /></td>
              <td>{{ u.latency ? `${u.latency}ms` : '-' }}</td>
              <td>{{ u.weight }}</td>
            </tr>
            <tr v-if="upstreamsStore.upstreams.length === 0">
              <td colspan="5" class="empty">{{ t('dashboard.empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.metrics-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 16px; margin-bottom: 24px; }
.panel { background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--radius); }
.panel-header { padding: 16px 20px; border-bottom: 1px solid var(--border); }
.panel-header h3 { font-size: 14px; font-weight: 600; }
.panel-body { padding: 20px; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 12px 16px; text-align: left; font-size: 13px; border-bottom: 1px solid var(--border); }
th { color: var(--text-secondary); font-weight: 600; font-size: 11px; text-transform: uppercase; }
.loading { padding: 20px; text-align: center; color: var(--text-secondary); }
.empty { text-align: center; color: var(--text-muted); padding: 24px; }
</style>
