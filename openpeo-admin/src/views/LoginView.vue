<template>
  <div class="admin-login-wrapper">
    <div class="login-card">
      <div class="login-header">
        <h1>OpenPeo</h1>
        <p>Admin Portal</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>Email Admin</label>
          <input type="email" v-model="email" required placeholder="admin@openpeo.com" />
        </div>
        
        <div class="form-group">
          <label>Password</label>
          <input type="password" v-model="password" required />
        </div>
        
        <button type="submit" class="btn-login" :disabled="isLoading">
          {{ isLoading ? 'Memverifikasi...' : 'Masuk ke Portal' }}
        </button>
        
        <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const isLoading = ref(false)
const errorMsg = ref('')

const handleLogin = async () => {
  isLoading.value = true
  errorMsg.value = ''
  
  const result = await authStore.login(email.value, password.value)
  
  if (result.success) {
    router.push('/')
  } else {
    errorMsg.value = result.message
  }
  isLoading.value = false
}
</script>

<style scoped>
.admin-login-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f3f4f6;
  font-family: 'Inter', sans-serif;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.05);
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #111827;
  margin: 0;
  letter-spacing: -0.5px;
}

.login-header p {
  color: #6b7280;
  margin: 5px 0 0;
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 2px;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  font-size: 0.85rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 1rem;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

input:focus {
  outline: none;
  border-color: #111827;
  box-shadow: 0 0 0 3px rgba(17,24,39,0.1);
}

.btn-login {
  width: 100%;
  padding: 14px;
  background-color: #111827;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-top: 10px;
}

.btn-login:hover {
  background-color: #374151;
}

.btn-login:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.error-msg {
  margin-top: 16px;
  padding: 12px;
  background-color: #fee2e2;
  color: #991b1b;
  border-radius: 6px;
  font-size: 0.85rem;
  text-align: center;
}
</style>
