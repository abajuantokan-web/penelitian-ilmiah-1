<template>
  <article class="product-card" @click="$emit('click', product)">
    <div class="product-card-image">
      <img :src="$getImageUrl(product.image_url)" alt="Product Image" class="w-full h-full object-cover">
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

const storeName = computed(() => {
  if (props.product.seller?.seller_profile?.store_name) {
    return props.product.seller.seller_profile.store_name
  }
  if (props.product.seller?.store_name) {
    return props.product.seller.store_name
  }
  return 'Toko Anonim'
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}
</script>
