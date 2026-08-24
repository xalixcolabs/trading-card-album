<template>
  <div class="px-5 pt-4 pb-32">
    <!-- Encabezado -->
    <div class="flex items-center gap-2">
      <button type="button" @click="router.back()"
        class="flex h-10 w-10 items-center justify-center rounded-xl bg-raise text-mist ring-1 ring-edge transition-transform active:scale-90"
        aria-label="Volver">
        <PhArrowLeft :size="20" />
      </button>
      <div class="min-w-0 flex-1">
        <h1 class="truncate text-xl font-bold tracking-tight text-ink">{{ albumTitle }}</h1>
        <p class="text-[13px] text-mist">{{ collectionCount }} de {{ totalCards }} tarjetas</p>
      </div>
      <button type="button" @click="inviteOpen = true"
        class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-raise text-mist ring-1 ring-edge transition-transform active:scale-90"
        aria-label="Invitar al álbum">
        <PhShareNetwork :size="20" />
      </button>
    </div>

    <!-- Progreso -->
    <div class="mt-4 h-1.5 w-full overflow-hidden rounded-full bg-raise">
      <div class="h-full rounded-full bg-accent transition-all duration-700" :style="{ width: progressPct + '%' }"></div>
    </div>

    <!-- Tu carta asignada -->
    <section class="mt-7">
      <p class="text-[11px] font-semibold uppercase tracking-[0.18em] text-faint">Tu carta asignada</p>

      <div v-if="assignedPending" class="mx-auto mt-4 w-[220px] animate-pulse">
        <div class="aspect-2/3 w-full rounded-card bg-raise ring-1 ring-edge"></div>
      </div>

      <div v-else-if="assignedCard" class="relative mx-auto mt-4 w-[220px]">
        <div class="absolute inset-0 -z-10 scale-110 rounded-full bg-accent/10 blur-2xl"></div>
        <Card :card="assignedCard" />
      </div>

      <p v-else class="mt-4 text-center text-sm text-mist">
        Aún no tienes una tarjeta asignada en este álbum.
      </p>
    </section>

    <!-- Álbum completo -->
    <section class="mt-8">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-bold tracking-tight text-ink">Tu colección</h2>
        <span class="font-mono text-xs font-semibold text-mist tabular">{{ collectionCount }}/{{ totalCards }}</span>
      </div>

      <div v-if="albumPending" class="mt-4 grid grid-cols-2 gap-4">
        <div v-for="i in 4" :key="i" class="aspect-2/3 w-full animate-pulse rounded-card bg-raise ring-1 ring-edge"></div>
      </div>

      <div v-else-if="slots.length" class="mt-4 grid grid-cols-2 gap-4">
        <template v-for="(card, index) in slots" :key="index">
          <Card v-if="card" :card="card" />
          <div v-else
            class="relative flex aspect-2/3 w-full flex-col items-center justify-center gap-2 rounded-card bg-raise/50 border border-dashed border-edge">
            <span class="text-4xl font-bold text-accent/70">?</span>
            <span class="text-[10px] font-semibold uppercase tracking-[0.14em] text-faint">Por descubrir</span>
          </div>
        </template>
      </div>

      <div v-else class="mt-4 overflow-hidden rounded-3xl border border-dashed border-edge bg-raise/60 p-6 text-center">
        <div class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-accent-soft text-accent">
          <PhStack :size="26" weight="duotone" />
        </div>
        <h2 class="mt-4 text-lg font-bold text-ink">Sin tarjetas aún</h2>
        <p class="mx-auto mt-1 max-w-[34ch] text-sm leading-relaxed text-mist">
          Este álbum todavía no tiene tarjetas registradas.
        </p>
      </div>

      <p v-if="collectionCount === 0 && slots.length" class="mt-5 text-center text-sm text-mist">
        Escanea el QR de otros desarrolladores para desbloquear tarjetas y completar tu álbum.
      </p>
    </section>

    <!-- Dock de acciones -->
    <div
      class="fixed bottom-0 left-1/2 z-40 w-full max-w-[480px] -translate-x-1/2 border-t border-edge bg-panel/90 px-4 py-3 backdrop-blur-xl safe-bottom">
      <div class="flex gap-3">
        <button type="button" @click="scanOpen = true"
          class="flex flex-1 items-center justify-center gap-2 rounded-2xl bg-accent px-4 py-3.5 text-sm font-semibold text-accent-ink shadow-glow transition-transform active:scale-[0.97]">
          <PhScan :size="20" weight="fill" />
          Escanear
        </button>
        <button type="button" @click="qrOpen = true"
          class="flex flex-1 items-center justify-center gap-2 rounded-2xl bg-raise px-4 py-3.5 text-sm font-semibold text-ink ring-1 ring-edge transition-transform active:scale-[0.97]">
          <PhQrCode :size="20" weight="fill" />
          Mi QR
        </button>
      </div>
    </div>

    <ScanQrModal :is-open="scanOpen" @close="scanOpen = false" @card-added-successfully="onCardAdded" />
    <MyQrModal :is-open="qrOpen" :album-id="albumId" @close="qrOpen = false" />
    <InviteSheet :is-open="inviteOpen" :album-id="albumId" @close="inviteOpen = false" />
  </div>
</template>

<script setup lang="ts">
import { PhArrowLeft, PhQrCode, PhScan, PhShareNetwork, PhStack } from '@phosphor-icons/vue'
import { getApiV1AlbumId, getApiV1AlbumIdAssignedCard } from '~/services/album/album'
import { getApiV1AuthMe } from '~/services/auth/auth'
import type { CardModelCard } from '~/models'

const route = useRoute()
const router = useRouter()
const albumId = route.params.id as string

const scanOpen = ref(false)
const qrOpen = ref(false)
const inviteOpen = ref(false)

const { data: profile } = useApiData(() => getApiV1AuthMe(), 'profile')
const { data: album, pending: albumPending, refresh: refreshAlbum } = useApiData(() => getApiV1AlbumId(albumId), `album-${albumId}`)
const { data: assignedCard, pending: assignedPending } = useApiData(
  () => getApiV1AlbumIdAssignedCard(albumId), `assigned-${albumId}`,
)

const albumTitle = computed(() => album.value?.title || 'Álbum')
const collected = computed(() => album.value?.cards ?? [])
const totalCards = computed(() => album.value?.total_cards || collected.value.length)
const collectionCount = computed(() => collected.value.length)

const slots = computed<(CardModelCard | null)[]>(() => {
  const total = totalCards.value
  const list: (CardModelCard | null)[] = new Array(total).fill(null)
  const sorted = [...collected.value].sort((a, b) => Number(a.number) - Number(b.number))
  for (const card of sorted) {
    const index = Number(card.number) - 1
    if (Number.isFinite(index) && index >= 0 && index < total) {
      list[index] = card
    } else {
      const empty = list.findIndex(slot => slot === null)
      if (empty !== -1) list[empty] = card
    }
  }
  return list
})

const progressPct = computed(() => {
  if (!totalCards.value) return 0
  return Math.min(100, Math.round((collectionCount.value / totalCards.value) * 100))
})

watch(profile, (value) => {
  if (value) useProfile(value)
}, { immediate: true })

const onCardAdded = async () => {
  await refreshAlbum()
  await refreshNuxtData('assigned-' + albumId)
}
</script>