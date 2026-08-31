










export const getImageUrl = (path) => {

  if (!path) return _svgPlaceholder()


  if (path.startsWith('http://') || path.startsWith('https://')) return path

  const cleanPath = path.startsWith('/') ? path.slice(1) : path



  if (/^images\/upload_/.test(cleanPath)) {
    const backendBase = import.meta.env.VITE_API_BASE_URL || 'https://penelitian-ilmiah-1-production.up.railway.app/'
    return `${backendBase.replace(/\/$/, '')}/${cleanPath}`
  }



  return `/${cleanPath}`
}


const _svgPlaceholder = () =>
  "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='80' height='80' viewBox='0 0 80 80'%3E%3Crect width='80' height='80' fill='%23f3f0eb'/%3E%3Ctext x='40' y='44' text-anchor='middle' font-family='sans-serif' font-size='10' fill='%23a8956e'%3EOpenPeo%3C/text%3E%3C/svg%3E"

