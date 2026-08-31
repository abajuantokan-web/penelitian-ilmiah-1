<template>
  <div class="profile-page">
    <div class="container profile-container">
      
      
      <aside class="profile-sidebar">
        
        <router-link to="/" class="back-home-btn">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="19" y1="12" x2="5" y2="12"></line>
            <polyline points="12 19 5 12 12 5"></polyline>
          </svg>
          Kembali ke Beranda
        </router-link>
        <div class="profile-user-info">
          <div class="profile-avatar">{{ authStore.userInitials }}</div>
          <h2 class="profile-name">Halo, {{ authStore.user?.name || 'Pengguna' }}</h2>
          <p class="profile-email">{{ authStore.user?.email }}</p>
        </div>

        <nav class="profile-nav" aria-label="Menu Profil">
          <button 
            :class="['profile-nav-btn', { active: currentTab === 'account' }]" 
            @click="currentTab = 'account'"
          >
            Informasi Akun
          </button>
          <button 
            :class="['profile-nav-btn', { active: currentTab === 'orders' }]" 
            @click="currentTab = 'orders'"
          >
            Riwayat Pesanan
          </button>
          <button 
            :class="['profile-nav-btn', { active: currentTab === 'settings' }]" 
            @click="currentTab = 'settings'"
          >
            Pengaturan
          </button>
        </nav>

        <div class="profile-logout-wrapper">
          
          <router-link 
            v-if="authStore.user?.role === 'seller'" 
            to="/seller/dashboard" 
            class="store-link-btn"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
              <polyline points="9 22 9 12 15 12 15 22"></polyline>
            </svg>
            Toko Saya
          </router-link>
          
          <router-link 
            v-else 
            to="/register-seller" 
            class="store-link-btn"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="12" y1="8" x2="12" y2="16"></line>
              <line x1="8" y1="12" x2="16" y2="12"></line>
            </svg>
            Buat Toko
          </router-link>
          <button class="profile-logout-btn" @click="handleLogout">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
            Keluar Akun
          </button>
        </div>
      </aside>

      
      <main class="profile-content">
        
        
        <div v-if="currentTab === 'account'" class="profile-tab fade-in">
          <h1 class="profile-title">Informasi Akun</h1>
          <p class="profile-subtitle">Perbarui detail profil dan informasi kontak Anda.</p>
          
          <form @submit.prevent="handleUpdateProfile" class="profile-form">
            <div class="form-group">
              <label for="name">Nama Lengkap</label>
              <input type="text" id="name" v-model="profileForm.name" required />
            </div>
            
            <div class="form-group">
              <label for="phone">Nomor Telepon</label>
              <input type="tel" id="phone" v-model="profileForm.phone" placeholder="Contoh: 081234567890" />
            </div>
            
            <div class="form-group">
              <label for="address">Alamat Pengiriman Utama</label>
              <textarea id="address" v-model="profileForm.address" rows="4" placeholder="Alamat lengkap beserta kode pos"></textarea>
            </div>
            
            <div class="form-actions">
              <button type="submit" class="btn-primary" :disabled="isUpdatingProfile">
                {{ isUpdatingProfile ? 'Menyimpan...' : 'Simpan Perubahan' }}
              </button>
              <span v-if="profileMsg" :class="['msg', profileStatus]">{{ profileMsg }}</span>
            </div>
          </form>
        </div>

        
        <div v-if="currentTab === 'orders'" class="profile-tab fade-in">
          <h1 class="profile-title">Riwayat Pesanan</h1>
          <p class="profile-subtitle">Lacak status pesanan pre-order dan pembelian Anda.</p>
          
          <div v-if="orderStore.isLoading" class="loading-state">
            <div class="spinner"></div>
            <p>Memuat pesanan...</p>
          </div>
          
          <div v-else-if="!orderStore.orders || orderStore.orders.length === 0" class="empty-state">
            <p>Anda belum memiliki riwayat pesanan.</p>
            <router-link to="/#products" class="btn-outline">Mulai Belanja</router-link>
          </div>
          
          <div v-else class="order-list">
            <div v-for="order in orderStore.orders" :key="order.id" class="order-card">
              
              
              <div class="order-header">
                <div class="order-id">
                  <span>Order #{{ order.id }}</span>
                  <span class="order-date">{{ formatDate(order.created_at) }}</span>
                </div>
                <div :class="['order-badge', getStatusClass(order.status)]">
                  {{ order.status }}
                </div>
              </div>
              
              <div v-if="order.order_items && order.order_items.length > 0" class="order-items-list">
                 <div v-for="item in order.order_items" :key="item.id" class="order-product border-b-subtle">
                   <div class="product-thumb">
                     <img :src="$getImageUrl(item.product?.image_url)" :alt="item.product?.name" loading="lazy" />
                   </div>
                   <div class="product-details">
                     <h3 class="product-name">{{ item.product?.name || 'Produk' }}</h3>
                     <p class="product-meta">Qty: {{ item.quantity }} &times; {{ formatPrice(item.price) }}</p>
                   </div>
                 </div>
                 <div v-if="order.custom_notes || order.note" class="custom-note mt-3">
                    <strong>Catatan:</strong> {{ order.custom_notes || order.note }}
                 </div>
              </div>
              
              
              <div>
                <div class="order-footer">
                  <div class="order-total-label">Total Pembayaran</div>
                  <div class="order-total-val">{{ formatPrice(order.total_price) }}</div>
                </div>
                <div class="flex justify-end mb-4 px-6" v-if="['diproses', 'dikirim', 'sedang_diproses', 'pesanan sedang diproses', 'diproses perajin'].includes(order.status?.toLowerCase())">
                  <button @click="promptCompleteOrder(order.id)" class="bg-zinc-900 text-white px-5 py-2 rounded-md text-sm font-medium hover:bg-zinc-800 transition shadow-sm" style="background-color: #18181b; color: white;">
                    Pesanan Selesai
                  </button>
                </div>
              </div>
              
            </div>
          </div>
        </div>

        
        <div v-if="currentTab === 'settings'" class="profile-tab fade-in">
          <h1 class="profile-title">Pengaturan Akun</h1>
          <p class="profile-subtitle">Ubah password dan kelola keamanan akun Anda.</p>
          
          <form @submit.prevent="handleChangePassword" class="profile-form">
            <div class="form-group">
              <label for="current_password">Password Saat Ini</label>
              <input type="password" id="current_password" v-model="passwordForm.current" required />
            </div>
            
            <div class="form-group">
              <label for="new_password">Password Baru</label>
              <input type="password" id="new_password" v-model="passwordForm.new" required minlength="6" />
            </div>
            
            <div class="form-group">
              <label for="confirm_password">Konfirmasi Password Baru</label>
              <input type="password" id="confirm_password" v-model="passwordForm.confirm" required minlength="6" />
            </div>
            
            <div class="form-actions">
              <button type="submit" class="btn-primary" :disabled="isChangingPassword">
                {{ isChangingPassword ? 'Memproses...' : 'Ubah Password' }}
              </button>
              <span v-if="passMsg" :class="['msg', passStatus]">{{ passMsg }}</span>
            </div>
          </form>
        </div>

      </main>
    </div>
  </div>

  <div v-if="showCompleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm transition-opacity">
    <div class="bg-white rounded-2xl shadow-2xl p-6 w-full max-w-md mx-4 transform transition-all scale-100">
      <div class="text-center">
        <div class="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-green-100 mb-4">
          <svg class="h-8 w-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h3 class="text-xl font-bold text-gray-900 mb-2">Konfirmasi Pesanan Diterima</h3>
        <p class="text-sm text-gray-500 mb-6">
          Apakah Anda yakin pesanan ini sudah diterima dengan baik? Dana akan langsung diteruskan ke saldo Penjual dan tindakan ini tidak dapat dibatalkan.
        </p>
      </div>
      <div class="flex gap-3 justify-center">
        <button @click="cancelCompleteOrder" class="px-5 py-2.5 bg-gray-100 text-gray-700 font-medium rounded-xl hover:bg-gray-200 transition-colors w-full">
          Batal
        </button>
        <button @click="confirmCompleteOrder" class="px-5 py-2.5 bg-black text-white font-medium rounded-xl hover:bg-gray-800 transition-colors w-full">
          Ya, Selesai
        </button>
      </div>
    </div>
  </div>

  <transition enter-active-class="transition ease-out duration-300 transform" enter-from-class="translate-y-10 opacity-0" enter-to-class="translate-y-0 opacity-100" leave-active-class="transition ease-in duration-200" leave-from-class="opacity-100" leave-to-class="opacity-0 translate-y-10">
    <div v-if="showToast" class="fixed bottom-6 right-6 z-[60] flex items-center w-full max-w-xs p-4 space-x-3 text-gray-900 bg-white rounded-2xl shadow-[0_8px_30px_rgb(0,0,0,0.12)] border border-gray-100">
      <div class="inline-flex items-center justify-center flex-shrink-0 w-10 h-10 bg-green-100 rounded-xl">
        <svg class="w-6 h-6 text-green-600" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
        </svg>
      </div>
      <div class="text-sm font-semibold">{{ toastMessage }}</div>
    </div>
  </transition>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useOrderStore } from '../stores/orders'
import { useCartStore } from '../stores/cart'
import axios from '../axios'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const orderStore = useOrderStore()
const cartStore = useCartStore()

const currentTab = ref('account')
const showCompleteModal = ref(false)
const orderToComplete = ref(null)

const showToast = ref(false)
const toastMessage = ref('')
const triggerToast = (msg) => {
  toastMessage.value = msg;
  showToast.value = true;
  setTimeout(() => { showToast.value = false; }, 3000);
};

watch(
  () => route.query.tab,
  (newTab) => {
    if (!newTab) return; 
    if (newTab === 'orders') {
      currentTab.value = 'orders'
    } else if (newTab === 'account') {
      currentTab.value = 'account'
    } else if (newTab === 'settings') {
      currentTab.value = 'settings'
    }
  },
  { immediate: true, deep: true }
)


const profileForm = ref({
  name: '',
  phone: '',
  address: ''
})
const isUpdatingProfile = ref(false)
const profileMsg = ref('')
const profileStatus = ref('')


const passwordForm = ref({
  current: '',
  new: '',
  confirm: ''
})
const isChangingPassword = ref(false)
const passMsg = ref('')
const passStatus = ref('')


const formatPrice = (price) => {
  if (!price) return 'Rp0'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const d = new Date(dateString)
  return new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }).format(d)
}

const getStatusClass = (status) => {
  switch(status?.toLowerCase()) {
    case 'pending':
    case 'menunggu pembayaran':
      return 'badge-warning'
    case 'diproses perajin':
    case 'diproses':
      return 'badge-info'
    case 'dikirim':
      return 'badge-primary'
    case 'selesai':
      return 'badge-success'
    default:
      return 'badge-default'
  }
}


const handleLogout = () => {
  authStore.logout()
  cartStore.clearLocalCart()
  orderStore.clearOrders()
  router.push('/')
}

const handleUpdateProfile = async () => {
  isUpdatingProfile.value = true
  profileMsg.value = ''
  
  const result = await authStore.updateProfile({
    name: profileForm.value.name,
    phone: profileForm.value.phone,
    address: profileForm.value.address
  })
  
  isUpdatingProfile.value = false
  profileMsg.value = result.message
  profileStatus.value = result.success ? 'msg-success' : 'msg-error'
  
  if (result.success) {
    setTimeout(() => { profileMsg.value = '' }, 3000)
  }
}

const handleChangePassword = async () => {
  if (passwordForm.value.new !== passwordForm.value.confirm) {
    passStatus.value = 'msg-error'
    passMsg.value = 'Konfirmasi password baru tidak cocok.'
    return
  }

  isChangingPassword.value = true
  passMsg.value = ''
  
  const result = await authStore.changePassword(passwordForm.value.current, passwordForm.value.new)
  
  isChangingPassword.value = false
  passMsg.value = result.message
  passStatus.value = result.success ? 'msg-success' : 'msg-error'
  
  if (result.success) {
    passwordForm.value = { current: '', new: '', confirm: '' }
    setTimeout(() => { passMsg.value = '' }, 3000)
  }
}

const promptCompleteOrder = (orderId) => {
  orderToComplete.value = orderId;
  showCompleteModal.value = true;
};

const cancelCompleteOrder = () => {
  showCompleteModal.value = false;
  orderToComplete.value = null;
};

const confirmCompleteOrder = async () => {
  if (!orderToComplete.value) return;
  
  try {
    const response = await axios.put(`/api/orders/${orderToComplete.value}/complete`, {}, {
      headers: { Authorization: `Bearer ${authStore.token || localStorage.getItem('token')}` }
    })
    if (response.data.success) {
      showCompleteModal.value = false;
      orderToComplete.value = null;
      triggerToast('Terima kasih! Pesanan telah selesai.');
      orderStore.fetchOrders() 
    }
  } catch (error) {
    console.error("Gagal menyelesaikan pesanan:", error)
    alert(error.response?.data?.message || 'Gagal menyelesaikan pesanan.')
    showCompleteModal.value = false;
    orderToComplete.value = null;
  }
};

onMounted(async () => {
  
  await authStore.fetchProfile()
  if (authStore.user) {
    profileForm.value.name = authStore.user.name || ''
    profileForm.value.phone = authStore.user.phone || ''
    profileForm.value.address = authStore.user.address || ''
  }
  
  
  orderStore.fetchOrders()
})
</script>

<style scoped>
.profile-page {
  padding: 120px 0 80px;
  min-height: 80vh;
  background-color: #fafafa;
}

.profile-container {
  display: flex;
  gap: 40px;
  align-items: flex-start;
}


.profile-sidebar {
  width: 280px;
  flex-shrink: 0;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.03);
  overflow: hidden;
}

.back-home-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 24px;
  color: #666;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.85rem;
  font-weight: 600;
  text-decoration: none;
  border-bottom: 1px solid #eee;
  transition: color 0.2s, background-color 0.2s;
}

.back-home-btn:hover {
  color: #1a1a1a;
  background-color: #f9f9f9;
}

.profile-user-info {
  padding: 32px 24px;
  text-align: center;
  border-bottom: 1px solid #eee;
}

.profile-avatar {
  width: 72px;
  height: 72px;
  background-color: #1a1a1a;
  color: #fff;
  font-family: 'Playfair Display', serif;
  font-size: 1.5rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  margin: 0 auto 16px;
}

.profile-name {
  font-size: 1.1rem;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 4px;
}

.profile-email {
  font-size: 0.85rem;
  color: #666;
  margin: 0;
}

.profile-nav {
  display: flex;
  flex-direction: column;
  padding: 16px 0;
}

.profile-nav-btn {
  background: none;
  border: none;
  text-align: left;
  padding: 14px 24px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.9rem;
  font-weight: 500;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
  border-left: 3px solid transparent;
}

.profile-nav-btn:hover {
  background-color: #f9f9f9;
  color: #1a1a1a;
}

.profile-nav-btn.active {
  background-color: #f5f5f5;
  color: #1a1a1a;
  border-left-color: #1a1a1a;
  font-weight: 600;
}

.store-link-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px;
  margin-bottom: 12px;
  background-color: #fff;
  border: 1px solid #1a1a1a;
  color: #1a1a1a;
  border-radius: 4px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.9rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s;
}

.store-link-btn:hover {
  background-color: #1a1a1a;
  color: #fff;
}

.profile-logout-wrapper {
  padding: 16px 24px 24px;
  border-top: 1px solid #eee;
}

.profile-logout-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: none;
  border: 1px solid #eee;
  padding: 10px;
  border-radius: 4px;
  color: #dc2626;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.profile-logout-btn:hover {
  background-color: #fef2f2;
  border-color: #fca5a5;
}


.profile-content {
  flex: 1;
  min-width: 0;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.03);
  padding: 40px;
}

.profile-title {
  font-family: 'Playfair Display', serif;
  font-size: 1.8rem;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px;
}

.profile-subtitle {
  font-size: 0.95rem;
  color: #666;
  margin: 0 0 32px;
}


.profile-form {
  max-width: 500px;
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  font-size: 0.85rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.9rem;
  transition: border-color 0.2s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #1a1a1a;
}

.form-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 32px;
}

.btn-primary {
  padding: 12px 28px;
  background-color: #1a1a1a;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background-color: #333;
}

.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-outline {
  display: inline-block;
  padding: 12px 28px;
  background-color: transparent;
  color: #1a1a1a;
  border: 1px solid #1a1a1a;
  border-radius: 4px;
  font-family: 'Montserrat', sans-serif;
  font-size: 0.9rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s;
}

.btn-outline:hover {
  background-color: #1a1a1a;
  color: #fff;
}

.msg {
  font-size: 0.85rem;
  font-weight: 500;
}
.msg-success { color: #16a34a; }
.msg-error { color: #dc2626; }


.order-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.order-card {
  border: 1px solid #eee;
  border-radius: 8px;
  overflow: hidden;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background-color: #fafafa;
  border-bottom: 1px solid #eee;
}

.order-id {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.order-id span:first-child {
  font-weight: 600;
  color: #1a1a1a;
  font-size: 0.9rem;
}

.order-date {
  font-size: 0.8rem;
  color: #666;
}

.order-badge {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.badge-warning { background-color: #fffbeb; color: #b45309; }
.badge-info { background-color: #eff6ff; color: #1d4ed8; }
.badge-primary { background-color: #f3e8ff; color: #7e22ce; }
.badge-success { background-color: #f0fdf4; color: #15803d; }
.badge-default { background-color: #f3f4f6; color: #374151; }

.order-product {
  display: flex;
  padding: 24px;
  gap: 20px;
}

.border-b-subtle {
  border-bottom: 1px solid #f9f9f9;
}
.order-product:last-child {
  border-bottom: none;
}

.product-thumb {
  width: 80px;
  height: 80px;
  border-radius: 4px;
  overflow: hidden;
  background-color: #f5f5f5;
  flex-shrink: 0;
}

.product-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-details {
  flex: 1;
}

.product-name {
  font-size: 1rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 6px;
}

.product-meta {
  font-size: 0.85rem;
  color: #666;
  margin: 0 0 12px;
}

.custom-note {
  background-color: #f9f9f9;
  padding: 10px 14px;
  border-radius: 4px;
  font-size: 0.85rem;
  color: #444;
  border-left: 3px solid #ddd;
}
.mt-3 {
  margin-top: 12px;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background-color: #fff;
  border-top: 1px solid #eee;
}

.order-total-label {
  font-size: 0.85rem;
  color: #666;
  font-weight: 500;
}

.order-total-val {
  font-family: 'Playfair Display', serif;
  font-size: 1.2rem;
  font-weight: 700;
  color: #1a1a1a;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #666;
}
.empty-state p {
  margin-bottom: 24px;
}

.spinner {
  width: 30px;
  height: 30px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #1a1a1a;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}
@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
.loading-state {
  text-align: center;
  padding: 40px;
  color: #666;
}
.loading-state p { margin-top: 12px; font-size: 0.9rem; }


.fade-in {
  animation: fadeIn 0.3s ease-in-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 768px) {
  .profile-container {
    flex-direction: column;
  }
  .profile-sidebar {
    width: 100%;
  }
  .profile-content {
    padding: 24px;
  }
  .order-product {
    flex-direction: column;
  }
  .product-thumb {
    width: 100%;
    height: 180px;
  }
}
</style>
