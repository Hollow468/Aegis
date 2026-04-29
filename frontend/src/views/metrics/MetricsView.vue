<script setup lang="ts">
import { ref } from 'vue'

const topRoutes = ref([
  { path: '/api/users', requests: '452K', avgLatency: '12.3ms', p99: '45.1ms', errors: '0.8%' },
  { path: '/api/orders', requests: '312K', avgLatency: '18.7ms', p99: '62.3ms', errors: '2.1%' },
  { path: '/api/files/*', requests: '198K', avgLatency: '34.2ms', p99: '120.5ms', errors: '0.3%' },
  { path: '/api/payments', requests: '87K', avgLatency: '28.4ms', p99: '95.7ms', errors: '1.5%' }
])

const statusCodes = ref([
  { code: '200', count: '892K', pct: '74.3%' },
  { code: '201', count: '156K', pct: '13.0%' },
  { code: '404', count: '89K', pct: '7.4%' },
  { code: '500', count: '42K', pct: '3.5%' },
  { code: '429', count: '21K', pct: '1.8%' }
])
</script>

<template>
  <div class="metrics-page">
    <div class="metrics-grid">
      <div class="panel">
        <div class="panel-header"><h3>Top Routes by Traffic</h3></div>
        <div class="panel-body">
          <table>
            <thead>
              <tr><th>Route</th><th>Requests</th><th>Avg Latency</th><th>P99</th><th>Errors</th></tr>
            </thead>
            <tbody>
              <tr v-for="r in topRoutes" :key="r.path">
                <td><code>{{ r.path }}</code></td>
                <td>{{ r.requests }}</td>
                <td>{{ r.avgLatency }}</td>
                <td>{{ r.p99 }}</td>
                <td>{{ r.errors }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="panel">
        <div class="panel-header"><h3>Status Code Distribution</h3></div>
        <div class="panel-body">
          <div v-for="s in statusCodes" :key="s.code" class="status-row">
            <span class="code">{{ s.code }}</span>
            <div class="bar-container">
              <div class="bar" :style="{ width: s.pct }"></div>
            </div>
            <span class="count">{{ s.count }} ({{ s.pct }})</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.metrics-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.panel { background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--radius); }
.panel-header { padding: 16px 20px; border-bottom: 1px solid var(--border); }
.panel-header h3 { font-size: 14px; font-weight: 600; }
.panel-body { padding: 20px; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 12px 16px; text-align: left; font-size: 13px; border-bottom: 1px solid var(--border); }
th { color: var(--text-secondary); font-weight: 600; font-size: 11px; text-transform: uppercase; background: var(--bg-tertiary); }

.status-row { display: flex; align-items: center; gap: 12px; padding: 10px 0; border-bottom: 1px solid var(--border); }
.code { font-weight: 700; font-size: 14px; width: 40px; }
.bar-container { flex: 1; height: 8px; background: var(--bg-tertiary); border-radius: 4px; overflow: hidden; }
.bar { height: 100%; background: var(--accent); border-radius: 4px; transition: width 0.3s; }
.count { font-size: 12px; color: var(--text-secondary); min-width: 100px; text-align: right; }
</style>
