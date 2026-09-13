<template>
  <section class="section-divider">
    <div class="section-divider__image">
      <!--
        Performance: serve modern formats (AVIF/WebP) for the two large divider
        images (divider1.png ~888 KB, divider2.png ~765 KB).
        loading="lazy" is intentional — these images are below the fold.
        vite-plugin-image-optimizer generates the .avif and .webp variants at build time.
      -->
      <picture>
        <source :srcset="avifSrc" type="image/avif">
        <source :srcset="webpSrc" type="image/webp">
        <img :src="src" :alt="alt" loading="lazy" width="1920" height="1080">
      </picture>
    </div>
    <div class="section-divider__overlay"></div>
    <div class="section-divider__content container reveal">
      <h2 class="section-divider__title">{{ title }}</h2>
      <p class="section-divider__subtitle">{{ subtitle }}</p>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  src: String,
  alt: String,
  title: String,
  subtitle: String
})

// Derive WebP and AVIF paths from the original src prop.
// e.g. "/images/divider1.png" → "/images/divider1.webp" / "/images/divider1.avif"
// This works for any image format (.png, .jpg, .jpeg) without hardcoding paths.
const webpSrc = computed(() => props.src?.replace(/\.(png|jpe?g)$/i, '.webp') ?? '')
const avifSrc = computed(() => props.src?.replace(/\.(png|jpe?g)$/i, '.avif') ?? '')
</script>

