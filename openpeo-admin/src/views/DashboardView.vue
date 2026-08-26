<template>
  <div class="admin-layout">
    
    <aside class="admin-sidebar">
      <div class="sidebar-header">
        <h2>OpenPeo Admin</h2>
      </div>
      <nav class="sidebar-nav">
        <button 
          :class="['nav-item', { active: currentTab === 'dashboard' }]" 
          @click="currentTab = 'dashboard'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
          Dashboard
        </button>
        <button 
          :class="['nav-item', { active: currentTab === 'users' }]" 
          @click="currentTab = 'users'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
          Manajemen Pengguna
        </button>
        <button 
          :class="['nav-item', { active: currentTab === 'products' }]" 
          @click="currentTab = 'products'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
          Moderasi Produk
        </button>
        <button 
          :class="['nav-item', { active: currentTab === 'transactions' }]" 
          @click="currentTab = 'transactions'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path></svg>
          Manajemen Transaksi
        </button>
      </nav>
      <div class="sidebar-footer">
        <div class="admin-info">
          <p class="admin-name">{{ authStore.user?.name?.replace('Admin ', '') }}</p>
          <p class="admin-email">{{ authStore.user?.email }}</p>
        </div>
        <button class="logout-btn" @click="handleLogout">
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
          Logout
        </button>
      </div>
    </aside>

    
    <main class="admin-main">
      <header class="topbar">
        <h1>{{ tabTitle }}</h1>
      </header>

      <div class="content-wrapper">
        
        <div v-if="currentTab === 'dashboard'" class="dashboard-stats">
          <div class="stat-card">
            <div class="stat-icon users">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
            </div>
            <div class="stat-info">
              <h3>Total Pembeli</h3>
              <p class="stat-value">{{ stats?.total_users || 0 }}</p>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon sellers">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
            </div>
            <div class="stat-info">
              <h3>Total Penjual</h3>
              <p class="stat-value">{{ stats?.total_sellers || 0 }}</p>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon products">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"></path><line x1="3" y1="6" x2="21" y2="6"></line><path d="M16 10a4 4 0 0 1-8 0"></path></svg>
            </div>
            <div class="stat-info">
              <h3>Total Produk</h3>
              <p class="stat-value">{{ stats?.total_products || 0 }}</p>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon transactions">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="1" x2="12" y2="23"></line><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path></svg>
            </div>
            <div class="stat-info">
              <h3>Total Transaksi</h3>
              <p class="stat-value">{{ formatPrice(stats?.total_revenue || 0) }}</p>
            </div>
          </div>
        </div>

        
        <div v-if="currentTab === 'dashboard'" class="dashboard-activity">
          <h2 class="section-title">Transaksi Terbaru</h2>
          <div class="data-table-container overflow-x-auto w-full">
            <table class="data-table">
              <thead>
                <tr>
                  <th style="width: 1%; white-space: nowrap;">Tanggal</th>
                  <th style="width: 40%;">Produk</th>
                  <th>Pembeli</th>
                  <th>Penjual</th>
                  <th style="width: 1%; white-space: nowrap;">Total</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="order in recentOrders" :key="order.id">
                  <td>{{ formatDate(order.created_at) }}</td>
                  <td>{{ order.order_items && order.order_items.length > 0 ? order.order_items[0].product?.name + (order.order_items.length > 1 ? ` (+${order.order_items.length - 1} item)` : '') : '-' }}</td>
                  <td>{{ order.customer?.name || '-' }}</td>
                  <td>{{ order.seller_profile?.store_name || 'Nama Toko Belum Diatur' }}</td>
                  <td>{{ formatPrice(order.total_price) }}</td>
                  <td>
                    <span :class="['status-badge', getOrderStatusClass(order.status)]">{{ order.status }}</span>
                  </td>
                </tr>
                <tr v-if="recentOrders.length === 0">
                  <td colspan="6" class="empty-state">Belum ada transaksi.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        
        <div v-if="currentTab === 'users'" class="data-table-container overflow-x-auto w-full">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Nama Lengkap</th>
                <th>Email</th>
                <th>Role</th>
                <th>No. Telepon</th>
                <th>Terdaftar</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in users" :key="user.id">
                <td>#{{ user.id }}</td>
                <td>
                  <strong>{{ user.name }}</strong>
                  <div v-if="user.role === 'seller'" class="seller-store">{{ user?.seller_profile?.store_name || user?.store_name || 'Nama Toko Belum Diatur' }}</div>
                </td>
                <td>{{ user.email }}</td>
                <td>
                  <span :class="['role-badge', user.role]">{{ user.role }}</span>
                </td>
                <td>{{ user.phone || '-' }}</td>
                <td>{{ formatDate(user.created_at) }}</td>
              </tr>
              <tr v-if="users.length === 0">
                <td colspan="6" class="empty-state">Tidak ada pengguna ditemukan.</td>
              </tr>
            </tbody>
          </table>
        </div>

        
        <div v-if="currentTab === 'products'" class="data-table-container overflow-x-auto w-full">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Produk</th>
                <th>Toko / Penjual</th>
                <th>Harga</th>
                <th>Status</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="product in products" :key="product.id">
                <td>#{{ product.id }}</td>
                <td>
                  <div class="product-cell">
                    <img :src="getImageUrl(product)" @error="onImageError" alt="" class="product-img" />
                    <div>
                      <strong>{{ product.name }}</strong>
                      <div class="product-category">{{ product.category }}</div>
                    </div>
                  </div>
                </td>
                <td>{{ product?.seller_profile?.store_name || 'Nama Toko Belum Diatur' }}</td>
                <td>{{ formatPrice(product.price) }}</td>
                <td>
                  <span :class="['status-badge', product.is_active ? 'active' : 'inactive']">
                    {{ product.is_active ? 'Aktif' : 'Draft' }}
                  </span>
                </td>
                <td>
                  <button class="btn-delete" @click="deleteProduct(product.id)" title="Hapus Permanen">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                  </button>
                </td>
              </tr>
              <tr v-if="products.length === 0">
                <td colspan="6" class="empty-state">Tidak ada produk ditemukan.</td>
              </tr>
            </tbody>
          </table>
        </div>

        
        <div v-if="currentTab === 'transactions'" class="data-table-container overflow-x-auto w-full">
          <table class="data-table">
            <thead>
              <tr>
                <th style="width: 1%; white-space: nowrap;">ID</th>
                <th style="width: 1%; white-space: nowrap;">Tanggal</th>
                <th>Pembeli</th>
                <th>Toko / Penjual</th>
                <th style="width: 35%;">Produk & Catatan</th>
                <th style="width: 1%; white-space: nowrap;">Total</th>
                <th>Status</th>
                <th style="width: 1%; white-space: nowrap;">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="order in orders" :key="order.id">
                <td>#{{ order.id }}</td>
                <td>{{ formatDate(order.created_at) }}</td>
                <td>
                  <strong>{{ order.customer?.name || '-' }}</strong>
                  <div class="seller-store">{{ order.customer?.email }}</div>
                </td>
                <td>{{ order.seller_profile?.store_name || 'Nama Toko Belum Diatur' }}</td>
                <td>
                  <strong>{{ order.order_items && order.order_items.length > 0 ? order.order_items.map(item => item.product?.name + ' (x' + item.quantity + ')').join(', ') : '-' }}</strong>
                  <div class="product-category" v-if="order.custom_notes">Catatan: {{ order.custom_notes }}</div>
                </td>
                <td>{{ formatPrice(order.total_price) }}</td>
                <td>
                  <span :class="['status-badge', getOrderStatusClass(order.status)]">{{ order.status }}</span>
                </td>
                <td>
                  <button 
                    v-if="order.status !== 'Dibatalkan Admin' && order.status !== 'Selesai'" 
                    class="btn-delete" 
                    @click="cancelOrder(order.id)" 
                    title="Batalkan Pesanan (Force Cancel)"
                  >
                    Batalkan
                  </button>
                  <span v-else class="action-disabled">-</span>
                </td>
              </tr>
              <tr v-if="orders.length === 0">
                <td colspan="8" class="empty-state">Tidak ada transaksi ditemukan.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'
import { BASE_URL } from '../axios'

const router = useRouter()
const authStore = useAuthStore()

const currentTab = ref('dashboard')
const stats = ref(null)
const users = ref([])
const products = ref([])
const logs = ref([])

const CATEGORY_FALLBACKS = {
  'Koleksi Tenun NTT': 'https://images.unsplash.com/photo-1596755094514-f87e34085b2c?q=80&w=800&auto=format&fit=crop',
  'Cita Rasa Lokal':   'https://images.unsplash.com/photo-1447933601403-0c6688de566e?q=80&w=800&auto=format&fit=crop',
  'Koleksi Aksesoris': 'https://images.unsplash.com/photo-1515562141207-7a88fb7ce338?q=80&w=800&auto=format&fit=crop',
}

const SVG_PLACEHOLDER = "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='400' height='400' viewBox='0 0 400 400'%3E%3Crect width='400' height='400' fill='%23f3f0eb'/%3E%3Ctext x='200' y='210' text-anchor='middle' font-family='sans-serif' font-size='14' fill='%23a8956e'%3EOpenPeo%3C/text%3E%3C/svg%3E"

const getImageUrl = (product) => {
  const raw = product?.image_url
  if (raw && (raw.startsWith('http://') || raw.startsWith('https://'))) {
    return raw
  }
  
  if (raw && (raw.startsWith('/images/') || raw.startsWith('images/'))) {
    return `${BASE_URL}${raw.startsWith('/') ? raw.slice(1) : raw}`
  }

  const categoryFallback = CATEGORY_FALLBACKS[product?.category]
  if (categoryFallback) return categoryFallback

  return SVG_PLACEHOLDER
}

const onImageError = (event) => {
  if (event.target.src !== SVG_PLACEHOLDER) {
    event.target.src = SVG_PLACEHOLDER
  }
}

const tabTitle = computed(() => {
  switch (currentTab.value) {
    case 'dashboard': return 'Platform Overview'
    case 'users': return 'Manajemen Pengguna'
    case 'products': return 'Moderasi Produk'
    case 'transactions': return 'Manajemen Transaksi'
    default: return ''
  }
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const fetchData = async () => {
  const token = authStore.token
  const headers = token ? { Authorization: `Bearer ${token}` } : {}
  try {
    const [statsRes, usersRes, productsRes, ordersRes, logsRes] = await Promise.all([
      axios.get(`${BASE_URL}api/admin/stats`, { headers }),
      axios.get(`${BASE_URL}api/admin/users`, { headers }),
      axios.get(`${BASE_URL}api/admin/products`, { headers }),
      axios.get(`${BASE_URL}api/admin/orders`, { headers }),
      axios.get(`${BASE_URL}api/admin/activity-logs`, { headers })
    ])
    
    if (statsRes.data.success) stats.value = statsRes.data.data
    if (usersRes.data.success) users.value = usersRes.data.data
    if (productsRes.data.success) products.value = productsRes.data.data
    if (ordersRes.data.success) orders.value = ordersRes.data.data
    if (logsRes.data.success) logs.value = logsRes.data.data
  } catch (error) {
    console.error('Failed to fetch admin data', error)
    if (error.response?.status === 401 || error.response?.status === 403) {
      handleLogout()
    }
  }
}

const deleteProduct = async (id) => {
  if (!confirm('Apakah Anda yakin ingin menghapus produk ini secara permanen?')) return
  
  try {
    const token = authStore.token
    const res = await axios.delete(`${BASE_URL}api/admin/products/${id}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    if (res.data.success) {
      products.value = products.value.filter(p => p.id !== id)
      
      if (stats.value) stats.value.total_products--
      alert('Produk berhasil dihapus')
    }
  } catch (error) {
    alert(error.response?.data?.message || 'Gagal menghapus produk')
  }
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(price)
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }).format(date)
}

const orders = ref([])

const recentOrders = computed(() => {
  return orders.value.slice(0, 5)
})

const getOrderStatusClass = (status) => {
  switch (status) {
    case 'Menunggu Pembayaran': return 'warning'
    case 'Diproses Perajin': return 'info'
    case 'Dikirim': return 'primary'
    case 'Selesai': return 'success'
    case 'Dibatalkan Admin': return 'danger'
    default: return 'inactive'
  }
}

const cancelOrder = async (id) => {
  if (!confirm('Apakah Anda yakin ingin membatalkan pesanan ini (Force Cancel)?')) return
  
  try {
    const token = authStore.token
    const res = await axios.put(`${BASE_URL}api/admin/orders/${id}/cancel`, {}, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    if (res.data.success) {
      
      const order = orders.value.find(o => o.id === id)
      if (order) order.status = 'Dibatalkan Admin'
      
      fetchData()
      alert('Pesanan berhasil dibatalkan')
    }
  } catch (error) {
    alert(error.response?.data?.message || 'Gagal membatalkan pesanan')
  }
}

let pollInterval = null

onMounted(() => {
  fetchData()
  pollInterval = setInterval(() => {
    fetchData()
  }, 3000)
})

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval)
  }
})
</script>

<style scoped>
.admin-layout {
  display: flex;
  width: 100vw;
  height: 100vh;
  background-color: #f3f4f6;
  font-family: 'Inter', sans-serif;
  overflow: hidden;
}


.admin-sidebar {
  width: 260px;
  background-color: #111827;
  color: white;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-header {
  padding: 24px;
  border-bottom: 1px solid #1f2937;
}

.sidebar-header h2 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.sidebar-nav {
  flex: 1;
  padding: 20px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: #9ca3af;
  background: none;
  border: none;
  border-radius: 6px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
}

.nav-item:hover {
  background-color: #1f2937;
  color: white;
}

.nav-item.active {
  background-color: #374151;
  color: white;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid #1f2937;
}

.admin-info {
  margin-bottom: 16px;
}

.admin-name {
  margin: 0;
  font-weight: 600;
  font-size: 0.9rem;
}

.admin-email {
  margin: 2px 0 0;
  font-size: 0.8rem;
  color: #9ca3af;
}

.logout-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 10px;
  background-color: transparent;
  border: 1px solid #4b5563;
  color: #d1d5db;
  border-radius: 6px;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s;
}

.logout-btn:hover {
  background-color: #374151;
  color: white;
}


.admin-main {
  flex: 1;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.topbar {
  background-color: white;
  padding: 20px 32px;
  border-bottom: 1px solid #e5e7eb;
}

.topbar h1 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #111827;
}

.content-wrapper {
  flex: 1;
  padding: 32px;
  overflow-y: auto;
}


.dashboard-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 24px;
}

.stat-card {
  background: white;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon.users { background-color: #dbeafe; color: #2563eb; }
.stat-icon.sellers { background-color: #fce7f3; color: #db2777; }
.stat-icon.products { background-color: #d1fae5; color: #059669; }

.stat-info h3 {
  margin: 0 0 4px;
  font-size: 0.85rem;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  color: #111827;
}


.overflow-x-auto {
  overflow-x: auto;
}

.w-full {
  width: 100%;
}

.data-table-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th, .data-table td {
  padding: 12px 12px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
  line-height: 1.4;
}

.data-table th {
  background-color: #f9fafb;
  font-size: 0.75rem;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.data-table td {
  font-size: 0.9rem;
  color: #374151;
  vertical-align: top;
}

.data-table tr:last-child td {
  border-bottom: none;
}


.seller-store {
  font-size: 0.8rem;
  color: #6b7280;
  margin-top: 2px;
}

.role-badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}
.role-badge.customer { background-color: #e0f2fe; color: #0369a1; }
.role-badge.seller { background-color: #fce7f3; color: #be185d; }

.product-cell {
  display: flex;
  align-items: center;
  gap: 16px;
}

.product-img {
  width: 48px;
  height: 48px;
  border-radius: 6px;
  object-fit: cover;
  background-color: #f3f4f6;
}

.product-category {
  font-size: 0.8rem;
  color: #6b7280;
  margin-top: 2px;
}

.status-badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
}
.status-badge.active { background-color: #d1fae5; color: #065f46; }
.status-badge.inactive { background-color: #f3f4f6; color: #4b5563; }
.status-badge.warning { background-color: #fef08a; color: #854d0e; }
.status-badge.info { background-color: #bfdbfe; color: #1e40af; }
.status-badge.primary { background-color: #e0e7ff; color: #3730a3; }
.status-badge.success { background-color: #dcfce3; color: #166534; }
.status-badge.danger { background-color: #fee2e2; color: #991b1b; }

.btn-delete {
  color: #ef4444;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.btn-delete:hover {
  background-color: #fee2e2;
}

.empty-state {
  text-align: center;
  padding: 40px !important;
  color: #6b7280;
}

.dashboard-activity {
  margin-top: 32px;
}
.section-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #111827;
  margin-bottom: 16px;
}
.stat-icon.transactions { background-color: #fef3c7; color: #d97706; }
.stat-value {
  font-size: 1.6rem;
}
.action-disabled {
  color: #9ca3af;
  font-size: 0.85rem;
}
</style>
