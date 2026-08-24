<template>
  <AdminShell active="albums">
    <h2 class="text-lg font-bold tracking-tight text-ink">Álbumes</h2>
    <p class="mt-1 text-sm text-mist">{{ albums?.length ?? 0 }} álbumes en total.</p>

    <div v-if="pending" class="mt-5 flex flex-col gap-3">
      <div v-for="i in 4" :key="i" class="h-20 animate-pulse rounded-2xl bg-raise ring-1 ring-edge"></div>
    </div>

    <div v-else class="mt-5 flex flex-col gap-3">
      <div v-for="album in albums" :key="album.id"
        class="flex items-center gap-3 rounded-2xl bg-raise p-4 ring-1 ring-edge">
        <div class="min-w-0 flex-1">
          <p class="truncate text-[15px] font-semibold text-ink">{{ album.title }}</p>
          <p class="mt-0.5 text-xs text-mist">
            {{ album.total_cards }} tarjetas · {{ album.participant_count }} participantes
          </p>
        </div>
        <button type="button" @click="removeAlbum(album)"
          class="shrink-0 rounded-lg px-3 py-2 text-xs font-semibold transition-colors"
          :class="confirmingId === album.id ? 'bg-danger text-white' : 'bg-panel text-faint ring-1 ring-edge'"
          :disabled="deleting">
          {{ confirmingId === album.id ? '¿Eliminar?' : 'Eliminar' }}
        </button>
      </div>

      <p v-if="albums && albums.length === 0" class="mt-4 text-center text-sm text-faint">Aún no hay álbumes.</p>
    </div>
  </AdminShell>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'admin' })
import { getApiV1AdminAlbums, deleteApiV1AdminAlbumsId } from '~/services/admin/admin'

const toast = useToast()
const confirmingId = ref<string | null>(null)
const deleting = ref(false)

const { data: albums, pending, refresh } = useApiData(() => getApiV1AdminAlbums(), 'admin-albums')

function removeAlbum(album: { id: string; title?: string }) {
  if (confirmingId.value !== album.id) {
    confirmingId.value = album.id
    setTimeout(() => {
      if (confirmingId.value === album.id) confirmingId.value = null
    }, 2500)
    return
  }
  performDelete(album)
}

async function performDelete(album: { id: string; title?: string }) {
  deleting.value = true
  try {
    await deleteApiV1AdminAlbumsId(album.id)
    toast.success({ title: 'Álbum eliminado', message: album.title || '' })
    confirmingId.value = null
    await refresh()
    await refreshNuxtData('albums')
  } catch {
    toast.error({ title: 'No se pudo eliminar', message: 'Inténtalo de nuevo.' })
  } finally {
    deleting.value = false
  }
}
</script>