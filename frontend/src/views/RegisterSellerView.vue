<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-header">
        <h1>Buka Toko Gratis</h1>
        <p>Jangkau pembeli di seluruh Indonesia dengan platform OpenPeo.</p>
      </div>

      <form @submit.prevent="handleRegisterSeller" class="auth-form">
        <div class="form-group">
          <label for="storeName">Nama Toko</label>
          <input type="text" id="storeName" v-model="form.storeName" required placeholder="Contoh: Toko Tenun Sumba" />
        </div>

        <div class="form-group">
          <label for="name">Nama Pemilik</label>
          <input type="text" id="name" v-model="form.name" required placeholder="Nama lengkap Anda" />
        </div>



        <div class="form-group">
          <label for="email">Alamat Email</label>
          <input type="email" id="email" v-model="form.email" required placeholder="email@contoh.com" pattern="[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$" title="Masukkan format email yang valid (contoh: user@domain.com)" />
        </div>

        <div class="form-group">
          <label for="password">Kata Sandi</label>
          <input type="password" id="password" v-model="form.password" required minlength="6" placeholder="Minimal 6 karakter" />
        </div>

        <div v-if="msg" :class="['auth-msg', statusClass]">
          {{ msg }}
        </div>

        <button type="submit" class="btn-primary auth-submit" :disabled="isLoading">
          {{ isLoading ? 'Mendaftar...' : 'Daftar Sekarang' }}
        </button>

        <div class="auth-footer">
          Sudah punya akun penjual? <router-link to="/login">Masuk di sini</router-link>
        </div>
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

const form = ref({
  storeName: '',
  name: '',
  email: '',
  password: ''
})

const isLoading = ref(false)
const msg = ref('')
const statusClass = ref('')

const handleRegisterSeller = async () => {
  isLoading.value = true
  msg.value = ''
  
  const payload = {
    store_name: form.value.storeName,
    name: form.value.name,
    email: form.value.email,
    password: form.value.password
  }

  const result = await authStore.registerSeller(payload)
  isLoading.value = false

  if (result.success) {
    statusClass.value = 'msg-success'
    msg.value = result.message
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } else {
    statusClass.value = 'msg-error'
    msg.value = result.message
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 120px 20px 80px;
  background-color: #fafafa;
}

.auth-container {
  width: 100%;
  max-width: 480px;
  background-color: #fff;
  padding: 48px;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.03);
  animation: slideUp 0.4s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.auth-header {
  text-align: center;
  margin-bottom: 40px;
}

.auth-header h1 {
  font-family: 'Playfair Display', serif;
  font-size: 2rem;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 12px;
}

.auth-header p {
  color: #666;
  font-size: 0.95rem;
  line-height: 1.5;
  margin: 0;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #1a1a1a;
}

.form-group input {
  padding: 14px 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.9rem;
  transition: border-color 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #1a1a1a;
}

.auth-msg {
  padding: 12px;
  border-radius: 4px;
  font-size: 0.85rem;
  text-align: center;
  font-weight: 500;
}

.msg-error {
  background-color: #fef2f2;
  color: #dc2626;
}

.msg-success {
  background-color: #f0fdf4;
  color: #16a34a;
}

.auth-submit {
  width: 100%;
  padding: 16px;
  margin-top: 8px;
  font-size: 0.95rem;
  font-weight: 600;
  border: none;
  border-radius: 4px;
  background-color: #1a1a1a;
  color: #fff;
  cursor: pointer;
  transition: background-color 0.2s;
}

.auth-submit:hover:not(:disabled) {
  background-color: #333;
}

.auth-submit:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.auth-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 0.85rem;
  color: #666;
}

.auth-footer a {
  color: #1a1a1a;
  font-weight: 600;
  text-decoration: none;
}

.auth-footer a:hover {
  text-decoration: underline;
}

@media (max-width: 480px) {
  .auth-container {
    padding: 32px 24px;
  }
}
</style>
