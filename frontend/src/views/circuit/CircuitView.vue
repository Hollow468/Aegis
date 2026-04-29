<script setup lang="ts">
import { onMounted } from 'vue'
import { useCircuitStore } from '../../stores/circuit'
import StateTag from '../../components/common/StateTag.vue'
import { ElMessage } from 'element-plus'

const store = useCircuitStore()

onMounted(() => { store.fetchBreakers() })

async function handleReset(route: string) {
  try {
    await store.resetBreaker(route)
    ElMessage.success('Circuit breaker reset')
  } catch {
    ElMessage.error('Failed to reset')
  }
}
</script>

<template>
  <div class="circuit-page">
    <div class="panel">
      <div class="panel-header"><h3>Circuit Breaker States</h3></div>
      <div class="panel-body">
        <div v-if="store.loading" class="loading">Loading...</div>
        <table v-else>
          <thead>
            <tr><th>Route</th><th>State</th><th>Error Rate</th><th>Min Requests</th><th>Window</th><th>Recovery</th><th>Actions</th></tr>
          </thead>
          <tbody>
            <tr v-for="b in store.breakers" :key="b.route">
              <td><code>{{ b.route }}</code></td>
              <td><StateTag :state="b.state" variant="circuit" /></td>
              <td>{{ (b.error_rate * 100).toFixed(1) }}%</td>
              <td>{{ b.min_requests }}</td>
              <td>{{ b.window_seconds }}s</td>
              <td>{{ b.recovery_time }}s</td>
              <td><button class="btn btn-secondary btn-sm" @click="handleReset(b.route)">Reset</button></td>
            </tr>
            <tr v-if="store.breakers.length === 0">
              <td colspan="7" class="empty">No circuit breakers configured</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="panel" style="margin-top:16px;">
      <div class="panel-header"><h3>State Machine</h3></div>
      <div class="panel-body state-machine">
        <div class="state-box"><div class="cb-state cb-closed">CLOSED</div><div class="desc">Normal operation</div></div>
        <div class="arrow">&rarr;</div>
        <div class="state-box"><div class="cb-state cb-open">OPEN</div><div class="desc">Rejecting requests</div></div>
        <div class="arrow">&rarr;</div>
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
.loading { padding: 20px; text-align: center; color: var(--text-secondary); }
.empty { text-align: center; color: var(--text-muted); padding: 24px; }

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
