<template>
  <div class="bg-gray-50 min-h-screen pt-24 pb-12">
    <div class="max-w-4xl mx-auto px-4">
      <h1 class="text-2xl font-bold mb-6 text-gray-800">Checkout</h1>

      <div v-if="checkoutItems.length === 0" class="bg-white p-6 rounded shadow text-center">
        <p class="text-gray-500 mb-4">Produk tidak ditemukan atau sesi telah berakhir.</p>
        <button @click="router.push('/')" class="btn-primary--solid-dark px-4 py-2 rounded">Kembali ke Beranda</button>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- Main Content -->
        <div class="md:col-span-2 space-y-6">
          
          <!-- Shipping Address -->
          <div class="bg-white p-6 rounded shadow">
            <div class="flex items-center justify-between mb-4 border-b pb-2">
              <h2 class="text-lg font-semibold flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-2 text-red-500"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
                Alamat Pengiriman
              </h2>
              <button @click="openEditModal" class="text-sm text-blue-600 font-medium hover:text-blue-800">Ubah</button>
            </div>
            
            <div class="text-sm text-gray-700">
              <p class="font-bold">{{ shippingName }} <span class="font-normal text-gray-500 ml-2">{{ shippingPhone }}</span></p>
              <p class="mt-1">{{ shippingAddress }}</p>
            </div>
          </div>

          <!-- Product Details -->
          <div class="bg-white p-6 rounded shadow">
            <h2 class="text-lg font-semibold mb-4 border-b pb-2">Pesanan Anda</h2>
            <div v-for="item in checkoutItems" :key="item.product.id" class="flex items-start gap-4 mb-4">
              <img :src="getImageUrl(item.product.image_url || item.product.image)" :alt="item.product.name" class="w-24 h-24 object-cover rounded border">
              <div class="flex-1">
                <h3 class="font-medium text-gray-800">{{ item.product.name }}</h3>
                <p class="text-sm text-gray-500 mb-2">Variant: {{ item.product.category }}</p>
                <div class="flex justify-between items-center mt-4">
                  <span class="text-sm text-gray-700">{{ formatPrice(item.product.price) }}</span>
                  <span class="text-sm text-gray-700">x{{ item.qty }}</span>
                </div>
              </div>
            </div>

            <!-- Optional Note -->
            <div class="mt-6 border-t pt-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">Pesan: (Opsional) Tinggalkan pesan</label>
              <input type="text" v-model="customerNote" placeholder="Contoh: Tolong sesuaikan ukuran lingkar dada 100cm" class="w-full border-gray-300 border rounded-md shadow-sm p-2 focus:ring-black focus:border-black text-sm">
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Shipping Option -->
          <div class="bg-white p-6 rounded shadow">
            <h2 class="text-lg font-semibold mb-4 border-b pb-2">Opsi Pengiriman</h2>
            <div class="flex justify-between items-center mb-2">
              <span class="text-sm font-medium text-gray-800">Reguler (JNE/J&T)</span>
              <span class="text-sm font-medium text-gray-800">Rp 0</span>
            </div>
            <p class="text-xs text-gray-500">Estimasi Tiba: {{ estimatedArrival }}</p>
            <p v-if="maxPoDuration > 0" class="text-xs text-orange-500 mt-1">*Termasuk masa Pre-Order maksimum {{ maxPoDuration }} hari</p>
          </div>

          <!-- Order Summary -->
          <div class="bg-white p-6 rounded shadow">
            <h2 class="text-lg font-semibold mb-4 border-b pb-2">Ringkasan Belanja</h2>
            
            <div class="space-y-3 mb-4">
              <div class="flex justify-between text-sm">
                <span class="text-gray-600">Total Harga ({{ totalQty }} barang)</span>
                <span class="text-gray-800">{{ formatPrice(totalAmount) }}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-gray-600">Total Ongkos Kirim</span>
                <span class="text-gray-800">Rp 0</span>
              </div>
            </div>
            
            <div class="border-t pt-4 flex justify-between items-center mb-6">
              <span class="font-bold text-gray-800">Total Tagihan</span>
              <span class="font-bold text-red-600 text-lg">{{ formatPrice(totalAmount) }}</span>
            </div>

            <button 
              @click="handleBuatPesanan" 
              :disabled="isProcessing"
              class="w-full bg-black text-white font-medium py-3 rounded hover:bg-gray-800 transition-colors disabled:opacity-50"
            >
              <span v-if="isProcessing">Memproses...</span>
              <span v-else>Buat Pesanan</span>
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Edit Address Modal -->
    <div v-if="isEditingAddress" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b flex justify-between items-center">
          <h3 class="text-lg font-semibold text-gray-800">Ubah Alamat Pengiriman</h3>
          <button @click="cancelEdit" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Penerima</label>
            <input type="text" v-model="editForm.name" class="w-full border-gray-300 border rounded-md p-2 focus:ring-black focus:border-black text-sm">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nomor Telepon</label>
            <input type="tel" v-model="editForm.phone" class="w-full border-gray-300 border rounded-md p-2 focus:ring-black focus:border-black text-sm">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Detail Alamat</label>
            <textarea v-model="editForm.address" rows="3" class="w-full border-gray-300 border rounded-md p-2 focus:ring-black focus:border-black text-sm"></textarea>
          </div>
        </div>
        <div class="px-6 py-4 bg-gray-50 border-t flex justify-end gap-3">
          <button @click="cancelEdit" class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-100">Batal</button>
          <button @click="saveEdit" class="px-4 py-2 bg-black text-white rounded-md text-sm font-medium hover:bg-gray-800">Simpan</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import { getImageUrl } from '../utils/imageUtils'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const cartStore = useCartStore()

const isFromCart = computed(() => route.query.source === 'cart')
const checkoutItems = ref([])
const customerNote = ref('')
const isProcessing = ref(false)

// Shipping Address State
const shippingName = ref('Customer')
const shippingPhone = ref('081234567890')
const shippingAddress = ref('Jl. Piet A. Tallo, Liliba, Kec. Oebobo, Kota Kupang, Nusa Tenggara Timur 85111')

// Edit Form State
const isEditingAddress = ref(false)
const editForm = reactive({
  name: '',
  phone: '',
  address: ''
})

const openEditModal = () => {
  // Populate form with current state values
  editForm.name = shippingName.value
  editForm.phone = shippingPhone.value
  editForm.address = shippingAddress.value
  isEditingAddress.value = true
}

const saveEdit = () => {
  // Update state with new input values
  shippingName.value = editForm.name
  shippingPhone.value = editForm.phone
  shippingAddress.value = editForm.address
  isEditingAddress.value = false
}

const cancelEdit = () => {
  // Close without saving
  isEditingAddress.value = false
}

onMounted(async () => {
  if (isFromCart.value) {
    if (cartStore.items.length === 0) {
      await cartStore.fetchCart()
    }
    if (cartStore.items.length === 0) {
      alert('Keranjang kosong.')
      router.push('/')
      return
    }
    checkoutItems.value = cartStore.items.map(item => ({
      product: item.product,
      qty: item.quantity
    }))
  } else {
    if (history.state && history.state.product) {
      try {
        const prod = JSON.parse(history.state.product)
        const q = history.state.qty || 1
        checkoutItems.value = [{ product: prod, qty: q }]
      } catch (e) {
        console.error("Gagal memuat data produk:", e)
      }
    }
  }
  
  // Set fallback defaults from user store if available on mount
  if (authStore.user) {
    shippingName.value = authStore.user.name || 'Customer'
    shippingPhone.value = authStore.user.phone || '081234567890'
  }
})

const totalAmount = computed(() => {
  return checkoutItems.value.reduce((sum, item) => sum + (item.product.price * item.qty), 0)
})

const totalQty = computed(() => {
  return checkoutItems.value.reduce((sum, item) => sum + item.qty, 0)
})

const formatPrice = (price) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(price)
}

const maxPoDuration = computed(() => {
  if (checkoutItems.value.length === 0) return 0
  return checkoutItems.value.reduce((max, item) => {
    return Math.max(max, parseInt(item.product.po_duration || 0))
  }, 0)
})

const estimatedArrival = computed(() => {
  if (checkoutItems.value.length === 0) return ''
  const poDays = maxPoDuration.value
  // Tambah 3 hari estimasi pengiriman setelah PO
  const totalDays = parseInt(poDays) + 3 
  const targetDate = new Date()
  targetDate.setDate(targetDate.getDate() + totalDays)
  
  return targetDate.toLocaleDateString('id-ID', { 
    weekday: 'long', 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
})

const handleBuatPesanan = async () => {
  if (!authStore.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: '/checkout' } })
    return
  }
  
  isProcessing.value = true
  
  try {
    let checkoutRes;
    
    // Create the DB record first based on Cart vs Direct Buy
    if (isFromCart.value) {
      const payloadItems = checkoutItems.value.map(item => ({
        product_id: item.product.id,
        quantity: item.qty,
        price: item.product.price
      }));

      checkoutRes = await axios.post('http://localhost:8081/api/orders', {
        items: payloadItems,
        note: customerNote.value
      }, {
        headers: { Authorization: `Bearer ${authStore.token}` }
      });
      // Clear local cart since it's converted to an order
      if (checkoutRes.data.success) {
        cartStore.clearLocalCart()
      }
    } else {
      checkoutRes = await axios.post('http://localhost:8081/api/orders/direct', {
        product_id: checkoutItems.value[0].product.id,
        quantity: checkoutItems.value[0].qty,
        price: checkoutItems.value[0].product.price,
        note: customerNote.value
      }, {
        headers: { Authorization: `Bearer ${authStore.token}` }
      });
    }

    let snap_token = checkoutRes.data.snap_token;
    
    // Fallback: if backend hasn't been restarted, it won't return snap_token
    // so we fetch it manually via /api/checkout
    if (!snap_token && checkoutRes.data.success) {
      const fallbackRes = await axios.post('http://localhost:8081/api/checkout', {
        total_amount: totalAmount.value,
        first_name: shippingName.value,
        last_name: '',
        email: authStore.user?.email || 'customer@example.com',
        phone: shippingPhone.value,
        note: customerNote.value
      });
      snap_token = fallbackRes.data.snap_token;
    }

    if (checkoutRes.data.success && snap_token) {
      window.snap.pay(snap_token, {
        onSuccess: async function (result) {
          console.log('✅ PAYMENT SUCCESS:', result)
          try {
            const refId = checkoutRes.data.payment_reference || checkoutRes.data.data.id;
            await axios.put(`http://localhost:8081/api/orders/${refId}/confirm-payment`, {}, {
              headers: { Authorization: `Bearer ${authStore.token}` }
            })
          } catch (e) {
            console.error("Failed to confirm payment", e)
          }
          router.push({ path: '/profile', query: { tab: 'orders' } })
        },
        onPending: async function (result) {
          console.log('⏳ PAYMENT PENDING:', result)
          try {
            const refId = checkoutRes.data.payment_reference || checkoutRes.data.data.id;
            await axios.put(`http://localhost:8081/api/orders/${refId}/confirm-payment`, {}, {
              headers: { Authorization: `Bearer ${authStore.token}` }
            })
          } catch (e) {
            console.error("Failed to confirm payment", e)
          }
          router.push({ path: '/profile', query: { tab: 'orders' } })
        },
        onError: function (result) {
          console.error('❌ PAYMENT FAILED:', result)
          isProcessing.value = false
        },
        onClose: function () {
          console.log('⚠️ PAYMENT POPUP CLOSED')
          isProcessing.value = false
        }
      })
    } else {
      isProcessing.value = false
    }
  } catch (error) {
    console.error("AXIOS ERROR:", error.response?.data || error.message);
    alert("Error: " + (error.response?.data?.message || error.message));
    isProcessing.value = false
  }
}
</script>

<style scoped>
/* Tailwind classes handled via CDN. Any additional overrides can go here. */
</style>
