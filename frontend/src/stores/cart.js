import { defineStore } from 'pinia'
import axios from 'axios'
import { useAuthStore } from './auth'

export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [],
    isLoading: false,
    isOpen: false
  }),
  
  getters: {
    totalItems: (state) => state.items.reduce((acc, item) => acc + item.quantity, 0),
    totalPrice: (state) => state.items.reduce((acc, item) => {
      const price = item.product?.price || 0
      return acc + (price * item.quantity)
    }, 0)
  },
  
  actions: {
    toggleDrawer() {
      this.isOpen = !this.isOpen
    },

    openDrawer() {
      this.isOpen = true
    },

    closeDrawer() {
      this.isOpen = false
    },

    async fetchCart() {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) return
      
      this.isLoading = true
      try {
        const response = await axios.get('http://localhost:8081/api/cart')
        if (response.data.success) {
          this.items = response.data.data || []
        }
      } catch (error) {
        console.error('Failed to fetch cart', error)
        this.items = []
      } finally {
        this.isLoading = false
      }
    },
    
    async addToCart(productId, quantity = 1) {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) return false
      
      try {
        const response = await axios.post('http://localhost:8081/api/cart', {
          product_id: productId,
          quantity: quantity
        })
        
        if (response.data.success) {
          await this.fetchCart() 
          return true
        }
      } catch (error) {
        console.error('Failed to add to cart', error)
        return false
      }
    },

    async removeFromCart(cartItemId) {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) return false

      try {
        const response = await axios.delete(`http://localhost:8081/api/cart/${cartItemId}`)
        if (response.data.success) {
          
          this.items = this.items.filter(item => item.id !== cartItemId)
          return true
        }
      } catch (error) {
        console.error('Failed to remove from cart', error)
        return false
      }
    },
    
    async checkout() {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) return { success: false, message: 'Unauthenticated' }
      
      this.isLoading = true
      try {
        const response = await axios.post('http://localhost:8081/api/orders', {})
        if (response.data.success) {
          this.clearLocalCart()
          return { success: true, message: 'Pesanan berhasil dibuat' }
        }
      } catch (error) {
        return { 
          success: false, 
          message: error.response?.data?.message || 'Gagal membuat pesanan' 
        }
      } finally {
        this.isLoading = false
      }
    },

    clearLocalCart() {
      this.items = []
      this.isOpen = false
    }
  }
})
