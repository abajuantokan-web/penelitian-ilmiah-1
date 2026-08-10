<template>
  <div id="app-root" >
    <!-- Header -->
    <header :class="headerClass">
      <div class="container header-inner">
        <div class="logo">
          <router-link to="/">OpenPeo</router-link>
        </div>

        <nav class="nav-links" aria-label="Main Navigation">
          <router-link to="/">Beranda</router-link>
          <router-link to="/koleksi">Koleksi</router-link>
          <router-link to="/#cerita-kami">Cerita Kami</router-link>
          <router-link to="/tentang">Tentang</router-link>
        </nav>

        <div class="header-icons">
          <!-- Search -->
          <button class="icon-btn" aria-label="Search" @click="goToSearch">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
          </button>
          
          <!-- Profile/Login -->
          <button 
            :class="['icon-btn', { 'icon-btn--active': authStore.isAuthenticated }]" 
            aria-label="Profile" 
            @click="handleProfileClick"
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
            <span v-if="authStore.isAuthenticated" class="profile-dot"></span>
          </button>

          <!-- Store / Seller / Order History -->
          <div class="store-menu-wrapper" @mouseleave="isStoreDropdownOpen = false">
            <template v-if="authStore.isAuthenticated && authStore.user?.role !== 'seller'">
              <button class="icon-btn relative" aria-label="Order History" @click="router.push({ path: '/profile', query: { tab: 'orders' } })">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                  <polyline points="14 2 14 8 20 8"></polyline>
                  <line x1="16" y1="13" x2="8" y2="13"></line>
                  <line x1="16" y1="17" x2="8" y2="17"></line>
                  <polyline points="10 9 9 9 8 9"></polyline>
                </svg>
              </button>
            </template>
            <template v-else>
              <button class="icon-btn relative" aria-label="Store" @mouseenter="isStoreDropdownOpen = true" @click="handleStoreClick">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
                  <polyline points="9 22 9 12 15 12 15 22"></polyline>
                </svg>
                <span v-if="pendingOrdersCount > 0" class="absolute -top-0.5 -right-0.5 flex h-2 w-2">
                  <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
                  <span class="relative inline-flex rounded-full h-2 w-2 bg-red-500"></span>
                </span>
              </button>
              
              <div v-if="isStoreDropdownOpen" class="store-dropdown">
                <div v-if="authStore.user?.role === 'seller'">
                  <p class="store-dropdown-title">Kelola Toko Anda</p>
                  <div class="store-dropdown-profile" v-if="authStore.user?.seller_profile?.store_name || authStore.user?.store_name">
                    <img v-if="authStore.user?.seller_profile?.store_logo || authStore.user?.store_logo" :src="$getImageUrl(authStore.user.seller_profile?.store_logo || authStore.user.store_logo)" alt="Logo" class="store-dropdown-logo" />
                    <span class="store-dropdown-name">{{ authStore.user.seller_profile?.store_name || authStore.user.store_name }}</span>
                  </div>
                  <router-link to="/seller/dashboard" class="btn-primary store-btn">Dashboard Seller</router-link>
                </div>
                <div v-else>
                  <p class="store-dropdown-title">Mulai Berjualan</p>
                  <p class="store-dropdown-desc">Jangkau pembeli di seluruh Indonesia.</p>
                  <router-link to="/register-seller" class="btn-primary store-btn">Buka Toko Gratis</router-link>
                </div>
              </div>
            </template>
          </div>

          <!-- Cart -->
          <button class="icon-btn cart-btn" aria-label="Cart" @click="cartStore.toggleDrawer">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
              <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/>
              <line x1="3" y1="6" x2="21" y2="6"/>
              <path d="M16 10a4 4 0 0 1-8 0"/>
            </svg>
            <span v-if="cartStore.totalItems > 0" class="cart-badge">{{ cartStore.totalItems }}</span>
          </button>
        </div>
      </div>
    </header>

    <!-- Router View for Pages -->
    <main class="main-content" :class="{ 'with-header-padding': route.path !== '/' }">
      <router-view />
    </main>

    <!-- Cart Drawer -->
    <CartDrawer />
    
    <!-- Live Chat Widget -->
    <LiveChat v-if="authStore.user?.role !== 'seller'" />

    <!-- Footer -->
    <footer class="site-footer" role="contentinfo">
      <div class="container">
        <div class="footer-grid">
          <div class="footer-brand">
            <p class="footer-brand-name">OpenPeo</p>
            <p class="footer-brand-desc">Platform pre-order premium untuk produk artisan Nusa Tenggara Timur. Menghubungkan pengrajin lokal dengan pecinta warisan budaya.</p>
          </div>
          <div class="footer-col">
            <p class="footer-col-title">Belanja</p>
            <div class="footer-links">
              <a href="#">Koleksi Baru</a>
              <a href="#">Tenun Sumba</a>
              <a href="#">Tenun Flores</a>
              <a href="#">Aksesoris</a>
            </div>
          </div>
          <div class="footer-col">
            <p class="footer-col-title">Bantuan</p>
            <div class="footer-links">
              <a href="#">FAQ</a>
              <a href="#">Pengiriman</a>
              <a href="#">Pengembalian</a>
              <a href="#">Panduan Ukuran</a>
            </div>
          </div>
          <div class="footer-col">
            <p class="footer-col-title">Berlangganan</p>
            <p style="font-size: 13px; margin-bottom: 16px; opacity: 0.8;">Dapatkan info koleksi terbaru.</p>
            <div class="newsletter">
              <input type="email" placeholder="Alamat Email" aria-label="Email untuk berlangganan">
              <button>→</button>
            </div>
          </div>
        </div>
        <div class="footer-bottom">
          <p>&copy; 2026 OpenPeo. Hak Cipta Dilindungi.</p>
          <div class="footer-legal">
            <a href="#">Kebijakan Privasi</a>
            <a href="#">Syarat &amp; Ketentuan</a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useCartStore } from './stores/cart'
import { useNotificationStore } from './stores/notification'
import { useDashboardStore } from './stores/dashboard'
import { useWebsocketStore } from './stores/websocket'
import CartDrawer from './components/CartDrawer.vue'
import LiveChat from './components/LiveChat.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const cartStore = useCartStore()
const notificationStore = useNotificationStore()
const dashboardStore = useDashboardStore()
const websocketStore = useWebsocketStore() // Initializes singleton WS

const isScrolled = ref(false)
const isStoreDropdownOpen = ref(false)
const pendingOrdersCount = computed(() => dashboardStore.pendingCount || 0)

const headerClass = computed(() => {
  return ['site-header', { 'scrolled': isScrolled.value || route.path !== '/' }]
})

const handleScroll = () => {
  isScrolled.value = window.scrollY > 50
}

const handleProfileClick = () => {
  if (authStore.isAuthenticated) {
    router.push({ path: '/profile', query: { tab: 'account' } })
  } else {
    router.push('/login')
  }
}

const handleStoreClick = () => {
  if (authStore.user?.role === 'seller') {
    router.push('/seller/dashboard')
  } else {
    router.push('/register-seller')
  }
}

const goToSearch = () => {
  router.push({ path: '/koleksi', query: { search: 1 } })
}

// IntersectionObserver for scroll reveal animations
let revealObserver = null

const initRevealObserver = () => {
  if (revealObserver) revealObserver.disconnect()

  revealObserver = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible')
        revealObserver.unobserve(entry.target)
      }
    })
  }, {
    threshold: 0.1,
    rootMargin: '0px 0px -40px 0px'
  })

  document.querySelectorAll('.reveal').forEach((el) => {
    revealObserver.observe(el)
  })
}

// Re-observe when route changes (new elements may appear)
watch(() => router.currentRoute.value, () => {
  nextTick(() => {
    setTimeout(initRevealObserver, 100)
  })
})

onMounted(() => {
  window.addEventListener('scroll', handleScroll)

  if (authStore.isAuthenticated) {
    cartStore.fetchCart()
    if (authStore.user?.role === 'seller') {
      dashboardStore.fetchInitialStats()
    }
  }

  // Initial reveal observer setup
  nextTick(() => {
    setTimeout(initRevealObserver, 300)
  })
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  if (revealObserver) revealObserver.disconnect()
})
</script>

<style>
/* Profile active dot */
.icon-btn {
  position: relative;
}

.profile-dot {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #16a34a;
  border: 2px solid transparent;
  pointer-events: none;
}

.site-header:not(.scrolled) .profile-dot {
  border-color: rgba(0, 0, 0, 0.2);
}

.site-header.scrolled .profile-dot {
  border-color: #fff;
}

/* Store Dropdown */
.store-menu-wrapper {
  position: relative;
}

.store-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  width: 260px;
  background-color: #fff;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 10px 25px rgba(0,0,0,0.05);
  z-index: 100;
  animation: fadeInDown 0.2s ease-out;
}

.store-dropdown-title {
  font-weight: 600;
  margin-bottom: 8px;
  color: #1a1a1a;
  font-size: 0.95rem;
}

.store-dropdown-profile {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  padding: 8px;
  background: #f9fafb;
  border-radius: 6px;
  border: 1px solid #f3f4f6;
}

.store-dropdown-logo {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  object-fit: cover;
}

.store-dropdown-name {
  font-size: 0.85rem;
  font-weight: 500;
  color: #374151;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.store-dropdown-desc {
  font-size: 0.85rem;
  color: #666;
  margin: 0 0 16px;
  line-height: 1.4;
}

.store-btn {
  display: block;
  text-align: center;
  width: 100%;
  text-decoration: none;
  padding: 10px;
}

@keyframes fadeInDown {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Enforce Global Navbar Fixed & Layout padding */
.site-header {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  width: 100% !important;
  z-index: 50 !important;
}

.with-header-padding {
  padding-top: 80px; /* mt-20 equivalent */
}
</style>
