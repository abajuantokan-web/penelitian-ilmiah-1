<template>
  <div class="pengaturan-toko">
    <form @submit.prevent="saveSettings" class="settings-form">
      
      
      <section class="settings-card fade-in">
        <h2 class="card-title">Profil Toko</h2>
        <p class="card-subtitle">Informasi dasar mengenai toko Anda yang akan dilihat oleh pembeli.</p>
        
        <div class="form-group mb-6">
          <label>Logo Toko</label>
          <div class="image-upload-flex">
            
            <input type="file" ref="logoInput" accept="image/*" class="hidden" @change="uploadLogo" style="display: none;" />
            
            <div class="logo-preview-container">
              <img :src="$getImageUrl(form.store_logo)" alt="Logo Toko" class="logo-preview" />
            </div>
            
            <button type="button" class="btn-upload" @click="triggerLogoUpload" title="Upload Logo">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                <circle cx="8.5" cy="8.5" r="1.5"></circle>
                <polyline points="21 15 16 10 5 21"></polyline>
              </svg>
              <span>{{ isUploading ? 'Mengupload...' : 'Ubah Logo' }}</span>
            </button>
          </div>
        </div>

        <div class="form-group">
          <label>Nama Toko</label>
          <input type="text" v-model="form.store_name" required placeholder="Masukkan nama toko Anda" class="w-full form-input" />
        </div>

        <div class="form-group mt-4">
          <label>Deskripsi Toko</label>
          <textarea v-model="form.description" rows="4" placeholder="Ceritakan keunikan dan produk unggulan toko Anda..." class="w-full form-input"></textarea>
        </div>
      </section>

      
      <section class="settings-card fade-in" style="animation-delay: 0.1s;">
        <h2 class="card-title">Kontak & Lokasi</h2>
        <p class="card-subtitle">Detail kontak dan alamat pengiriman asal produk Anda.</p>
        
        <div class="grid-2-cols mt-4">
          <div class="form-group">
            <label>Nomor Telepon / WhatsApp</label>
            <input type="tel" v-model="form.phone" placeholder="Contoh: 08123456789" class="w-full form-input" />
          </div>
          
          <div class="form-group">
            <label>Region Asal</label>
            <select v-model="form.region" class="w-full form-input" required>
              <option value="" disabled>Pilih Region</option>
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

        <div class="form-group mt-4">
          <label>Alamat Lengkap</label>
          <textarea v-model="form.address" rows="3" placeholder="Alamat lengkap toko atau workshop Anda..." class="w-full form-input"></textarea>
        </div>
      </section>

      
      <section class="settings-card fade-in" style="animation-delay: 0.2s;">
        <h2 class="card-title">Informasi Rekening Bank</h2>
        <p class="card-subtitle">Rekening yang akan digunakan untuk pencairan dana hasil penjualan.</p>
        
        <div class="grid-2-cols mt-4">
          <div class="form-group">
            <label>Nama Bank</label>
            <input type="text" v-model="form.bank_name" placeholder="Contoh: BCA, Mandiri, BNI" class="w-full form-input" required />
          </div>
          
          <div class="form-group">
            <label>Nomor Rekening</label>
            <input type="text" v-model="form.bank_account_number" placeholder="Contoh: 1234567890" class="w-full form-input" required />
          </div>
        </div>

        <div class="form-group mt-4">
          <label>Nama Pemilik Rekening</label>
          <input type="text" v-model="form.bank_account_name" placeholder="Sesuai dengan nama di buku tabungan" class="w-full form-input" required />
        </div>
      </section>

      
      <div class="settings-action">
        <span v-if="saveMessage" :class="['save-message', saveStatus]">{{ saveMessage }}</span>
        <button type="submit" class="btn-primary--solid-dark save-btn" :disabled="isSaving">
          {{ isSaving ? 'Menyimpan...' : 'Simpan Pengaturan' }}
        </button>
      </div>

    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from '../../axios'
import { useAuthStore } from '../../stores/auth'
import { useProductStore } from '../../stores/products'

const authStore = useAuthStore()
const productStore = useProductStore()

const logoInput = ref(null)
const isUploading = ref(false)
const isSaving = ref(false)
const saveMessage = ref('')
const saveStatus = ref('')

const form = ref({
  store_name: '',
  description: '',
  store_logo: '',
  phone: '',
  address: '',
  region: '',
  bank_name: '',
  bank_account_number: '',
  bank_account_name: ''
})

onMounted(async () => {
  
  const sp = authStore.user?.seller_profile
  if (sp?.store_name) {
    form.value.store_name = sp.store_name
  } else if (authStore.user?.store_name) {
    form.value.store_name = authStore.user.store_name
  }
  await fetchProfile()
})

const fetchProfile = async () => {
  try {
    const response = await axios.get(`/api/seller/profile`)
    
    if (response.data && response.data.id) {
      form.value = {
        store_name: response.data.store_name || form.value.store_name,
        description: response.data.description || '',
        store_logo: response.data.store_logo || '',
        phone: response.data.phone || '',
        address: response.data.address || '',
        region: response.data.region || '',
        bank_name: response.data.bank_name || '',
        bank_account_number: response.data.bank_account_number || '',
        bank_account_name: response.data.bank_account_name || ''
      }
    }
  } catch (error) {
    console.error('Failed to fetch seller profile:', error)
  }
}

const triggerLogoUpload = () => {
  if (logoInput.value) {
    logoInput.value.click()
  }
}

const uploadLogo = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  isUploading.value = true
  const formData = new FormData()
  formData.append('image', file)
  
  try {
    const response = await axios.post(`/api/upload`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    
    if (response.data.success) {
      form.value.store_logo = response.data.url
    } else {
      alert('Gagal mengupload logo: ' + response.data.message)
    }
  } catch (error) {
    console.error('Upload error:', error)
    alert('Terjadi kesalahan saat mengupload logo')
  } finally {
    isUploading.value = false
    event.target.value = ''
  }
}

const saveSettings = async () => {
  isSaving.value = true
  saveMessage.value = ''
  
  try {
    await axios.put(`/api/seller/profile`, form.value)
    
    authStore.updateSellerProfileLocally(form.value)
    
    await productStore.fetchProducts()
    
    saveMessage.value = 'Pengaturan berhasil disimpan!'
    saveStatus.value = 'success'
  } catch (error) {
    console.error('Failed to save settings:', error)
    saveMessage.value = 'Gagal menyimpan pengaturan.'
    saveStatus.value = 'error'
  } finally {
    isSaving.value = false
    setTimeout(() => {
      saveMessage.value = ''
    }, 3000)
  }
}
</script>

<style scoped>
.pengaturan-toko {
  max-width: 800px;
  margin: 0 auto;
}

.settings-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 32px;
  margin-bottom: 24px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
}

.card-title {
  font-family: 'Playfair Display', serif;
  font-size: 1.5rem;
  color: #111827;
  margin: 0 0 8px 0;
}

.card-subtitle {
  color: #6b7280;
  font-size: 0.95rem;
  margin: 0 0 24px 0;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 0.9rem;
  font-weight: 600;
  color: #374151;
}

.form-input {
  width: 100%;
  box-sizing: border-box;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-family: inherit;
  font-size: 0.95rem;
  transition: border-color 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #1a1a1a;
}

.grid-2-cols {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

@media (min-width: 768px) {
  .grid-2-cols {
    grid-template-columns: 1fr 1fr;
  }
}

.image-upload-flex {
  display: flex;
  align-items: center;
  gap: 20px;
}

.logo-preview-container {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 1px solid #e5e7eb;
  overflow: hidden;
  background-color: #f9fafb;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.btn-upload {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background-color: #fff;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  cursor: pointer;
  color: #374151;
  font-weight: 500;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.btn-upload:hover {
  background-color: #f9fafb;
}

.settings-action {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 32px;
  padding-bottom: 40px;
}

.save-btn {
  padding: 12px 32px;
  font-size: 1rem;
  border-radius: 6px;
  cursor: pointer;
  background-color: #111827;
  color: #fff;
  border: none;
  font-weight: 600;
  transition: opacity 0.2s;
}

.save-btn:hover {
  opacity: 0.9;
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.save-message {
  font-weight: 500;
  font-size: 0.95rem;
}

.save-message.success {
  color: #10b981; 
}

.save-message.error {
  color: #ef4444; 
}

.mb-6 { margin-bottom: 24px; }
.mt-4 { margin-top: 16px; }

.fade-in {
  animation: fadeIn 0.4s ease forwards;
  opacity: 0;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
