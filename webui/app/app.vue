<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
useHead({
  title: 'Trading Card Album',
  meta: [
    { name: 'description', content: 'Colecciona, intercambia y comparte tarjetas coleccionables en el DevFest 2026.' },
    { name: 'theme-color', content: '#0b0d11' },
    { name: 'apple-mobile-web-app-capable', content: 'yes' },
    { name: 'apple-mobile-web-app-status-bar-style', content: 'black-translucent' },
  ],
  htmlAttrs: { lang: 'es' },
})
</script>