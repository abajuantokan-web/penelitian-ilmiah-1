import axios from 'axios'

// =============================================================
//  URL Backend Production (Railway)
//  Ubah jika alamat backend berubah.
// =============================================================
export const BASE_URL = 'https://penelitian-ilmiah-1-production.up.railway.app/'

// Instance axios dengan baseURL
const apiClient = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Sisipkan token admin secara otomatis
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('admin_token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

export default apiClient
