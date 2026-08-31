<template>
  <div class="admin-dashboard">
    
    <nav class="navbar glass-strong">
      <div class="navbar-container">
        <div class="navbar-brand">
          <div class="brand-icon">
            <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="20" cy="20" r="18" stroke="url(#brandGrad)" stroke-width="2.5" fill="none"/>
              <path d="M12 20C12 15.58 15.58 12 20 12C24.42 12 28 15.58 28 20C28 24.42 24.42 28 20 28" stroke="url(#brandGrad)" stroke-width="2.5" stroke-linecap="round"/>
              <circle cx="20" cy="20" r="4" fill="url(#brandGrad)"/>
              <defs>
                <linearGradient id="brandGrad" x1="0" y1="0" x2="40" y2="40">
                  <stop offset="0%" stop-color="#d35400"/>
                  <stop offset="50%" stop-color="#f5a623"/>
                  <stop offset="100%" stop-color="#ffd166"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="brand-text">
            <span class="brand-name text-gradient">OpenPeo</span>
            <span class="brand-tagline">Panel Admin</span>
          </div>
        </div>

        <div class="navbar-links">
          <router-link to="/" class="nav-link">Beranda</router-link>
          <span class="nav-link active">Dashboard Admin</span>
        </div>

        <div class="navbar-actions">
          <div class="user-profile-nav">
            <span class="user-greeting">Halo, <strong>{{ currentUser?.name || currentUser?.username }}</strong></span>
            <span class="badge badge-amber">Admin</span>
          </div>
          <button class="btn btn-secondary btn-sm" @click="logout">Keluar</button>
        </div>
      </div>
    </nav>

    
    <div class="dashboard-container">
      
      <div v-if="loading && salesDataLoading" class="loading-state">
        <div class="spinner"></div>
        <p>Memuat data panel admin...</p>
      </div>

      <div v-else-if="error" class="error-state glass">
        <div class="error-icon">⚠️</div>
        <h3>Gagal Memuat Data</h3>
        <p>{{ error }}</p>
        <button class="btn btn-primary btn-sm" @click="initializeDashboard">Coba Lagi</button>
      </div>

      <div v-else class="dashboard-content">
        
        <header class="dashboard-header">
          <div>
            <h1 class="dashboard-title">Panel <span class="text-gradient">Kendali Admin</span></h1>
            <p class="dashboard-subtitle">Kelola produk marketplace dan pantau ringkasan log transaksi secara real-time.</p>
          </div>
          <div class="header-actions">
            
            <div class="tab-navigation glass">
              <button 
                class="tab-btn" 
                :class="{ active: activeTab === 'summary' }" 
                @click="activeTab = 'summary'"
              >
                📊 Ringkasan
              </button>
              <button 
                class="tab-btn" 
                :class="{ active: activeTab === 'products' }" 
                @click="activeTab = 'products'"
              >
                🛍️ Kelola Produk
              </button>
            </div>
            <button class="btn btn-secondary btn-refresh" @click="initializeDashboard" :disabled="loading">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ 'spin-active': loading }">
                <path d="M21.5 2v6h-6M21.34 15.57a10 10 0 1 1-.57-8.38l5.67-5.67"/>
              </svg>
              Segarkan
            </button>
          </div>
        </header>

        
        <div v-if="activeTab === 'summary'" class="tab-pane animate-fade-in">
          
          <div class="metrics-grid">
            
            <div class="metric-card glass card-glow-amber">
              <div class="metric-header">
                <span class="metric-icon">💰</span>
                <span class="metric-badge positive">Keuangan</span>
              </div>
              <div class="metric-body">
                <span class="metric-label">Total Pendapatan</span>
                <h2 class="metric-value text-gradient">{{ formatCurrency(salesData.total_revenue) }}</h2>
              </div>
              <div class="metric-footer">
                <span class="metric-subtext">Akumulasi nilai transaksi</span>
              </div>
            </div>

            
            <div class="metric-card glass card-glow-indigo">
              <div class="metric-header">
                <span class="metric-icon">📦</span>
                <span class="metric-badge info">Pesanan</span>
              </div>
              <div class="metric-body">
                <span class="metric-label">Jumlah Pre-Order</span>
                <h2 class="metric-value text-gradient">{{ salesData.total_pre_orders }}</h2>
              </div>
              <div class="metric-footer">
                <span class="metric-subtext">Total pesanan masuk</span>
              </div>
            </div>

            
            <div class="metric-card glass card-glow-amber">
              <div class="metric-header">
                <span class="metric-icon">🛍️</span>
                <span class="metric-badge positive">Katalog</span>
              </div>
              <div class="metric-body">
                <span class="metric-label">Produk Aktif</span>
                <h2 class="metric-value text-gradient">{{ salesData.total_products }}</h2>
              </div>
              <div class="metric-footer">
                <span class="metric-subtext">Tayang di katalog belanja</span>
              </div>
            </div>

            
            <div class="metric-card glass card-glow-indigo">
              <div class="metric-header">
                <span class="metric-icon">👥</span>
                <span class="metric-badge info">User</span>
              </div>
              <div class="metric-body">
                <span class="metric-label">Total Pelanggan</span>
                <h2 class="metric-value text-gradient">{{ salesData.total_customers }}</h2>
              </div>
              <div class="metric-footer">
                <span class="metric-subtext">Pengguna terdaftar</span>
              </div>
            </div>
          </div>

          
          <div class="charts-grid">
            
            <div class="chart-box glass">
              <h3 class="chart-title">Tren Pendapatan Harian (7 Hari Terakhir)</h3>
              <div class="chart-container">
                <div v-if="!salesData.daily_sales || salesData.daily_sales.length === 0" class="no-chart-data">
                  Belum ada data penjualan 7 hari terakhir.
                </div>
                <div v-else class="svg-chart-wrapper">
                  <svg viewBox="0 0 600 180" class="line-chart">
                    <defs>
                      <linearGradient id="chartGrad" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stop-color="var(--color-amber)" stop-opacity="0.25"/>
                        <stop offset="100%" stop-color="var(--color-amber)" stop-opacity="0"/>
                      </linearGradient>
                      <linearGradient id="lineGrad" x1="0" y1="0" x2="1" y2="0">
                        <stop offset="0%" stop-color="var(--color-burnt-orange)"/>
                        <stop offset="100%" stop-color="var(--color-amber)"/>
                      </linearGradient>
                    </defs>
                    <line x1="30" y1="30" x2="570" y2="30" stroke="rgba(255,255,255,0.03)" stroke-width="1" />
                    <line x1="30" y1="80" x2="570" y2="80" stroke="rgba(255,255,255,0.03)" stroke-width="1" />
                    <line x1="30" y1="130" x2="570" y2="130" stroke="rgba(255,255,255,0.03)" stroke-width="1" />
                    <line x1="30" y1="155" x2="570" y2="155" stroke="rgba(255,255,255,0.08)" stroke-width="1.5" />
                    <path :d="chartAreaPath" fill="url(#chartGrad)" />
                    <path :d="chartPath" fill="none" stroke="url(#lineGrad)" stroke-width="3" stroke-linecap="round" />
                    <g class="chart-dots">
                      <circle 
                        v-for="(p, i) in chartPoints" 
                        :key="i" 
                        :cx="p.x" 
                        :cy="p.y" 
                        r="5" 
                        fill="var(--color-amber-light)" 
                        stroke="var(--color-bg)" 
                        stroke-width="1.5"
                      >
                        <title>{{ p.date }}: {{ p.revenue }} ({{ p.orders }} order)</title>
                      </circle>
                    </g>
                  </svg>
                  <div class="chart-x-labels">
                    <span v-for="(p, i) in chartPoints" :key="i" class="x-label">{{ p.shortDate }}</span>
                  </div>
                </div>
              </div>
            </div>

            
            <div class="chart-box glass">
              <h3 class="chart-title">Pendapatan Berdasarkan Asal Wilayah</h3>
              <div class="region-breakdown">
                <div v-if="!salesData.revenue_by_region || salesData.revenue_by_region.length === 0" class="no-chart-data">
                  Belum ada catatan pendapatan wilayah.
                </div>
                <div v-else class="region-list">
                  <div v-for="r in salesData.revenue_by_region" :key="r.region" class="region-item">
                    <div class="region-info">
                      <span class="region-name">🏝️ {{ r.region }}</span>
                      <span class="region-val">{{ formatCurrency(r.revenue) }} ({{ r.count }} PO)</span>
                    </div>
                    <div class="region-progress-bg">
                      <div class="region-progress-fill" :style="{ width: `${getRegionPercentage(r.revenue)}%` }"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          
          <div class="recent-orders-box glass">
            <div class="box-header">
              <h3 class="box-title">Log Transaksi Pre-Order Terbaru</h3>
              <div class="box-stats">
                <span class="badge badge-indigo">Menunggu: {{ salesData.pending_orders }}</span>
                <span class="badge badge-amber">Selesai/Kirim: {{ salesData.completed_orders }}</span>
              </div>
            </div>

            <div class="table-wrapper">
              <table class="orders-table">
                <thead>
                  <tr>
                    <th>ID Order</th>
                    <th>Pelanggan</th>
                    <th>Produk NTT</th>
                    <th>Daerah</th>
                    <th>Kuantitas</th>
                    <th>Total Biaya</th>
                    <th>Tanggal Masuk</th>
                    <th>Status</th>
                    <th>Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-if="!salesData.recent_orders || salesData.recent_orders.length === 0">
                    <td colspan="8" class="empty-table">Belum ada transaksi terdaftar.</td>
                  </tr>
                  <tr v-for="order in salesData.recent_orders" :key="order.id">
                    <td><code>#PO-{{ order.id }}</code></td>
                    <td>
                      <div class="user-cell">
                        <div class="user-avatar-sm">{{ order.customer?.name?.[0] || 'C' }}</div>
                        <div>
                          <div class="user-cell-name">{{ order.customer?.name || 'Customer' }}</div>
                          <div class="user-cell-sub">{{ order.customer?.phone || '-' }}</div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div class="product-cell">
                        <img :src="order.product?.image_url" alt="" class="product-img-sm" @error="onImageError" />
                        <span class="product-cell-name">{{ order.product?.name || 'Produk' }}</span>
                      </div>
                    </td>
                    <td><span class="region-tag">{{ order.product?.region || 'NTT' }}</span></td>
                    <td>{{ order.quantity }} pcs</td>
                    <td class="total-price-cell">{{ formatCurrency(order.total_price) }}</td>
                    <td class="date-cell">{{ formatDateLong(order.created_at) }}</td>
                    <td>
                      <span class="status-badge" :class="order.status">
                        {{ order.status.toUpperCase() }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button class="btn-action delete-action" @click="handleDeleteOrder(order.id)" title="Hapus Log Transaksi">
                        🗑️ Hapus
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        
        <div v-if="activeTab === 'products'" class="tab-pane animate-fade-in">
          <div class="product-crud-box glass">
            <div class="crud-header">
              <div class="crud-search-bar">
                <svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
                </svg>
                <input 
                  v-model="productSearch" 
                  type="text" 
                  class="input search-input" 
                  placeholder="Cari nama produk katalog..."
                />
              </div>
              <button class="btn btn-primary btn-add-product" @click="openAddProductModal">
                ➕ Tambah Produk Baru
              </button>
            </div>

            <div class="table-wrapper">
              <table class="orders-table">
                <thead>
                  <tr>
                    <th>Foto</th>
                    <th>Nama Produk</th>
                    <th>Kategori</th>
                    <th>Asal Daerah</th>
                    <th>Harga</th>
                    <th>Stok Tersedia</th>
                    <th>Min. Order</th>
                    <th>Durasi PO</th>
                    <th>Status</th>
                    <th class="actions-col">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-if="filteredProducts.length === 0">
                    <td colspan="10" class="empty-table">Tidak ada produk katalog ditemukan.</td>
                  </tr>
                  <tr v-for="prod in filteredProducts" :key="prod.id">
                    <td>
                      <img :src="prod.image_url" alt="" class="product-img-sm crud-img" @error="onImageError" />
                    </td>
                    <td class="product-name-cell">
                      <strong>{{ prod.name }}</strong>
                    </td>
                    <td><span class="badge badge-indigo">{{ prod.category }}</span></td>
                    <td><span class="region-tag">{{ prod.region }}</span></td>
                    <td class="price-val-cell">{{ formatCurrency(prod.price) }}</td>
                    <td>
                      <span class="stock-display" :class="{ 'out-of-stock': prod.stock <= 0 }">
                        {{ prod.stock }} pcs
                      </span>
                    </td>
                    <td>{{ prod.min_order }} pcs</td>
                    <td>{{ prod.po_duration }} hari</td>
                    <td>
                      <span class="status-badge" :class="prod.is_active ? 'completed' : 'cancelled'">
                        {{ prod.is_active ? 'AKTIF' : 'NON-AKTIF' }}
                      </span>
                    </td>
                    <td class="actions-cell">
                      <button class="btn-action edit-action" @click="openEditProductModal(prod)" title="Edit Produk">
                        ✏️ Edit
                      </button>
                      <button class="btn-action delete-action" @click="handleDeleteProduct(prod.id)" title="Hapus Produk">
                        🗑️ Hapus
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>

    
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-card glass-strong animate-scale-in">
        <div class="modal-header">
          <h3>{{ isEditing ? '✏️ Edit Detail Produk' : '✨ Tambah Produk Katalog Baru' }}</h3>
          <button class="modal-close" @click="showModal = false">&times;</button>
        </div>
        <form @submit.prevent="handleSaveProduct" class="modal-form">
          <div class="form-row">
            <div class="form-group flex-2">
              <label class="form-label" for="prod-name">Nama Produk</label>
              <input id="prod-name" v-model="formProduct.name" type="text" class="input" placeholder="Kain Tenun Ikat Sumba..." required />
            </div>
            <div class="form-group flex-1">
              <label class="form-label" for="prod-price">Harga (IDR)</label>
              <input id="prod-price" v-model.number="formProduct.price" type="number" class="input" min="1" required />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label" for="prod-region">Asal Daerah NTT</label>
              <select id="prod-region" v-model="formProduct.region" class="input">
                <option value="Sumba">Sumba</option>
                <option value="Manggarai">Manggarai</option>
                <option value="Flores">Flores</option>
                <option value="Kupang">Kupang</option>
                <option value="Timor">Timor</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label" for="prod-category">Kategori</label>
              <select id="prod-category" v-model="formProduct.category" class="input">
                <option value="Tenun">Tenun</option>
                <option value="Kuliner">Kuliner</option>
                <option value="Kerajinan">Kerajinan</option>
                <option value="Aksesoris">Aksesoris</option>
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label" for="prod-stock">Stok Tersedia (Ditentukan Admin)</label>
              <input id="prod-stock" v-model.number="formProduct.stock" type="number" class="input" min="0" required />
            </div>
            <div class="form-group">
              <label class="form-label" for="prod-minorder">Minimal Pembelian (Min Order)</label>
              <input id="prod-minorder" v-model.number="formProduct.min_order" type="number" class="input" min="1" required />
            </div>
            <div class="form-group">
              <label class="form-label" for="prod-duration">Durasi Pre-Order (Hari)</label>
              <input id="prod-duration" v-model.number="formProduct.po_duration" type="number" class="input" min="1" required />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label" for="prod-image">URL Foto/Gambar Produk</label>
            <input id="prod-image" v-model="formProduct.image_url" type="text" class="input" placeholder="https://images.unsplash.com/..." required />
          </div>

          <div class="form-group">
            <label class="form-label" for="prod-desc">Deskripsi Produk</label>
            <textarea id="prod-desc" v-model="formProduct.description" class="input textarea" rows="3" placeholder="Deskripsikan keunikan motif adat..." required></textarea>
          </div>

          <div class="form-group checkbox-group">
            <input id="prod-active" v-model="formProduct.is_active" type="checkbox" class="checkbox-input" />
            <label for="prod-active" class="checkbox-label">Tampilkan Produk (Aktif dalam katalog belanja)</label>
          </div>

          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="showModal = false">Batal</button>
            <button type="submit" class="btn btn-primary" :disabled="formSubmitting">
              {{ formSubmitting ? 'Menyimpan...' : 'Simpan Perubahan' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    
    <ChatBox :visible="chatOpen" @close="chatOpen = false" @update-unread="unreadCount = $event" />
    <button v-show="!chatOpen" class="chat-toggle-btn btn btn-primary" @click="chatOpen = true">
      <span class="chat-toggle-badge" v-if="unreadCount > 0">{{ unreadCount }}</span>
      <span class="chat-toggle-icon">💬</span>
      <span>Chat Pelanggan</span>
    </button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ChatBox from '../components/ChatBox.vue'

const router = useRouter()
const salesData = ref({
  total_revenue: 0,
  total_pre_orders: 0,
  total_products: 0,
  total_customers: 0,
  pending_orders: 0,
  completed_orders: 0,
  revenue_by_region: [],
  recent_orders: [],
  orders_by_status: [],
  daily_sales: []
})

const products = ref([])
const activeTab = ref('summary')
const loading = ref(true)
const salesDataLoading = ref(true)
const error = ref('')
const chatOpen = ref(false)
const unreadCount = ref(0)


const showModal = ref(false)
const isEditing = ref(false)
const formSubmitting = ref(false)
const productSearch = ref('')
const formProduct = ref({
  id: null,
  name: '',
  description: '',
  price: 150000,
  image_url: '',
  region: 'Sumba',
  category: 'Tenun',
  min_order: 1,
  po_duration: 14,
  stock: 10,
  is_active: true
})

const currentUser = computed(() => {
  return JSON.parse(localStorage.getItem('openpeo_user') || 'null')
})


const filteredProducts = computed(() => {
  if (!productSearch.value.trim()) return products.value
  const query = productSearch.value.toLowerCase()
  return products.value.filter(p => 
    p.name.toLowerCase().includes(query) ||
    p.category.toLowerCase().includes(query) ||
    p.region.toLowerCase().includes(query)
  )
})


async function fetchProducts() {
  try {
    const response = await fetch('http://localhost:8080/api/products?limit=100')
    const data = await response.json()
    if (data.success && data.data) {
      products.value = data.data
    }
  } catch (err) {
    console.warn('Gagal memuat daftar produk', err)
  }
}


async function fetchSalesData() {
  salesDataLoading.value = true
  try {
    const response = await fetch('http://localhost:8080/api/admin/sales-data')
    const data = await response.json()
    if (data.success && data.data) {
      salesData.value = data.data
    } else {
      error.value = data.message || 'Gagal memproses data aggregator'
    }
  } catch (err) {
    error.value = 'Koneksi API gagal. Pastikan backend engine berjalan di port 8080.'
  } finally {
    salesDataLoading.value = false
  }
}


async function initializeDashboard() {
  loading.value = true
  error.value = ''
  await Promise.all([fetchSalesData(), fetchProducts()])
  loading.value = false
}


function openAddProductModal() {
  isEditing.value = false
  formProduct.value = {
    id: null,
    name: '',
    description: '',
    price: 100000,
    image_url: 'https://images.unsplash.com/photo-1558171813-4c088753af8f?w=400',
    region: 'Sumba',
    category: 'Tenun',
    min_order: 1,
    po_duration: 14,
    stock: 10,
    is_active: true
  }
  showModal.value = true
}

function openEditProductModal(prod) {
  isEditing.value = true
  formProduct.value = {
    id: prod.id,
    name: prod.name,
    description: prod.description,
    price: prod.price,
    image_url: prod.image_url,
    region: prod.region,
    category: prod.category,
    min_order: prod.min_order,
    po_duration: prod.po_duration,
    stock: prod.stock,
    is_active: prod.is_active
  }
  showModal.value = true
}

async function handleSaveProduct() {
  if (formProduct.value.price <= 0 || formProduct.value.stock < 0) {
    alert('Harga dan Stok harus bernilai valid!')
    return
  }

  formSubmitting.value = true
  
  const url = isEditing.value 
    ? `http://localhost:8080/api/products/${formProduct.value.id}`
    : 'http://localhost:8080/api/products'
    
  const method = isEditing.value ? 'PUT' : 'POST'
  const payload = {
    ...formProduct.value,
    vendor_id: currentUser.value ? currentUser.value.id : 1
  }

  try {
    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    const data = await response.json()
    if (data.success) {
      showModal.value = false
      await initializeDashboard()
    } else {
      alert(data.message || 'Gagal menyimpan perubahan produk')
    }
  } catch (err) {
    alert('Koneksi ke server gagal, pastikan API server menyala')
  } finally {
    formSubmitting.value = false
  }
}

async function handleDeleteProduct(id) {
  if (!confirm('Apakah Anda yakin ingin menghapus produk ini secara permanen dari database?')) return
  
  try {
    const response = await fetch(`http://localhost:8080/api/products/${id}?user_id=${currentUser.value.id}`, {
      method: 'DELETE'
    })
    const data = await response.json()
    if (data.success) {
      await initializeDashboard()
    } else {
      alert(data.message || 'Gagal menghapus produk')
    }
  } catch (err) {
    alert('Koneksi ke server gagal, pastikan API server menyala')
  }
}

async function handleDeleteOrder(id) {
  if (!confirm('Apakah Anda yakin ingin menghapus log transaksi pre-order ini?')) return
  
  try {
    const response = await fetch(`http://localhost:8080/api/orders/${id}?user_id=${currentUser.value.id}`, {
      method: 'DELETE'
    })
    const data = await response.json()
    if (data.success) {
      await initializeDashboard()
    } else {
      alert(data.message || 'Gagal menghapus log transaksi')
    }
  } catch (err) {
    alert('Koneksi ke server gagal, pastikan API server menyala')
  }
}


function onImageError(e) {
  e.target.src = 'https://images.unsplash.com/photo-1558171813-4c088753af8f?w=400'
}


function formatCurrency(val) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(val)
}


const maxRegionRevenue = computed(() => {
  if (!salesData.value?.revenue_by_region?.length) return 1
  return Math.max(...salesData.value.revenue_by_region.map(r => r.revenue), 1)
})

function getRegionPercentage(rev) {
  return (rev / maxRegionRevenue.value) * 100
}


const chartPoints = computed(() => {
  if (!salesData.value?.daily_sales?.length) return []
  const data = salesData.value.daily_sales
  const width = 600
  const height = 180
  const paddingX = 40
  const paddingY = 25
  const chartWidth = width - paddingX * 2
  const chartHeight = height - paddingY * 2
  
  const maxVal = Math.max(...data.map(d => d.revenue), 100000)
  
  return data.map((d, i) => {
    const x = paddingX + (i / (data.length - 1 || 1)) * chartWidth
    const y = height - paddingY - (d.revenue / maxVal) * chartHeight
    return {
      x,
      y,
      shortDate: formatShortDate(d.date),
      date: formatDateLong(d.date),
      revenue: formatCurrency(d.revenue),
      orders: d.orders,
      raw: d
    }
  })
})

const chartPath = computed(() => {
  const points = chartPoints.value
  if (!points.length) return ''
  return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ')
})

const chartAreaPath = computed(() => {
  const points = chartPoints.value
  if (!points.length) return ''
  const first = points[0]
  const last = points[points.length - 1]
  const height = 180
  const paddingY = 25
  const linePath = points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ')
  return `${linePath} L ${last.x} ${height - paddingY} L ${first.x} ${height - paddingY} Z`
})


function formatShortDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

function formatDateLong(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

function logout() {
  localStorage.removeItem('openpeo_user')
  router.push('/')
}

onMounted(() => {
  initializeDashboard()
})
</script>

<style scoped>
.admin-dashboard {
  min-height: 100vh;
  background: var(--color-bg);
  padding-top: 5.5rem;
  padding-bottom: var(--space-3xl);
}

.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  padding: 0.75rem 0;
}

.navbar-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 var(--space-xl);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.brand-icon svg {
  width: 36px;
  height: 36px;
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
}

.brand-name {
  font-family: var(--font-display);
  font-size: 1.3rem;
  font-weight: 800;
}

.brand-tagline {
  font-size: 0.65rem;
  color: var(--color-text-muted);
  letter-spacing: 0.15em;
  text-transform: uppercase;
}

.navbar-links {
  display: flex;
  gap: var(--space-lg);
}

.nav-link {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
  padding: 0.5rem 0;
  position: relative;
  transition: color var(--transition-fast);
}

.nav-link::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 0;
  height: 2px;
  background: var(--gradient-amber);
  border-radius: 1px;
  transition: width var(--transition-base);
}

.nav-link:hover, .nav-link.active {
  color: var(--color-text-primary);
}

.nav-link:hover::after, .nav-link.active::after {
  width: 100%;
}

.navbar-actions {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.user-profile-nav {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.user-greeting {
  font-size: 0.85rem;
  color: var(--color-text-secondary);
}

.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.8rem;
}

.dashboard-container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 var(--space-xl);
}

.loading-state, .error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 6rem var(--space-xl);
  text-align: center;
  border-radius: var(--radius-xl);
}

.error-state {
  background: rgba(231, 76, 60, 0.05);
  border: 1px solid rgba(231, 76, 60, 0.12);
  max-width: 500px;
  margin: 4rem auto;
}

.error-icon {
  font-size: 2.5rem;
  margin-bottom: var(--space-md);
}

.error-state h3 {
  font-size: 1.25rem;
  margin-bottom: var(--space-xs);
}

.error-state p {
  color: var(--color-text-secondary);
  font-size: 0.9rem;
  margin-bottom: var(--space-lg);
}

.dashboard-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-xl);
  flex-wrap: wrap;
  gap: var(--space-md);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

.dashboard-title {
  font-family: var(--font-display);
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 0.35rem;
}

.dashboard-subtitle {
  color: var(--color-text-secondary);
  font-size: 0.92rem;
}

.btn-refresh {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.spin-active {
  animation: spin-slow 1s linear infinite;
}


.tab-navigation {
  display: flex;
  padding: 0.3rem;
  border-radius: var(--radius-lg);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.tab-btn {
  padding: 0.5rem 1.2rem;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  background: transparent;
  border-radius: var(--radius-md);
  border: none;
  transition: all var(--transition-fast);
}

.tab-btn:hover {
  color: var(--color-text-primary);
  background: rgba(255, 255, 255, 0.03);
}

.tab-btn.active {
  color: var(--color-bg);
  background: var(--gradient-amber);
  box-shadow: var(--shadow-sm);
}

.tab-pane {
  animation: fadeInUp 0.4s ease-out;
}


.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: var(--space-lg);
  margin-bottom: var(--space-2xl);
}

.metric-card {
  padding: var(--space-xl);
  border-radius: var(--radius-xl);
  border: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  transition: transform var(--transition-base), box-shadow var(--transition-base);
}

.metric-card:hover {
  transform: translateY(-4px);
}

.card-glow-amber:hover {
  box-shadow: 0 8px 30px rgba(245, 166, 35, 0.08), inset 0 0 12px rgba(245, 166, 35, 0.04);
  border-color: rgba(245, 166, 35, 0.2);
}

.card-glow-indigo:hover {
  box-shadow: 0 8px 30px rgba(99, 102, 241, 0.08), inset 0 0 12px rgba(99, 102, 241, 0.04);
  border-color: rgba(99, 102, 241, 0.2);
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.metric-icon {
  font-size: 1.8rem;
}

.metric-badge {
  font-size: 0.68rem;
  font-weight: 700;
  padding: 0.15rem 0.45rem;
  border-radius: var(--radius-sm);
  text-transform: uppercase;
}

.metric-badge.positive {
  background: rgba(46, 204, 113, 0.15);
  color: var(--color-success);
}

.metric-badge.info {
  background: rgba(52, 152, 219, 0.15);
  color: #3498db;
}

.metric-body {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.metric-label {
  font-size: 0.82rem;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.metric-value {
  font-family: var(--font-display);
  font-size: 1.65rem;
  font-weight: 800;
}

.metric-footer {
  font-size: 0.72rem;
  color: var(--color-text-muted);
  border-top: 1px solid rgba(255, 255, 255, 0.04);
  padding-top: var(--space-xs);
  margin-top: auto;
}


.charts-grid {
  display: grid;
  grid-template-columns: 3fr 2fr;
  gap: var(--space-lg);
  margin-bottom: var(--space-2xl);
}

@media (max-width: 992px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
}

.chart-box {
  padding: var(--space-lg);
  border-radius: var(--radius-xl);
}

.chart-title {
  font-size: 0.95rem;
  font-weight: 700;
  margin-bottom: var(--space-lg);
  color: var(--color-text-secondary);
  border-left: 3px solid var(--color-amber);
  padding-left: var(--space-sm);
}

.chart-container {
  height: 220px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.no-chart-data {
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.svg-chart-wrapper {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.line-chart {
  width: 100%;
  height: 180px;
  overflow: visible;
}

.chart-dots circle {
  transition: r var(--transition-fast), fill var(--transition-fast);
  cursor: pointer;
}

.chart-dots circle:hover {
  r: 8px;
  fill: var(--color-burnt-orange);
}

.chart-x-labels {
  display: flex;
  justify-content: space-between;
  padding: 0 40px;
  font-size: 0.68rem;
  color: var(--color-text-muted);
  font-weight: 500;
}


.region-breakdown {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 200px;
}

.region-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.region-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.region-info {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
}

.region-name {
  font-weight: 600;
}

.region-val {
  color: var(--color-text-secondary);
  font-size: 0.76rem;
}

.region-progress-bg {
  height: 6px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.region-progress-fill {
  height: 100%;
  background: var(--gradient-amber);
  border-radius: var(--radius-full);
}


.recent-orders-box, .product-crud-box {
  padding: var(--space-xl);
  border-radius: var(--radius-xl);
  margin-bottom: var(--space-xl);
}

.box-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-lg);
  flex-wrap: wrap;
  gap: var(--space-sm);
}

.box-title {
  font-size: 1.05rem;
  font-weight: 700;
  border-left: 3px solid var(--color-indigo-glow);
  padding-left: var(--space-sm);
}

.box-stats {
  display: flex;
  gap: var(--space-xs);
}

.table-wrapper {
  overflow-x: auto;
  border-radius: var(--radius-md);
  border: 1px solid rgba(255, 255, 255, 0.04);
}

.orders-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.82rem;
}

.orders-table th, .orders-table td {
  padding: var(--space-md);
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}

.orders-table th {
  background: rgba(255, 255, 255, 0.02);
  font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.02em;
}

.orders-table tr:hover {
  background: rgba(255, 255, 255, 0.01);
}

.orders-table code {
  color: var(--color-amber-light);
  background: rgba(245, 166, 35, 0.08);
  padding: 0.15rem 0.4rem;
  border-radius: var(--radius-sm);
}

.empty-table {
  text-align: center;
  padding: var(--space-2xl) !important;
  color: var(--color-text-muted);
}


.user-cell, .product-cell {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
}

.user-avatar-sm {
  width: 28px;
  height: 28px;
  border-radius: var(--radius-full);
  background: var(--color-indigo);
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 700;
}

.user-cell-name {
  font-weight: 600;
}

.user-cell-sub {
  font-size: 0.72rem;
  color: var(--color-text-muted);
}

.product-img-sm {
  width: 32px;
  height: 32px;
  object-fit: cover;
  border-radius: var(--radius-sm);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.product-cell-name, .product-name-cell {
  font-weight: 500;
}

.product-cell-name {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.region-tag {
  background: rgba(255, 255, 255, 0.03);
  padding: 0.15rem 0.4rem;
  border-radius: var(--radius-sm);
  font-size: 0.75rem;
  color: var(--color-text-secondary);
}

.total-price-cell, .price-val-cell {
  font-weight: 700;
}

.date-cell {
  color: var(--color-text-secondary);
}


.status-badge {
  font-size: 0.68rem;
  padding: 0.2rem 0.5rem;
  border-radius: var(--radius-sm);
  font-weight: 700;
  letter-spacing: 0.02em;
  display: inline-block;
}

.status-badge.pending {
  background: rgba(245, 166, 35, 0.1);
  color: var(--color-amber);
  border: 1px solid rgba(245, 166, 35, 0.2);
}

.status-badge.completed {
  background: rgba(46, 204, 113, 0.1);
  color: var(--color-success);
  border: 1px solid rgba(46, 204, 113, 0.2);
}

.status-badge.shipped {
  background: rgba(52, 152, 219, 0.1);
  color: #3498db;
  border: 1px solid rgba(52, 152, 219, 0.2);
}

.status-badge.cancelled {
  background: rgba(231, 76, 60, 0.1);
  color: var(--color-error);
  border: 1px solid rgba(231, 76, 60, 0.2);
}


.crud-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-lg);
  gap: var(--space-md);
  flex-wrap: wrap;
}

.crud-search-bar {
  position: relative;
  min-width: 280px;
}

.crud-img {
  width: 45px;
  height: 45px;
  border-radius: var(--radius-sm);
}

.actions-col {
  width: 140px;
}

.actions-cell {
  display: flex;
  gap: 0.4rem;
}

.btn-action {
  padding: 0.35rem 0.65rem;
  font-size: 0.72rem;
  font-weight: 600;
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: none;
  transition: opacity var(--transition-fast);
}

.btn-action:hover {
  opacity: 0.85;
}

.edit-action {
  background: #3498db;
  color: white;
}

.delete-action {
  background: var(--color-error);
  color: white;
}

.stock-display.out-of-stock {
  color: var(--color-error);
  font-weight: 700;
}


.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 2000;
  background: rgba(11, 10, 18, 0.75);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-md);
}

.modal-card {
  width: 100%;
  max-width: 650px;
  border-radius: var(--radius-2xl);
  overflow: hidden;
  box-shadow: var(--shadow-lg);
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-lg) var(--space-xl);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.modal-header h3 {
  font-family: var(--font-display);
  font-size: 1.25rem;
  font-weight: 700;
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.8rem;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: color var(--transition-fast);
}

.modal-close:hover {
  color: var(--color-text-primary);
}

.modal-form {
  padding: var(--space-xl);
  display: flex;
  flex-direction: column;
  gap: var(--space-lg);
  max-height: 80vh;
  overflow-y: auto;
}

.form-row {
  display: flex;
  gap: var(--space-md);
  flex-wrap: wrap;
}

.form-row .form-group {
  flex: 1;
  min-width: 140px;
}

.form-row .flex-2 {
  flex: 2;
}

.form-row .flex-1 {
  flex: 1;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.form-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.textarea {
  resize: vertical;
}

.checkbox-group {
  flex-direction: row;
  align-items: center;
  gap: var(--space-sm);
  padding: 0.3rem 0;
}

.checkbox-input {
  width: 16px;
  height: 16px;
  accent-color: var(--color-amber);
}

.checkbox-label {
  font-size: 0.85rem;
  color: var(--color-text-secondary);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-md);
  margin-top: var(--space-md);
  padding-top: var(--space-md);
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}


.chat-toggle-btn {
  position: fixed;
  bottom: var(--space-xl);
  right: var(--space-xl);
  z-index: 1000;
  border-radius: var(--radius-full);
  padding: 0.8rem 1.5rem;
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  box-shadow: var(--shadow-lg);
  animation: float 4s ease-in-out infinite;
}

.chat-toggle-icon {
  font-size: 1.2rem;
}

.chat-toggle-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: var(--color-error);
  color: white;
  font-size: 0.65rem;
  font-weight: 700;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(255,255,255,0.05);
  border-top-color: var(--color-amber);
  border-radius: 50%;
  animation: spin-slow 0.8s linear infinite;
  margin-bottom: var(--space-sm);
}

@keyframes spin-slow {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(15px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-scale-in {
  animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

@media (max-width: 768px) {
  .navbar-links {
    display: none;
  }
  .crud-header {
    flex-direction: column;
    align-items: stretch;
  }
  .crud-search-bar {
    min-width: 100%;
  }
}
</style>
