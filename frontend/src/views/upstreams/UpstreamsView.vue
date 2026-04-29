<script setup lang="ts">
import { ref } from 'vue'

const upstreams = ref([
  { address: 'http://localhost:9001', service: 'user-service', healthy: true, latency: '8.5ms', weight: 100, connections: 12 },
  { address: 'http://localhost:9002', service: 'file-service', healthy: true, latency: '12.1ms', weight: 100, connections: 8 },
  { address: 'http://localhost:9003', service: 'user-service', healthy: false, latency: '-', weight: 100, connections: 0 }
])
</script>

<template>
  <div class="upstreams-page">
    <div class="panel">
      <div class="panel-header">
        <h3>Upstream Servers</h3>
        <button class="btn btn-primary">+ Register</button>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr><th>Address</th><th>Service</th><th>Status</th><th>Latency</th><th>Weight</th><th>Connections</th><th>Actions</th></tr>
          </thead>
          <tbody>
            <tr v-for="u in upstreams" :key="u.address">
              <td><code>{{ u.address }}</code></td>
              <td>{{ u.service }}</td>
              <td>
                <span class="status-indicator">
                  <span class="dot" :class="u.healthy ? 'dot-green' : 'dot-red'"></span>
                  {{ u.healthy ? 'Healthy' : 'Down' }}
                </span>
              </td>
              <td>{{ u.latency }}</td>
              <td>{{ u.weight }}</td>
              <td>{{ u.connections }}</td>
              <td><button class="btn btn-danger btn-sm">Deregister</button></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.panel { background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--radius); }
.panel-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid var(--border); }
.panel-header h3 { font-size: 14px; font-weight: 600; }
.panel-body { padding: 20px; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 12px 16px; text-align: left; font-size: 13px; border-bottom: 1px solid var(--border); }
th { color: var(--text-secondary); font-weight: 600; font-size: 11px; text-transform: uppercase; background: var(--bg-tertiary); }

.status-indicator { display: inline-flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 600; }
.dot { width: 8px; height: 8px; border-radius: 50%; }
.dot-green { background: var(--success); animation: pulse 2s infinite; }
.dot-red { background: var(--danger); animation: pulse 1s infinite; }
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 16px; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer; border: none; }
.btn-primary { background: var(--accent); color: #fff; }
.btn-danger { background: rgba(239,68,68,0.15); color: var(--danger); }
.btn-sm { padding: 4px 10px; font-size: 12px; }
</style>
