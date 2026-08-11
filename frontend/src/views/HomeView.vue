<template>
  <div>
    <Hero />
    
    <ProductGrid 
      id="products" 
      title="Koleksi Tenun NTT" 
      :products="productStore.tenunProducts.slice(0, 8)"
      @product-click="openProductModal"
    />
    
    <ImageDivider 
      src="/images/divider1.png"
      alt="Petani Kopi Flores"
      title="Dari Alam ke Cangkir Anda"
      subtitle="Kopi Flores Bajawa Single Origin"
    />
    
    <ProductGrid 
      title="Cita Rasa Lokal" 
      :products="productStore.foodProducts.slice(0, 8)"
      :altBg="true"
      @product-click="openProductModal"
    />
    
    <ImageDivider 
      src="/images/divider2.png"
      alt="Aksesoris Tradisional NTT"
      title="Warisan dalam Genggaman"
      subtitle="Koleksi Aksesoris Premium"
    />
    
    <ProductGrid 
      title="Koleksi Aksesoris" 
      :products="productStore.accessoriesProducts.slice(0, 8)"
      @product-click="openProductModal"
    />

    
    <section id="cerita-kami" class="story-section py-24">
      <div class="container story-container">
        <div class="story-image">
          
          <img src="/images/divider1.png" alt="Pengrajin Tenun NTT" />
        </div>
        <div class="story-content">
          <h2 class="story-title">Menghubungkan Warisan & Masa Depan</h2>
          <p class="story-text">OpenPeo hadir bukan sekadar sebagai platform belanja, melainkan sebagai jembatan yang menghubungkan karya tangan autentik para pengrajin Nusa Tenggara Timur dengan apresiasi global. Kami percaya setiap helaian benang dan setiap biji kopi memiliki cerita yang layak didengar dan dilestarikan.</p>
          <router-link to="/koleksi" class="btn-primary story-btn">Eksplorasi Koleksi</router-link>
        </div>
      </div>
    </section>

    
    <FloatingButtons />

    
    <ProductDetailModal 
      :is-open="isModalOpen" 
      :product="selectedProduct"
      @close="isModalOpen = false"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useProductStore } from '../stores/products'
import { useAuthStore } from '../stores/auth'
import Hero from '../components/Hero.vue'
import ProductGrid from '../components/ProductGrid.vue'
import ImageDivider from '../components/ImageDivider.vue'
import FloatingButtons from '../components/FloatingButtons.vue'
import ProductDetailModal from '../components/ProductDetailModal.vue'

const productStore = useProductStore()
const authStore = useAuthStore()

const isModalOpen = ref(false)
const selectedProduct = ref(null)

const openProductModal = (product) => {
  selectedProduct.value = product
  isModalOpen.value = true
}

onMounted(async () => {
  await productStore.fetchProducts()
  
  const pendingAction = localStorage.getItem('pendingProductAction')
  if (pendingAction && authStore.isAuthenticated) {
    const { productId } = JSON.parse(pendingAction)
    localStorage.removeItem('pendingProductAction')
    
    const targetProduct = productStore.allProducts.find(p => p.id === productId)
    if (targetProduct) {
       openProductModal(targetProduct)
    }
  }
})
</script>

<style scoped>
.story-section {
  background-color: #fff;
  padding-top: 96px;
  padding-bottom: 96px;
}

.story-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 64px;
  align-items: center;
}

.story-image {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
  aspect-ratio: 4/5;
}

.story-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.7s ease;
}

.story-image:hover img {
  transform: scale(1.05);
}

.story-content {
  padding-right: 40px;
}

.story-title {
  font-family: 'Playfair Display', serif;
  font-size: 2.75rem;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 24px;
  line-height: 1.2;
}

.story-text {
  font-family: 'Montserrat', sans-serif;
  font-size: 1.1rem;
  color: #555;
  line-height: 1.8;
  margin-bottom: 40px;
}

.story-btn {
  display: inline-block;
  padding: 14px 32px;
  font-size: 0.95rem;
  letter-spacing: 1px;
}

@media (max-width: 992px) {
  .story-container {
    grid-template-columns: 1fr;
    gap: 40px;
  }
  
  .story-content {
    padding-right: 0;
    text-align: center;
    order: -1;
  }
  
  .story-image {
    aspect-ratio: 16/9;
  }
}
</style>
