import { defineStore } from 'pinia'
import axios from 'axios'
import { useAuthStore } from './auth'

export const useSellerStore = defineStore('seller', {
  state: () => ({
    products: [],
    isLoading: false,
  }),
  
  getters: {
    productCount: (state) => state.products.length
  },
  
  actions: {
    async fetchProducts() {
      const auth = useAuthStore()
      if (!auth.isAuthenticated || auth.user?.role !== 'seller') return
      
      this.isLoading = true
      try {
        const response = await axios.get('http://localhost:8081/api/seller/products')
        if (response.data.success) {
          this.products = response.data.data || []
        }
      } catch (error) {
        console.error('Failed to fetch seller products', error)
        this.products = []
      } finally {
        this.isLoading = false
      }
    },
    
    async createProduct(productData) {
      try {
        const response = await axios.post('http://localhost:8081/api/seller/products', productData)
        if (response.data.success) {
          await this.fetchProducts()
          return { success: true, message: response.data.message }
        }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.message || 'Gagal menambah produk'
        }
      }
    },
    
    async updateProduct(id, productData) {
      try {
        const response = await axios.put(`http://localhost:8081/api/seller/products/${id}`, productData)
        if (response.data.success) {
          await this.fetchProducts()
          return { success: true, message: response.data.message }
        }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.message || 'Gagal mengubah produk'
        }
      }
    },
    
    async deleteProduct(id) {
      try {
        const response = await axios.delete(`http://localhost:8081/api/seller/products/${id}`)
        if (response.data.success) {
          await this.fetchProducts()
          return { success: true, message: response.data.message }
        }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.message || 'Gagal menghapus produk'
        }
      }
    }
  }
})
