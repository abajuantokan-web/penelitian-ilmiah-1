/**
 * Resolves a product image_url to a fully-qualified URL.
 *
 * Routing strategy:
 *   1. Absolute URL (http/https)  → pass through as-is (Unsplash, CDN, etc.)
 *   2. Uploaded file path          → backend server (e.g. images/upload_*.png)
 *      These are files POSTed to /api/upload and stored in backend/images/.
 *   3. Static product image path  → Vite frontend server (frontend/public/images/)
 *      e.g. "images/kalung-timor.jpg" → "/images/kalung-timor.jpg"
 *   4. Empty / null               → inline SVG placeholder (never broken icon)
 */
export const getImageUrl = (path) => {
  // ── Case 4: empty ────────────────────────────────────────────────────────────
  if (!path) return _svgPlaceholder()

  // ── Case 1: already an absolute URL ─────────────────────────────────────────
  if (path.startsWith('http://') || path.startsWith('https://')) return path

  const cleanPath = path.startsWith('/') ? path.slice(1) : path

  // ── Case 2: user-uploaded file (stored in backend /images/) ─────────────────
  // These files are created by /api/upload and have the pattern upload_*.ext
  if (/^images\/upload_/.test(cleanPath)) {
    const backendBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081'
    return `${backendBase.replace(/\/$/, '')}/${cleanPath}`
  }

  // ── Case 3: static product image in frontend/public/images/ ─────────────────
  // Vite serves frontend/public/ at root '/', so /images/foo.png loads directly.
  return `/${cleanPath}`
}

/** Returns a lightweight inline SVG so broken-icon never appears */
const _svgPlaceholder = () =>
  "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='80' height='80' viewBox='0 0 80 80'%3E%3Crect width='80' height='80' fill='%23f3f0eb'/%3E%3Ctext x='40' y='44' text-anchor='middle' font-family='sans-serif' font-size='10' fill='%23a8956e'%3EOpenPeo%3C/text%3E%3C/svg%3E"

