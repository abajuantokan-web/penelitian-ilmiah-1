import { defineStore } from 'pinia'
import axios from 'axios'

export const useProductStore = defineStore('products', {
  state: () => ({
    allProducts: [],
    isLoading: false
  }),
  
  getters: {
    tenunProducts: (state) => state.allProducts.filter(p => p.category === 'Koleksi Tenun NTT'),
    foodProducts: (state) => state.allProducts.filter(p => p.category === 'Cita Rasa Lokal'),
    accessoriesProducts: (state) => state.allProducts.filter(p => p.category === 'Koleksi Aksesoris')
  },
  
  actions: {
    async fetchProducts() {
      this.isLoading = true
      try {
        const response = await axios.get('http://localhost:8081/api/products?limit=100')
        if (response.data.success) {
          this.allProducts = response.data.data
        }
      } catch (error) {
        console.error('Failed to fetch products', error)
      } finally {
        this.isLoading = false
      }
    }
  }
})
