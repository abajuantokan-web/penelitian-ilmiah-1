<template>
  <div class="seller-layout-container">
    <div class="seller-layout max-w-7xl mx-auto w-full">
    <aside class="seller-sidebar">
      <div class="sidebar-header">
        <h2 class="relative inline-block">
          {{ authStore.user?.store_name || 'Toko Saya' }}
        </h2>
        <p>Seller Dashboard</p>
      </div>
      <nav class="sidebar-nav">
        <a href="#" :class="{ active: currentTab === 'dashboard' }" @click.prevent="currentTab = 'dashboard'">Dashboard Overview</a>
        <a href="#" :class="{ active: currentTab === 'products' }" @click.prevent="currentTab = 'products'">Produk Saya</a>
        <a href="#" :class="{ active: currentTab === 'orders' }" @click.prevent="currentTab = 'orders'">
          <span>Pesanan Masuk</span>
          <span v-if="pendingCount > 0" class="badge">{{ pendingCount }}</span>
        </a>
        <a href="#" :class="{ active: currentTab === 'keuangan' }" @click.prevent="currentTab = 'keuangan'">Keuangan</a>
        <a href="#" :class="{ active: currentTab === 'chat' }" @click.prevent="currentTab = 'chat'">
          <span>Pesan</span>
          <span v-if="unreadChatCount > 0" class="badge">{{ unreadChatCount }}</span>
        </a>
        <a href="#" :class="{ active: currentTab === 'settings' }" @click.prevent="currentTab = 'settings'">Pengaturan Toko</a>
      </nav>
    </aside>

    <main class="seller-content">
      <div class="content-header">
        <h1 v-if="currentTab === 'products'">Produk Saya</h1>
        <h1 v-else-if="currentTab === 'dashboard'">Dashboard</h1>
        <h1 v-else-if="currentTab === 'orders'">Pesanan Masuk</h1>
        <h1 v-else-if="currentTab === 'keuangan'">Keuangan</h1>
        <h1 v-else-if="currentTab === 'chat'">Pesan Pembeli</h1>
        <h1 v-else-if="currentTab === 'settings'">Pengaturan Toko</h1>
        
        <button v-if="currentTab === 'products'" class="btn-dark" @click="openAddModal">
          + Tambah Produk
        </button>
      </div>

      <!-- PRODUCTS TAB -->
      <div v-if="currentTab === 'products'">
        <div v-if="sellerStore.isLoading" class="loading-state">Memuat produk...</div>
        <div v-else-if="sellerStore.products.length === 0" class="empty-state">
          Belum ada produk di toko Anda.
        </div>
        <div v-else class="table-container">
          <table class="products-table">
            <thead>
              <tr>
                <th>Produk</th>
                <th>Harga</th>
                <th>Stok</th>
                <th class="text-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="product in sellerStore.products" :key="product.id">
                <td class="td-product">
                  <img :src="$getImageUrl(product.image_url)" :alt="product.name" class="table-img" />
                  <span class="table-product-name">{{ product.name }}</span>
                </td>
                <td>Rp {{ product.price.toLocaleString('id-ID') }}</td>
                <td>{{ product.stock }}</td>
                <td class="td-actions">
                  <button class="btn-outline btn-sm" @click="openEditModal(product)">Edit</button>
                  <button class="btn-danger btn-sm" @click="handleDelete(product.id)">Hapus</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ORDERS TAB -->
      <div v-else-if="currentTab === 'orders'">
        <div v-if="isLoadingOrders" class="loading-state">Memuat pesanan...</div>
        <div v-else-if="sellerOrders.length === 0" class="empty-state">
          Belum ada pesanan masuk.
        </div>
        <div v-else class="table-container">
          <table class="products-table">
            <thead>
              <tr>
                <th>Produk</th>
                <th>Pembeli</th>
                <th>Total</th>
                <th>Status</th>
                <th class="text-right">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="order in sellerOrders" :key="order.id">
                <td class="td-product">
                  <div v-if="order.order_items && order.order_items.length > 0">
                    <div v-for="item in order.order_items" :key="item.id" class="flex items-center gap-3 mb-2 last:mb-0">
                      <img :src="$getImageUrl(item.product?.image_url)" :alt="item.product?.name" class="table-img" style="width: 32px; height: 32px; border-radius: 4px;" />
                      <div>
                        <span class="table-product-name" style="display:block; font-size: 0.85rem;">{{ item.product?.name }}</span>
                        <span style="font-size: 0.75rem; color: #6b7280;">x{{ item.quantity }}</span>
                      </div>
                    </div>
                  </div>
                  <div v-else class="flex items-center gap-3">
                    <img :src="$getImageUrl(order.product?.image_url)" :alt="order.product?.name" class="table-img" />
                    <div>
                      <span class="table-product-name" style="display:block">{{ order.product?.name }}</span>
                      <span style="font-size: 0.75rem; color: #6b7280;">x{{ order.quantity }}</span>
                    </div>
                  </div>
                </td>
                <td>
                  <div style="font-weight: 500;">{{ order.customer?.name || 'Anonim' }}</div>
                  <div style="font-size: 0.75rem; color: #6b7280;">{{ order.note }}</div>
                </td>
                <td>Rp {{ order.total_price.toLocaleString('id-ID') }}</td>
                <td>
                  <span :class="['status-badge', getOrderStatusClass(order.status)]">{{ order.status }}</span>
                </td>
                <td class="td-actions">
                  <button 
                    v-if="order.status === 'Menunggu Konfirmasi'" 
                    class="btn-dark btn-sm" 
                    @click="processOrder(order.id)"
                    :disabled="isProcessingOrder === order.id">
                    {{ isProcessingOrder === order.id ? 'Memproses...' : 'Konfirmasi Pesanan' }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- PENGATURAN TOKO TAB -->
      <div v-else-if="currentTab === 'settings'">
        <PengaturanToko />
      </div>

      <!-- KEUANGAN TAB -->
      <div v-else-if="currentTab === 'keuangan'">
        <KeuanganView />
      </div>
      <!-- DASHBOARD OVERVIEW TAB -->
      <div v-else-if="currentTab === 'dashboard'">
        <div v-if="dashboardStore.isLoading" class="loading-state">Memuat data analitik...</div>
        <div v-else class="dashboard-overview">
          <!-- Stats Cards -->
          <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <div class="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 flex items-center gap-5 transition-transform hover:-translate-y-1">
              <div class="w-14 h-14 rounded-xl flex items-center justify-center bg-zinc-50 border border-zinc-200 text-zinc-800">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"></circle>
                  <polyline points="12 6 12 12 16 14"></polyline>
                </svg>
              </div>
              <div>
                <h3 class="text-gray-500 font-medium text-sm mb-1">Menunggu Konfirmasi</h3>
                <p class="text-zinc-900 font-semibold text-3xl">{{ pendingCount }}</p>
              </div>
            </div>
            
            <div class="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 flex items-center gap-5 transition-transform hover:-translate-y-1">
              <div class="w-14 h-14 rounded-xl flex items-center justify-center bg-zinc-50 border border-zinc-200 text-zinc-800">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="1" y="3" width="15" height="13"></rect>
                  <polygon points="16 8 20 8 23 11 23 16 16 16 16 8"></polygon>
                  <circle cx="5.5" cy="18.5" r="2.5"></circle>
                  <circle cx="18.5" cy="18.5" r="2.5"></circle>
                </svg>
              </div>
              <div>
                <h3 class="text-gray-500 font-medium text-sm mb-1">Sedang Diproses</h3>
                <p class="text-zinc-900 font-semibold text-3xl">{{ processingCount }}</p>
              </div>
            </div>

            <div class="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 flex items-center gap-5 transition-transform hover:-translate-y-1">
              <div class="w-14 h-14 rounded-xl flex items-center justify-center bg-zinc-50 border border-zinc-200 text-zinc-800">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                  <polyline points="22 4 12 14.01 9 11.01"></polyline>
                </svg>
              </div>
              <div>
                <h3 class="text-gray-500 font-medium text-sm mb-1">Pesanan Selesai</h3>
                <p class="text-zinc-900 font-semibold text-3xl">{{ completedCount }}</p>
              </div>
            </div>

            <div class="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 flex items-center gap-5 transition-transform hover:-translate-y-1">
              <div class="w-14 h-14 rounded-xl flex items-center justify-center bg-zinc-50 border border-zinc-200 text-zinc-800">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="2" y="6" width="20" height="12" rx="2"></rect>
                  <circle cx="12" cy="12" r="2"></circle>
                  <path d="M6 12h.01M18 12h.01"></path>
                </svg>
              </div>
              <div>
                <h3 class="text-gray-500 font-medium text-sm mb-1">Total Pendapatan</h3>
                <p class="text-zinc-900 font-semibold text-3xl">Rp {{ totalRevenue.toLocaleString('id-ID') }}</p>
              </div>
            </div>
          </div>

          <!-- Sales Chart Section -->
          <div class="bg-white shadow-sm border border-gray-100 rounded-2xl p-6 mb-8">
            <div class="flex justify-between items-center mb-6">
              <h2 class="section-title !mb-0" style="margin-bottom: 0;">Grafik Pendapatan</h2>
              <select v-model="chartFilter" class="border border-gray-200 rounded-lg px-3 py-1.5 text-sm bg-zinc-50 text-zinc-700 focus:outline-none focus:border-zinc-400">
                <option value="7_days">7 Hari Terakhir</option>
                <option value="30_days">30 Hari Terakhir</option>
                <option value="12_months">12 Bulan Terakhir</option>
              </select>
            </div>
            <div class="w-full" style="height: 320px;">
              <LineChart v-if="!dashboardStore.isLoading && chartDataConfig.labels.length" :data="chartDataConfig" :options="chartOptions" />
            </div>
          </div>

          <!-- Recent Orders -->
          <div class="w-full mb-8">
            <h2 class="section-title">Pesanan Terbaru</h2>
            <div class="table-container">
              <table class="products-table">
                <thead>
                  <tr>
                    <th>ID Pesanan</th>
                    <th>Pembeli</th>
                    <th>Total Harga</th>
                    <th>Status</th>
                    <th>Tanggal</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="order in recentOrders" :key="order.id">
                    <td>#{{ order.id }}</td>
                    <td>{{ order.customer?.name || 'Anonim' }}</td>
                    <td>Rp {{ order.total_price.toLocaleString('id-ID') }}</td>
                    <td>
                      <span :class="['status-badge', getOrderStatusClass(order.status)]">{{ order.status }}</span>
                    </td>
                    <td>{{ new Date(order.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) }}</td>
                  </tr>
                  <tr v-if="recentOrders.length === 0">
                    <td colspan="5" class="text-center py-8 text-gray-500">Belum ada pesanan terbaru.</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="currentTab !== 'chat'" class="mock-tab">
        <p>Fitur <strong>{{ currentTab }}</strong> sedang dalam pengembangan.</p>
      </div>

      <!-- CHAT TAB (Always mounted with v-show to keep WS alive) -->
      <div v-show="currentTab === 'chat'">
        <SellerChat @new-message-received="handleNewMessage" @update-unread="updateUnreadCount" />
      </div>
    </main>

    <!-- PRODUCT MODAL -->
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ isEditing ? 'Edit Produk' : 'Tambah Produk Baru' }}</h2>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        
        <form @submit.prevent="handleSubmitProduct" class="modal-form">
          <div class="custom-modal-flex">
            <!-- Left Column: Image Preview -->
            <div class="custom-modal-flex-left">
              <div class="image-preview" style="aspect-ratio: 4/5; border-radius: 8px; border: 1px solid #e5e7eb; overflow: hidden; height: 100%;">
                <img :src="$getImageUrl(form.image_url)" alt="Preview" style="width: 100%; height: 100%; object-fit: cover;" />
              </div>
            </div>
            
            <!-- Right Column: Form Fields -->
            <div class="custom-modal-flex-right">
              <div class="inner-form-grid">
                <div class="form-group">
                  <label>Nama Produk</label>
                  <input type="text" v-model="form.name" required />
                </div>
                <div class="form-group">
                  <label>Harga (Rp)</label>
                  <input type="number" v-model="form.price" required />
                </div>
              </div>
              
              <div class="form-group full-width-field">
                <label>Deskripsi</label>
                <textarea v-model="form.description" rows="3"></textarea>
              </div>
              
              <div class="inner-form-grid">
                <div class="form-group">
                  <label>Kategori</label>
                  <select v-model="form.category" required>
                    <option value="Koleksi Tenun NTT">Koleksi Tenun NTT</option>
                    <option value="Cita Rasa Lokal">Cita Rasa Lokal</option>
                    <option value="Koleksi Aksesoris">Koleksi Aksesoris</option>
                  </select>
                </div>
                <div class="form-group">
                  <label>Region</label>
                  <select v-model="form.region" required>
                    <option value="Sumba">Sumba</option>
                    <option value="Sabu">Sabu</option>
                    <option value="Amarasi">Amarasi</option>
                    <option value="Rote">Rote</option>
                    <option value="Ende">Ende</option>
                    <option value="Manggarai">Manggarai</option>
                    <option value="Alor">Alor</option>
                    <option value="Kupang">Kupang</option>
                    <option value="Timor">Timor</option>
                    <option value="Flores">Flores</option>
                    <option value="NTT">NTT</option>
                  </select>
                </div>
              </div>
              
              <div class="inner-form-grid">
                <div class="form-group">
                  <label>Stok</label>
                  <input type="number" v-model="form.stock" />
                </div>
                <div class="form-group">
                  <label>Image</label>
                  <div class="image-upload-flex">
                    <input type="text" v-model="form.image_url" placeholder="Contoh: images/product-1.png" style="flex: 1;" />
                    <button type="button" class="btn-upload" @click="triggerFileInput" title="Upload Foto">
                      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                        <circle cx="8.5" cy="8.5" r="1.5"></circle>
                        <polyline points="21 15 16 10 5 21"></polyline>
                      </svg>
                    </button>
                    <input type="file" ref="fileInput" accept="image/*" capture="environment" @change="uploadImage" style="display: none;" />
                  </div>
                </div>
              </div>

              <div class="inner-form-grid">
                <div class="form-group">
                  <label>Durasi Pre-Order (Hari)</label>
                  <input type="number" v-model="form.pre_order_duration" required min="1" />
                </div>
              </div>
            </div>
          </div>
          
          <div class="modal-footer flex justify-end space-x-4" style="display: flex; justify-content: flex-end; gap: 16px;">
            <button type="button" class="btn-cancel bg-gray-100 text-gray-800" style="padding: 10px 20px; border-radius: 4px; background-color: #f3f4f6; color: #1f2937; border: 1px solid #d1d5db;" @click="closeModal">Batal</button>
            <button type="submit" class="btn-primary bg-black text-white" style="padding: 10px 20px; border-radius: 4px; background-color: #000; color: #fff; border: none; font-weight: 500;" :disabled="isSubmitting">
              {{ isSubmitting ? 'Menyimpan...' : 'Simpan Produk' }}
            </button>
          </div>
        </form>
      </div>
    </div>
    <!-- Chat Toast Notification -->
    <div class="chat-toast-container" :class="{ 'toast-visible': showChatToast }">
      <div class="chat-toast-content">
        <div class="chat-toast-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
        </div>
        <div class="chat-toast-text">
          <strong>Pesan Baru</strong>
          <p>{{ toastSenderName }}</p>
        </div>
        <button class="chat-toast-close" @click="showChatToast = false">&times;</button>
      </div>
    </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Line as LineChart } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Filler
} from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale, Filler)
import { useAuthStore } from '../../stores/auth'
import { useSellerStore } from '../../stores/seller'
import { useProductStore } from '../../stores/products'
import PengaturanToko from './PengaturanToko.vue'
import KeuanganView from './KeuanganView.vue'
import SellerChat from './SellerChat.vue'

const router = useRouter()
const authStore = useAuthStore()
const productStore = useProductStore()
const sellerStore = useSellerStore()
import { useNotificationStore } from '../../stores/notification'
import { useChatStore } from '../../stores/chat'
import { useDashboardStore } from '../../stores/dashboard'

const notificationStore = useNotificationStore()
const chatStore = useChatStore()
const dashboardStore = useDashboardStore()

// Bound computed properties for Dashboard Overview
const pendingCount = computed(() => dashboardStore.pendingCount)
const processingCount = computed(() => dashboardStore.processingCount)
const completedCount = computed(() => dashboardStore.completedCount)
const totalRevenue = computed(() => dashboardStore.totalRevenue)
const recentOrders = computed(() => dashboardStore.recentOrders)

const currentTab = ref('products')
watch(currentTab, (newTab) => {
  chatStore.isSellerChatOpen = (newTab === 'chat')
  if (newTab === 'orders') {
    notificationStore.resetSellerNewOrders()
  } else if (newTab === 'chat') {
    notificationStore.resetSellerUnreadChats()
  } else if (newTab === 'dashboard') {
    dashboardStore.fetchInitialStats()
    dashboardStore.fetchChartData(chartFilter.value)
  }
})

const chartFilter = ref('7_days')
watch(chartFilter, (newRange) => {
  dashboardStore.fetchChartData(newRange)
})

const chartDataConfig = computed(() => ({
  labels: dashboardStore.chartLabels,
  datasets: [
    {
      label: 'Pendapatan (Rp)',
      backgroundColor: 'rgba(24, 24, 27, 0.05)',
      borderColor: '#18181b',
      borderWidth: 2,
      pointBackgroundColor: '#ffffff',
      pointBorderColor: '#18181b',
      pointHoverBackgroundColor: '#18181b',
      pointHoverBorderColor: '#ffffff',
      fill: true,
      tension: 0.4,
      data: dashboardStore.chartData
    }
  ]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        label: (context) => {
          let label = context.dataset.label || ''
          if (label) label += ': '
          if (context.parsed.y !== null) {
            label += 'Rp ' + context.parsed.y.toLocaleString('id-ID')
          }
          return label
        }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      grid: { color: '#f3f4f6', drawBorder: false },
      ticks: {
        callback: (value) => {
          if (value === 0) return '0'
          if (value >= 1000000) return (value / 1000000) + 'M'
          if (value >= 1000) return (value / 1000) + 'K'
          return value
        }
      }
    },
    x: {
      grid: { display: false, drawBorder: false }
    }
  }
}

const isModalOpen = ref(false)
const isEditing = ref(false)
const isSubmitting = ref(false)
const isUploading = ref(false)
const editingId = ref(null)
const fileInput = ref(null)

const form = ref({
  name: '',
  price: 0,
  description: '',
  category: 'Koleksi Tenun NTT',
  region: 'Sumba',
  stock: 0,
  image_url: '',
  pre_order_duration: 7
})

// Orders State
const sellerOrders = ref([])
const isLoadingOrders = ref(false)
const isProcessingOrder = ref(null)

const fetchSellerOrders = async () => {
  isLoadingOrders.value = true
  try {
    const response = await axios.get('http://localhost:8081/api/orders/seller', {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    console.log("Seller Orders API Response:", response.data)
    if (response.data.success) {
      sellerOrders.value = response.data.data || []
    }
  } catch (error) {
    console.error("Gagal memuat pesanan:", error)
  } finally {
    isLoadingOrders.value = false
  }
}

const getOrderStatusClass = (status) => {
  switch (status) {
    case 'Menunggu Konfirmasi': return 'bg-yellow-100 text-yellow-800'
    case 'Pesanan Sedang Diproses': return 'bg-blue-100 text-blue-800'
    case 'Dikirim': return 'bg-purple-100 text-purple-800'
    case 'Selesai': return 'bg-gray-100 text-gray-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

const processOrder = async (id) => {
  isProcessingOrder.value = id
  try {
    const response = await axios.put(`http://localhost:8081/api/orders/seller/${id}/process`, {}, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    if (response.data.success) {
      await fetchSellerOrders()
    }
  } catch (error) {
    alert("Gagal memproses pesanan: " + (error.response?.data?.message || error.message))
  } finally {
    isProcessingOrder.value = null
  }
}

const unreadChatCount = ref(0)
const showChatToast = ref(false)
const toastSenderName = ref('')
let toastTimeout = null

const playNotificationSound = () => {
  try {
    const AudioContext = window.AudioContext || window.webkitAudioContext;
    const ctx = new AudioContext();
    const osc = ctx.createOscillator();
    const gainNode = ctx.createGain();
    osc.type = 'sine';
    osc.frequency.setValueAtTime(800, ctx.currentTime);
    osc.frequency.exponentialRampToValueAtTime(1200, ctx.currentTime + 0.1);
    gainNode.gain.setValueAtTime(0.2, ctx.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.1);
    osc.connect(gainNode);
    gainNode.connect(ctx.destination);
    osc.start();
    osc.stop(ctx.currentTime + 0.1);
  } catch (e) {
    console.log("Audio not supported", e);
  }
}

const handleNewMessage = ({ message, senderName }) => {
  toastSenderName.value = senderName
  showChatToast.value = true
  playNotificationSound()
  
  if (toastTimeout) clearTimeout(toastTimeout)
  toastTimeout = setTimeout(() => {
    showChatToast.value = false
  }, 4000)
}

const updateUnreadCount = (count) => {
  unreadChatCount.value = count
}

onMounted(async () => {
  if (!authStore.isAuthenticated || authStore.user?.role !== 'seller') {
    router.push('/')
    return
  }
  await dashboardStore.fetchInitialStats()
  await dashboardStore.fetchChartData(chartFilter.value)
  await sellerStore.fetchProducts()
  await fetchSellerOrders()
})

const resetForm = () => {
  form.value = {
    name: '',
    price: 0,
    description: '',
    category: 'Koleksi Tenun NTT',
    region: 'Sumba',
    stock: 0,
    image_url: '',
    pre_order_duration: 7
  }
  isEditing.value = false
  editingId.value = null
}

const openAddModal = () => {
  resetForm()
  isModalOpen.value = true
}

const openEditModal = (product) => {
  form.value = { ...product }
  isEditing.value = true
  editingId.value = product.id
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
  resetForm()
}

const handleSubmitProduct = async () => {
  isSubmitting.value = true
  const payload = {
    name: form.value.name,
    price: Number(form.value.price),
    description: form.value.description,
    category: form.value.category,
    region: form.value.region,
    stock: Number(form.value.stock),
    image_url: form.value.image_url,
    pre_order_duration: Number(form.value.pre_order_duration)
  }

  if (isEditing.value) {
    await sellerStore.updateProduct(editingId.value, payload)
  } else {
    await sellerStore.createProduct(payload)
  }
  
  // Trigger global data synchronization so Homepage and Koleksi reflect changes instantly
  await productStore.fetchProducts()
  
  isSubmitting.value = false
  closeModal()
}

const triggerFileInput = () => {
  if (fileInput.value) {
    fileInput.value.click()
  }
}

const uploadImage = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  isUploading.value = true
  const formData = new FormData()
  formData.append('image', file)
  
  try {
    const response = await axios.post('http://localhost:8081/api/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    if (response.data.success) {
      form.value.image_url = response.data.url
    } else {
      alert('Gagal mengupload gambar: ' + response.data.message)
    }
  } catch (error) {
    console.error('Upload error:', error)
    alert('Terjadi kesalahan saat mengupload gambar')
  } finally {
    isUploading.value = false
    // reset input so same file can be uploaded again if needed
    event.target.value = ''
  }
}

const handleDelete = async (id) => {
  if (confirm('Yakin ingin menghapus produk ini?')) {
    await sellerStore.deleteProduct(id)
  }
}
</script>

<style scoped>
.seller-layout-container {
  width: 100%;
  display: flex;
  justify-content: center;
  background-color: #f9fafb;
}

.seller-layout {
  display: flex;
  align-items: flex-start;
  min-height: calc(100vh - 80px);
  width: 100%;
  max-width: 80rem; /* max-w-7xl */
}

.seller-sidebar {
  width: 260px;
  background-color: #fff;
  border-right: 1px solid #e5e7eb;
  padding: 32px 0;
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - 80px);
  position: sticky;
  top: 80px;
}

.sidebar-header {
  padding: 0 24px 24px;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 24px;
}

.sidebar-header h2 {
  font-family: 'Playfair Display', serif;
  font-size: 1.25rem;
  margin: 0 0 4px;
  color: #111827;
}

.sidebar-header p {
  font-size: 0.85rem;
  color: #6b7280;
  margin: 0;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 0 16px;
}

.sidebar-nav a {
  padding: 12px 16px;
  border-radius: 6px;
  color: #4b5563;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar-nav a .badge {
  background-color: #ef4444;
  color: white;
  font-size: 0.75rem;
  font-weight: bold;
  padding: 2px 8px;
  border-radius: 12px;
}

.sidebar-nav a:hover {
  background-color: #f3f4f6;
  color: #111827;
}

.sidebar-nav a.active {
  background-color: #1a1a1a;
  color: #fff;
}

.seller-content {
  flex: 1;
  padding: 32px 40px;
  min-height: calc(100vh - 80px);
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.content-header h1 {
  font-family: 'Playfair Display', serif;
  font-size: 2rem;
  margin: 0;
  color: #111827;
}

.btn-dark {
  background-color: #111827;
  color: #ffffff;
  padding: 10px 20px;
  border-radius: 6px;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-dark:hover {
  background-color: #1f2937;
}

/* Dashboard Stats styled via Tailwind CSS directly in template */

.section-title {
  font-family: 'Playfair Display', serif;
  font-size: 1.5rem;
  margin-bottom: 20px;
  color: #111827;
}

.mt-8 {
  margin-top: 32px;
}

.table-container {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.products-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.products-table th, .products-table td {
  padding: 16px 20px;
  border-bottom: 1px solid #e5e7eb;
}

.products-table th {
  background-color: #f9fafb;
  font-weight: 600;
  color: #374151;
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.products-table tbody tr:last-child td {
  border-bottom: none;
}

.products-table tbody tr:hover {
  background-color: #f9fafb;
}

.td-product {
  display: flex;
  align-items: center;
  gap: 16px;
}

.table-img {
  width: 48px;
  height: 48px;
  border-radius: 6px;
  object-fit: cover;
  border: 1px solid #e5e7eb;
}

.table-product-name {
  font-weight: 500;
  color: #111827;
}

.text-right {
  text-align: right;
}

.td-actions {
  text-align: right;
  white-space: nowrap;
}

.td-actions button {
  margin-left: 8px;
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  font-weight: 500;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s;
}

.td-actions .btn-outline {
  background: transparent;
  border: 1px solid #d1d5db;
  color: #374151;
}

.td-actions .btn-outline:hover {
  background-color: #f3f4f6;
}

.td-actions .btn-danger {
  background: transparent;
  color: #ef4444;
}

.td-actions .btn-danger:hover {
  background-color: #fef2f2;
}

.empty-state, .mock-tab, .loading-state {
  text-align: center;
  padding: 64px 20px;
  background-color: #fff;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  color: #6b7280;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: #fff;
  width: 100%;
  max-width: 900px; /* Widened modal as requested */
  border-radius: 8px;
  padding: 32px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
  max-height: 90vh;
  overflow-y: auto;
}

/* Modal Flexbox Styles */
.custom-modal-flex {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  align-items: flex-start;
}

.custom-modal-flex-left {
  width: 100%;
  flex-shrink: 0;
}

.custom-modal-flex-right {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  min-width: 0; /* CRITICAL for preventing flex blowout */
}

.inner-form-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
  width: 100%;
}

.image-upload-flex {
  display: flex;
  gap: 8px;
  align-items: center;
}

.image-upload-flex input[type="text"] {
  width: 100%;
}

.btn-upload {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px;
  background-color: #f3f4f6;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  cursor: pointer;
  color: #374151;
  transition: all 0.2s;
}

.btn-upload:hover {
  background-color: #e5e7eb;
}

.full-width-field {
  width: 100%;
}

@media (min-width: 768px) {
  .custom-modal-flex {
    flex-direction: row;
  }
  .custom-modal-flex-left {
    width: 33.333333%;
  }
  .custom-modal-flex-right {
    width: 66.666667%;
  }
  .inner-form-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.modal-header h2 {
  font-family: 'Playfair Display', serif;
  margin: 0;
  font-size: 1.5rem;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6b7280;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #374151;
}

.form-group input, .form-group select, .form-group textarea {
  width: 100%;
  box-sizing: border-box;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-family: inherit;
  max-width: 100%;
}

.form-group input:focus, .form-group select:focus, .form-group textarea:focus {
  outline: none;
  border-color: #1a1a1a;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.btn-outline {
  padding: 10px 20px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background: none;
  cursor: pointer;
  font-weight: 500;
}

.btn-outline:hover {
  background-color: #f3f4f6;
}

/* Chat Toast Notification */
.chat-toast-container {
  position: fixed;
  bottom: -100px;
  right: 24px;
  transition: bottom 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  z-index: 9999;
}

.chat-toast-container.toast-visible {
  bottom: 24px;
}

.chat-toast-content {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  border-left: 4px solid #111827;
  min-width: 300px;
}

.chat-toast-icon {
  background: #f3f4f6;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #111827;
}

.chat-toast-text {
  flex: 1;
}

.chat-toast-text strong {
  display: block;
  font-size: 0.95rem;
  color: #111827;
  margin-bottom: 2px;
}

.chat-toast-text p {
  margin: 0;
  font-size: 0.85rem;
  color: #6b7280;
}

.chat-toast-close {
  background: none;
  border: none;
  font-size: 1.2rem;
  color: #9ca3af;
  cursor: pointer;
  padding: 4px;
}

.chat-toast-close:hover {
  color: #111827;
}
</style>
