<template>
  <AppSheet :is-open="isOpen" @close="$emit('close')">
    <!-- Selector de tarjeta (asignar / regalar) -->
    <div v-if="mode !== 'view'" class="pt-3">
      <div class="flex items-center gap-2">
        <button type="button" @click="backToView"
          class="flex h-9 w-9 items-center justify-center rounded-xl bg-raise text-mist ring-1 ring-edge transition-transform active:scale-90"
          aria-label="Volver">
          <PhArrowLeft :size="18" />
        </button>
        <div class="min-w-0 flex-1">
          <h3 class="truncate text-base font-bold tracking-tight text-ink">
            {{ mode === 'assign' ? 'Cambiar tarjeta asignada' : 'Regalar tarjeta' }}
          </h3>
          <p class="truncate text-xs text-mist">{{ detail?.user?.name || 'Usuario' }}</p>
        </div>
        <span class="inline-flex h-9 w-9 items-center justify-center rounded-xl bg-accent-soft text-accent">
          <PhSwap v-if="mode === 'assign'" :size="18" weight="duotone" />
          <PhGift v-else :size="18" weight="duotone" />
        </span>
      </div>

      <p class="mt-3 text-[12px] leading-snug text-mist">
        {{ mode === 'assign'
          ? 'El usuario pasará a compartir esta tarjeta vía QR.'
          : 'Se sumará a su colección sin cambiar su tarjeta asignada.' }}
      </p>

      <div class="mt-4 flex flex-col gap-1.5">
        <label for="action-album" class="text-[13px] font-semibold text-mist">Álbum</label>
        <select id="action-album" v-model="selectedAlbumId"
          class="rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge focus:outline-none focus:ring-2 focus:ring-accent">
          <option value="" disabled>Selecciona un álbum</option>
          <option v-for="album in albums" :key="album.id" :value="album.id">{{ album.title }}</option>
        </select>
      </div>

      <div class="mt-4">
        <div v-if="cardsLoading" class="grid grid-cols-3 gap-3">
          <div v-for="i in 6" :key="i" class="aspect-2/3 animate-pulse rounded-lg bg-raise ring-1 ring-edge"></div>
        </div>
        <div v-else-if="cards.length" class="grid grid-cols-3 gap-3">
          <button v-for="card in cards" :key="card.id" type="button" @click="selectedCardId = card.id!"
            class="group relative overflow-hidden rounded-lg bg-raise ring-1 transition-all active:scale-95"
            :class="selectedCardId === card.id ? 'ring-2 ring-accent' : 'ring-edge'">
            <div class="aspect-2/3 w-full overflow-hidden">
              <img :src="card.image_url" :alt="card.name" loading="lazy" class="h-full w-full object-cover" />
            </div>
            <div class="px-1.5 py-1 text-left">
              <p class="truncate text-[11px] font-semibold text-mist">{{ card.name }}</p>
              <p class="font-mono text-[10px] text-faint tabular">#{{ card.number }}</p>
            </div>
            <span v-if="selectedCardId === card.id"
              class="absolute right-1.5 top-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-accent text-accent-ink">
              <PhCheck :size="12" weight="bold" />
            </span>
          </button>
        </div>
        <p v-else-if="selectedAlbumId" class="mt-3 text-center text-sm text-faint">Este álbum no tiene tarjetas.</p>
        <p v-else class="mt-3 text-center text-sm text-faint">Elige un álbum para ver sus tarjetas.</p>
      </div>

      <button type="button" :disabled="!selectedCardId || submitting" @click="confirmAction"
        class="mt-5 w-full rounded-2xl bg-accent px-4 py-3.5 text-sm font-semibold text-accent-ink shadow-glow transition-transform active:scale-[0.97] disabled:opacity-50">
        {{ submitting ? 'Guardando…' : (mode === 'assign' ? 'Asignar tarjeta' : 'Regalar tarjeta') }}
      </button>
    </div>

    <!-- Detalle del usuario -->
    <template v-else>
      <div v-if="loading" class="flex flex-col gap-3 pt-4">
        <div v-for="i in 3" :key="i" class="h-16 animate-pulse rounded-2xl bg-raise ring-1 ring-edge"></div>
      </div>

      <div v-else-if="detail" class="pt-3">
        <div class="flex items-center gap-3">
          <div class="h-12 w-12 overflow-hidden rounded-2xl ring-1 ring-edge">
            <UserAvatar :name="detail.user?.name" :picture="detail.user?.picture" size-class="text-lg" />
          </div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-[15px] font-semibold text-ink">{{ detail.user?.name || 'Sin nombre' }}</p>
            <p class="truncate text-[13px] text-mist">{{ detail.user?.email }}</p>
          </div>
          <span v-if="detail.user?.is_admin"
            class="shrink-0 rounded-full bg-accent-soft px-2.5 py-1 text-[11px] font-semibold text-accent">Admin</span>
        </div>

        <div class="mt-6">
          <h4 class="text-[11px] font-semibold uppercase tracking-[0.18em] text-faint">Álbumes ({{ detail.albums?.length ?? 0 }})</h4>
          <div v-if="detail.albums?.length" class="mt-3 flex flex-col gap-2">
            <div v-for="album in detail.albums" :key="album.id"
              class="flex items-center justify-between rounded-xl bg-raise px-4 py-3 ring-1 ring-edge-soft">
              <span class="truncate text-sm font-semibold text-ink">{{ album.title }}</span>
              <span class="font-mono text-xs text-mist tabular">{{ album.total_cards }}</span>
            </div>
          </div>
          <p v-else class="mt-3 text-sm text-faint">Sin álbumes.</p>
        </div>

        <div class="mt-6">
          <h4 class="text-[11px] font-semibold uppercase tracking-[0.18em] text-faint">Tarjetas recolectadas ({{ detail.cards?.length ?? 0 }})</h4>
          <div v-if="detail.cards?.length" class="mt-3 grid grid-cols-3 gap-3">
            <div v-for="card in detail.cards" :key="card.id" class="text-center">
              <div class="mx-auto aspect-2/3 w-full overflow-hidden rounded-lg bg-raise ring-1 ring-edge">
                <img :src="card.image_url" :alt="card.name" loading="lazy" class="h-full w-full object-cover" />
              </div>
              <p class="mt-1.5 truncate text-[11px] font-semibold text-mist">{{ card.name }}</p>
            </div>
          </div>
          <p v-else class="mt-3 text-sm text-faint">Aún no recolecta tarjetas.</p>
        </div>

        <div class="mt-6 grid grid-cols-2 gap-2">
          <button type="button" @click="openAction('assign')"
            class="inline-flex items-center justify-center gap-1.5 rounded-xl bg-panel px-3 py-3 text-xs font-semibold text-mist ring-1 ring-edge transition-transform active:scale-95">
            <PhSwap :size="15" /> Asignar tarjeta
          </button>
          <button type="button" @click="openAction('gift')"
            class="inline-flex items-center justify-center gap-1.5 rounded-xl bg-accent px-3 py-3 text-xs font-semibold text-accent-ink shadow-glow transition-transform active:scale-95">
            <PhGift :size="15" /> Regalar tarjeta
          </button>
        </div>
      </div>
    </template>
  </AppSheet>
</template>

<script setup lang="ts">
import { PhArrowLeft, PhCheck, PhGift, PhSwap } from '@phosphor-icons/vue'
import {
  getApiV1AdminAlbums,
  getApiV1AdminAlbumsIdCards,
  getApiV1AdminUsersId,
  postApiV1AdminUsersIdCards,
  postApiV1AdminUsersIdGiftCard,
} from '~/services/admin/admin'
import type { AdminDtoAlbum, AdminDtoUserDetail, CardModelCard } from '~/models'

const props = defineProps<{
  isOpen: boolean
  userId: string
}>()

defineEmits<{
  (e: 'close'): void
}>()

const toast = useToast()
const detail = ref<AdminDtoUserDetail | null>(null)
const loading = ref(false)

const mode = ref<'view' | 'assign' | 'gift'>('view')
const albums = ref<AdminDtoAlbum[]>([])
const selectedAlbumId = ref('')
const cards = ref<CardModelCard[]>([])
const cardsLoading = ref(false)
const selectedCardId = ref('')
const submitting = ref(false)

async function loadDetail() {
  if (!props.userId) {
    detail.value = null
    return
  }
  loading.value = true
  try {
    const response = await getApiV1AdminUsersId(props.userId)
    detail.value = response.data
  } catch {
    detail.value = null
  } finally {
    loading.value = false
  }
}

async function loadAlbums() {
  if (albums.value.length) return
  try {
    const response = await getApiV1AdminAlbums()
    albums.value = response.data ?? []
  } catch {
    albums.value = []
  }
}

async function loadCards() {
  cards.value = []
  selectedCardId.value = ''
  if (!selectedAlbumId.value) return
  cardsLoading.value = true
  try {
    const response = await getApiV1AdminAlbumsIdCards(selectedAlbumId.value)
    cards.value = response.data ?? []
  } catch {
    cards.value = []
  } finally {
    cardsLoading.value = false
  }
}

function openAction(action: 'assign' | 'gift') {
  mode.value = action
  selectedAlbumId.value = ''
  selectedCardId.value = ''
  cards.value = []
  loadAlbums()
}

function backToView() {
  mode.value = 'view'
  selectedAlbumId.value = ''
  selectedCardId.value = ''
  cards.value = []
}

async function confirmAction() {
  if (!selectedCardId.value) return
  const selected = cards.value.find(card => card.id === selectedCardId.value)
  submitting.value = true
  try {
    if (mode.value === 'assign') {
      await postApiV1AdminUsersIdCards(props.userId, { card_id: selectedCardId.value })
      toast.success({ title: 'Tarjeta asignada', message: selected?.name || '' })
    } else {
      await postApiV1AdminUsersIdGiftCard(props.userId, { card_id: selectedCardId.value })
      toast.success({ title: 'Tarjeta regalada', message: selected?.name || '' })
    }
    backToView()
    await loadDetail()
  } catch {
    toast.error({ title: 'No se pudo completar', message: 'Inténtalo de nuevo.' })
  } finally {
    submitting.value = false
  }
}

watch(() => props.isOpen, (isOpen) => {
  if (!isOpen) return
  backToView()
  loadDetail()
}, { immediate: true })

watch(() => props.userId, () => {
  if (props.isOpen && props.userId) loadDetail()
})

watch(selectedAlbumId, loadCards)
</script>
