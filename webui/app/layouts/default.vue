<template>
  <div class="min-h-[100dvh] w-full bg-base">
    <div
      class="grain relative mx-auto min-h-[100dvh] w-full max-w-[480px] border-x border-edge bg-panel">
      <div class="safe-top">
        <slot />
      </div>
      <div :class="showTabBar ? 'h-28' : ''"></div>
    </div>

    <TabBar v-if="showTabBar" :active="tabActive" @add="addOpen = true" />
    <AddAlbumSheet :is-open="addOpen" @close="addOpen = false" @created="onAlbumCreated" />
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const addOpen = ref(false)

const showTabBar = computed(() => {
  return !route.path.startsWith('/album') && route.path !== '/login'
})

const tabActive = computed<'home' | 'contacts' | 'profile'>(() => {
  if (route.path.startsWith('/profile')) return 'profile'
  if (route.path.startsWith('/contactos')) return 'contacts'
  return 'home'
})

const onAlbumCreated = async () => {
  await refreshNuxtData('albums')
}
</script>