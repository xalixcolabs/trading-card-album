<template>
  <AdminShell active="cards">
    <h2 class="text-lg font-bold tracking-tight text-ink">Tarjetas</h2>
    <p class="mt-1 text-sm text-mist">Elige un álbum para ver sus tarjetas.</p>

    <div v-if="pending" class="mt-5 flex flex-col gap-3">
      <div v-for="i in 4" :key="i" class="h-20 animate-pulse rounded-2xl bg-raise ring-1 ring-edge"></div>
    </div>

    <div v-else class="mt-5 flex flex-col gap-3">
      <button v-for="album in albums" :key="album.id" type="button"
        @click="router.push(`/admin/albums/${album.id}/cards`)"
        class="flex items-center gap-3 rounded-2xl bg-raise p-4 text-left ring-1 ring-edge transition-transform active:scale-[0.98]">
        <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-accent-soft text-accent">
          <PhStack :size="20" weight="duotone" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-[15px] font-semibold text-ink">{{ album.title }}</p>
          <p class="mt-0.5 text-xs text-mist">{{ album.total_cards }} tarjetas</p>
        </div>
        <PhCaretRight :size="18" class="shrink-0 text-faint" />
      </button>

      <p v-if="albums && albums.length === 0" class="mt-4 text-center text-sm text-faint">Aún no hay álbumes.</p>
    </div>
  </AdminShell>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'admin' })
import { PhCaretRight, PhStack } from '@phosphor-icons/vue'
import { getApiV1AdminAlbums } from '~/services/admin/admin'

const router = useRouter()
const { data: albums, pending } = useApiData(() => getApiV1AdminAlbums(), 'admin-albums')
</script>