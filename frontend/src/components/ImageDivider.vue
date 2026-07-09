<template>
  <section class="section-divider">
    <div class="section-divider__image">
      <img :src="src" :alt="alt" loading="lazy">
    </div>
    <div class="section-divider__overlay"></div>
    <div class="section-divider__content container reveal">
      <h2 class="section-divider__title">{{ title }}</h2>
      <p class="section-divider__subtitle">{{ subtitle }}</p>
    </div>
  </section>
</template>

<script setup>
defineProps({
  src: String,
  alt: String,
  title: String,
  subtitle: String
})
</script>
