<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRateLimitStore } from '../../stores/ratelimit'
import { useConfirm } from '../../composables/useConfirm'
import { useI18n } from 'vue-i18n'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import { ElMessage } from 'element-plus'

const { t } = useI18n()
const store = useRateLimitStore()
const { visible, title, message, confirm, handleConfirm, handleCancel } = useConfirm()

const showAdd = ref(false)
const form = ref({ route: '', strategy: 'token_bucket', rate: 100, burst: 200 })

onMounted(() => { store.fetchLimits() })

async function handleDelete(route: string) {
  const ok = await confirm(t('ratelimit.confirm.title'), t('ratelimit.confirm.message', { route }))
  if (!ok) return
  try {
    await store.removeLimit(route)
    ElMessage.success(t('ratelimit.success.deleted'))
  } catch {
    ElMessage.error(t('ratelimit.error.deleteFailed'))
  }
}

async function handleAdd() {
  if (!form.value.route) {
    ElMessage.warning(t('ratelimit.error.routeRequired'))
    return
  }
  try {
    await store.addLimit(form.value)
    showAdd.value = false
    form.value = { route: '', strategy: 'token_bucket', rate: 100, burst: 200 }
    ElMessage.success(t('ratelimit.success.added'))
  } catch {
    ElMessage.error(t('ratelimit.error.addFailed'))
  }
}
</script>

<template>
  <div class="ratelimit-page">
    <div class="panel">
      <div class="panel-header">
        <h3>{{ t('ratelimit.header') }}</h3>
        <button class="btn btn-primary" @click="showAdd = true">{{ t('ratelimit.addRule') }}</button>
      </div>
      <div class="panel-body">
        <div v-if="store.loading" class="loading">{{ t('common.loading') }}</div>
        <table v-else>
          <thead>
            <tr>
              <th>{{ t('common.route') }}</th>
              <th>{{ t('ratelimit.col.strategy') }}</th>
              <th>{{ t('ratelimit.col.rate') }}</th>
              <th>{{ t('ratelimit.col.burst') }}</th>
              <th>{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="l in store.limits" :key="l.route">
              <td><code>{{ l.route }}</code></td>
              <td>{{ l.strategy }}</td>
              <td>{{ l.rate }}/s</td>
              <td>{{ l.burst }}</td>
              <td><button class="btn btn-danger btn-sm" @click="handleDelete(l.route)">{{ t('common.delete') }}</button></td>
            </tr>
            <tr v-if="store.limits.length === 0">
              <td colspan="5" class="empty">{{ t('ratelimit.empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <el-dialog v-model="showAdd" :title="t('ratelimit.dialog.title')" width="420px">
      <div class="form-grid">
        <div class="form-item">
          <label>{{ t('ratelimit.form.route') }}</label>
          <el-input v-model="form.route" :placeholder="t('ratelimit.form.routePlaceholder')" />
        </div>
        <div class="form-item">
          <label>{{ t('ratelimit.form.strategy') }}</label>
          <el-select v-model="form.strategy">
            <el-option :label="t('ratelimit.form.strategyOptions.tokenBucket')" value="token_bucket" />
            <el-option :label="t('ratelimit.form.strategyOptions.redis')" value="redis" />
          </el-select>
        </div>
        <div class="form-item">
          <label>{{ t('ratelimit.form.rate') }}</label>
          <el-input-number v-model="form.rate" :min="1" />
        </div>
        <div class="form-item">
          <label>{{ t('ratelimit.form.burst') }}</label>
          <el-input-number v-model="form.burst" :min="1" />
        </div>
      </div>
      <template #footer>
        <el-button @click="showAdd = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :disabled="store.loading" @click="handleAdd">{{ t('common.add') }}</el-button>
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
