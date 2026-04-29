<script setup lang="ts">
import { useMetricsStore } from '../../stores/metrics'
import { usePolling } from '../../composables/usePolling'
import { useI18n } from 'vue-i18n'
import QpsChart from '../../components/charts/QpsChart.vue'
import LatencyChart from '../../components/charts/LatencyChart.vue'
import StatusCodeChart from '../../components/charts/StatusCodeChart.vue'

const { t } = useI18n()
const store = useMetricsStore()
usePolling(() => store.fetchSummary(), 2000)
</script>

<template>
  <div class="metrics-page">
    <div v-if="store.loading && !store.summary.qps" class="loading">{{ t('common.loading') }}</div>
    <template v-else>
      <div class="charts-grid">
        <div class="panel">
          <div class="panel-header"><h3>{{ t('metrics.requestsPerSec') }}</h3></div>
          <div class="panel-body">
            <QpsChart :value="store.summary.qps" />
          </div>
        </div>
        <div class="panel">
          <div class="panel-header"><h3>{{ t('metrics.latency') }}</h3></div>
          <div class="panel-body">
            <LatencyChart :avg="store.summary.avg_latency" :p99="store.summary.p99_latency" />
          </div>
        </div>
      </div>

      <div class="panel" style="margin-top:16px;">
        <div class="panel-header"><h3>{{ t('metrics.statusDistribution') }}</h3></div>
        <div class="panel-body">
          <StatusCodeChart :s2xx="store.summary.status_2xx" :s4xx="store.summary.status_4xx" :s5xx="store.summary.status_5xx" />
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.charts-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.panel { background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--radius); }
.panel-header { padding: 16px 20px; border-bottom: 1px solid var(--border); }
.panel-header h3 { font-size: 14px; font-weight: 600; }
.panel-body { padding: 16px; }
.loading { padding: 40px; text-align: center; color: var(--text-secondary); }
</style>
