<template>
  <div class="koleksi-page">
    <div class="koleksi-header">
      <div class="container">
        <h1>Koleksi Lengkap OpenPeo</h1>
        <p>Jelajahi seluruh koleksi warisan Nusa Tenggara Timur, mulai dari tenun autentik, cita rasa lokal, hingga aksesoris tradisional.</p>
        
        <div class="search-container">
          <input 
            type="text" 
            v-model="searchQuery" 
            ref="searchInput"
            class="search-input" 
            placeholder="Cari produk autentik NTT..."
            aria-label="Cari produk"
          />
          <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
        </div>
      </div>
    </div>
    
    <div v-if="filteredTenun.length === 0 && filteredFood.length === 0 && filteredAccessories.length === 0" class="no-results">
      <p>Pencarian untuk "{{ searchQuery }}" tidak ditemukan.</p>
    </div>

    <ProductGrid 
      v-if="filteredTenun.length > 0"
      id="tenun" 
      title="Koleksi Tenun NTT" 
      :products="filteredTenun"
      @product-click="openProductModal"
    />
    
    <ImageDivider 
      v-if="filteredFood.length > 0 && !searchQuery"
      src="/images/divider1.png"
      alt="Petani Kopi Flores"
      title="Dari Alam ke Cangkir Anda"
      subtitle="Kopi Flores Bajawa Single Origin"
    />
    
    <ProductGrid 
      v-if="filteredFood.length > 0"
      id="kuliner"
      title="Cita Rasa Lokal" 
      :products="filteredFood"
      :altBg="true"
      @product-click="openProductModal"
    />
    
    <ImageDivider 
      v-if="filteredAccessories.length > 0 && !searchQuery"
      src="/images/divider2.png"
      alt="Aksesoris Tradisional NTT"
      title="Warisan dalam Genggaman"
      subtitle="Koleksi Aksesoris Premium"
    />
    
    <ProductGrid 
      v-if="filteredAccessories.length > 0"
      id="aksesoris"
      title="Koleksi Aksesoris" 
      :products="filteredAccessories"
      @product-click="openProductModal"
    />

    <!-- Floating Buttons -->
    <FloatingButtons />

    <!-- Product Detail Modal -->
    <ProductDetailModal 
      :is-open="isModalOpen" 
      :product="selectedProduct"
      @close="isModalOpen = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useProductStore } from '../stores/products'
import ProductGrid from '../components/ProductGrid.vue'
import ImageDivider from '../components/ImageDivider.vue'
import FloatingButtons from '../components/FloatingButtons.vue'
import ProductDetailModal from '../components/ProductDetailModal.vue'

const route = useRoute()
const productStore = useProductStore()

const isModalOpen = ref(false)
const selectedProduct = ref(null)

const searchQuery = ref('')
const searchInput = ref(null)

const filteredTenun = computed(() => filterProducts(productStore.tenunProducts))
const filteredFood = computed(() => filterProducts(productStore.foodProducts))
const filteredAccessories = computed(() => filterProducts(productStore.accessoriesProducts))

const filterProducts = (products) => {
  if (!searchQuery.value) return products
  const query = searchQuery.value.toLowerCase()
  return products.filter(p => 
    p.name.toLowerCase().includes(query) || 
    p.description.toLowerCase().includes(query)
  )
}

const openProductModal = (product) => {
  selectedProduct.value = product
  isModalOpen.value = true
}

onMounted(() => {
  productStore.fetchProducts()
  
  if (route.query.search) {
    nextTick(() => {
      if (searchInput.value) searchInput.value.focus()
    })
  }
})
</script>

<style scoped>
.koleksi-page {
  background-color: #fafafa;
  min-height: 100vh;
}

.koleksi-header {
  padding: 60px 20px;
  text-align: center;
  background-color: #ffffff;
  border-bottom: 1px solid #eaeaea;
}

.koleksi-header h1 {
  font-family: 'Playfair Display', serif;
  font-size: 2.5rem;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 16px;
}

.koleksi-header p {
  font-family: 'Montserrat', sans-serif;
  color: #666;
  font-size: 1rem;
  line-height: 1.6;
  max-width: 600px;
  margin: 0 auto 32px;
}

.search-container {
  position: relative;
  max-width: 480px;
  margin: 0 auto;
}

.search-input {
  width: 100%;
  padding: 16px 24px 16px 52px;
  border-radius: 40px;
  border: 1px solid #eaeaea;
  font-size: 1rem;
  font-family: 'Inter', sans-serif;
  color: #333;
  outline: none;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0,0,0,0.02);
}

.search-input:focus {
  border-color: #d4af37;
  box-shadow: 0 4px 20px rgba(212, 175, 55, 0.15);
}

.search-icon {
  position: absolute;
  left: 20px;
  top: 50%;
  transform: translateY(-50%);
  width: 20px;
  height: 20px;
  fill: none;
  stroke: #9ca3af;
  stroke-width: 2;
  pointer-events: none;
}

.no-results {
  text-align: center;
  padding: 60px 20px;
  font-size: 1.1rem;
  color: #666;
  min-height: 40vh;
}

@media (max-width: 768px) {
  .koleksi-header h1 {
    font-size: 2rem;
  }
}
</style>
