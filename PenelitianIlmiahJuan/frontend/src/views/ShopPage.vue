<template>
  <div class="shop-page">
    <!-- Navbar -->
    <nav class="navbar glass-strong">
      <div class="navbar-container">
        <div class="navbar-brand">
          <div class="brand-icon">
            <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="20" cy="20" r="18" stroke="url(#brandGrad)" stroke-width="2.5" fill="none"/>
              <path d="M12 20C12 15.58 15.58 12 20 12C24.42 12 28 15.58 28 20C28 24.42 24.42 28 20 28" stroke="url(#brandGrad)" stroke-width="2.5" stroke-linecap="round"/>
              <circle cx="20" cy="20" r="4" fill="url(#brandGrad)"/>
              <defs>
                <linearGradient id="brandGrad" x1="0" y1="0" x2="40" y2="40">
                  <stop offset="0%" stop-color="#d35400"/>
                  <stop offset="50%" stop-color="#f5a623"/>
                  <stop offset="100%" stop-color="#ffd166"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="brand-text">
            <span class="brand-name text-gradient">OpenPeo</span>
            <span class="brand-tagline">Pasar NTT</span>
          </div>
        </div>

        <div class="navbar-links">
          <router-link to="/" class="nav-link">Beranda</router-link>
          <span class="nav-link active">Katalog Belanja</span>
        </div>

        <div class="navbar-actions">
          <div class="user-profile-nav">
            <span class="user-greeting">Halo, <strong>{{ currentUser?.name || currentUser?.username }}</strong></span>
            <span class="badge badge-indigo">Customer</span>
          </div>
          <button class="btn btn-secondary btn-sm" @click="logout">Keluar</button>
        </div>
      </div>
    </nav>

    <!-- Content -->
    <div class="shop-content">
      <header class="shop-header">
        <h1 class="shop-title">Eksplorasi <span class="text-gradient">Produk Autentik</span> NTT</h1>
        <p class="shop-subtitle">
          Selamat datang di platform marketplace NTT. Dapatkan kain tenun premium, kopi cita rasa tinggi, madu murni, dan kerajinan khas daerah NTT langsung dari produsen lokal dengan sistem pre-order terpercaya.
        </p>
      </header>

      <!-- Products Grid -->
      <ProductList @order="openOrderModal" />
    </div>

    <!-- Live Chat Floating -->
    <ChatBox :visible="chatOpen" @close="chatOpen = false" />
    
    <!-- Float chat toggle button if chat is closed -->
    <button v-show="!chatOpen" class="chat-toggle-btn btn btn-primary" @click="chatOpen = true">
      <span class="chat-toggle-icon">💬</span>
      <span>Tanya Penjual</span>
    </button>

    <!-- Order Modal -->
    <OrderModal v-if="selectedProduct" :product="selectedProduct" @close="selectedProduct = null" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import ProductList from '../components/ProductList.vue'
import OrderModal from '../components/OrderModal.vue'
import ChatBox from '../components/ChatBox.vue'

const router = useRouter()
const selectedProduct = ref(null)
const chatOpen = ref(true)

const currentUser = computed(() => {
  return JSON.parse(localStorage.getItem('openpeo_user') || 'null')
})

function openOrderModal(product) {
  selectedProduct.value = product
}

function logout() {
  localStorage.removeItem('openpeo_user')
  router.push('/')
}
</script>

<style scoped>
.shop-page {
  min-height: 100vh;
  background: var(--color-bg);
  padding-top: 5.5rem;
}

.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  padding: 0.75rem 0;
}

.navbar-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 var(--space-xl);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.brand-icon svg {
  width: 36px;
  height: 36px;
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
}

.brand-name {
  font-family: var(--font-display);
  font-size: 1.3rem;
  font-weight: 800;
}

.brand-tagline {
  font-size: 0.65rem;
  color: var(--color-text-muted);
  letter-spacing: 0.15em;
  text-transform: uppercase;
}

.navbar-links {
  display: flex;
  gap: var(--space-lg);
}

.nav-link {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
  padding: 0.5rem 0;
  position: relative;
  transition: color var(--transition-fast);
}

.nav-link::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 0;
  height: 2px;
  background: var(--gradient-amber);
  border-radius: 1px;
  transition: width var(--transition-base);
}

.nav-link:hover, .nav-link.active {
  color: var(--color-text-primary);
}

.nav-link:hover::after, .nav-link.active::after {
  width: 100%;
}

.navbar-actions {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.user-profile-nav {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.user-greeting {
  font-size: 0.85rem;
  color: var(--color-text-secondary);
}

.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.8rem;
}

.shop-content {
  max-width: 1360px;
  margin: 0 auto;
  padding: var(--space-xl);
}

.shop-header {
  text-align: center;
  max-width: 800px;
  margin: 0 auto var(--space-3xl);
  animation: fadeInUp 0.6s ease-out;
}

.shop-title {
  font-family: var(--font-display);
  font-size: clamp(2rem, 4vw, 3rem);
  font-weight: 800;
  margin-bottom: var(--space-md);
}

.shop-subtitle {
  color: var(--color-text-secondary);
  font-size: 0.98rem;
  line-height: 1.6;
}

.chat-toggle-btn {
  position: fixed;
  bottom: var(--space-xl);
  right: var(--space-xl);
  z-index: 1000;
  border-radius: var(--radius-full);
  padding: 0.8rem 1.5rem;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  box-shadow: var(--shadow-lg);
  animation: float 4s ease-in-out infinite;
}

.chat-toggle-icon {
  font-size: 1.2rem;
}

@media (max-width: 768px) {
  .navbar-links {
    display: none;
  }
}
</style>
