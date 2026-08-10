<template>
  <div id="app-root" >
    <!-- Header -->
    <header :class="headerClass">
      <div class="container header-inner">
        <div class="logo">
          <router-link to="/">OpenPeo</router-link>
        </div>

        <nav class="nav-links" aria-label="Main Navigation">
          <router-link to="/">Dashboard Seller</router-link>
        </nav>

        <div class="header-icons">
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
        </div>
      </div>
    </header>

    <!-- Router View for Pages -->
    <main class="main-content" :class="{ 'with-header-padding': route.path !== '/' }">
      <router-view />
    </main>



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
import { useDashboardStore } from './stores/dashboard'
import { useWebsocketStore } from './stores/websocket'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const dashboardStore = useDashboardStore()
const websocketStore = useWebsocketStore() // Initializes singleton WS

const isScrolled = ref(false)

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

.icon-btn--active {
  /* subtle visual for logged-in state */
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
