<template>
  <div>
    <div @click="openHero"
      class="group relative aspect-2/3 w-full rounded-2xl overflow-hidden shadow-md cursor-pointer transition-all duration-300 hover:scale-[1.03] hover:shadow-xl active:scale-[0.98] border border-base-200/80 bg-base-200 text-base-content">
      <img :src="`${card.image_url}`" :alt="card.name"
        class="w-full h-full object-cover select-none transition-transform duration-500 group-hover:scale-105" />

      <div
        class="absolute inset-0 bg-linear-to-t from-black/80 via-black/20 to-transparent opacity-60 group-hover:opacity-40 transition-opacity">
      </div>

      <div
        class="absolute bottom-2 left-2 right-2 flex justify-between items-center bg-black/60 backdrop-blur-xs py-1 px-2.5 rounded-lg border border-white/10">
        <span class="text-[10px] font-bold text-primary font-mono">#{{ card.number }}</span>
        <span class="text-xs font-semibold text-primary-content truncate max-w-[80%]">{{ card.name }}</span>
      </div>
    </div>

    <Teleport to="body">
      <Transition name="fade">
        <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm"
          @click.self="closeHero">
          <div
            class="relative w-full max-w-[320px] aspect-[2.5/3.6] bg-base-100 text-base-content rounded-2xl shadow-[0_0_35px_rgba(0,0,0,0.3)] dark:shadow-[0_0_35px_rgba(0,0,0,0.8)] border border-base-300 flex flex-col p-3.5 overflow-hidden select-none animate-card-float">
            <div
              class="absolute inset-0 pointer-events-none opacity-40 mix-blend-color-dodge animate-shine-sweep bg-holo-shine">
            </div>

            <button @click="closeHero"
              class="absolute top-2 right-2 z-10 w-8 h-8 rounded-full bg-black/60 hover:bg-black/80 text-white flex items-center justify-center border border-white/10 hover:border-white/20 transition-all cursor-pointer">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>

            <div
              class="flex-1 flex flex-col justify-between border-2 border-base-300 rounded-xl p-3 h-full bg-base-100 shadow-[inset_0_0_12px_rgba(0,0,0,0.02)]">

              <div
                class="flex justify-between items-center bg-base-200 px-3 py-2 rounded-lg border border-base-300 shadow-[inset_0_1px_2px_rgba(0,0,0,0.03)]">
                <span class="text-[13px] font-black text-base-content leading-tight truncate max-w-47.5">{{ card.name
                  }}</span>
                <span class="text-[11px] font-black font-mono text-primary">#{{ card.number }}</span>
              </div>

              <div
                class="my-3 relative aspect-square w-full overflow-hidden rounded-lg border border-base-300 bg-base-200/50 flex items-center justify-center shadow-inner">
                <img :src="`${card.image_url}`" :alt="card.name" class="w-full h-full object-cover" />
              </div>

              <div
                class="flex-1 flex items-center justify-center bg-base-200/30 p-3 rounded-lg border border-base-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.02)]">
                <p class="text-[11px] text-base-content/85 italic leading-relaxed text-center font-medium">
                  {{ card.description || 'Una carta misteriosa por descubrir.' }}
                </p>
              </div>
            </div>

            <div
              class="absolute inset-0 pointer-events-none bg-linear-to-tr from-transparent via-white/5 to-transparent mix-blend-overlay">
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">

defineProps<{
  card: {
    id: string,
    album_id: string,
    number: string,
    name: string,
    description: string,
    image_url: string,
    created_at: number,
    updated_at: number,
  }
}>()

const isOpen = useState(() => false)

const openHero = () => {
  isOpen.value = true
}

const closeHero = () => {
  isOpen.value = false
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes float-card {
  0% {
    transform: perspective(1000px) rotateX(3deg) rotateY(-6deg) scale3d(1, 1, 1);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.3);
  }

  25% {
    transform: perspective(1000px) rotateX(-3deg) rotateY(-3deg) scale3d(1, 1, 1);
    box-shadow: 0 25px 40px rgba(0, 0, 0, 0.25);
  }

  50% {
    transform: perspective(1000px) rotateX(3deg) rotateY(6deg) scale3d(1, 1, 1);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.3);
  }

  75% {
    transform: perspective(1000px) rotateX(4deg) rotateY(-3deg) scale3d(1, 1, 1);
    box-shadow: 0 20px 35px rgba(0, 0, 0, 0.25);
  }

  100% {
    transform: perspective(1000px) rotateX(3deg) rotateY(-6deg) scale3d(1, 1, 1);
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.3);
  }
}

.animate-card-float {
  animation: float-card 7s infinite ease-in-out;
}

@keyframes shine-sweep {
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

.animate-shine-sweep {
  animation: shine-sweep 7s infinite ease-in-out;
  background-size: 200% 200%;
}

.bg-holo-shine {
  background-image: linear-gradient(115deg,
      transparent 15%,
      rgba(255, 255, 255, 0.15) 25%,
      rgba(0, 229, 255, 0.12) 40%,
      rgba(255, 0, 128, 0.12) 60%,
      rgba(255, 255, 255, 0.15) 75%,
      transparent 85%);
}
</style>