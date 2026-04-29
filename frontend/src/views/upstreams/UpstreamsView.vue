<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useUpstreamsStore } from '../../stores/upstreams'
import { useConfirm } from '../../composables/useConfirm'
import { useI18n } from 'vue-i18n'
import StatusDot from '../../components/common/StatusDot.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import { ElMessage } from 'element-plus'

const { t } = useI18n()
const upstreamsStore = useUpstreamsStore()
const { visible, title, message, confirm, handleConfirm, handleCancel } = useConfirm()

const showRegister = ref(false)
const regForm = ref({ service: '', address: '', weight: 100 })

onMounted(() => { upstreamsStore.fetchUpstreams() })

async function handleDeregister(service: string, address: string) {
  const ok = await confirm(t('upstreams.confirm.title'), t('upstreams.confirm.message', { address, service }))
  if (!ok) return
  try {
    await upstreamsStore.deregister(service, address)
    ElMessage.success(t('upstreams.success.deregistered'))
  } catch {
    ElMessage.error(t('upstreams.error.deregisterFailed'))
  }
}

async function handleRegister() {
  const { service, address, weight } = regForm.value
  if (!service || !address) {
    ElMessage.warning(t('upstreams.error.required'))
    return
  }
  try {
    await upstreamsStore.register(service, address, weight)
    showRegister.value = false
    regForm.value = { service: '', address: '', weight: 100 }
    ElMessage.success(t('upstreams.success.registered'))
  } catch {
    ElMessage.error(t('upstreams.error.registerFailed'))
  }
}
</script>

<template>
  <div class="upstreams-page">
    <div class="panel">
      <div class="panel-header">
        <h3>{{ t('upstreams.header') }}</h3>
        <button class="btn btn-primary" @click="showRegister = true">{{ t('upstreams.register') }}</button>
      </div>
      <div class="panel-body">
        <div v-if="upstreamsStore.loading" class="loading">{{ t('common.loading') }}</div>
        <table v-else>
          <thead>
            <tr>
              <th>{{ t('common.address') }}</th>
              <th>{{ t('common.service') }}</th>
              <th>{{ t('common.status') }}</th>
              <th>{{ t('common.latency') }}</th>
              <th>{{ t('common.weight') }}</th>
              <th>{{ t('upstreams.col.connections') }}</th>
              <th>{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in upstreamsStore.upstreams" :key="u.address">
              <td><code>{{ u.address }}</code></td>
              <td>{{ u.service }}</td>
              <td><StatusDot :status="u.healthy ? 'healthy' : 'down'" :label="u.healthy ? t('common.healthy') : t('common.down')" /></td>
              <td>{{ u.latency ? `${u.latency}ms` : '-' }}</td>
              <td>{{ u.weight }}</td>
              <td>{{ u.connections }}</td>
              <td><button class="btn btn-danger btn-sm" @click="handleDeregister(u.service, u.address)">{{ t('common.deregister') }}</button></td>
            </tr>
            <tr v-if="upstreamsStore.upstreams.length === 0">
              <td colspan="7" class="empty">{{ t('upstreams.empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <el-dialog v-model="showRegister" :title="t('upstreams.dialog.title')" width="420px">
      <div class="form-grid">
        <div class="form-item">
          <label>{{ t('upstreams.form.serviceName') }}</label>
          <el-input v-model="regForm.service" :placeholder="t('upstreams.form.servicePlaceholder')" />
        </div>
        <div class="form-item">
          <label>{{ t('upstreams.form.address') }}</label>
          <el-input v-model="regForm.address" :placeholder="t('upstreams.form.addressPlaceholder')" />
        </div>
        <div class="form-item">
          <label>{{ t('upstreams.form.weight') }}</label>
          <el-input-number v-model="regForm.weight" :min="1" :max="1000" />
        </div>
      </div>
      <template #footer>
        <el-button @click="showRegister = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :disabled="upstreamsStore.loading" @click="handleRegister">{{ t('common.register') }}</el-button>
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
