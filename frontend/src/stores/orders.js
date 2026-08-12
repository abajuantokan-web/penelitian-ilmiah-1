import { defineStore } from 'pinia'
import apiClient from '../axios'
import { useAuthStore } from './auth'

export const useOrderStore = defineStore('orders', {
  state: () => ({
    orders: [],
    isLoading: false
  }),

  getters: {
    orderCount: (state) => state.orders.length
  },

  actions: {
    async fetchOrders() {
      const auth = useAuthStore()
      if (!auth.isAuthenticated) return

      this.isLoading = true
      try {
        const response = await apiClient.get('/api/user/orders')
        if (response.data.success) {
          this.orders = response.data.data || []
        }
      } catch (error) {
        console.error('Failed to fetch orders', error)
        this.orders = []
      } finally {
        this.isLoading = false
      }
    },

    clearOrders() {
      this.orders = []
    }
  }
})
