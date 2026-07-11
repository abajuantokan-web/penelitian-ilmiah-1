<template>
  <article class="product-card" @click="$emit('click', product)">
    <div class="product-card-image">
      <img
        :src="resolvedImageUrl"
        :alt="product.name || 'Product Image'"
        class="w-full h-full object-cover object-center transition-transform duration-300"
        loading="lazy"
        @error="onImageError"
      >
    </div>
    <div class="product-card-info">
      <h3 class="product-card-name">{{ product.name }}</h3>
      <div class="product-seller">
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
          <polyline points="9 22 9 12 15 12 15 22"></polyline>
        </svg>
        <span>{{ storeName }}</span>
      </div>
      <p class="product-card-price">{{ formatPrice(product.price) }}</p>
    </div>
  </article>
</template>

<style scoped>
.product-seller {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.85rem;
  color: #6b7280; /* gray-500 */
  margin-top: 4px;
  margin-bottom: 8px;
}
</style>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  product: {
    type: Object,
    required: true
  }
})

defineEmits(['click'])

// ---------------------------------------------------------------------------
// Image resolution — three-tier priority:
//   1. The product's own image_url (local path → proxied through $getImageUrl,
//      absolute http(s) URL → passed through as-is).
//   2. A high-quality Unsplash fallback chosen by category if the primary
//      URL is empty or a bare local path.
//   3. Last-resort SVG data URI rendered entirely in the browser so there is
//      never a broken-image icon, even without a network connection.
// ---------------------------------------------------------------------------

/** Category → fallback Unsplash URL map */
const CATEGORY_FALLBACKS = {
  'Koleksi Tenun NTT': 'https://images.unsplash.com/photo-1596755094514-f87e34085b2c?q=80&w=800&auto=format&fit=crop',
  'Cita Rasa Lokal':   'https://images.unsplash.com/photo-1447933601403-0c6688de566e?q=80&w=800&auto=format&fit=crop',
  'Koleksi Aksesoris': 'https://images.unsplash.com/photo-1515562141207-7a88fb7ce338?q=80&w=800&auto=format&fit=crop',
}

/** Minimal SVG placeholder — shown only when the network image also fails */
const SVG_PLACEHOLDER =
  "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='400' height='400' viewBox='0 0 400 400'%3E%3Crect width='400' height='400' fill='%23f3f0eb'/%3E%3Ctext x='200' y='210' text-anchor='middle' font-family='sans-serif' font-size='14' fill='%23a8956e'%3EOpenPeo%3C/text%3E%3C/svg%3E"

/** The resolved URL fed to :src */
const resolvedImageUrl = computed(() => {
  const raw = props.product?.image_url

  // ── Priority 1: absolute remote URL (Unsplash, CDN, etc.) ────────────────
  if (raw && (raw.startsWith('http://') || raw.startsWith('https://'))) {
    return raw
  }

  // ── Priority 2: local public path e.g. /images/kopi-flores.png ───────────
  // Vite serves frontend/public/ at root '/', so these load directly.
  if (raw && (raw.startsWith('/images/') || raw.startsWith('images/'))) {
    return raw
  }

  // ── Priority 3: category-based Unsplash fallback ─────────────────────────
  const categoryFallback = CATEGORY_FALLBACKS[props.product?.category]
  if (categoryFallback) return categoryFallback

  // ── Priority 4: last-resort inline SVG (works fully offline) ─────────────
  return SVG_PLACEHOLDER
})

/** Called by @error when the browser fails to load the resolved URL */
const onImageError = (event) => {
  // Prevent infinite retry loops if the fallback itself 404s
  if (event.target.src !== SVG_PLACEHOLDER) {
    event.target.src = SVG_PLACEHOLDER
  }
}

// ---------------------------------------------------------------------------

const storeName = computed(() => {
  if (props.product.seller_profile?.store_name) {
    return props.product.seller_profile.store_name
  }
  return 'Toko Tidak Diketahui'
  return 'Toko Anonim'
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}
</script>
