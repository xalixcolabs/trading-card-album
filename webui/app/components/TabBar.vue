<template>
  <nav class="fixed bottom-0 left-1/2 z-40 w-full max-w-[480px] -translate-x-1/2 px-3 safe-bottom">
    <div
      class="flex items-center rounded-2xl border border-edge bg-panel/90 px-2 py-2 shadow-[0_-8px_30px_-12px_rgba(0,0,0,0.6)] backdrop-blur-xl">

      <NuxtLink to="/"
        class="flex flex-1 flex-col items-center gap-1 rounded-xl py-1.5 transition-transform active:scale-95"
        :class="active === 'home' ? 'text-ink' : 'text-faint'">
        <PhHouse :size="21" :weight="active === 'home' ? 'fill' : 'regular'" />
        <span class="text-[10px] font-semibold tracking-wide">Inicio</span>
      </NuxtLink>

      <button type="button" @click="$emit('join')"
        class="-mt-7 flex h-13 w-13 shrink-0 items-center justify-center rounded-2xl bg-accent text-accent-ink shadow-glow ring-4 ring-panel transition-transform active:scale-90"
        aria-label="Unirse a un álbum">
        <PhPlus :size="26" weight="bold" />
      </button>

      <NuxtLink to="/contactos"
        class="flex flex-1 flex-col items-center gap-1 rounded-xl py-1.5 transition-transform active:scale-95"
        :class="active === 'contacts' ? 'text-ink' : 'text-faint'">
        <PhUsers :size="21" :weight="active === 'contacts' ? 'fill' : 'regular'" />
        <span class="text-[10px] font-semibold tracking-wide">Contactos</span>
      </NuxtLink>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { PhHouse, PhPlus, PhUsers } from '@phosphor-icons/vue'

defineProps<{
  active: 'home' | 'contacts'
}>()

defineEmits<{
  (e: 'join'): void
}>()
</script>