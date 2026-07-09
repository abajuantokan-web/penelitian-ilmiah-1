<template>
  <section :class="['section-products', { 'section-products--alt': altBg }]" :id="id" :aria-label="title">
    <div class="container">
      <div class="section-header reveal">
        <h2 class="section-title">{{ title }}</h2>
        <router-link to="/koleksi" class="section-link">Lihat Semua</router-link>
      </div>
      
      <div class="product-grid">
        <ProductCard 
          v-for="product in products" 
          :key="product.id" 
          :product="product"
          @click="$emit('product-click', product)"
        />
      </div>
    </div>
  </section>
</template>

<script setup>
import ProductCard from './ProductCard.vue'

defineProps({
  id: String,
  title: String,
  products: Array,
  altBg: {
    type: Boolean,
    default: false
  }
})

defineEmits(['product-click'])
</script>
