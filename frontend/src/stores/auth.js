import { defineStore } from 'pinia'
import apiClient from '../axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('openpeo_user')) || null,
    token: localStorage.getItem('openpeo_token') || null,
  }),
  
  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin',
    userInitials: (state) => {
      if (!state.user?.name) return '?'
      return state.user.name
        .split(' ')
        .map(w => w[0])
        .join('')
        .toUpperCase()
        .slice(0, 2)
    }
  },
  
  actions: {
    async login(email, password) {
      try {
        const response = await apiClient.post('/api/login', {
          email,
          password
        })
        
        if (response.data.success) {
          this.token = response.data.token
          this.user = response.data.data
          
          localStorage.setItem('openpeo_token', this.token)
          localStorage.setItem('openpeo_user', JSON.stringify(this.user))
          
          
          return { success: true }
        }
      } catch (error) {
        return { 
          success: false, 
          message: error.response?.data?.message || 'Login failed' 
        }
      }
    },

    async fetchProfile() {
      if (!this.isAuthenticated) return
      try {
        const response = await apiClient.get('/api/user/profile')
        if (response.data.success) {
          this.user = response.data.data
          localStorage.setItem('openpeo_user', JSON.stringify(this.user))
        }
      } catch (error) {
        console.error('Failed to fetch profile', error)
      }
    },

    updateSellerProfileLocally(profile) {
      if (this.user && this.user.role === 'seller') {
        if (!this.user.seller_profile) this.user.seller_profile = {}
        this.user.seller_profile.store_name = profile.store_name || this.user.seller_profile?.store_name
        this.user.seller_profile.store_logo = profile.store_logo || this.user.seller_profile?.store_logo
        
        this.user.store_name = this.user.seller_profile.store_name
        this.user.store_logo = this.user.seller_profile.store_logo
        localStorage.setItem('openpeo_user', JSON.stringify(this.user))
      }
    },

    async updateProfile(data) {
      try {
        const response = await apiClient.put('/api/user/profile', data)
        if (response.data.success) {
          this.user = response.data.data
          localStorage.setItem('openpeo_user', JSON.stringify(this.user))
          return { success: true, message: response.data.message }
        }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.message || 'Gagal memperbarui profil'
        }
      }
    },

    async changePassword(currentPassword, newPassword) {
      try {
        const response = await apiClient.put('/api/user/password', {
          current_password: currentPassword,
          new_password: newPassword,
        })
        if (response.data.success) {
          return { success: true, message: response.data.message }
        }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.message || 'Gagal mengubah password'
        }
      }
    },

    async registerSeller(data) {
      try {
        const response = await apiClient.post('/api/register-seller', data)
        if (response.data.success) {
          return { success: true, message: response.data.message }
        }
      } catch (error) {
        return {
          success: false,
          message: error.response?.data?.message || 'Gagal mendaftar toko'
        }
      }
    },
    
    logout() {
      import('./notification').then(module => {
        const notificationStore = module.useNotificationStore()
        notificationStore.resetAll()
      }).catch(err => console.error(err))

      this.token = null
      this.user = null
      localStorage.removeItem('openpeo_token')
      localStorage.removeItem('openpeo_user')
      // Token dihapus dari localStorage, interceptor apiClient tidak akan menyisipkannya lagi
    }
  }
})
