<script setup lang="ts">
import { ref } from 'vue'

const breakers = ref([
  { route: '/api/users', state: 'closed', error_rate: '2.1%', min_requests: 10, window: '60s', recovery: '30s' },
  { route: '/api/users/{id}', state: 'closed', error_rate: '1.8%', min_requests: 20, window: '30s', recovery: '15s' },
  { route: '/api/orders', state: 'open', error_rate: '67.3%', min_requests: 10, window: '60s', recovery: '30s' },
  { route: '/api/payments', state: 'half-open', error_rate: '45.2%', min_requests: 10, window: '60s', recovery: '30s' }
])

const stateClass: Record<string, string> = { closed: 'cb-closed', open: 'cb-open', 'half-open': 'cb-half-open' }
</script>

<template>
  <div class="circuit-page">
    <div class="panel">
      <div class="panel-header"><h3>Circuit Breaker States</h3></div>
      <div class="panel-body">
        <table>
          <thead>
            <tr><th>Route</th><th>State</th><th>Error Rate</th><th>Min Requests</th><th>Window</th><th>Recovery</th><th>Actions</th></tr>
          </thead>
          <tbody>
            <tr v-for="b in breakers" :key="b.route">
              <td><code>{{ b.route }}</code></td>
              <td><span class="cb-state" :class="stateClass[b.state]">{{ b.state.toUpperCase() }}</span></td>
              <td>{{ b.error_rate }}</td>
              <td>{{ b.min_requests }}</td>
              <td>{{ b.window }}</td>
              <td>{{ b.recovery }}</td>
              <td><button class="btn btn-secondary btn-sm">Reset</button></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="panel" style="margin-top:16px;">
      <div class="panel-header"><h3>State Machine</h3></div>
      <div class="panel-body state-machine">
        <div class="state-box"><div class="cb-state cb-closed">CLOSED</div><div class="desc">Normal operation</div></div>
        <div class="arrow">→</div>
        <div class="state-box"><div class="cb-state cb-open">OPEN</div><div class="desc">Rejecting requests</div></div>
        <div class="arrow">→</div>
        <div class="state-box"><div class="cb-state cb-half-open">HALF-OPEN</div><div class="desc">Testing recovery</div></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.panel { background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--radius); }
.panel-header { padding: 16px 20px; border-bottom: 1px solid var(--border); }
.panel-header h3 { font-size: 14px; font-weight: 600; }
.panel-body { padding: 20px; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 12px 16px; text-align: left; font-size: 13px; border-bottom: 1px solid var(--border); }
th { color: var(--text-secondary); font-weight: 600; font-size: 11px; text-transform: uppercase; background: var(--bg-tertiary); }

.cb-state { display: inline-block; padding: 6px 14px; border-radius: 6px; font-size: 12px; font-weight: 700; letter-spacing: 0.5px; }
.cb-closed { background: rgba(34,197,94,0.15); color: var(--success); border: 1px solid rgba(34,197,94,0.3); }
.cb-open { background: rgba(239,68,68,0.15); color: var(--danger); border: 1px solid rgba(239,68,68,0.3); }
.cb-half-open { background: rgba(245,158,11,0.15); color: var(--warning); border: 1px solid rgba(245,158,11,0.3); }

.state-machine { display: flex; align-items: center; justify-content: center; gap: 40px; padding: 40px; }
.state-box { text-align: center; }
.desc { font-size: 11px; color: var(--text-muted); margin-top: 8px; }
.arrow { font-size: 24px; color: var(--text-muted); }

.btn { display: inline-flex; padding: 8px 16px; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer; border: none; }
.btn-secondary { background: var(--bg-tertiary); color: var(--text-primary); border: 1px solid var(--border); }
.btn-sm { padding: 4px 10px; font-size: 12px; }
</style>
