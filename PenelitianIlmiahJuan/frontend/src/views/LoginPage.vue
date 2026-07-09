<template>
  <div class="login-page">
    <!-- Background effects -->
    <div class="login-bg">
      <div class="login-orb login-orb-1"></div>
      <div class="login-orb login-orb-2"></div>
      <div class="login-grid"></div>
    </div>

    <div class="login-container" ref="containerRef">
      <!-- Back to home -->
      <router-link to="/" class="back-link">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
        Kembali ke Beranda
      </router-link>

      <!-- Login/Register Card -->
      <div class="login-card glass-strong" ref="cardRef">
        <!-- Brand -->
        <div class="login-brand">
          <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg" class="brand-svg">
            <circle cx="20" cy="20" r="18" stroke="url(#lgGrad)" stroke-width="2.5" fill="none"/>
            <path d="M12 20C12 15.58 15.58 12 20 12C24.42 12 28 15.58 28 20C28 24.42 24.42 28 20 28" stroke="url(#lgGrad)" stroke-width="2.5" stroke-linecap="round"/>
            <circle cx="20" cy="20" r="4" fill="url(#lgGrad)"/>
            <defs><linearGradient id="lgGrad" x1="0" y1="0" x2="40" y2="40"><stop offset="0%" stop-color="#d35400"/><stop offset="50%" stop-color="#f5a623"/><stop offset="100%" stop-color="#ffd166"/></linearGradient></defs>
          </svg>
          <h1 class="login-title">
            <span class="text-gradient">OpenPeo</span>
          </h1>
          <p class="login-subtitle">{{ isRegister ? 'Daftar Akun Pembeli Baru' : 'Masuk ke akun Anda' }}</p>
        </div>

        <!-- Success message -->
        <div v-if="successMsg" class="success-message">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" style="margin-right: 4px;"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>
          {{ successMsg }}
        </div>

        <!-- Error message -->
        <div v-if="errorMsg" class="error-message">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" style="margin-right: 4px;"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/></svg>
          {{ errorMsg }}
        </div>

        <template v-if="!isRegister">
          <!-- Role selector pills -->
          <div class="role-selector">
            <button
              v-for="r in roles"
              :key="r.value"
              class="role-pill"
              :class="{ active: selectedRole === r.value }"
              @click="selectRole(r)"
            >
              <span class="role-icon">{{ r.icon }}</span>
              <span>{{ r.label }}</span>
            </button>
          </div>

          <!-- Login Form -->
          <form @submit.prevent="handleLogin" class="login-form">
            <div class="form-group">
              <label class="form-label" for="login-username">Username</label>
              <div class="input-wrapper">
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/>
                </svg>
                <input
                  id="login-username"
                  v-model="username"
                  type="text"
                  class="input"
                  :placeholder="usernamePlaceholder"
                  required
                  autocomplete="username"
                />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="login-password">Password</label>
              <div class="input-wrapper">
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                <input
                  id="login-password"
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  class="input"
                  placeholder="Masukkan password"
                  required
                  autocomplete="current-password"
                />
                <button type="button" class="toggle-password" @click="showPassword = !showPassword">
                  <svg v-if="!showPassword" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>
                  </svg>
                  <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                    <line x1="1" y1="1" x2="23" y2="23"/>
                  </svg>
                </button>
              </div>
            </div>

            <!-- Submit -->
            <button type="submit" class="btn btn-primary btn-login" :disabled="loading">
              <span v-if="loading" class="spinner"></span>
              {{ loading ? 'Memproses...' : 'Masuk' }}
            </button>
          </form>

          <div class="switch-auth">
            Belum punya akun? <a href="#" @click.prevent="toggleAuthMode(true)">Daftar sekarang</a>
          </div>

          <!-- Demo credentials hint -->
          <div class="demo-hint">
            <p class="hint-title">🔑 Demo Credentials</p>
            <div class="hint-grid">
              <div class="hint-item">
                <code>admin_flores</code>
                <span class="hint-role badge badge-amber">Admin</span>
              </div>
              <div class="hint-item">
                <code>pembeli_flores</code>
                <span class="hint-role badge badge-indigo">Customer</span>
              </div>
            </div>
            <p class="hint-password">Password: <code>password123</code></p>
          </div>
        </template>

        <template v-else>
          <!-- Register Form -->
          <form @submit.prevent="handleRegister" class="login-form">
            <div class="form-group">
              <label class="form-label" for="reg-name">Nama Lengkap *</label>
              <div class="input-wrapper">
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/>
                </svg>
                <input
                  id="reg-name"
                  v-model="regName"
                  type="text"
                  class="input"
                  placeholder="Nama lengkap Anda"
                  required
                />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="reg-username">Username *</label>
              <div class="input-wrapper">
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="4"/><path d="M16 8v5a3 3 0 0 0 6 0v-1a10 10 0 1 0-3.92 7.94"/>
                </svg>
                <input
                  id="reg-username"
                  v-model="regUsername"
                  type="text"
                  class="input"
                  placeholder="Username unik"
                  required
                />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="reg-email">Email *</label>
              <div class="input-wrapper">
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/>
                </svg>
                <input
                  id="reg-email"
                  v-model="regEmail"
                  type="email"
                  class="input"
                  placeholder="alamat@email.com"
                  required
                />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="reg-phone">No. Telepon</label>
              <div class="input-wrapper">
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
                </svg>
                <input
                  id="reg-phone"
                  v-model="regPhone"
                  type="tel"
                  class="input"
                  placeholder="Nomor telepon aktif"
                />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="reg-address">Alamat Pengiriman</label>
              <textarea
                id="reg-address"
                v-model="regAddress"
                class="input textarea"
                placeholder="Alamat pengiriman lengkap untuk pre-order"
                rows="2"
              ></textarea>
            </div>

            <div class="form-group">
              <label class="form-label" for="reg-password">Password *</label>
              <div class="input-wrapper">
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                <input
                  id="reg-password"
                  v-model="regPassword"
                  :type="showPassword ? 'text' : 'password'"
                  class="input"
                  placeholder="Password minimal 6 karakter"
                  required
                />
              </div>
            </div>

            <!-- Submit -->
            <button type="submit" class="btn btn-primary btn-login" :disabled="loading">
              <span v-if="loading" class="spinner"></span>
              {{ loading ? 'Mendaftar...' : 'Daftar' }}
            </button>
          </form>

          <div class="switch-auth">
            Sudah punya akun? <a href="#" @click.prevent="toggleAuthMode(false)">Masuk di sini</a>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { animate } from 'animejs'

const router = useRouter()
const route = useRoute()

const API_BASE = 'http://localhost:8080/api'

const roles = [
  { value: 'admin', label: 'Admin', icon: '👑', placeholder: 'admin_flores' },
  { value: 'customer', label: 'Customer', icon: '🛒', placeholder: 'pembeli_flores' },
]

// Auth states
const isRegister = ref(false)
const selectedRole = ref('customer')
const username = ref('')
const password = ref('password123')
const showPassword = ref(false)
const loading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')
const containerRef = ref(null)
const cardRef = ref(null)

// Registration states
const regName = ref('')
const regUsername = ref('')
const regEmail = ref('')
const regPhone = ref('')
const regAddress = ref('')
const regPassword = ref('')

const usernamePlaceholder = computed(() => {
  const role = roles.find(r => r.value === selectedRole.value)
  return role ? role.placeholder : 'Username'
})

function selectRole(role) {
  selectedRole.value = role.value
  username.value = role.placeholder
  errorMsg.value = ''
  successMsg.value = ''
}

function toggleAuthMode(registerMode) {
  isRegister.value = registerMode
  errorMsg.value = ''
  successMsg.value = ''
}

async function handleLogin() {
  if (!username.value.trim() || !password.value.trim()) {
    errorMsg.value = 'Username dan password harus diisi'
    return
  }

  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''

  try {
    const response = await fetch(`${API_BASE}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: username.value.trim(),
        password: password.value.trim()
      })
    })

    const data = await response.json()

    if (data.success) {
      // Save session to localStorage
      localStorage.setItem('openpeo_user', JSON.stringify(data.data))

      // Redirect based on role
      const redirect = route.query.redirect
      if (redirect) {
        router.push(redirect)
      } else if (data.data.role === 'admin') {
        router.push('/admin/dashboard')
      } else {
        router.push('/shop')
      }
    } else {
      errorMsg.value = data.message || 'Login gagal'
      shakeCard()
    }
  } catch (error) {
    errorMsg.value = 'Tidak dapat terhubung ke server. Pastikan backend berjalan.'
    shakeCard()
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  if (!regName.value.trim() || !regUsername.value.trim() || !regEmail.value.trim() || !regPassword.value.trim()) {
    errorMsg.value = 'Mohon isi semua bidang wajib (*)'
    return
  }

  if (regPassword.value.length < 6) {
    errorMsg.value = 'Password harus minimal 6 karakter'
    return
  }

  loading.value = true
  errorMsg.value = ''
  successMsg.value = ''

  try {
    const response = await fetch(`${API_BASE}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: regName.value.trim(),
        username: regUsername.value.trim(),
        email: regEmail.value.trim(),
        password: regPassword.value.trim(),
        phone: regPhone.value.trim(),
        address: regAddress.value.trim()
      })
    })

    const data = await response.json()

    if (data.success) {
      successMsg.value = 'Pendaftaran berhasil! Silakan masuk dengan akun baru Anda.'
      // Prefill login form
      username.value = regUsername.value.trim()
      password.value = regPassword.value.trim()
      selectedRole.value = 'customer'
      
      // Clear fields
      regName.value = ''
      regUsername.value = ''
      regEmail.value = ''
      regPhone.value = ''
      regAddress.value = ''
      regPassword.value = ''
      
      // Go to login mode
      isRegister.value = false
    } else {
      errorMsg.value = data.message || 'Pendaftaran gagal'
      shakeCard()
    }
  } catch (error) {
    errorMsg.value = 'Tidak dapat terhubung ke server. Pastikan backend berjalan.'
    shakeCard()
  } finally {
    loading.value = false
  }
}

function shakeCard() {
  if (cardRef.value) {
    animate(cardRef.value, {
      translateX: [0, -10, 10, -10, 10, 0],
      duration: 400,
      ease: 'easeInOutQuad'
    })
  }
}

onMounted(() => {
  if (containerRef.value) {
    animate(containerRef.value, {
      opacity: [0, 1],
      translateY: [40, 0],
      duration: 700,
      ease: 'outExpo'
    })
  }
})
</script>

<style scoped>
.login-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden; padding: var(--space-xl); }
.login-bg { position: absolute; inset: 0; pointer-events: none; }
.login-orb { position: absolute; border-radius: 50%; filter: blur(120px); opacity: 0.25; }
.login-orb-1 { width: 500px; height: 500px; background: radial-gradient(circle, var(--color-indigo-glow), transparent 70%); top: -150px; right: -150px; animation: float 8s ease-in-out infinite; }
.login-orb-2 { width: 400px; height: 400px; background: radial-gradient(circle, var(--color-burnt-orange), transparent 70%); bottom: -100px; left: -100px; animation: float 10s ease-in-out infinite reverse; }
.login-grid { position: absolute; inset: 0; background-image: linear-gradient(rgba(255,255,255,0.015) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.015) 1px, transparent 1px); background-size: 60px 60px; mask-image: radial-gradient(ellipse at center, black 20%, transparent 70%); }

.login-container { position: relative; z-index: 1; width: 100%; max-width: 440px; }
.back-link { display: inline-flex; align-items: center; gap: 0.3rem; color: var(--color-text-muted); font-size: 0.82rem; margin-bottom: var(--space-lg); transition: color var(--transition-fast); }
.back-link:hover { color: var(--color-text-primary); }

.login-card { padding: var(--space-2xl); border-radius: var(--radius-2xl); }
.login-brand { text-align: center; margin-bottom: var(--space-xl); }
.brand-svg { width: 56px; height: 56px; margin: 0 auto var(--space-md); }
.login-title { font-family: var(--font-display); font-size: 1.8rem; font-weight: 800; margin-bottom: var(--space-xs); }
.login-subtitle { color: var(--color-text-muted); font-size: 0.9rem; }

.role-selector { display: flex; gap: var(--space-xs); margin-bottom: var(--space-xl); }
.role-pill { flex: 1; display: flex; align-items: center; justify-content: center; gap: 0.4rem; padding: 0.6rem; border-radius: var(--radius-md); background: rgba(255,255,255,0.03); color: var(--color-text-secondary); font-size: 0.82rem; font-weight: 500; border: 1px solid transparent; transition: all var(--transition-fast); }
.role-pill:hover { background: rgba(255,255,255,0.06); }
.role-pill.active { background: var(--color-amber-glow); color: var(--color-amber); border-color: rgba(245,166,35,0.3); }
.role-icon { font-size: 1.1rem; }

.login-form { display: flex; flex-direction: column; gap: var(--space-lg); }
.form-group { display: flex; flex-direction: column; gap: var(--space-xs); }
.form-label { font-size: 0.82rem; font-weight: 600; color: var(--color-text-secondary); }
.input-wrapper { position: relative; }
.input-icon { position: absolute; left: 0.85rem; top: 50%; transform: translateY(-50%); color: var(--color-text-muted); pointer-events: none; }
.input-wrapper .input { padding-left: 2.5rem; padding-right: 2.5rem; }
.toggle-password { position: absolute; right: 0.75rem; top: 50%; transform: translateY(-50%); background: none; color: var(--color-text-muted); transition: color var(--transition-fast); }
.toggle-password:hover { color: var(--color-text-primary); }

.error-message { display: flex; align-items: center; gap: var(--space-sm); padding: var(--space-sm) var(--space-md); background: rgba(231,76,60,0.1); border: 1px solid rgba(231,76,60,0.2); border-radius: var(--radius-md); color: var(--color-error); font-size: 0.82rem; margin-bottom: var(--space-sm); }
.success-message { display: flex; align-items: center; gap: var(--space-sm); padding: var(--space-sm) var(--space-md); background: rgba(46,204,113,0.1); border: 1px solid rgba(46,204,113,0.2); border-radius: var(--radius-md); color: var(--color-success); font-size: 0.82rem; margin-bottom: var(--space-sm); }

.switch-auth { margin-top: var(--space-lg); font-size: 0.85rem; color: var(--color-text-secondary); text-align: center; }
.switch-auth a { color: var(--color-amber); font-weight: 600; text-decoration: none; transition: color var(--transition-fast); }
.switch-auth a:hover { text-decoration: underline; color: var(--color-amber-light); }

.btn-login { width: 100%; padding: 0.85rem; font-size: 0.95rem; }
.spinner { width: 16px; height: 16px; border: 2px solid rgba(11,10,18,0.3); border-top-color: #0b0a12; border-radius: 50%; animation: spin-slow 0.6s linear infinite; margin-right: var(--space-sm); }

.demo-hint { margin-top: var(--space-xl); padding: var(--space-md); background: rgba(255,255,255,0.02); border-radius: var(--radius-md); border: 1px dashed rgba(255,255,255,0.08); }
.hint-title { font-size: 0.78rem; font-weight: 600; color: var(--color-text-muted); margin-bottom: var(--space-sm); }
.hint-grid { display: flex; flex-direction: column; gap: var(--space-xs); margin-bottom: var(--space-sm); }
.hint-item { display: flex; align-items: center; justify-content: space-between; }
.hint-item code { font-size: 0.78rem; color: var(--color-amber-light); background: rgba(245,166,35,0.1); padding: 0.15rem 0.5rem; border-radius: var(--radius-sm); }
.hint-role { font-size: 0.65rem; padding: 0.15rem 0.5rem; }
.hint-password { font-size: 0.75rem; color: var(--color-text-muted); }
.hint-password code { font-size: 0.75rem; color: var(--color-amber-light); background: rgba(245,166,35,0.1); padding: 0.1rem 0.4rem; border-radius: var(--radius-sm); }

.textarea { resize: vertical; min-height: 60px; font-family: inherit; font-size: 0.875rem; }
</style>
