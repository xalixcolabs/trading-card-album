<template>
  <AppSheet :is-open="isOpen" @close="$emit('close')">
    <div class="pt-4 text-center">
      <div class="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-accent-soft text-accent">
        <PhScan :size="26" weight="fill" />
      </div>
      <h3 class="mt-3 text-lg font-bold tracking-tight text-ink">Escanear tarjeta</h3>
      <p class="mt-1 text-sm text-mist">Apunta al QR de otro desarrollador para sumar su tarjeta a tu colección.</p>
    </div>

    <div
      class="relative mt-5 aspect-square w-full overflow-hidden rounded-2xl border border-edge bg-black shadow-inner">
      <QrcodeStream @error="onError" @detect="onDetect" />
      <div class="pointer-events-none absolute inset-6 rounded-xl border-2 border-dashed border-accent/40"></div>
      <div class="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/40 to-transparent"></div>
    </div>

    <p v-if="scanning" class="mt-4 text-center text-xs text-mist">Registrando interacción…</p>
  </AppSheet>
</template>

<script setup lang="ts">
import { PhScan } from '@phosphor-icons/vue'
import type { DetectedBarcode } from 'nuxt-qrcode'
import { postApiV1AlbumNewCard } from '~/services/album/album'

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'card-added-successfully'): void
}>()

const toast = useToast()
const scanning = ref(false)

async function onDetect(detectedCodes: DetectedBarcode[]) {
  if (scanning.value) return
  scanning.value = true
  const qrPayload = detectedCodes.map(payload => payload.rawValue).join()
  try {
    const newCardPayload = JSON.parse(qrPayload)
    await postApiV1AlbumNewCard(newCardPayload)
    toast.success({
      title: 'Nueva tarjeta',
      message: 'Añadida a tu colección',
    })
    emit('card-added-successfully')
  } catch {
    toast.error({
      title: 'No se pudo registrar',
      message: 'Revisa que el código sea válido',
    })
  }
  scanning.value = false
  emit('close')
}

function onError(_err: Error) { }
</script>