<script setup lang="ts">
import { ref } from 'vue'

const routes = ref([
  { path: '/api/users', method: 'GET', match_type: 'exact', balancer: 'Round Robin', rate_limit: '100/s', upstreams: 1 },
  { path: '/api/v1/user', method: 'GET', match_type: 'exact', balancer: 'Weighted RR', rate_limit: '100/s', upstreams: 1 },
  { path: '/api/files/*', method: 'GET', match_type: 'prefix', balancer: 'Random', rate_limit: '-', upstreams: 1 },
  { path: '/api/users/{id}', method: 'GET', match_type: 'regex', balancer: 'Consistent Hash', rate_limit: '50/s', upstreams: 1 }
])

const methodClass: Record<string, string> = { GET: 'tag-get', POST: 'tag-post', PUT: 'tag-put', DELETE: 'tag-delete' }
const matchClass: Record<string, string> = { exact: 'tag-exact', prefix: 'tag-prefix', regex: 'tag-regex' }
</script>

<template>
  <div class="routes-page">
    <div class="panel">
      <div class="panel-header">
        <h3>Route Configuration</h3>
        <button class="btn btn-primary">+ Add Route</button>
      </div>
      <div class="panel-body">
        <table>
          <thead>
            <tr><th>Path</th><th>Method</th><th>Match Type</th><th>Balancer</th><th>Rate Limit</th><th>Upstreams</th><th>Actions</th></tr>
          </thead>
          <tbody>
            <tr v-for="r in routes" :key="r.path">
              <td><code>{{ r.path }}</code></td>
              <td><span class="tag" :class="methodClass[r.method]">{{ r.method }}</span></td>
              <td><span class="tag" :class="matchClass[r.match_type]">{{ r.match_type }}</span></td>
              <td>{{ r.balancer }}</td>
              <td>{{ r.rate_limit }}</td>
              <td>{{ r.upstreams }}</td>
              <td>
                <button class="btn btn-secondary btn-sm">Edit</button>
                <button class="btn btn-danger btn-sm">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
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

.tag { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; }
.tag-get { background: rgba(34,197,94,0.15); color: var(--success); }
.tag-post { background: rgba(59,130,246,0.15); color: var(--accent); }
.tag-put { background: rgba(245,158,11,0.15); color: var(--warning); }
.tag-delete { background: rgba(239,68,68,0.15); color: var(--danger); }
.tag-exact { background: rgba(59,130,246,0.1); color: var(--accent); }
.tag-prefix { background: rgba(6,182,212,0.1); color: var(--info); }
.tag-regex { background: rgba(168,85,247,0.1); color: #a855f7; }

.btn { display: inline-flex; align-items: center; gap: 6px; padding: 8px 16px; border-radius: 6px; font-size: 13px; font-weight: 600; cursor: pointer; border: none; }
.btn-primary { background: var(--accent); color: #fff; }
.btn-secondary { background: var(--bg-tertiary); color: var(--text-primary); border: 1px solid var(--border); }
.btn-danger { background: rgba(239,68,68,0.15); color: var(--danger); }
.btn-sm { padding: 4px 10px; font-size: 12px; }
</style>
