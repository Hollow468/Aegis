<script setup lang="ts">
import { ref } from 'vue'

const limits = ref([
  { route: '/api/users', strategy: 'Token Bucket', rate: '100/s', burst: 200, current: '12/s', status: 'ok' },
  { route: '/api/users/{id}', strategy: 'Token Bucket', rate: '50/s', burst: 100, current: '8/s', status: 'ok' },
  { route: '/api/orders', strategy: 'Sliding Window', rate: '200/s', burst: 300, current: '187/s', status: 'warning' },
  { route: '/api/payments', strategy: 'Sliding Window', rate: '30/s', burst: 50, current: '29/s', status: 'warning' }
])

const statusClass: Record<string, string> = { ok: 'status-ok', warning: 'status-warning', exceeded: 'status-exceeded' }
</script>

<template>
  <div class="ratelimit-page">
    <div class="panel">
      <div class="panel-header">
        <h3>Rate Limiting Rules</h3>
        <button class="btn btn-primary">+ Add Rule</button>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr><th>Route</th><th>Strategy</th><th>Rate</th><th>Burst</th><th>Current</th><th>Status</th><th>Actions</th></tr>
          </thead>
          <tbody>
            <tr v-for="l in limits" :key="l.route">
              <td><code>{{ l.route }}</code></td>
              <td>{{ l.strategy }}</td>
              <td>{{ l.rate }}</td>
              <td>{{ l.burst }}</td>
              <td>{{ l.current }}</td>
              <td><span class="status-tag" :class="statusClass[l.status]">{{ l.status.toUpperCase() }}</span></td>
              <td>
                <button class="btn btn-secondary btn-sm">Edit</button>
                <button class="btn btn-danger btn-sm">Delete</button>
              </td>
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

.status-tag { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; }
.status-ok { background: rgba(34,197,94,0.15); color: var(--success); }
.status-warning { background: rgba(245,158,11,0.15); color: var(--warning); }
.status-exceeded { background: rgba(239,68,68,0.15); color: var(--danger); }

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 16px; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer; border: none; }
.btn-primary { background: var(--accent); color: #fff; }
.btn-secondary { background: var(--bg-tertiary); color: var(--text-primary); border: 1px solid var(--border); }
.btn-danger { background: rgba(239,68,68,0.15); color: var(--danger); }
.btn-sm { padding: 4px 10px; font-size: 12px; }
</style>
