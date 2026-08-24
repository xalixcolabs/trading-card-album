<template>
  <AdminShell active="albums">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-lg font-bold tracking-tight text-ink">Álbumes</h2>
        <p class="mt-1 text-sm text-mist">{{ albums?.length ?? 0 }} álbumes en total.</p>
      </div>
      <button type="button" @click="openCreate"
        class="inline-flex shrink-0 items-center gap-1.5 rounded-xl bg-accent px-3.5 py-2.5 text-xs font-semibold text-accent-ink shadow-glow transition-transform active:scale-95">
        <PhPlus :size="15" weight="bold" />
        Nuevo
      </button>
    </div>

    <div v-if="pending" class="mt-5 flex flex-col gap-3">
      <div v-for="i in 4" :key="i" class="h-20 animate-pulse rounded-2xl bg-raise ring-1 ring-edge"></div>
    </div>

    <div v-else class="mt-5 flex flex-col gap-3">
      <div v-for="album in albums" :key="album.id"
        class="flex items-center gap-3 rounded-2xl bg-raise p-4 ring-1 ring-edge transition-transform active:scale-[0.98]">
        <button type="button" @click="router.push(`/admin/albums/${album.id}/cards`)"
          class="flex min-w-0 flex-1 items-center gap-3 text-left">
          <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-accent-soft text-accent">
            <PhStack :size="20" weight="duotone" />
          </div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-[15px] font-semibold text-ink">{{ album.title }}</p>
            <p class="mt-0.5 text-xs text-mist">
              {{ album.total_cards }} tarjetas · {{ album.participant_count }} participantes
            </p>
          </div>
        </button>
        <button type="button" @click.stop="openEdit(album)"
          class="shrink-0 rounded-lg bg-panel p-2.5 text-mist ring-1 ring-edge transition-transform active:scale-90"
          aria-label="Editar álbum">
          <PhPencilSimple :size="16" />
        </button>
        <button type="button" @click.stop="removeAlbum(album)"
          class="shrink-0 rounded-lg px-3 py-2 text-xs font-semibold transition-colors"
          :class="confirmingId === album.id ? 'bg-danger text-white' : 'bg-panel text-faint ring-1 ring-edge'"
          :disabled="deleting">
          {{ confirmingId === album.id ? '¿Eliminar?' : 'Eliminar' }}
        </button>
      </div>

      <p v-if="albums && albums.length === 0" class="mt-4 text-center text-sm text-faint">Aún no hay álbumes.</p>
    </div>

    <AdminAlbumSheet :is-open="sheetOpen" :mode="sheetMode" :album="editingAlbum" @close="sheetOpen = false" @saved="onSaved" />
  </AdminShell>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'admin' })
import { PhPencilSimple, PhPlus, PhStack } from '@phosphor-icons/vue'
import { getApiV1AdminAlbums, deleteApiV1AdminAlbumsId } from '~/services/admin/admin'
import type { AdminDtoAlbum } from '~/models'

const router = useRouter()
const toast = useToast()
const confirmingId = ref<string | null>(null)
const deleting = ref(false)
const sheetOpen = ref(false)
const sheetMode = ref<'create' | 'edit'>('create')
const editingAlbum = ref<AdminDtoAlbum | undefined>()

const { data: albums, pending, refresh } = useApiData(() => getApiV1AdminAlbums(), 'admin-albums')

const openCreate = () => {
  editingAlbum.value = undefined
  sheetMode.value = 'create'
  sheetOpen.value = true
}

const openEdit = (album: AdminDtoAlbum) => {
  editingAlbum.value = album
  sheetMode.value = 'edit'
  sheetOpen.value = true
}

function removeAlbum(album: AdminDtoAlbum) {
  if (confirmingId.value !== album.id) {
    confirmingId.value = album.id
    setTimeout(() => {
      if (confirmingId.value === album.id) confirmingId.value = null
    }, 2500)
    return
  }
  performDelete(album)
}

async function performDelete(album: AdminDtoAlbum) {
  deleting.value = true
  try {
    await deleteApiV1AdminAlbumsId(album.id!)
    toast.success({ title: 'Álbum eliminado', message: album.title || '' })
    confirmingId.value = null
    await onSaved()
  } catch {
    toast.error({ title: 'No se pudo eliminar', message: 'Inténtalo de nuevo.' })
  } finally {
    deleting.value = false
  }
}

const onSaved = async () => {
  await refresh()
  await refreshNuxtData('albums')
}
</script>