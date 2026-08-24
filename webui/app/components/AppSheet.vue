<template>
  <Teleport to="body">
    <Transition name="sheet">
      <div v-if="isOpen" class="fixed inset-0 z-50 flex items-end justify-center">
        <div class="absolute inset-0 bg-black/65 backdrop-blur-[2px]" @click="$emit('close')"></div>

        <div
          class="grain relative max-h-[88dvh] w-full max-w-[480px] overflow-y-auto rounded-t-[1.5rem] border-t border-edge bg-panel shadow-[0_-18px_50px_-12px_rgba(0,0,0,0.8)]">
          <div class="mx-auto mt-3 h-1 w-11 shrink-0 rounded-full bg-edge"></div>
          <div class="px-5 pb-8 safe-bottom">
            <slot />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
defineProps<{
  isOpen: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()
</script>

<style scoped>
.sheet-enter-active,
.sheet-leave-active {
  transition: transform 0.32s cubic-bezier(0.32, 0.72, 0, 1), opacity 0.22s ease;
}

.sheet-enter-active > div:last-child,
.sheet-leave-active > div:last-child {
  transition: transform 0.32s cubic-bezier(0.32, 0.72, 0, 1);
}

.sheet-enter-from,
.sheet-leave-to {
  opacity: 0;
}

.sheet-enter-from > div:last-child,
.sheet-leave-to > div:last-child {
  transform: translateY(100%);
}
</style>