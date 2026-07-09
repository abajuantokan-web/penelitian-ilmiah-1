<template>
  <section class="hero" id="hero" aria-label="Hero banner">
    <div class="hero-image">
      <img src="/images/hero.png"
           alt="Seorang perempuan mengenakan kain tenun ikat tradisional NTT dengan latar bebatuan alam"
           loading="eager"
           fetchpriority="high">
    </div>
    <div class="hero-overlay" aria-hidden="true"></div>

    <div class="hero-motif hero-motif--top-left" aria-hidden="true">
      <svg viewBox="0 0 48 48">
        <line x1="24" y1="0" x2="24" y2="48"/>
        <line x1="16" y1="4" x2="16" y2="44"/>
        <line x1="32" y1="4" x2="32" y2="44"/>
        <line x1="8" y1="12" x2="8" y2="36"/>
        <line x1="40" y1="12" x2="40" y2="36"/>
      </svg>
    </div>
    <div class="hero-motif hero-motif--bottom-right" aria-hidden="true">
      <svg viewBox="0 0 48 48">
        <circle cx="24" cy="24" r="20" stroke-dasharray="4 4"/>
        <circle cx="24" cy="24" r="12" stroke-dasharray="2 6"/>
        <circle cx="24" cy="24" r="4"/>
      </svg>
    </div>

    <div class="hero-content container">
      <h1 class="hero-heading reveal">Bawa Produk Autentik<br>Nusa Tenggara Timur<br>ke Rumah Anda</h1>
      <p class="hero-subheading reveal reveal-delay-1">Platform e-commerce pre-order yang menghubungkan Anda langsung dengan berbagai UMKM kreatif NTT. Temukan ragam warisan budaya, mulai dari keindahan tenun ikat legendaris hingga kriya dan kuliner lokal pilihan.</p>
      <a href="#products" class="btn-primary reveal reveal-delay-2" id="hero-cta">Belanja Sekarang</a>
    </div>
  </section>
</template>

<script setup>
</script>
