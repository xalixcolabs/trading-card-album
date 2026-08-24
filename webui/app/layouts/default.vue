<template>
  <div class="min-h-[100dvh] w-full bg-base">
    <div
      class="grain relative mx-auto min-h-[100dvh] w-full max-w-[480px] border-x border-edge bg-panel">
      <div class="safe-top">
        <slot />
      </div>
      <div :class="showTabBar ? 'h-28' : ''"></div>
    </div>

    <TabBar v-if="showTabBar" :active="tabActive" @join="joinOpen = true" />
    <JoinAlbumSheet :is-open="joinOpen" @close="joinOpen = false" @joined="onJoined" />
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const joinOpen = ref(false)

const showTabBar = computed(() => {
  return !route.path.startsWith('/album') && !route.path.startsWith('/admin') && route.path !== '/login'
})

const tabActive = computed<'home' | 'contacts'>(() => {
  if (route.path.startsWith('/contactos')) return 'contacts'
  return 'home'
})

const onJoined = async (albumId: string) => {
  await refreshNuxtData('albums')
  await navigateTo(`/album/${albumId}`)
}
</script>