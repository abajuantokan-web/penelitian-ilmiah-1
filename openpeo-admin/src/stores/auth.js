import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('admin_user')) || null,
    token: localStorage.getItem('admin_token') || null,
  }),
  getters: {
    isAuthenticated: (state) => !!state.token && state.user?.role === 'admin'
  },
  actions: {
    async login(email, password) {
      try {
        const response = await axios.post('http://localhost:8081/api/login', {
          email,
          password
        })
        
        if (response.data.success) {
          if (response.data.data.role !== 'admin') {
             return { success: false, message: 'Akses ditolak: Hanya untuk admin' }
          }
          
          this.token = response.data.token
          this.user = response.data.data
          
          localStorage.setItem('admin_token', this.token)
          localStorage.setItem('admin_user', JSON.stringify(this.user))
          
          axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
          
          return { success: true }
        }
      } catch (error) {
        return { 
          success: false, 
          message: error.response?.data?.message || 'Login failed' 
        }
      }
    },
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
      delete axios.defaults.headers.common['Authorization']
    }
  }
})
