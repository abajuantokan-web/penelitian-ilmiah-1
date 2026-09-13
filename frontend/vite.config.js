import { fileURLToPath, URL } from 'node:url'
import { readdirSync, statSync } from 'node:fs'
import { join, extname, basename } from 'node:path'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { ViteImageOptimizer } from 'vite-plugin-image-optimizer'
import tailwindcss from '@tailwindcss/vite'

/**
 * Custom Vite plugin: generates .webp and .avif sibling files from every
 * PNG/JPG in the dist output directory after the bundle is written.
 *
 * This is distinct from ViteImageOptimizer which only optimises *existing*
 * formats in-place. We need to *create new files* for the <picture> srcset
 * sources used in Hero.vue and ImageDivider.vue.
 *
 * The plugin runs as a `closeBundle` async hook so it runs after all assets
 * are already written to `dist/`.
 */
function viteImageConverter() {
  return {
    name: 'vite-image-converter',
    apply: 'build',  // only runs during `vite build`, never in dev server
    async closeBundle() {
      // Dynamic import of sharp so vite dev server startup is never slowed
      const sharp = (await import('sharp')).default

      const imagesDir = join(process.cwd(), 'dist', 'images')

      // Guard: images dir may not exist on first run
      let files
      try {
        files = readdirSync(imagesDir)
      } catch {
        return
      }

      const targetFormats = [
        { ext: '.webp', options: { quality: 80, effort: 4 } },
        { ext: '.avif', options: { quality: 65, effort: 4 } },
      ]

      const tasks = []

      for (const file of files) {
        const srcExt = extname(file).toLowerCase()
        if (!['.png', '.jpg', '.jpeg'].includes(srcExt)) continue

        const srcPath = join(imagesDir, file)
        const nameWithoutExt = basename(file, srcExt)

        // Skip empty files (e.g., the corrupt kalung-timor.jpg)
        try {
          if (statSync(srcPath).size === 0) continue
        } catch {
          continue
        }

        for (const { ext, options } of targetFormats) {
          // Standard full-size generation
          const destPath = join(imagesDir, nameWithoutExt + ext)
          const task = sharp(srcPath)
            [ext === '.webp' ? 'webp' : 'avif'](options)
            .toFile(destPath)
            .then(() => {
              const orig = statSync(srcPath).size
              const dest = statSync(destPath).size
              const saving = (((orig - dest) / orig) * 100).toFixed(0)
              console.log(
                `  [image-converter] ${file} → ${basename(destPath)}` +
                ` (${(dest / 1024).toFixed(0)} KB, -${saving}%)`
              )
            })
            .catch((err) => {
              console.warn(`  [image-converter] skipped ${file}: ${err.message}`)
            })
          tasks.push(task)

          // Responsive sizes generation exclusively for hero image
          if (nameWithoutExt === 'hero') {
            const sizes = [480, 800, 1200];
            for (const width of sizes) {
              const responsiveDestPath = join(imagesDir, `${nameWithoutExt}-${width}${ext}`)
              const resizeTask = sharp(srcPath)
                .resize({ width, withoutEnlargement: true })
                [ext === '.webp' ? 'webp' : 'avif'](options)
                .toFile(responsiveDestPath)
                .then(() => {
                  const dest = statSync(responsiveDestPath).size
                  console.log(
                    `  [image-converter] ${file} → ${basename(responsiveDestPath)}` +
                    ` (${(dest / 1024).toFixed(0)} KB) [Responsive]`
                  )
                })
                .catch((err) => {
                  console.warn(`  [image-converter] skipped responsive ${file} (${width}): ${err.message}`)
                })
              tasks.push(resizeTask)
            }
          }
        }
      }

      if (tasks.length > 0) {
        console.log(`\n✨ [vite-image-converter] generating WebP + AVIF for ${tasks.length / 2} images...\n`)
        await Promise.all(tasks)
        console.log('\n✅ [vite-image-converter] done\n')
      }
    },
  }
}

export default defineConfig(({ mode }) => ({
  clearScreen: false,
  plugins: [
    vue(),
    tailwindcss(),
    // Only load devtools in development — keeps production bundle lean
    mode === 'development' && vueDevTools(),

    // Step 1 (build): Generate .webp and .avif siblings from PNG/JPG in dist/images/
    // These are the files served by <picture> srcset in Hero.vue + ImageDivider.vue
    viteImageConverter(),

    // Step 2 (build): Losslessly compress the original PNG/JPG fallback files
    ViteImageOptimizer({
      png:  { quality: 80 },
      jpg:  { quality: 80 },
      jpeg: { quality: 80 },
      // webp/avif entries optimise *existing* .webp/.avif files (our generated ones)
      webp: { quality: 80, effort: 4 },
      avif: { quality: 65, effort: 4 },
      logStats: true,
    }),
  ].filter(Boolean),

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

  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },

  build: {
    chunkSizeWarningLimit: 500,
    rollupOptions: {
      output: {
        // Vendor-split: isolate vue/pinia/vue-router into a separately cached chunk.
        // Vite 8 (rolldown) requires manualChunks to be a function, not a plain object.
        manualChunks(id) {
          if (id.includes('node_modules/vue/') ||
              id.includes('node_modules/pinia/') ||
              id.includes('node_modules/vue-router/')) {
            return 'vendor-vue'
          }
        },
      },
    },
  },
}))

