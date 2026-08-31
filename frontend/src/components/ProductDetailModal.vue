<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="close">
    <div class="modal-content">
      <button class="modal-close" @click="close" aria-label="Tutup">&times;</button>
      
      <div v-if="product" class="pdp-grid">
        
        <div class="pdp-gallery">
          <div class="pdp-thumbnails">
            <div class="pdp-thumb active">
              <img :src="productImageUrl" :alt="product.name">
            </div>
            <div class="pdp-thumb" v-for="i in 3" :key="i">
               <img :src="productImageUrl" :alt="product.name">
            </div>
          </div>
          <div class="pdp-main-image">
            <img :src="productImageUrl" :alt="product.name">
          </div>
        </div>

        
        <div class="pdp-info">
          <p class="pdp-brand">{{ product.region }} &bull; {{ product.category }}</p>
          <h2 class="pdp-title">{{ product.name }}</h2>
          <p class="pdp-price">{{ formatPrice(product.price) }}</p>
          
          <div class="pdp-preorder-badge">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right:8px"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            Estimasi Pre-order: {{ product.po_duration || 14 }} hari kerja
          </div>

          <div class="pdp-quantity-selector">
            <label>Kuantitas:</label>
            <div class="qty-controls">
              <button @click="quantity > 1 ? quantity-- : null">-</button>
              <input type="number" v-model.number="quantity" min="1" readonly>
              <button @click="quantity++">+</button>
            </div>
          </div>

          <div class="pdp-actions">
            <button class="btn-primary--dark" @click="handleAddToCart">
              <span v-if="isAdding">Menambahkan...</span>
              <span v-else>Add to Cart</span>
            </button>
            <button class="btn-primary--solid-dark" @click="handleBuyNow">
              Buy It Now
            </button>
          </div>
          
          <button class="btn-outline w-full mt-4" style="display:flex; justify-content:center; align-items:center; gap:8px;" @click="handleChat">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
            Chat Penjual
          </button>

          <div class="pdp-description">
            <p>{{ product.description }}</p>
            <p v-if="isTenun" class="pdp-custom-note">
              <strong>Tersedia layanan custom size (request ukuran sesuai keinginan Anda).</strong>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import { useChatStore } from '../stores/chat'
import { getImageUrl } from '../utils/imageUtils'


const props = defineProps({
  isOpen: Boolean,
  product: Object
})

const emit = defineEmits(['close'])
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const cartStore = useCartStore()
const chatStore = useChatStore()

const quantity = ref(1)
const isAdding = ref(false)


watch(() => props.isOpen, (newVal) => {
  if (newVal) quantity.value = 1
})

const productImageUrl = computed(() => getImageUrl(props.product?.image_url))

const isTenun = computed(() => props.product?.category === 'Koleksi Tenun NTT' || props.product?.category === 'Tenun')

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

const close = () => {
  emit('close')
}

const handleAddToCart = async () => {
  if (!authStore.isAuthenticated) {
    localStorage.setItem('pendingProductAction', JSON.stringify({ productId: props.product.id }));
    close()
    router.push('/login')
    return
  }
  
  isAdding.value = true
  const success = await cartStore.addToCart(props.product.id, quantity.value)
  isAdding.value = false
  
  if (success) {
    close()
  } else {
    alert('Gagal menambahkan ke keranjang.')
  }
}

const handleBuyNow = () => {
  console.log("Tombol Buy It Now ditekan untuk produk:", props.product?.name);

  if (!authStore.isAuthenticated) {
    localStorage.setItem('pendingProductAction', JSON.stringify({ productId: props.product.id }));
    close()
    router.push('/login')
    return
  }
  
  close()
  router.push({
    name: 'checkout',
    state: {
      product: JSON.stringify(props.product),
      qty: quantity.value
    }
  })
}

const handleChat = () => {
  if (!authStore.isAuthenticated) {
    localStorage.setItem('pendingProductAction', JSON.stringify({ productId: props.product.id }));
    close()
    router.push('/login')
    return
  }
  
  
  if (authStore.user?.id === props.product.seller_profile?.user_id) {
    alert('Ini adalah produk toko Anda sendiri.')
    return
  }
  
  if (props.product.seller_profile?.user_id) {
    const sellerName = props.product.seller_profile?.store_name || 'Penjual'
    chatStore.openChat(props.product.seller_profile.user_id, sellerName)
    close() 
  } else {
    alert('Informasi penjual tidak tersedia.')
  }
}
</script>
