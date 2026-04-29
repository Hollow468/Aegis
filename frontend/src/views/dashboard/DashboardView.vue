<script setup lang="ts">
import { ref } from 'vue'

const metrics = ref([
  { label: 'Requests / sec', value: '1,247', trend: '▲ 12.3%', trendType: 'up' },
  { label: 'Avg Latency', value: '12.4ms', trend: '▼ 3.1%', trendType: 'down' },
  { label: 'Total Requests', value: '1.2M', trend: '▲ 8.7%', trendType: 'up' },
  { label: 'In Flight', value: '23', trend: 'stable', trendType: 'neutral' }
])

const upstreams = ref([
  { address: 'http://localhost:9001', healthy: true, latency: '8.5ms', weight: 100 },
  { address: 'http://localhost:9002', healthy: true, latency: '12.1ms', weight: 100 },
  { address: 'http://localhost:9003', healthy: false, latency: '-', weight: 100 }
])
</script>

<template>
  <div class="dashboard">
    <div class="metrics-grid">
      <div v-for="m in metrics" :key="m.label" class="metric-card">
        <div class="label">{{ m.label }}</div>
        <div class="value">{{ m.value }}</div>
        <div class="trend" :class="m.trendType">{{ m.trend }}</div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-header">
        <h3>Upstream Health</h3>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr><th>Address</th><th>Status</th><th>Latency</th><th>Weight</th></tr>
          </thead>
          <tbody>
            <tr v-for="u in upstreams" :key="u.address">
              <td><code>{{ u.address }}</code></td>
              <td>
                <span class="status-indicator">
                  <span class="dot" :class="u.healthy ? 'dot-green' : 'dot-red'"></span>
                  {{ u.healthy ? 'Healthy' : 'Down' }}
                </span>
              </td>
              <td>{{ u.latency }}</td>
              <td>{{ u.weight }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.metric-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 20px;
}

.metric-card .label {
  font-size: 12px;
  color: var(--text-secondary);
  text-transform: uppercase;
  margin-bottom: 8px;
}

.metric-card .value { font-size: 28px; font-weight: 700; }

.metric-card .trend { font-size: 12px; }
.trend.up { color: var(--success); }
.trend.down { color: var(--danger); }
.trend.neutral { color: var(--text-secondary); }

.panel {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.panel-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}

.panel-header h3 { font-size: 14px; font-weight: 600; }
.panel-body { padding: 20px; }

table { width: 100%; border-collapse: collapse; }
th, td { padding: 12px 16px; text-align: left; font-size: 13px; border-bottom: 1px solid var(--border); }
th { color: var(--text-secondary); font-weight: 600; font-size: 11px; text-transform: uppercase; }

.status-indicator { display: inline-flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 600; }
.dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-green { background: var(--success); animation: pulse 2s infinite; }
.dot-red { background: var(--danger); animation: pulse 1s infinite; }

@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
</style>
