<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useUpstreamsStore } from '../../stores/upstreams'
import { useConfirm } from '../../composables/useConfirm'
import StatusDot from '../../components/common/StatusDot.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import { ElMessage } from 'element-plus'

const upstreamsStore = useUpstreamsStore()
const { visible, title, message, confirm, handleConfirm, handleCancel } = useConfirm()

const showRegister = ref(false)
const regForm = ref({ service: '', address: '', weight: 100 })

onMounted(() => { upstreamsStore.fetchUpstreams() })

async function handleDeregister(service: string, address: string) {
  const ok = await confirm('Deregister Upstream', `Remove ${address} from ${service}?`)
  if (!ok) return
  try {
    await upstreamsStore.deregister(service, address)
    ElMessage.success('Upstream deregistered')
  } catch {
    ElMessage.error('Failed to deregister')
  }
}

async function handleRegister() {
  const { service, address, weight } = regForm.value
  if (!service || !address) {
    ElMessage.warning('Service and address are required')
    return
  }
  try {
    await upstreamsStore.register(service, address, weight)
    showRegister.value = false
    regForm.value = { service: '', address: '', weight: 100 }
    ElMessage.success('Upstream registered')
  } catch {
    ElMessage.error('Failed to register')
  }
}
</script>

<template>
  <div class="upstreams-page">
    <div class="panel">
      <div class="panel-header">
        <h3>Upstream Servers</h3>
        <button class="btn btn-primary" @click="showRegister = true">+ Register</button>
      </div>
      <div class="panel-body">
        <div v-if="upstreamsStore.loading" class="loading">Loading...</div>
        <table v-else>
          <thead>
            <tr><th>Address</th><th>Service</th><th>Status</th><th>Latency</th><th>Weight</th><th>Connections</th><th>Actions</th></tr>
          </thead>
          <tbody>
            <tr v-for="u in upstreamsStore.upstreams" :key="u.address">
              <td><code>{{ u.address }}</code></td>
              <td>{{ u.service }}</td>
              <td><StatusDot :status="u.healthy ? 'healthy' : 'down'" :label="u.healthy ? 'Healthy' : 'Down'" /></td>
              <td>{{ u.latency ? `${u.latency}ms` : '-' }}</td>
              <td>{{ u.weight }}</td>
              <td>{{ u.connections }}</td>
              <td><button class="btn btn-danger btn-sm" @click="handleDeregister(u.service, u.address)">Deregister</button></td>
            </tr>
            <tr v-if="upstreamsStore.upstreams.length === 0">
              <td colspan="7" class="empty">No upstreams registered</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Register Dialog -->
    <el-dialog v-model="showRegister" title="Register Upstream" width="420px">
      <div class="form-grid">
        <div class="form-item">
          <label>Service Name</label>
          <el-input v-model="regForm.service" placeholder="user-service" />
        </div>
        <div class="form-item">
          <label>Address</label>
          <el-input v-model="regForm.address" placeholder="http://localhost:9001" />
        </div>
        <div class="form-item">
          <label>Weight</label>
          <el-input-number v-model="regForm.weight" :min="1" :max="1000" />
        </div>
      </div>
      <template #footer>
        <el-button @click="showRegister = false">Cancel</el-button>
        <el-button type="primary" :disabled="upstreamsStore.loading" @click="handleRegister">Register</el-button>
      </template>
    </el-dialog>

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
.form-grid { display: flex; flex-direction: column; gap: 16px; }
.form-item label { display: block; font-size: 12px; font-weight: 600; color: var(--text-secondary); margin-bottom: 6px; }
.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 16px; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer; border: none; }
.btn-primary { background: var(--accent); color: #fff; }
.btn-danger { background: rgba(239,68,68,0.15); color: var(--danger); }
.btn-sm { padding: 4px 10px; font-size: 12px; }
</style>
