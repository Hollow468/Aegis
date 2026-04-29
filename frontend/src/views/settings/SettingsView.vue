<script setup lang="ts">
import { onMounted } from 'vue'
import { useSettingsStore } from '../../stores/settings'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const store = useSettingsStore()
onMounted(() => { store.fetchSettings() })
</script>

<template>
  <div class="settings-page">
    <div v-if="store.loading" class="loading">{{ t('common.loading') }}</div>
    <template v-else-if="store.settings">
      <div class="panel">
        <div class="panel-header"><h3>{{ t('settings.serverConfig') }}</h3></div>
        <div class="panel-body">
          <div class="form-grid">
            <div class="form-item"><label>{{ t('settings.form.port') }}</label><input type="text" :value="store.settings.server.port" readonly /></div>
            <div class="form-item"><label>{{ t('settings.form.host') }}</label><input type="text" :value="store.settings.server.host" readonly /></div>
            <div class="form-item"><label>{{ t('settings.form.readTimeout') }}</label><input type="text" :value="store.settings.server.readTimeout" readonly /></div>
            <div class="form-item"><label>{{ t('settings.form.writeTimeout') }}</label><input type="text" :value="store.settings.server.writeTimeout" readonly /></div>
          </div>
        </div>
      </div>

      <div class="panel" style="margin-top:16px;">
        <div class="panel-header"><h3>{{ t('settings.externalServices') }}</h3></div>
        <div class="panel-body">
          <div class="form-grid">
            <div class="form-item"><label>{{ t('settings.form.etcdEndpoints') }}</label><input type="text" :value="store.settings.etcd.endpoints" readonly /></div>
            <div class="form-item"><label>{{ t('settings.form.etcdPrefix') }}</label><input type="text" :value="store.settings.etcd.prefix" readonly /></div>
            <div class="form-item"><label>{{ t('settings.form.redisAddress') }}</label><input type="text" :value="store.settings.redis.addr" readonly /></div>
            <div class="form-item"><label>{{ t('settings.form.redisDb') }}</label><input type="text" :value="store.settings.redis.db" readonly /></div>
          </div>
        </div>
      </div>

      <div class="panel" style="margin-top:16px;">
        <div class="panel-header"><h3>{{ t('settings.authLogging') }}</h3></div>
        <div class="panel-body">
          <div class="form-grid">
            <div class="form-item"><label>{{ t('settings.form.tokenExpiry') }}</label><input type="text" :value="store.settings.jwt.expire" readonly /></div>
            <div class="form-item"><label>{{ t('settings.form.logLevel') }}</label><input type="text" :value="store.settings.log.level" readonly /></div>
            <div class="form-item"><label>{{ t('settings.form.logFormat') }}</label><input type="text" :value="store.settings.log.format" readonly /></div>
          </div>
        </div>
      </div>
    </template>
    <div v-else class="empty">{{ t('settings.error.loadFailed') }}</div>
  </div>
</template>

<style scoped>
.panel { background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--radius); }
.panel-header { padding: 16px 20px; border-bottom: 1px solid var(--border); }
.panel-header h3 { font-size: 14px; font-weight: 600; }
.panel-body { padding: 20px; }
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.form-item { display: flex; flex-direction: column; gap: 6px; }
.form-item label { font-size: 12px; font-weight: 600; color: var(--text-secondary); text-transform: uppercase; }
.form-item input {
  padding: 10px 12px; border-radius: 6px; border: 1px solid var(--border);
  background: var(--bg-tertiary); color: var(--text-primary); font-size: 13px;
}
.loading { padding: 40px; text-align: center; color: var(--text-secondary); }
.empty { padding: 40px; text-align: center; color: var(--text-muted); }
</style>
