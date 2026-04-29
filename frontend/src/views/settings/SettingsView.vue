<script setup lang="ts">
import { onMounted } from 'vue'
import { useSettingsStore } from '../../stores/settings'

const store = useSettingsStore()
onMounted(() => { store.fetchSettings() })
</script>

<template>
  <div class="settings-page">
    <div v-if="store.loading" class="loading">Loading...</div>
    <template v-else-if="store.settings">
      <div class="panel">
        <div class="panel-header"><h3>Server Configuration</h3></div>
        <div class="panel-body">
          <div class="form-grid">
            <div class="form-item"><label>Port</label><input type="text" :value="store.settings.server.port" readonly /></div>
            <div class="form-item"><label>Host</label><input type="text" :value="store.settings.server.host" readonly /></div>
            <div class="form-item"><label>Read Timeout</label><input type="text" :value="store.settings.server.readTimeout" readonly /></div>
            <div class="form-item"><label>Write Timeout</label><input type="text" :value="store.settings.server.writeTimeout" readonly /></div>
          </div>
        </div>
      </div>

      <div class="panel" style="margin-top:16px;">
        <div class="panel-header"><h3>External Services</h3></div>
        <div class="panel-body">
          <div class="form-grid">
            <div class="form-item"><label>etcd Endpoints</label><input type="text" :value="store.settings.etcd.endpoints" readonly /></div>
            <div class="form-item"><label>etcd Prefix</label><input type="text" :value="store.settings.etcd.prefix" readonly /></div>
            <div class="form-item"><label>Redis Address</label><input type="text" :value="store.settings.redis.addr" readonly /></div>
            <div class="form-item"><label>Redis DB</label><input type="text" :value="store.settings.redis.db" readonly /></div>
          </div>
        </div>
      </div>

      <div class="panel" style="margin-top:16px;">
        <div class="panel-header"><h3>Authentication & Logging</h3></div>
        <div class="panel-body">
          <div class="form-grid">
            <div class="form-item"><label>Token Expiry</label><input type="text" :value="store.settings.jwt.expire" readonly /></div>
            <div class="form-item"><label>Log Level</label><input type="text" :value="store.settings.log.level" readonly /></div>
            <div class="form-item"><label>Log Format</label><input type="text" :value="store.settings.log.format" readonly /></div>
          </div>
        </div>
      </div>
    </template>
    <div v-else class="empty">Failed to load settings</div>
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
