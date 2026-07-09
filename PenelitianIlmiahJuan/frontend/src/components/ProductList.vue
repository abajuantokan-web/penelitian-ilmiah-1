<template>
  <div class="product-list">
    <!-- Section Header -->
    <div class="list-header">
      <div class="header-text">
        <span class="badge badge-amber">Katalog Produk</span>
        <h2 class="section-title text-display">
          Produk Pre-Order <span class="text-gradient">Terbaik NTT</span>
        </h2>
        <p class="section-desc">
          Telusuri koleksi produk tradisional dari berbagai daerah di Nusa Tenggara Timur
        </p>
      </div>
    </div>

    <!-- Region Filter Tabs -->
    <div class="filter-bar">
      <div class="filter-tabs">
        <button
          v-for="tab in regionTabs"
          :key="tab.value"
          class="filter-tab"
          :class="{ active: activeRegion === tab.value }"
          @click="setRegion(tab.value)"
        >
          <span class="tab-icon">{{ tab.icon }}</span>
          <span class="tab-label">{{ tab.label }}</span>
        </button>
      </div>

      <div class="filter-search">
        <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <path d="m21 21-4.35-4.35"/>
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          class="input search-input"
          placeholder="Cari produk..."
          @input="debouncedSearch"
        />
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="loading-grid">
        <div v-for="n in 6" :key="n" class="skeleton-card">
          <div class="skeleton-image skeleton-animate"></div>
          <div class="skeleton-body">
            <div class="skeleton-line skeleton-animate" style="width: 80%"></div>
            <div class="skeleton-line skeleton-animate" style="width: 60%"></div>
            <div class="skeleton-line skeleton-animate" style="width: 40%"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Product Grid -->
    <div v-else-if="products.length > 0" class="product-grid" ref="gridRef">
      <ProductCard
        v-for="(product, index) in products"
        :key="product.id"
        :product="product"
        :index="index"
        @order="(p) => $emit('order', p)"
      />
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state">
      <div class="empty-icon">🏝️</div>
      <h3>Belum ada produk</h3>
      <p>Produk dari daerah ini belum tersedia. Coba jelajahi daerah lain!</p>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="pagination">
      <button
        class="btn btn-secondary btn-page"
        :disabled="currentPage <= 1"
        @click="goToPage(currentPage - 1)"
      >
        ← Sebelumnya
      </button>
      <span class="page-info">
        Halaman {{ currentPage }} dari {{ totalPages }}
      </span>
      <button
        class="btn btn-secondary btn-page"
        :disabled="currentPage >= totalPages"
        @click="goToPage(currentPage + 1)"
      >
        Selanjutnya →
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import ProductCard from './ProductCard.vue'

const emit = defineEmits(['order'])

const API_BASE = 'http://localhost:8080/api'

// Region filter tabs with NTT regions
const regionTabs = [
  { value: '',           label: 'Semua',      icon: '🌏' },
  { value: 'Sumba',      label: 'Sumba',      icon: '🏝️' },
  { value: 'Manggarai',  label: 'Manggarai',  icon: '⛰️'  },
  { value: 'Flores',     label: 'Flores',     icon: '🌺' },
  { value: 'Kupang',     label: 'Kupang',     icon: '🏙️'  },
  { value: 'Timor',      label: 'Timor',      icon: '🌿' },
]

const products = ref([])
const loading = ref(true)
const activeRegion = ref('')
const searchQuery = ref('')
const currentPage = ref(1)
const totalPages = ref(1)
const gridRef = ref(null)

let searchTimeout = null

async function fetchProducts() {
  loading.value = true

  try {
    const params = new URLSearchParams()
    params.set('page', currentPage.value.toString())
    params.set('limit', '12')

    if (activeRegion.value) {
      params.set('region', activeRegion.value)
    }
    if (searchQuery.value.trim()) {
      params.set('search', searchQuery.value.trim())
    }

    const response = await fetch(`${API_BASE}/products?${params}`)
    const data = await response.json()

    if (data.success) {
      products.value = data.data || []
      totalPages.value = data.meta?.total_pages || 1
    } else {
      products.value = []
    }
  } catch (error) {
    console.warn('Backend not available, loading demo data:', error.message)
    loadDemoProducts()
  } finally {
    loading.value = false
  }
}

// Fallback demo data when backend is not running
function loadDemoProducts() {
  const demo = [
    {
      id: 1, vendor_id: 2,
      vendor: { name: 'Maria Tenun Sumba' },
      name: 'Kain Tenun Ikat Sumba Pahikung',
      description: 'Kain tenun ikat tradisional Sumba Timur bermotif Pahikung yang melambangkan kebangsawanan dan kekuatan.',
      price: 1500000,
      image_url: 'https://images.unsplash.com/photo-1558171813-4c088753af8f?w=400',
      region: 'Sumba', category: 'Tenun', min_order: 2, po_duration: 30
    },
    {
      id: 2, vendor_id: 2,
      vendor: { name: 'Maria Tenun Sumba' },
      name: 'Selimut Tenun Sumba Hinggi Kombu',
      description: 'Selimut tenun Hinggi Kombu khas Sumba dengan motif hewan dan geometris.',
      price: 2500000,
      image_url: 'https://images.unsplash.com/photo-1606722590583-6951b5ea92ad?w=400',
      region: 'Sumba', category: 'Tenun', min_order: 1, po_duration: 45
    },
    {
      id: 3, vendor_id: 3,
      vendor: { name: 'Yohanes Manggarai Craft' },
      name: 'Songke Manggarai — Kain Adat',
      description: 'Kain tenun songke khas Manggarai dengan motif wela mpuu (bunga penuh).',
      price: 850000,
      image_url: 'https://images.unsplash.com/photo-1621600411688-4be93cd68504?w=400',
      region: 'Manggarai', category: 'Tenun', min_order: 3, po_duration: 21
    },
    {
      id: 4, vendor_id: 3,
      vendor: { name: 'Yohanes Manggarai Craft' },
      name: 'Kopi Flores Bajawa Single Origin',
      description: 'Kopi arabika single origin dari dataran tinggi Bajawa dengan cita rasa coklat dan karamel.',
      price: 120000,
      image_url: 'https://images.unsplash.com/photo-1447933601403-0c6688de566e?w=400',
      region: 'Flores', category: 'Kuliner', min_order: 5, po_duration: 14
    },
    {
      id: 5, vendor_id: 2,
      vendor: { name: 'Maria Tenun Sumba' },
      name: 'Manik-Manik Sumba Handmade Necklace',
      description: 'Kalung manik-manik handmade khas Sumba dengan warna-warna tanah.',
      price: 350000,
      image_url: 'https://images.unsplash.com/photo-1611085583191-a3b181a88401?w=400',
      region: 'Sumba', category: 'Aksesoris', min_order: 3, po_duration: 14
    },
    {
      id: 6, vendor_id: 3,
      vendor: { name: 'Yohanes Manggarai Craft' },
      name: 'Madu Hutan Timor Asli',
      description: 'Madu hutan asli dari pegunungan Timor. Dipanen langsung dari sarang lebah liar.',
      price: 175000,
      image_url: 'https://images.unsplash.com/photo-1587049352846-4a222e784d38?w=400',
      region: 'Timor', category: 'Kuliner', min_order: 4, po_duration: 10
    },
    {
      id: 7, vendor_id: 3,
      vendor: { name: 'Yohanes Manggarai Craft' },
      name: 'Patung Ukir Kayu Sandalwood Kupang',
      description: 'Patung ukiran kayu cendana buatan tangan pengrajin Kupang.',
      price: 950000,
      image_url: 'https://images.unsplash.com/photo-1513519245088-0e12902e35ca?w=400',
      region: 'Kupang', category: 'Kerajinan', min_order: 1, po_duration: 30
    },
    {
      id: 8, vendor_id: 2,
      vendor: { name: 'Maria Tenun Sumba' },
      name: 'Syal Tenun Ende Lio',
      description: 'Syal tenun ikat Ende Lio dengan motif bunga dan geometris.',
      price: 450000,
      image_url: 'https://images.unsplash.com/photo-1601924921557-45e16393d8e1?w=400',
      region: 'Flores', category: 'Tenun', min_order: 2, po_duration: 21
    },
  ]

  let filtered = demo
  if (activeRegion.value) {
    filtered = filtered.filter(p => p.region === activeRegion.value)
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase()
    filtered = filtered.filter(p =>
      p.name.toLowerCase().includes(q) ||
      p.description.toLowerCase().includes(q)
    )
  }

  products.value = filtered
  totalPages.value = 1
}

function setRegion(region) {
  activeRegion.value = region
  currentPage.value = 1
  fetchProducts()
}

function debouncedSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    fetchProducts()
  }, 400)
}

function goToPage(page) {
  currentPage.value = page
  fetchProducts()
}

onMounted(() => {
  fetchProducts()
})
</script>

<style scoped>
.product-list {
  width: 100%;
}

/* Header */
.list-header {
  text-align: center;
  margin-bottom: var(--space-2xl);
}

.section-title {
  font-size: clamp(1.8rem, 3.5vw, 2.6rem);
  margin: var(--space-md) 0;
  line-height: 1.2;
}

.section-desc {
  color: var(--color-text-secondary);
  font-size: 1rem;
  max-width: 550px;
  margin: 0 auto;
}

/* Filter Bar */
.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-lg);
  margin-bottom: var(--space-2xl);
  flex-wrap: wrap;
}

.filter-tabs {
  display: flex;
  gap: var(--space-xs);
  flex-wrap: wrap;
}

.filter-tab {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.55rem 1rem;
  border-radius: var(--radius-full);
  background: transparent;
  color: var(--color-text-secondary);
  font-size: 0.82rem;
  font-weight: 500;
  border: 1px solid transparent;
  transition: all var(--transition-fast);
}

.filter-tab:hover {
  background: rgba(255, 255, 255, 0.04);
  color: var(--color-text-primary);
}

.filter-tab.active {
  background: var(--color-amber-glow);
  color: var(--color-amber);
  border-color: rgba(245, 166, 35, 0.3);
}

.tab-icon {
  font-size: 1rem;
}

.filter-search {
  position: relative;
  min-width: 220px;
}

.search-icon {
  position: absolute;
  left: 0.85rem;
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-muted);
  pointer-events: none;
}

.search-input {
  padding-left: 2.5rem;
  border-radius: var(--radius-full);
  font-size: 0.85rem;
}

/* Product Grid */
.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-lg);
}

/* Skeleton Loading */
.loading-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-lg);
}

.skeleton-card {
  background: var(--gradient-card);
  border-radius: var(--radius-xl);
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.skeleton-image {
  aspect-ratio: 4 / 3;
  background: var(--color-surface);
}

.skeleton-body {
  padding: var(--space-lg);
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.skeleton-line {
  height: 14px;
  border-radius: var(--radius-sm);
  background: var(--color-surface);
}

.skeleton-animate {
  background: linear-gradient(90deg, var(--color-surface) 25%, var(--color-surface-hover) 50%, var(--color-surface) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: var(--space-3xl) var(--space-xl);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: var(--space-md);
}

.empty-state h3 {
  font-size: 1.3rem;
  margin-bottom: var(--space-sm);
}

.empty-state p {
  color: var(--color-text-secondary);
  font-size: 0.95rem;
}

/* Pagination */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-md);
  margin-top: var(--space-2xl);
}

.btn-page {
  padding: 0.5rem 1rem;
  font-size: 0.82rem;
}

.btn-page:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-info {
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

/* Responsive */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-tabs {
    justify-content: center;
  }

  .product-grid,
  .loading-grid {
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  }
}
</style>
