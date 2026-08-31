import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'


export default defineConfig({
  plugins: [
    vue(),
  ],
  server: {
    host: true,
    port: 5154,
    strictPort: true,
    hmr: {
      host: "127.0.0.1",
    },
    proxy: {
      "/api": {
        target: "https://penelitian-ilmiah-1-production.up.railway.app",
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
