<script setup lang="ts">
import { onMounted } from 'vue'
import { useCircuitStore } from '../../stores/circuit'
import { useI18n } from 'vue-i18n'
import StateTag from '../../components/common/StateTag.vue'
import { ElMessage } from 'element-plus'

const { t } = useI18n()
const store = useCircuitStore()

onMounted(() => { store.fetchBreakers() })

async function handleReset(route: string) {
  try {
    await store.resetBreaker(route)
    ElMessage.success(t('circuit.success.reset'))
  } catch {
    ElMessage.error(t('circuit.error.resetFailed'))
  }
}
</script>

<template>
  <div class="circuit-page">
    <div class="panel">
      <div class="panel-header"><h3>{{ t('circuit.header') }}</h3></div>
      <div class="panel-body">
        <div v-if="store.loading" class="loading">{{ t('common.loading') }}</div>
        <table v-else>
          <thead>
            <tr>
              <th>{{ t('common.route') }}</th>
              <th>{{ t('circuit.col.state') }}</th>
              <th>{{ t('circuit.col.errorRate') }}</th>
              <th>{{ t('circuit.col.minRequests') }}</th>
              <th>{{ t('circuit.col.window') }}</th>
              <th>{{ t('circuit.col.recovery') }}</th>
              <th>{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="b in store.breakers" :key="b.route">
              <td><code>{{ b.route }}</code></td>
              <td><StateTag :state="b.state" variant="circuit" /></td>
              <td>{{ (b.error_rate * 100).toFixed(1) }}%</td>
              <td>{{ b.min_requests }}</td>
              <td>{{ b.window_seconds }}s</td>
              <td>{{ b.recovery_time }}s</td>
              <td><button class="btn btn-secondary btn-sm" @click="handleReset(b.route)">{{ t('common.reset') }}</button></td>
            </tr>
            <tr v-if="store.breakers.length === 0">
              <td colspan="7" class="empty">{{ t('circuit.empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="panel" style="margin-top:16px;">
      <div class="panel-header"><h3>{{ t('circuit.stateMachine') }}</h3></div>
      <div class="panel-body state-machine">
        <div class="state-box"><div class="cb-state cb-closed">{{ t('circuit.state.closed') }}</div><div class="desc">{{ t('circuit.state.closedDesc') }}</div></div>
        <div class="arrow">&rarr;</div>
        <div class="state-box"><div class="cb-state cb-open">{{ t('circuit.state.open') }}</div><div class="desc">{{ t('circuit.state.openDesc') }}</div></div>
        <div class="arrow">&rarr;</div>
        <div class="state-box"><div class="cb-state cb-half-open">{{ t('circuit.state.halfOpen') }}</div><div class="desc">{{ t('circuit.state.halfOpenDesc') }}</div></div>
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
