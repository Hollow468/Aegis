<script setup lang="ts">
import { ref } from 'vue'

const settings = ref({
  server: { port: 8080, host: '0.0.0.0', readTimeout: '30s', writeTimeout: '30s' },
  etcd: { endpoints: 'http://localhost:2379', prefix: '/aegis/services' },
  redis: { addr: 'localhost:6379', db: 0 },
  jwt: { secret: '********', expire: '24h' },
  log: { level: 'info', format: 'json' }
})
</script>

<template>
  <div class="settings-page">
    <div class="panel">
      <div class="panel-header"><h3>Server Configuration</h3></div>
      <div class="panel-body">
        <div class="form-grid">
          <div class="form-item">
            <label>Port</label>
            <input type="text" :value="settings.server.port" readonly />
          </div>
          <div class="form-item">
            <label>Host</label>
            <input type="text" :value="settings.server.host" readonly />
          </div>
          <div class="form-item">
            <label>Read Timeout</label>
            <input type="text" :value="settings.server.readTimeout" readonly />
          </div>
          <div class="form-item">
            <label>Write Timeout</label>
            <input type="text" :value="settings.server.writeTimeout" readonly />
          </div>
        </div>
      </div>
    </div>

    <div class="panel" style="margin-top:16px;">
      <div class="panel-header"><h3>External Services</h3></div>
      <div class="panel-body">
        <div class="form-grid">
          <div class="form-item">
            <label>etcd Endpoints</label>
            <input type="text" :value="settings.etcd.endpoints" readonly />
          </div>
          <div class="form-item">
            <label>etcd Prefix</label>
            <input type="text" :value="settings.etcd.prefix" readonly />
          </div>
          <div class="form-item">
            <label>Redis Address</label>
            <input type="text" :value="settings.redis.addr" readonly />
          </div>
          <div class="form-item">
            <label>Redis DB</label>
            <input type="text" :value="settings.redis.db" readonly />
          </div>
        </div>
      </div>
    </div>

    <div class="panel" style="margin-top:16px;">
      <div class="panel-header"><h3>Authentication & Logging</h3></div>
      <div class="panel-body">
        <div class="form-grid">
          <div class="form-item">
            <label>JWT Secret</label>
            <input type="password" :value="settings.jwt.secret" readonly />
          </div>
          <div class="form-item">
            <label>Token Expiry</label>
            <input type="text" :value="settings.jwt.expire" readonly />
          </div>
          <div class="form-item">
            <label>Log Level</label>
            <input type="text" :value="settings.log.level" readonly />
          </div>
          <div class="form-item">
            <label>Log Format</label>
            <input type="text" :value="settings.log.format" readonly />
          </div>
        </div>
      </div>
    </div>
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
</style>
