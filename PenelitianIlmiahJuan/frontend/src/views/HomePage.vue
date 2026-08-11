<template>
  <div class="home-page">
    
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
          <a href="#beranda" class="nav-link active">Beranda</a>
          <a href="#produk" class="nav-link">Produk</a>
          <a href="#tentang" class="nav-link">Tentang</a>
        </div>

        <div class="navbar-actions">
          <router-link to="/login" class="btn btn-secondary btn-sm" v-if="!currentUser">
            Masuk
          </router-link>
          <template v-else>
            <span class="user-greeting">Halo, <strong>{{ currentUser.name || currentUser.username }}</strong></span>
            <router-link :to="currentUser.role === 'admin' ? '/admin/dashboard' : '/shop'" class="btn btn-primary btn-sm">
              {{ currentUser.role === 'admin' ? 'Panel Admin' : 'Katalog' }}
            </router-link>
            <button class="btn btn-secondary btn-sm" @click="logout">Keluar</button>
          </template>
        </div>
      </div>
    </nav>

    
    <section id="beranda" class="hero">
      <div class="hero-bg">
        <div class="hero-orb hero-orb-1"></div>
        <div class="hero-orb hero-orb-2"></div>
        <div class="hero-orb hero-orb-3"></div>
        <div class="hero-grid"></div>
      </div>
      <div class="hero-content">
        <div class="hero-badge badge badge-amber">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"/></svg>
          Marketplace Pre-Order #1 NTT
        </div>
        <h1 class="hero-title">
          <span class="hero-title-line">Warisan Budaya</span>
          <span class="hero-title-line text-gradient">Nusa Tenggara Timur</span>
          <span class="hero-title-line">di Ujung Jari Anda</span>
        </h1>
        <p class="hero-subtitle">
          Temukan keindahan tenun ikat, kopi Flores, madu Timor, dan kerajinan tangan
          autentik langsung dari pengrajin lokal NTT melalui sistem pre-order eksklusif.
        </p>
        <div class="hero-actions">
          <router-link :to="!currentUser ? '/login' : (currentUser.role === 'admin' ? '/admin/dashboard' : '/shop')" class="btn btn-primary btn-lg">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/></svg>
            Mulai Belanja
          </router-link>
          <a href="#tentang" class="btn btn-secondary btn-lg">Pelajari Lebih Lanjut</a>
        </div>
        <div class="hero-stats">
          <div class="stat-item"><span class="stat-value text-gradient">500+</span><span class="stat-label">Produk NTT</span></div>
          <div class="stat-divider"></div>
          <div class="stat-item"><span class="stat-value text-gradient">200+</span><span class="stat-label">Pengrajin Lokal</span></div>
          <div class="stat-divider"></div>
          <div class="stat-item"><span class="stat-value text-gradient">22</span><span class="stat-label">Kabupaten NTT</span></div>
        </div>
      </div>
    </section>

    
    <section id="produk" class="section-products">
      <ProductList @order="openOrderModal" />
    </section>

    
    <section id="tentang" class="section-about">
      <div class="about-container">
        <div class="about-header">
          <span class="badge badge-indigo">Tentang OpenPeo</span>
          <h2 class="section-title text-display">Menghubungkan <span class="text-gradient">NTT</span> dengan Dunia</h2>
          <p class="section-desc">OpenPeo adalah platform marketplace khusus yang menghubungkan pengrajin dan produsen lokal NTT dengan pembeli di seluruh Indonesia melalui sistem pre-order yang adil dan transparan.</p>
        </div>
        <div class="about-features">
          <div class="feature-card glass"><div class="feature-icon">🏝️</div><h3>Autentik NTT</h3><p>Produk langsung dari sumber terpercaya di Sumba, Manggarai, Flores, Kupang, dan Timor.</p></div>
          <div class="feature-card glass"><div class="feature-icon">📦</div><h3>Pre-Order Fair</h3><p>Sistem pre-order yang melindungi pengrajin dan memastikan kualitas terjaga.</p></div>
          <div class="feature-card glass"><div class="feature-icon">💬</div><h3>Chat Langsung</h3><p>Komunikasi real-time dengan vendor untuk diskusi produk dan pesanan.</p></div>
        </div>
      </div>
    </section>

    <footer class="footer">
      <div class="footer-container">
        <div class="footer-brand"><span class="brand-name text-gradient">OpenPeo</span><p>© 2026 OpenPeo. Platform Marketplace Produk Tradisional NTT.</p></div>
        <div class="footer-links"><span>Dibuat dengan ❤️ dari Nusa Tenggara Timur</span></div>
      </div>
    </footer>

    <OrderModal v-if="selectedProduct" :product="selectedProduct" @close="selectedProduct = null" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import ProductList from '../components/ProductList.vue'
import OrderModal from '../components/OrderModal.vue'

const router = useRouter()
const selectedProduct = ref(null)

const currentUser = computed(() => {
  return JSON.parse(localStorage.getItem('openpeo_user') || 'null')
})

function openOrderModal(product) { selectedProduct.value = product }

function logout() {
  localStorage.removeItem('openpeo_user')
  router.push('/')
}
</script>

<style scoped>
.navbar { position: fixed; top: 0; left: 0; right: 0; z-index: 1000; padding: 0.75rem 0; }
.navbar-container { max-width: 1280px; margin: 0 auto; padding: 0 var(--space-xl); display: flex; align-items: center; justify-content: space-between; }
.navbar-brand { display: flex; align-items: center; gap: var(--space-sm); }
.brand-icon svg { width: 36px; height: 36px; }
.brand-text { display: flex; flex-direction: column; line-height: 1.1; }
.brand-name { font-family: var(--font-display); font-size: 1.3rem; font-weight: 800; }
.brand-tagline { font-size: 0.65rem; color: var(--color-text-muted); letter-spacing: 0.15em; text-transform: uppercase; }
.navbar-links { display: flex; gap: var(--space-lg); }
.nav-link { color: var(--color-text-secondary); font-size: 0.875rem; font-weight: 500; padding: 0.5rem 0; position: relative; transition: color var(--transition-fast); }
.nav-link::after { content: ''; position: absolute; bottom: 0; left: 0; width: 0; height: 2px; background: var(--gradient-amber); border-radius: 1px; transition: width var(--transition-base); }
.nav-link:hover, .nav-link.active { color: var(--color-text-primary); }
.nav-link:hover::after, .nav-link.active::after { width: 100%; }
.navbar-actions { display: flex; align-items: center; gap: var(--space-sm); }
.user-greeting { font-size: 0.82rem; color: var(--color-text-secondary); }
.btn-sm { padding: 0.5rem 1rem; font-size: 0.8rem; }
.btn-lg { padding: 0.875rem 1.75rem; font-size: 0.95rem; }
.hero { min-height: 100vh; display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden; padding: 6rem var(--space-xl) var(--space-3xl); }
.hero-bg { position: absolute; inset: 0; pointer-events: none; }
.hero-orb { position: absolute; border-radius: 50%; filter: blur(100px); opacity: 0.3; }
.hero-orb-1 { width: 600px; height: 600px; background: radial-gradient(circle, var(--color-indigo-glow) 0%, transparent 70%); top: -200px; right: -200px; animation: float 8s ease-in-out infinite; }
.hero-orb-2 { width: 500px; height: 500px; background: radial-gradient(circle, var(--color-burnt-orange) 0%, transparent 70%); bottom: -150px; left: -150px; animation: float 10s ease-in-out infinite reverse; }
.hero-orb-3 { width: 300px; height: 300px; background: radial-gradient(circle, var(--color-amber) 0%, transparent 70%); top: 40%; left: 50%; transform: translateX(-50%); opacity: 0.15; animation: float 6s ease-in-out infinite 2s; }
.hero-grid { position: absolute; inset: 0; background-image: linear-gradient(rgba(255,255,255,0.02) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.02) 1px, transparent 1px); background-size: 60px 60px; mask-image: radial-gradient(ellipse at center, black 20%, transparent 70%); -webkit-mask-image: radial-gradient(ellipse at center, black 20%, transparent 70%); }
.hero-content { position: relative; z-index: 1; text-align: center; max-width: 800px; }
.hero-badge { margin-bottom: var(--space-xl); display: inline-flex; gap: var(--space-xs); animation: fadeInUp 0.6s ease-out; }
.hero-title { font-family: var(--font-display); font-size: clamp(2.2rem, 5vw, 3.8rem); font-weight: 800; line-height: 1.15; margin-bottom: var(--space-lg); animation: fadeInUp 0.6s ease-out 0.1s both; }
.hero-title-line { display: block; }
.hero-subtitle { font-size: 1.05rem; color: var(--color-text-secondary); line-height: 1.7; max-width: 600px; margin: 0 auto var(--space-xl); animation: fadeInUp 0.6s ease-out 0.2s both; }
.hero-actions { display: flex; gap: var(--space-md); justify-content: center; flex-wrap: wrap; margin-bottom: var(--space-3xl); animation: fadeInUp 0.6s ease-out 0.3s both; }
.hero-stats { display: flex; align-items: center; justify-content: center; gap: var(--space-xl); animation: fadeInUp 0.6s ease-out 0.4s both; }
.stat-item { display: flex; flex-direction: column; align-items: center; }
.stat-value { font-family: var(--font-display); font-size: 1.8rem; font-weight: 800; }
.stat-label { font-size: 0.8rem; color: var(--color-text-muted); margin-top: 0.2rem; }
.stat-divider { width: 1px; height: 40px; background: rgba(255,255,255,0.1); }
.section-products { padding: var(--space-3xl) var(--space-xl); max-width: 1360px; margin: 0 auto; }
.section-about { padding: var(--space-3xl) var(--space-xl); background: var(--color-bg-elevated); border-top: 1px solid rgba(255,255,255,0.04); }
.about-container { max-width: 1100px; margin: 0 auto; }
.about-header { text-align: center; margin-bottom: var(--space-3xl); }
.section-title { font-size: clamp(1.8rem, 3.5vw, 2.6rem); margin: var(--space-md) 0; line-height: 1.2; }
.section-desc { color: var(--color-text-secondary); max-width: 650px; margin: 0 auto; font-size: 1rem; line-height: 1.7; }
.about-features { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: var(--space-lg); }
.feature-card { padding: var(--space-xl); border-radius: var(--radius-xl); text-align: center; transition: all var(--transition-base); }
.feature-card:hover { transform: translateY(-4px); box-shadow: var(--shadow-glow); }
.feature-icon { font-size: 2.5rem; margin-bottom: var(--space-md); }
.feature-card h3 { font-size: 1.15rem; margin-bottom: var(--space-sm); font-weight: 700; }
.feature-card p { color: var(--color-text-secondary); font-size: 0.9rem; line-height: 1.6; }
.footer { padding: var(--space-xl); border-top: 1px solid rgba(255,255,255,0.04); }
.footer-container { max-width: 1280px; margin: 0 auto; display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: var(--space-md); }
.footer-brand .brand-name { font-size: 1.1rem; margin-bottom: 0.25rem; display: block; }
.footer-brand p, .footer-links span { font-size: 0.8rem; color: var(--color-text-muted); }
@media (max-width: 768px) { .navbar-links { display: none; } .hero-stats { gap: var(--space-md); } .stat-value { font-size: 1.4rem; } .footer-container { flex-direction: column; text-align: center; } }
</style>
