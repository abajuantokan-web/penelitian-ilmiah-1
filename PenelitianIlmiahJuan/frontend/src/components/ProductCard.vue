<template>
  <div
    class="product-card"
    ref="cardRef"
    @mouseenter="onMouseEnter"
    @mousemove="onMouseMove"
    @mouseleave="onMouseLeave"
  >
    
    <div class="card-glare" ref="glareRef"></div>

    
    <div class="card-image-wrapper">
      <img
        :src="product.image_url"
        :alt="product.name"
        class="card-image"
        loading="lazy"
        @error="onImageError"
      />
      <div class="card-image-overlay"></div>

      
      <span class="card-region badge badge-amber">
        {{ product.region }}
      </span>

      
      <span v-if="product.category" class="card-category badge badge-indigo">
        {{ product.category }}
      </span>
    </div>

    
    <div class="card-body">
      <h3 class="card-title">{{ product.name }}</h3>

      <p class="card-description">
        {{ truncatedDescription }}
      </p>

      <div class="card-meta">
        <div class="meta-item">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
          </svg>
          <span>{{ product.vendor?.name || 'Vendor NTT' }}</span>
        </div>
        <div class="meta-item">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
            <line x1="16" y1="2" x2="16" y2="6"/>
            <line x1="8" y1="2" x2="8" y2="6"/>
            <line x1="3" y1="10" x2="21" y2="10"/>
          </svg>
          <span>PO {{ product.po_duration }} hari</span>
        </div>
        <div class="meta-item">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          <span>Stok: <strong>{{ product.stock }} pcs</strong></span>
        </div>
      </div>

      <div class="card-footer">
        <div class="card-price">
          <span class="price-label">Harga</span>
          <span class="price-value">{{ formatPrice(product.price) }}</span>
        </div>
        <button 
          class="btn btn-card" 
          :class="product.stock > 0 ? 'btn-primary' : 'btn-secondary'"
          :disabled="product.stock <= 0"
          @click.stop="$emit('order', product)"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="9" cy="21" r="1"/>
            <circle cx="20" cy="21" r="1"/>
            <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
          </svg>
          {{ product.stock > 0 ? 'Pre-Order' : 'Habis' }}
        </button>
      </div>

      <div class="card-min-order">
        Min. order: {{ product.min_order }} pcs
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { animate } from 'animejs'

const props = defineProps({
  product: {
    type: Object,
    required: true
  },
  index: {
    type: Number,
    default: 0
  }
})

defineEmits(['order'])

const cardRef = ref(null)
const glareRef = ref(null)
const isHovering = ref(false)

const truncatedDescription = computed(() => {
  if (!props.product.description) return ''
  return props.product.description.length > 100
    ? props.product.description.substring(0, 100) + '...'
    : props.product.description
})

function formatPrice(price) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(price)
}

function onImageError(e) {
  e.target.src = 'data:image/svg+xml;base64,' + btoa(`
    <svg xmlns="http://www.w3.org/2000/svg" width="400" height="300" viewBox="0 0 400 300">
      <rect fill="#1a1832" width="400" height="300"/>
      <text fill="#6e6a8e" font-family="sans-serif" font-size="16" text-anchor="middle" x="200" y="150">Gambar Produk NTT</text>
    </svg>
  `)
}


function onMouseEnter() {
  isHovering.value = true
  if (cardRef.value) {
    cardRef.value.style.transition = 'none'
  }
}

function onMouseMove(e) {
  if (!cardRef.value || !isHovering.value) return

  const card = cardRef.value
  const rect = card.getBoundingClientRect()

  
  const x = (e.clientX - rect.left) / rect.width
  const y = (e.clientY - rect.top) / rect.height

  
  const rotateX = (0.5 - y) * 24
  const rotateY = (x - 0.5) * 24

  
  card.style.transform = `perspective(1000px) rotateX(${rotateX}deg) rotateY(${rotateY}deg) scale3d(1.03, 1.03, 1.03)`

  
  if (glareRef.value) {
    glareRef.value.style.background = `radial-gradient(circle at ${x * 100}% ${y * 100}%, rgba(255,255,255,0.15) 0%, transparent 60%)`
    glareRef.value.style.opacity = '1'
  }
}

function onMouseLeave() {
  isHovering.value = false
  if (!cardRef.value) return

  
  animate(cardRef.value, {
    rotateX: 0,
    rotateY: 0,
    scale: 1,
    duration: 600,
    ease: 'outElastic(1, 0.5)'
  })

  if (glareRef.value) {
    glareRef.value.style.opacity = '0'
  }
}


onMounted(() => {
  if (cardRef.value) {
    
    cardRef.value.style.opacity = '0'
    cardRef.value.style.transform = 'translateY(60px) scale(0.95)'

    
    animate(cardRef.value, {
      opacity: [0, 1],
      translateY: [60, 0],
      scale: [0.95, 1],
      duration: 800,
      delay: props.index * 120,
      ease: 'outExpo'
    })
  }
})
</script>

<style scoped>
.product-card {
  position: relative;
  background: var(--gradient-card);
  border-radius: var(--radius-xl);
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.06);
  transform-style: preserve-3d;
  will-change: transform;
  cursor: pointer;
  transition: box-shadow var(--transition-base);
}

.product-card:hover {
  box-shadow: var(--shadow-lg), var(--shadow-glow);
}


.card-glare {
  position: absolute;
  inset: 0;
  z-index: 3;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
  border-radius: var(--radius-xl);
}


.card-image-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 4 / 3;
  overflow: hidden;
}

.card-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-slow);
}

.product-card:hover .card-image {
  transform: scale(1.08);
}

.card-image-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    transparent 40%,
    rgba(11, 10, 18, 0.9) 100%
  );
}

.card-region {
  position: absolute;
  top: var(--space-md);
  left: var(--space-md);
  z-index: 2;
}

.card-category {
  position: absolute;
  top: var(--space-md);
  right: var(--space-md);
  z-index: 2;
}


.card-body {
  padding: var(--space-lg);
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.card-title {
  font-size: 1.05rem;
  font-weight: 700;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-description {
  font-size: 0.82rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.card-meta {
  display: flex;
  gap: var(--space-md);
  flex-wrap: wrap;
  margin-top: var(--space-xs);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.meta-item svg {
  color: var(--color-amber);
  flex-shrink: 0;
}


.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: var(--space-sm);
  padding-top: var(--space-sm);
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.price-label {
  display: block;
  font-size: 0.7rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.price-value {
  font-size: 1.1rem;
  font-weight: 800;
  background: var(--gradient-amber);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.btn-card {
  padding: 0.5rem 1rem;
  font-size: 0.8rem;
  border-radius: var(--radius-md);
}

.card-min-order {
  font-size: 0.7rem;
  color: var(--color-text-muted);
  text-align: right;
}
</style>
