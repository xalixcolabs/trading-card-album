<template>
  <div class="px-5 pt-4 pb-24">
    <div class="flex items-center gap-2">
      <button type="button" @click="router.push('/')"
        class="flex h-10 w-10 items-center justify-center rounded-xl bg-raise text-mist ring-1 ring-edge transition-transform active:scale-90"
        aria-label="Volver">
        <PhArrowLeft :size="20" />
      </button>
      <div class="flex min-w-0 flex-1 items-center gap-2">
        <span class="flex h-10 w-10 items-center justify-center rounded-xl bg-accent-soft text-accent">
          <PhShieldStar :size="20" weight="fill" />
        </span>
        <h1 class="truncate text-xl font-bold tracking-tight text-ink">Administración</h1>
      </div>
    </div>

    <nav class="no-scrollbar -mx-5 mt-5 flex gap-2 overflow-x-auto px-5">
      <NuxtLink to="/admin" class="shrink-0 rounded-xl px-3.5 py-2 text-xs font-semibold ring-1 transition-colors"
        :class="active === 'overview' ? 'bg-accent text-accent-ink ring-accent' : 'bg-raise text-mist ring-edge'">
        Resumen
      </NuxtLink>
      <NuxtLink to="/admin/albums" class="shrink-0 rounded-xl px-3.5 py-2 text-xs font-semibold ring-1 transition-colors"
        :class="active === 'albums' ? 'bg-accent text-accent-ink ring-accent' : 'bg-raise text-mist ring-edge'">
        Álbumes
      </NuxtLink>
      <NuxtLink to="/admin/users" class="shrink-0 rounded-xl px-3.5 py-2 text-xs font-semibold ring-1 transition-colors"
        :class="active === 'users' ? 'bg-accent text-accent-ink ring-accent' : 'bg-raise text-mist ring-edge'">
        Usuarios
      </NuxtLink>
      <NuxtLink to="/admin/cards" class="shrink-0 rounded-xl px-3.5 py-2 text-xs font-semibold ring-1 transition-colors"
        :class="active === 'cards' ? 'bg-accent text-accent-ink ring-accent' : 'bg-raise text-mist ring-edge'">
        Tarjetas
      </NuxtLink>
    </nav>

    <div class="mt-6">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { PhArrowLeft, PhShieldStar } from '@phosphor-icons/vue'

defineProps<{
  active: 'overview' | 'albums' | 'users' | 'cards'
}>()

const router = useRouter()
</script>