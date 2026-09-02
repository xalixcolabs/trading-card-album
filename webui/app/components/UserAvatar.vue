<template>
  <div v-if="picture" class="overflow-hidden rounded-[inherit]">
    <img :src="picture" :alt="name || 'Avatar'" loading="lazy" referrerpolicy="no-referrer" class="h-full w-full object-cover" />
  </div>
  <div v-else
    class="flex h-full w-full items-center justify-center rounded-[inherit] bg-raise text-lg font-bold text-accent ring-1 ring-edge"
    :class="sizeClass">
    {{ initial }}
  </div>
</template>

<script setup lang="ts">
const { name, picture, sizeClass } = defineProps<{
  name?: string
  picture?: string
  sizeClass?: string
}>()

const initial = computed(() => (name || '?').trim().charAt(0).toUpperCase())
</script>