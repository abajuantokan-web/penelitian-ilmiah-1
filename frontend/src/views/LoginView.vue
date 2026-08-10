<template>
  <div class="login-page">
    
    <div class="login-visual">
      <div class="login-visual__overlay"></div>
      <div class="login-visual__content">
        <h1 class="login-visual__brand">OpenPeo</h1>
        <p class="login-visual__tagline">Platform Pre-order Premium<br>Produk Autentik NTT</p>
        <div class="login-visual__motif">
          <svg viewBox="0 0 48 48" aria-hidden="true">
            <circle cx="24" cy="24" r="20" stroke-dasharray="4 4"/>
            <circle cx="24" cy="24" r="12" stroke-dasharray="2 6"/>
            <circle cx="24" cy="24" r="4"/>
          </svg>
        </div>
      </div>
    </div>

    
    <div class="login-form-panel">
      <div class="login-container">
        <div class="login-header">
          <router-link to="/" class="login-back-link">← Kembali ke Beranda</router-link>
          <h2 class="login-title">{{ isRegisterMode ? 'Buat Akun Baru' : 'Selamat Datang' }}</h2>
          <p class="login-subtitle">{{ isRegisterMode ? 'Bergabunglah dengan komunitas pecinta warisan NTT' : 'Masuk ke akun OpenPeo Anda' }}</p>
        </div>

        
        <form v-if="!isRegisterMode" @submit.prevent="handleLogin" class="login-form" id="login-form">
          <div v-if="errorMsg" class="error-msg">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
            {{ errorMsg }}
          </div>

          <div v-if="successMsg" class="success-msg">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            {{ successMsg }}
          </div>

          <div class="form-group">
            <label for="login-email">Email</label>
            <input 
              type="email" 
              id="login-email" 
              v-model="email" 
              required
              autocomplete="email"
              placeholder="Masukkan email aktif"
              pattern="[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$"
              title="Masukkan format email yang valid (contoh: user@domain.com)"
            >
          </div>

          <div class="form-group">
            <label for="login-password">Password</label>
            <input 
              type="password" 
              id="login-password" 
              v-model="password" 
              required
              autocomplete="current-password"
              placeholder="Masukkan password"
            >
          </div>

          <button type="submit" class="login-submit-btn" :disabled="isLoading" id="login-submit">
            <span v-if="isLoading" class="spinner"></span>
            <span>{{ isLoading ? 'Memproses...' : 'Masuk' }}</span>
          </button>
        </form>

        
        <form v-else @submit.prevent="handleRegister" class="login-form" id="register-form">
          <div v-if="errorMsg" class="error-msg">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
            {{ errorMsg }}
          </div>

          <div class="form-group">
            <label for="reg-name">Nama Lengkap</label>
            <input type="text" id="reg-name" v-model="regName" required placeholder="Nama lengkap Anda">
          </div>



          <div class="form-group">
            <label for="reg-email">Email</label>
            <input type="email" id="reg-email" v-model="regEmail" required placeholder="Alamat email aktif" pattern="[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$" title="Masukkan format email yang valid (contoh: user@domain.com)">
          </div>

          <div class="form-group">
            <label for="reg-password">Password</label>
            <input type="password" id="reg-password" v-model="regPassword" required placeholder="Minimal 6 karakter" minlength="6">
          </div>

          <button type="submit" class="login-submit-btn" :disabled="isLoading" id="register-submit">
            <span v-if="isLoading" class="spinner"></span>
            <span>{{ isLoading ? 'Mendaftarkan...' : 'Daftar Sekarang' }}</span>
          </button>
        </form>

        <div class="login-footer">
          <p v-if="!isRegisterMode">
            Belum punya akun? 
            <a href="#" @click.prevent="toggleMode" id="switch-to-register">Daftar sekarang</a>
          </p>
          <p v-else>
            Sudah punya akun? 
            <a href="#" @click.prevent="toggleMode" id="switch-to-login">Masuk di sini</a>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const cartStore = useCartStore()

const isRegisterMode = ref(false)
const isLoading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')


const email = ref('')
const password = ref('')


const regName = ref('')
const regEmail = ref('')
const regPassword = ref('')

const toggleMode = () => {
  isRegisterMode.value = !isRegisterMode.value
  errorMsg.value = ''
  successMsg.value = ''
}

const handleLogin = async () => {
  isLoading.value = true
  errorMsg.value = ''
  const result = await authStore.login(email.value, password.value)
  
  isLoading.value = false
  
  if (result.success) {
    // Sync cart from backend after login
    await cartStore.fetchCart()
    
    // Redirect to the seller dashboard (which is now the root route)
    router.push('/');
  } else {
    errorMsg.value = result.message || 'Login gagal. Periksa kembali email dan password Anda.'
  }
}

const handleRegister = async () => {
  isLoading.value = true
  errorMsg.value = ''
  successMsg.value = ''
  
  try {
    const response = await axios.post('http://localhost:8081/api/register', {
      name: regName.value,
      email: regEmail.value,
      password: regPassword.value,
    })

    if (response.status === 201 || response.data.success) {
      successMsg.value = 'Pendaftaran berhasil! Silakan masuk.'
      isRegisterMode.value = false
      email.value = regEmail.value
    } else {
      errorMsg.value = response.data.message || 'Pendaftaran gagal.'
    }
  } catch (error) {
    errorMsg.value = error.response?.data?.message || 'Pendaftaran gagal: Email sudah digunakan atau format salah'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  min-height: 100vh;
  background-color: #fafafa;
}

/* Left Visual Panel */
.login-visual {
  position: relative;
  flex: 0 0 45%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #0a0a0a;
  overflow: hidden;
}

.login-visual__overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(10,10,10,0.95) 0%, rgba(30,30,30,0.9) 100%);
  z-index: 1;
}

.login-visual__content {
  position: relative;
  z-index: 2;
  text-align: center;
  padding: 40px;
}

.login-visual__brand {
  font-family: 'Playfair Display', serif;
  font-size: 3rem;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 3px;
  margin-bottom: 20px;
}

.login-visual__tagline {
  font-family: 'Montserrat', sans-serif;
  font-size: 0.95rem;
  color: rgba(255,255,255,0.55);
  letter-spacing: 1.5px;
  line-height: 1.8;
  text-transform: uppercase;
}

.login-visual__motif {
  margin-top: 48px;
  opacity: 0.15;
}

.login-visual__motif svg {
  width: 64px;
  height: 64px;
  stroke: #ffffff;
  fill: none;
  stroke-width: 1;
}

/* Right Form Panel */
.login-form-panel {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.login-container {
  width: 100%;
  max-width: 380px;
}

.login-back-link {
  display: inline-block;
  font-size: 0.8rem;
  color: #999;
  letter-spacing: 0.5px;
  margin-bottom: 32px;
  transition: color 0.3s;
  text-decoration: none;
}

.login-back-link:hover {
  color: #1a1a1a;
}

.login-title {
  font-family: 'Playfair Display', serif;
  font-size: 2rem;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px;
  letter-spacing: 0.5px;
}

.login-subtitle {
  font-size: 0.85rem;
  color: #888;
  margin: 0 0 40px;
  letter-spacing: 0.3px;
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 0.75rem;
  font-weight: 500;
  color: #555;
  text-transform: uppercase;
  letter-spacing: 1.5px;
}

.form-group input {
  width: 100%;
  padding: 14px 16px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.9rem;
  color: #1a1a1a;
  background: #fff;
  transition: border-color 0.3s, box-shadow 0.3s;
  box-sizing: border-box;
}

.form-group input::placeholder {
  color: #bbb;
}

.form-group input:focus {
  outline: none;
  border-color: #1a1a1a;
  box-shadow: 0 0 0 3px rgba(26,26,26,0.05);
}

.login-submit-btn {
  width: 100%;
  padding: 16px 32px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.85rem;
  font-weight: 500;
  letter-spacing: 1.5px;
  text-transform: uppercase;
  color: #ffffff;
  background-color: #1a1a1a;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 8px;
}

.login-submit-btn:hover:not(:disabled) {
  background-color: #333;
}

.login-submit-btn:active:not(:disabled) {
  transform: scale(0.99);
}

.login-submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-msg {
  display: flex;
  align-items: center;
  gap: 10px;
  background-color: #fef2f2;
  color: #dc2626;
  padding: 14px 16px;
  border-radius: 4px;
  margin-bottom: 24px;
  font-size: 0.85rem;
  border-left: 3px solid #dc2626;
}

.success-msg {
  display: flex;
  align-items: center;
  gap: 10px;
  background-color: #f0fdf4;
  color: #16a34a;
  padding: 14px 16px;
  border-radius: 4px;
  margin-bottom: 24px;
  font-size: 0.85rem;
  border-left: 3px solid #16a34a;
}

.login-footer {
  margin-top: 32px;
  text-align: center;
  font-size: 0.85rem;
  color: #888;
}

.login-footer a {
  color: #1a1a1a;
  font-weight: 500;
  text-decoration: none;
  border-bottom: 1px solid #1a1a1a;
  transition: opacity 0.3s;
}

.login-footer a:hover {
  opacity: 0.7;
}

/* Responsive */
@media (max-width: 768px) {
  .login-page {
    flex-direction: column;
  }
  
  .login-visual {
    flex: 0 0 200px;
  }
  
  .login-visual__brand {
    font-size: 2rem;
  }
  
  .login-visual__tagline {
    font-size: 0.8rem;
  }
  
  .login-visual__motif {
    display: none;
  }
  
  .login-form-panel {
    padding: 32px 20px;
  }
}
</style>
