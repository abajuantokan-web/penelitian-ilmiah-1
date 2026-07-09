<template>
  <div class="premium-placeholder">
    <div class="placeholder-content">
      <h3 class="placeholder-title">{{ name }}</h3>
      <div class="placeholder-brand">OpenPeo</div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  name: {
    type: String,
    required: true
  }
})
</script>

<style scoped>
.premium-placeholder {
  width: 100%;
  aspect-ratio: 3/4;
  background-color: #060606;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.placeholder-content {
  padding: 20px;
  text-align: center;
  position: relative;
  z-index: 2;
}

.placeholder-title {
  color: #FFFFFF;
  font-family: 'Montserrat', sans-serif;
  font-weight: 300;
  font-size: 1.2rem;
  letter-spacing: 1px;
  line-height: 1.4;
  margin: 0;
}

.placeholder-brand {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  color: rgba(255, 255, 255, 0.2);
  font-family: 'Montserrat', sans-serif;
  font-weight: 600;
  font-size: 0.8rem;
  letter-spacing: 2px;
  text-transform: uppercase;
  z-index: 1;
  width: 100%;
  text-align: center;
}
</style>
