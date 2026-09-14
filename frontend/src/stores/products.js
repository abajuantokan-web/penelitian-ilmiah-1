import { defineStore } from 'pinia'
import apiClient from '../axios'

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
    /**
     * Fetch products from the API.
     *
     * @param {number|null} customLimit - Optional limit. When null/undefined the
     *   API's own default page size is used — no aggressive ?limit=100 is sent.
     *   Pass a small number (4–8) for initial "above-the-fold" loads, and a
     *   larger number only when the user explicitly requests more.
     */
    async fetchProducts(customLimit = null) {
      this.isLoading = true
      try {
        // Build params conditionally so the URL stays clean when no limit is needed
        const params = {}
        if (customLimit !== null && customLimit !== undefined) {
          params.limit = customLimit
        }

        const response = await apiClient.get('/api/products', { params })
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
