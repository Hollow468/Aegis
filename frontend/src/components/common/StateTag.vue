<script setup lang="ts">
defineProps<{
  state: string
  variant?: 'method' | 'match' | 'circuit' | 'status' | 'default'
}>()

const classMap: Record<string, Record<string, string>> = {
  method: { GET: 'tag-get', POST: 'tag-post', PUT: 'tag-put', DELETE: 'tag-delete' },
  match: { exact: 'tag-exact', prefix: 'tag-prefix', regex: 'tag-regex' },
  circuit: { closed: 'tag-closed', open: 'tag-open', 'half-open': 'tag-half-open' },
  status: { ok: 'tag-ok', warning: 'tag-warning', exceeded: 'tag-exceeded' },
  default: {}
}

function getClass(variant: string, state: string): string {
  return classMap[variant]?.[state] || 'tag-default'
}
</script>

<template>
  <span class="state-tag" :class="getClass(variant || 'default', state)">
    {{ state.toUpperCase() }}
  </span>
</template>

<style scoped>
.state-tag { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; }

/* Method tags */
.tag-get { background: rgba(34,197,94,0.15); color: var(--success); }
.tag-post { background: rgba(59,130,246,0.15); color: var(--accent); }
.tag-put { background: rgba(245,158,11,0.15); color: var(--warning); }
.tag-delete { background: rgba(239,68,68,0.15); color: var(--danger); }

/* Match type tags */
.tag-exact { background: rgba(59,130,246,0.1); color: var(--accent); }
.tag-prefix { background: rgba(6,182,212,0.1); color: var(--info); }
.tag-regex { background: rgba(168,85,247,0.1); color: #a855f7; }

/* Circuit breaker tags */
.tag-closed { background: rgba(34,197,94,0.15); color: var(--success); border: 1px solid rgba(34,197,94,0.3); }
.tag-open { background: rgba(239,68,68,0.15); color: var(--danger); border: 1px solid rgba(239,68,68,0.3); }
.tag-half-open { background: rgba(245,158,11,0.15); color: var(--warning); border: 1px solid rgba(245,158,11,0.3); }

/* Status tags */
.tag-ok { background: rgba(34,197,94,0.15); color: var(--success); }
.tag-warning { background: rgba(245,158,11,0.15); color: var(--warning); }
.tag-exceeded { background: rgba(239,68,68,0.15); color: var(--danger); }

.tag-default { background: var(--bg-tertiary); color: var(--text-secondary); }
</style>
