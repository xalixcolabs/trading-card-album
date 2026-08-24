<template>
  <div class="min-h-[100dvh] w-full bg-base">
    <div
      class="grain relative mx-auto min-h-[100dvh] w-full max-w-[480px] border-x border-edge bg-panel">
      <div class="safe-top">
        <slot />
      </div>
      <div :class="showTabBar ? 'h-28' : ''"></div>
    </div>

    <TabBar v-if="showTabBar" :active="tabActive" @scan="scanOpen = true" />

    <ScanQrModal :is-open="scanOpen" @close="scanOpen = false" @card-added-successfully="onCardAdded" />
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const scanOpen = ref(false)

const showTabBar = computed(() => {
  return !route.path.startsWith('/album') && route.path !== '/login'
})

const tabActive = computed<'home' | 'contacts' | 'profile'>(() => {
  if (route.path.startsWith('/profile')) return 'profile'
  if (route.path.startsWith('/contactos')) return 'contacts'
  return 'home'
})

const onCardAdded = async () => {
  await refreshNuxtData()
}
</script>