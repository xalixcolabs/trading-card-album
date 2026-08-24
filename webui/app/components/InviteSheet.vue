<template>
  <AppSheet :is-open="isOpen" @close="$emit('close')">
    <div class="flex items-center justify-between pt-3">
      <div>
        <h3 class="text-lg font-bold tracking-tight text-ink">Invitar al álbum</h3>
        <p class="text-sm text-mist">Comparte el QR o el código para que otros se unan.</p>
      </div>
      <div class="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-accent-soft text-accent">
        <PhShareNetwork :size="24" weight="duotone" />
      </div>
    </div>

    <div class="mx-auto mt-5 flex aspect-square w-full max-w-[240px] items-center justify-center rounded-2xl bg-white p-4 ring-1 ring-edge">
      <img v-if="qrUrl" :src="qrUrl" alt="Código QR para unirse al álbum" class="h-full w-full select-none" />
      <div v-else class="h-24 w-24 animate-pulse rounded-xl bg-black/10"></div>
    </div>

    <div class="mt-5 flex items-center gap-2 rounded-xl bg-raise p-2 pl-4 ring-1 ring-edge">
      <span class="min-w-0 flex-1 truncate font-mono text-sm text-ink">{{ albumId }}</span>
      <button type="button" @click="copyCode"
        class="inline-flex shrink-0 items-center gap-1.5 rounded-lg bg-accent px-3 py-2 text-xs font-semibold text-accent-ink transition-transform active:scale-95">
        <PhCopy :size="14" weight="bold" />
        Copiar
      </button>
    </div>

    <p class="mt-4 text-center text-sm leading-relaxed text-mist">
      Quien escanee este código recibirá una tarjeta asignada y empezará su colección.
    </p>
  </AppSheet>
</template>

<script setup lang="ts">
import { PhCopy, PhShareNetwork } from '@phosphor-icons/vue'
import { getApiV1AlbumIdJoinQr } from '~/services/album/album'

const { albumId } = defineProps<{
  albumId: string
  isOpen: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()

const toast = useToast()
const qrUrl = ref<string | null>(null)

const { data: qrImage, refresh: refreshQr } = useApiData(
  () => getApiV1AlbumIdJoinQr(albumId) as any,
  `join-qr-${albumId}`,
)

watch(qrImage, (blob) => {
  if (blob) {
    qrUrl.value = URL.createObjectURL(blob as any)
  }
}, { immediate: true })

async function copyCode() {
  try {
    await navigator.clipboard.writeText(albumId)
    toast.success({ title: 'Código copiado', message: 'Comparte el código o el QR.' })
  } catch {
    toast.error({ title: 'No se pudo copiar', message: 'Copia el código manualmente.' })
  }
}
</script>