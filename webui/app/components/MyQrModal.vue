<template>
  <AppSheet :is-open="isOpen" @close="$emit('close')">
    <div class="flex items-center justify-between pt-3">
      <div>
        <h3 class="text-lg font-bold tracking-tight text-ink">Mi código QR</h3>
        <p class="text-sm text-mist">{{ profile?.name || 'Participante' }}</p>
      </div>
      <button type="button" @click="refreshQr"
        class="inline-flex items-center gap-1.5 rounded-xl bg-raise px-3 py-2 text-xs font-semibold text-mist ring-1 ring-edge transition-transform active:scale-95">
        <PhArrowsClockwise :size="15" />
        Refrescar
      </button>
    </div>

    <p class="mt-4 inline-flex items-center gap-1.5 rounded-full bg-accent-soft px-3 py-1 text-[11px] font-semibold text-accent">
      <PhDotOutline :size="14" weight="fill" />
      Código de un solo uso
    </p>

    <div class="mx-auto mt-5 flex aspect-square w-full max-w-[260px] items-center justify-center rounded-2xl bg-white p-4 ring-1 ring-edge">
      <img v-if="imageUrl" :src="imageUrl" alt="Mi código QR para compartir" class="h-full w-full select-none" />
      <div v-else class="h-24 w-24 animate-pulse rounded-xl bg-black/10"></div>
    </div>

    <p class="mt-5 text-center text-sm leading-relaxed text-mist">
      Muestra este código a otros desarrolladores para que agreguen tu tarjeta a su colección.
    </p>
  </AppSheet>
</template>

<script setup lang="ts">
import { PhArrowsClockwise, PhDotOutline } from '@phosphor-icons/vue'
import { getApiV1AlbumIdShareAssignedCard } from '~/services/album/album'

const props = defineProps<{
  albumId: string
  isOpen: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()

const toast = useToast()
const profile = useProfile()
const imageUrl = ref<string | null>(null)
let eventSource: EventSource | null = null

const { data: qrImage, refresh: refreshQr } = useApiData(
  () => getApiV1AlbumIdShareAssignedCard(props.albumId, { qr: 't' }) as any,
  `qr-${props.albumId}`,
)

watch(qrImage, (blob) => {
  if (blob) {
    imageUrl.value = URL.createObjectURL(blob as any)
  }
}, { immediate: true })

watch(() => props.isOpen, (open) => {
  if (open) {
    connectQrEvents()
  } else {
    eventSource?.close()
    eventSource = null
  }
})

function connectQrEvents() {
  if (eventSource) return
  eventSource = new EventSource(`/api/v1/album/${props.albumId}/qr_events`)
  eventSource.addEventListener('message', async (event) => {
    try {
      const payload = JSON.parse(event.data)
      if (payload?.album_id === props.albumId) {
        toast.info({ title: 'QR escaneado', message: 'Tu tarjeta fue agregada; tu código se renovó.' })
        await refreshQr()
      }
    } catch { }
  })
}
</script>