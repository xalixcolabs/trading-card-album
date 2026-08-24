<template>
  <div class="px-5 pt-4 pb-24">
    <div class="flex items-center gap-2">
      <button type="button" @click="router.push('/admin/albums')"
        class="flex h-10 w-10 items-center justify-center rounded-xl bg-raise text-mist ring-1 ring-edge transition-transform active:scale-90"
        aria-label="Volver a álbumes">
        <PhArrowLeft :size="20" />
      </button>
      <div class="min-w-0 flex-1">
        <h1 class="truncate text-xl font-bold tracking-tight text-ink">{{ albumTitle }}</h1>
        <p class="text-[13px] text-mist">{{ cards?.length ?? 0 }} tarjetas</p>
      </div>
      <button type="button" @click="openCreate"
        class="inline-flex shrink-0 items-center gap-1.5 rounded-xl bg-accent px-3.5 py-2.5 text-xs font-semibold text-accent-ink shadow-glow transition-transform active:scale-95">
        <PhPlus :size="15" weight="bold" />
        Nueva
      </button>
    </div>

    <div v-if="pending" class="mt-5 flex flex-col gap-3">
      <div v-for="i in 4" :key="i" class="h-16 animate-pulse rounded-2xl bg-raise ring-1 ring-edge"></div>
    </div>

    <div v-else class="mt-5 flex flex-col gap-3">
      <div v-for="card in cards" :key="card.id"
        class="flex items-center gap-3 rounded-2xl bg-raise p-3 ring-1 ring-edge">
        <div class="h-14 w-10 shrink-0 overflow-hidden rounded-lg bg-panel ring-1 ring-edge">
          <img :src="card.image_url" :alt="card.name" loading="lazy" class="h-full w-full object-cover" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-[15px] font-semibold text-ink">{{ card.name }}</p>
          <p class="font-mono text-xs text-mist tabular">#{{ card.number }}</p>
        </div>
        <button type="button" @click="openEdit(card)"
          class="shrink-0 rounded-lg bg-panel p-2.5 text-mist ring-1 ring-edge transition-transform active:scale-90"
          aria-label="Editar tarjeta">
          <PhPencilSimple :size="16" />
        </button>
        <button type="button" @click="removeCard(card)"
          class="shrink-0 rounded-lg px-3 py-2 text-xs font-semibold transition-colors"
          :class="confirmingId === card.id ? 'bg-danger text-white' : 'bg-panel text-faint ring-1 ring-edge'"
          :disabled="deleting">
          {{ confirmingId === card.id ? '¿Eliminar?' : 'Eliminar' }}
        </button>
      </div>

      <p v-if="cards && cards.length === 0" class="mt-4 text-center text-sm text-faint">
        Este álbum aún no tiene tarjetas.
      </p>
    </div>

    <AdminCardSheet :is-open="sheetOpen" :mode="sheetMode" :card="editingCard" :albums="[]" :fixed-album-id="albumId"
      @close="sheetOpen = false" @saved="onSaved" />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'admin' })
import { PhArrowLeft, PhPencilSimple, PhPlus } from '@phosphor-icons/vue'
import { getApiV1AdminAlbumsIdCards, deleteApiV1AdminCardsId, getApiV1AdminAlbums } from '~/services/admin/admin'
import type { CardModelCard } from '~/models'

const route = useRoute()
const router = useRouter()
const albumId = route.params.id as string

const toast = useToast()
const confirmingId = ref<string | null>(null)
const deleting = ref(false)
const sheetOpen = ref(false)
const sheetMode = ref<'create' | 'edit'>('create')
const editingCard = ref<CardModelCard | undefined>()

const { data: cards, pending, refresh } = useApiData(() => getApiV1AdminAlbumsIdCards(albumId), `admin-cards-${albumId}`)
const { data: albums } = useApiData(() => getApiV1AdminAlbums(), 'admin-albums')

const albumTitle = computed(() => {
  const found = albums.value?.find(album => album.id === albumId)
  return found?.title || 'Tarjetas del álbum'
})

const openCreate = () => {
  editingCard.value = undefined
  sheetMode.value = 'create'
  sheetOpen.value = true
}

const openEdit = (card: CardModelCard) => {
  editingCard.value = card
  sheetMode.value = 'edit'
  sheetOpen.value = true
}

function removeCard(card: CardModelCard) {
  if (confirmingId.value !== card.id) {
    confirmingId.value = card.id
    setTimeout(() => {
      if (confirmingId.value === card.id) confirmingId.value = null
    }, 2500)
    return
  }
  performDelete(card)
}

async function performDelete(card: CardModelCard) {
  deleting.value = true
  try {
    await deleteApiV1AdminCardsId(card.id!)
    toast.success({ title: 'Tarjeta eliminada', message: card.name || '' })
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
  await refreshNuxtData('admin-albums')
}
</script>