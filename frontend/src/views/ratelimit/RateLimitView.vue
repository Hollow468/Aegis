<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRateLimitStore } from '../../stores/ratelimit'
import { useConfirm } from '../../composables/useConfirm'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import { ElMessage } from 'element-plus'

const store = useRateLimitStore()
const { visible, title, message, confirm, handleConfirm, handleCancel } = useConfirm()

const showAdd = ref(false)
const form = ref({ route: '', strategy: 'token_bucket', rate: 100, burst: 200 })

onMounted(() => { store.fetchLimits() })

async function handleDelete(route: string) {
  const ok = await confirm('Delete Rule', `Remove rate limit for "${route}"?`)
  if (!ok) return
  try {
    await store.removeLimit(route)
    ElMessage.success('Rule deleted')
  } catch {
    ElMessage.error('Failed to delete rule')
  }
}

async function handleAdd() {
  if (!form.value.route) {
    ElMessage.warning('Route is required')
    return
  }
  try {
    await store.addLimit(form.value)
    showAdd.value = false
    form.value = { route: '', strategy: 'token_bucket', rate: 100, burst: 200 }
    ElMessage.success('Rule added')
  } catch {
    ElMessage.error('Failed to add rule')
  }
}
</script>

<template>
  <div class="ratelimit-page">
    <div class="panel">
      <div class="panel-header">
        <h3>Rate Limiting Rules</h3>
        <button class="btn btn-primary" @click="showAdd = true">+ Add Rule</button>
      </div>
      <div class="panel-body">
        <div v-if="store.loading" class="loading">Loading...</div>
        <table v-else>
          <thead>
            <tr><th>Route</th><th>Strategy</th><th>Rate</th><th>Burst</th><th>Actions</th></tr>
          </thead>
          <tbody>
            <tr v-for="l in store.limits" :key="l.route">
              <td><code>{{ l.route }}</code></td>
              <td>{{ l.strategy }}</td>
              <td>{{ l.rate }}/s</td>
              <td>{{ l.burst }}</td>
              <td><button class="btn btn-danger btn-sm" @click="handleDelete(l.route)">Delete</button></td>
            </tr>
            <tr v-if="store.limits.length === 0">
              <td colspan="5" class="empty">No rate limit rules configured</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <el-dialog v-model="showAdd" title="Add Rate Limit Rule" width="420px">
      <div class="form-grid">
        <div class="form-item">
          <label>Route</label>
          <el-input v-model="form.route" placeholder="/api/users" />
        </div>
        <div class="form-item">
          <label>Strategy</label>
          <el-select v-model="form.strategy">
            <el-option label="Token Bucket" value="token_bucket" />
            <el-option label="Redis Sliding Window" value="redis" />
          </el-select>
        </div>
        <div class="form-item">
          <label>Rate (req/s)</label>
          <el-input-number v-model="form.rate" :min="1" />
        </div>
        <div class="form-item">
          <label>Burst</label>
          <el-input-number v-model="form.burst" :min="1" />
        </div>
      </div>
      <template #footer>
        <el-button @click="showAdd = false">Cancel</el-button>
        <el-button type="primary" :disabled="store.loading" @click="handleAdd">Add</el-button>
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
