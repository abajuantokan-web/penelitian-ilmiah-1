import axios from 'axios'

// =============================================================
//  GANTI URL INI SESUAI DENGAN ALAMAT BACKEND YANG DIPAKAI
//  Contoh production : 'https://api.namadomain.com'
//  Contoh lokal      : 'http://localhost:8081'
// =============================================================
export const BASE_URL = 'https://penelitian-ilmiah-1-production.up.railway.app/'

// Instance axios dengan baseURL sudah terset otomatis
const apiClient = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Interceptor: sisipkan token Authorization dari localStorage secara otomatis
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('openpeo_token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

export default apiClient
