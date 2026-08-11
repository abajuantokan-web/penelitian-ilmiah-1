<template>
  <div class="keuangan-view">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
      <div class="bg-zinc-900 text-white rounded-xl p-6 shadow-md flex flex-col justify-between">
        <h2 class="text-zinc-400 text-sm font-medium tracking-wider mb-2">SALDO AKTIF</h2>
        <p class="text-3xl font-bold mb-6">Rp {{ activeBalance.toLocaleString('id-ID') }}</p>
        <button class="bg-white text-zinc-900 px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-zinc-100 transition w-max" @click="openWithdrawModal">Tarik Saldo</button>
      </div>
      <div class="bg-white border border-zinc-200 text-zinc-800 rounded-xl p-6 shadow-sm flex flex-col justify-between">
        <h2 class="text-zinc-500 text-sm font-medium tracking-wider mb-2">SALDO TERTAHAN (ESCROW)</h2>
        <p class="text-3xl font-bold mb-4">Rp {{ pendingBalance.toLocaleString('id-ID') }}</p>
        <p class="text-xs text-zinc-400 leading-relaxed">Dana akan masuk ke Saldo Aktif setelah pembeli menyelesaikan pesanan.</p>
      </div>
    </div>
    
    <div class="bg-white border border-zinc-200 rounded-xl p-6 shadow-sm overflow-hidden">
      <h3 class="text-xl font-serif font-bold text-zinc-900 mb-6">Riwayat Transaksi</h3>
      <div v-if="transactions.length === 0" class="text-zinc-500">Belum ada transaksi.</div>
      <table v-else class="w-full text-left border-collapse">
        <thead>
          <tr>
            <th class="border-b border-zinc-200 py-3 px-4 text-sm font-semibold text-zinc-700 uppercase tracking-wider">Tanggal</th>
            <th class="border-b border-zinc-200 py-3 px-4 text-sm font-semibold text-zinc-700 uppercase tracking-wider">Tipe</th>
            <th class="border-b border-zinc-200 py-3 px-4 text-sm font-semibold text-zinc-700 uppercase tracking-wider">Jumlah</th>
            <th class="border-b border-zinc-200 py-3 px-4 text-sm font-semibold text-zinc-700 uppercase tracking-wider">Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="tx in transactions" :key="tx.id">
            <td class="border-b border-zinc-100 py-4 px-4 text-sm text-zinc-600">{{ new Date(tx.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) }}</td>
            <td class="border-b border-zinc-100 py-4 px-4 text-sm text-zinc-600">{{ tx.type === 'withdrawal' ? 'Penarikan' : 'Pemasukan' }}</td>
            <td class="border-b border-zinc-100 py-4 px-4 text-sm text-zinc-600">Rp {{ tx.amount.toLocaleString('id-ID') }}</td>
            <td class="border-b border-zinc-100 py-4 px-4 text-sm text-zinc-600">
              <span v-if="tx.status === 'completed'" class="bg-green-100 text-green-700 px-3 py-1 rounded-full text-xs font-semibold">Completed</span>
              <span v-else-if="tx.status === 'failed'" class="bg-red-100 text-red-700 px-3 py-1 rounded-full text-xs font-semibold">Failed</span>
              <span v-else class="bg-orange-100 text-orange-700 px-3 py-1 rounded-full text-xs font-semibold">Processing</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content !max-w-md">
        <div class="modal-header">
          <h2>Tarik Saldo</h2>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        
        <form @submit.prevent="submitWithdrawal" class="modal-form" style="padding: 24px;">
          <div v-if="errorMessage" class="error-alert mb-4">
            {{ errorMessage }}
          </div>
          
          <div class="form-group mb-6">
            <label style="font-size: 0.9rem; font-weight: 500; color: #374151;">Nominal Penarikan (Rp)</label>
            <input type="number" v-model="withdrawAmount" class="w-full form-input mt-1" style="padding: 10px; border: 1px solid #d1d5db; border-radius: 6px; width: 100%; box-sizing: border-box;" required min="50000" />
            <p style="font-size: 0.75rem; color: #6b7280; margin-top: 4px;">Minimal penarikan Rp 50.000</p>
          </div>
          
          <div style="background-color: #fafafa; border: 1px solid #f3f4f6; border-radius: 8px; padding: 16px; margin-bottom: 24px;">
            <p style="font-size: 0.75rem; color: #6b7280; margin-bottom: 4px; margin-top: 0;">Akan dicairkan ke:</p>
            <p style="font-size: 0.875rem; font-weight: 600; color: #18181b; margin: 0;">{{ bankName || 'Belum diatur' }} - {{ bankAccount || '-' }}</p>
            <p style="font-size: 0.75rem; color: #3f3f46; margin-top: 2px; margin-bottom: 0;">a.n {{ bankAccountName || '-' }}</p>
          </div>

          <div style="display: flex; justify-content: flex-end; gap: 12px;">
            <button type="button" style="background-color: #f3f4f6; color: #1f2937; padding: 8px 16px; border-radius: 6px; border: 1px solid #e5e7eb; cursor: pointer;" @click="closeModal">Batal</button>
            <button type="submit" style="background-color: #000; color: #fff; padding: 8px 16px; border-radius: 6px; border: none; font-weight: 500; cursor: pointer;" :disabled="isSubmitting">
              {{ isSubmitting ? 'Memproses...' : 'Tarik Sekarang' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()

const activeBalance = ref(0)
const pendingBalance = ref(0)
const bankName = ref('')
const bankAccount = ref('')
const bankAccountName = ref('')
const transactions = ref([])

const isModalOpen = ref(false)
const withdrawAmount = ref(0)
const isSubmitting = ref(false)
const errorMessage = ref('')

const fetchWallet = async () => {
  try {
    const response = await axios.get('http://localhost:8081/api/seller/wallet', {
      headers: { Authorization: `Bearer ${authStore.token || localStorage.getItem('token')}` }
    })
    if (response.data.success) {
      activeBalance.value = response.data.data.active_balance
      pendingBalance.value = response.data.data.pending_balance
      bankName.value = response.data.data.bank_name
      bankAccount.value = response.data.data.bank_account_number
      bankAccountName.value = response.data.data.bank_account_name
      transactions.value = response.data.data.transactions || []
    }
  } catch (error) {
    console.error('Failed to fetch wallet:', error)
  }
}

onMounted(() => {
  fetchWallet()
})

const openWithdrawModal = () => {
  errorMessage.value = ''
  if (!bankAccount.value || !bankName.value) {
    errorMessage.value = "Silakan atur rekening bank di Pengaturan Toko terlebih dahulu."
  }
  withdrawAmount.value = activeBalance.value > 0 ? activeBalance.value : 0
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const submitWithdrawal = async () => {
  isSubmitting.value = true
  errorMessage.value = ''
  try {
    const response = await axios.post('http://localhost:8081/api/seller/wallet/withdraw', {
      amount: Number(withdrawAmount.value)
    }, {
      headers: { Authorization: `Bearer ${authStore.token || localStorage.getItem('token')}` }
    })
    if (response.data.success) {
      alert(response.data.message)
      closeModal()
      fetchWallet() 
    }
  } catch (error) {
    errorMessage.value = error.response?.data?.message || 'Gagal melakukan penarikan.'
  } finally {
    isSubmitting.value = false
  }
}

const getTxStatusClass = (status) => {
  switch (status) {
    case 'processing': return 'bg-orange-100 text-orange-800'
    case 'completed': return 'bg-green-100 text-green-800'
    case 'failed': return 'bg-red-100 text-red-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}
</script>

<style scoped>
.keuangan-view {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
  border: 1px solid #e5e7eb;
}
.balance-amount {
  font-size: 2.5rem;
  font-weight: 700;
  color: #111827;
  margin: 8px 0;
}
.error-alert {
  background-color: #fee2e2;
  color: #991b1b;
  padding: 12px;
  border-radius: 6px;
  font-size: 0.9rem;
}
</style>
