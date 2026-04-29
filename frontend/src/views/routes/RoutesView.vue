<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoutesStore } from '../../stores/routes'
import { useConfirm } from '../../composables/useConfirm'
import StateTag from '../../components/common/StateTag.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import { ElMessage } from 'element-plus'

const routesStore = useRoutesStore()
const { visible, title, message, confirm, handleConfirm, handleCancel } = useConfirm()

onMounted(() => { routesStore.fetchRoutes() })

async function handleDelete(path: string) {
  const ok = await confirm('Delete Route', `Are you sure you want to delete route "${path}"?`)
  if (!ok) return
  try {
    await routesStore.removeRoute(path)
    ElMessage.success('Route deleted')
  } catch {
    ElMessage.error('Failed to delete route')
  }
}
</script>

<template>
  <div class="routes-page">
    <div class="panel">
      <div class="panel-header">
        <h3>Route Configuration</h3>
        <button class="btn btn-primary">+ Add Route</button>
      </div>
      <div class="panel-body">
        <div v-if="routesStore.loading" class="loading">Loading...</div>
        <table v-else>
          <thead>
            <tr><th>Path</th><th>Method</th><th>Match Type</th><th>Balancer</th><th>Rate Limit</th><th>Upstreams</th><th>Actions</th></tr>
          </thead>
          <tbody>
            <tr v-for="r in routesStore.routes" :key="r.path">
              <td><code>{{ r.path }}</code></td>
              <td><StateTag :state="r.method" variant="method" /></td>
              <td><StateTag :state="r.match_type" variant="match" /></td>
              <td>{{ r.balancer }}</td>
              <td>{{ r.rate_limit ? `${r.rate_limit.limit}/s` : '-' }}</td>
              <td>{{ r.upstreams?.length ?? 0 }}</td>
              <td>
                <button class="btn btn-secondary btn-sm">Edit</button>
                <button class="btn btn-danger btn-sm" @click="handleDelete(r.path)">Delete</button>
              </td>
            </tr>
            <tr v-if="routesStore.routes.length === 0">
              <td colspan="7" class="empty">No routes configured</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <ConfirmDialog :visible="visible" :title="title" :message="message" @confirm="handleConfirm" @cancel="handleCancel" />
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
.loading { padding: 20px; text-align: center; color: var(--text-secondary); }
.empty { text-align: center; color: var(--text-muted); padding: 24px; }
.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 16px; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer; border: none; }
.btn-primary { background: var(--accent); color: #fff; }
.btn-secondary { background: var(--bg-tertiary); color: var(--text-primary); border: 1px solid var(--border); }
.btn-danger { background: rgba(239,68,68,0.15); color: var(--danger); }
.btn-sm { padding: 4px 10px; font-size: 12px; }
</style>
