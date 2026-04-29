<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()
const username = ref('')
const password = ref('')

async function handleLogin() {
  if (!username.value || !password.value) {
    ElMessage.warning('Please enter username and password')
    return
  }
  try {
    await authStore.login(username.value, password.value)
    router.push('/')
  } catch {
    ElMessage.error('Login failed')
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="logo">AEGIS</h1>
      <p class="subtitle">API Gateway Management</p>
      <form @submit.prevent="handleLogin">
        <div class="form-item">
          <label>Username</label>
          <input v-model="username" type="text" placeholder="admin" />
        </div>
        <div class="form-item">
          <label>Password</label>
          <input v-model="password" type="password" placeholder="password" />
        </div>
        <button type="submit" class="btn-login" :disabled="authStore.loading">
          {{ authStore.loading ? 'Signing in...' : 'Sign In' }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh; display: flex; align-items: center; justify-content: center;
  background: var(--bg-primary);
}
.login-card {
  background: var(--bg-secondary); border: 1px solid var(--border); border-radius: 12px;
  padding: 40px; width: 380px; text-align: center;
}
.logo { font-size: 32px; font-weight: 800; letter-spacing: 4px; margin-bottom: 4px; }
.subtitle { font-size: 13px; color: var(--text-secondary); margin-bottom: 32px; }
form { display: flex; flex-direction: column; gap: 16px; }
.form-item { text-align: left; }
.form-item label { display: block; font-size: 12px; font-weight: 600; color: var(--text-secondary); margin-bottom: 6px; }
.form-item input {
  width: 100%; padding: 10px 12px; border-radius: 6px; border: 1px solid var(--border);
  background: var(--bg-tertiary); color: var(--text-primary); font-size: 14px; box-sizing: border-box;
}
.btn-login {
  padding: 12px; border-radius: 6px; border: none; font-size: 14px; font-weight: 600;
  cursor: pointer; width: 100%; background: var(--accent); color: #fff; margin-top: 8px;
}
.btn-login:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
