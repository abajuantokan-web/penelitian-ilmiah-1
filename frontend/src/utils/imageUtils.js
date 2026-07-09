export const getImageUrl = (path) => {
  if (!path) return '/images/placeholder.jpg'
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  
  const cleanPath = path.startsWith('/') ? path.slice(1) : path
  // In a production app, use VITE_API_BASE_URL from env, but for now fallback to localhost:8081
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081'
  return `${baseUrl.replace(/\/$/, '')}/${cleanPath}`
}
