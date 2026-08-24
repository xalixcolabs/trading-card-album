<template>
  <div>
    <button type="button" @click="isOpen = true"
      class="group relative block aspect-2/3 w-full cursor-pointer overflow-hidden rounded-card bg-raise text-left ring-1 ring-edge transition-all duration-300 hover:ring-accent/40 active:scale-[0.97]">
      <img :src="card.image_url" :alt="card.name" loading="lazy"
        class="h-full w-full select-none object-cover transition-transform duration-500 group-hover:scale-105" />

      <div class="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/75 via-black/10 to-transparent"></div>

      <div
        class="absolute bottom-2 left-2 right-2 flex items-center justify-between rounded-lg border border-white/10 bg-black/45 px-2 py-1.5 backdrop-blur-sm">
        <span class="font-mono text-[10px] font-bold text-accent tabular">#{{ card.number }}</span>
        <span class="truncate pl-2 text-xs font-semibold text-white">{{ card.name }}</span>
      </div>
    </button>

    <AppSheet :is-open="isOpen" @close="isOpen = false">
      <div class="pt-4 text-center">
        <div class="relative mx-auto w-full max-w-[260px]">
          <div
            class="relative aspect-2/3 w-full overflow-hidden rounded-card border border-edge bg-raise shadow-glow">
            <img :src="card.image_url" :alt="card.name" class="h-full w-full object-cover" />

            <div class="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/55 via-transparent to-transparent"></div>

            <div class="pointer-events-none absolute inset-0 holo-shine"></div>

            <div class="absolute bottom-3 left-3 right-3 text-left">
              <span class="font-mono text-xs font-bold text-accent tabular">#{{ card.number }}</span>
              <h3 class="text-lg font-bold leading-tight text-white">{{ card.name }}</h3>
            </div>
          </div>

          <div
            class="pointer-events-none absolute -inset-px rounded-card ring-1 ring-inset ring-white/10">
          </div>
        </div>

        <p class="mx-auto mt-5 max-w-[34ch] text-sm leading-relaxed text-mist">
          {{ card.description || 'Una carta misteriosa por descubrir.' }}
        </p>
      </div>
    </AppSheet>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  card: {
    id: string
    album_id: string
    number: string
    name: string
    description: string
    image_url: string
    created_at: number
    updated_at: number
  }
}>()

const isOpen = ref(false)
</script>

<style scoped>
.holo-shine {
  background-image: linear-gradient(115deg,
      transparent 15%,
      rgba(255, 255, 255, 0.14) 25%,
      rgba(45, 212, 191, 0.14) 40%,
      rgba(255, 255, 255, 0.14) 75%,
      transparent 85%);
  background-size: 200% 200%;
  animation: holo 8s ease-in-out infinite;
}

@media (prefers-reduced-motion: reduce) {
  .holo-shine {
    animation: none;
  }
}

@keyframes holo {
  0% {
    background-position: 0% 0%;
  }

  50% {
    background-position: 100% 100%;
  }

  100% {
    background-position: 0% 0%;
  }
}
</style>