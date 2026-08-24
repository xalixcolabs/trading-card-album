<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-xs"
        @click.self="$emit('close')">
        <div
          class="relative w-full max-w-85 bg-base-100 text-base-content rounded-2xl p-6 shadow-2xl border border-base-300 flex flex-col">

          <button @click="$emit('close')"
            class="absolute top-3 right-3 btn btn-xs btn-circle btn-ghost cursor-pointer text-base-content/70">✕</button>

          <div class="text-center mb-4">
            <h3 class="font-extrabold text-base uppercase tracking-wider text-secondary">Escanear QR</h3>
            <p class="text-xs text-base-content/60 font-medium">Registrar interacción</p>
          </div>

          <div
            class="relative w-full aspect-square bg-black rounded-xl overflow-hidden border-2 border-base-300 flex items-center justify-center shadow-inner">
            <QrcodeStream @error="onError" @detect="onDetect" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import type { DetectedBarcode } from 'nuxt-qrcode'
import { postApiV1AlbumNewCard } from '~/services/album/album';

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void,
  (e: 'card-added-successfully'): void,
}>()

const toast = useToast()

async function onDetect(detectedCodes: DetectedBarcode[]) {
  const qrPayload = detectedCodes.map(payload => payload.rawValue).join()
  const newCardPayload = JSON.parse(qrPayload)
  try {
    await postApiV1AlbumNewCard(newCardPayload)
    toast.success({
      title: 'Nueva tarjeta',
      message: 'agregada'
    })
    emit('card-added-successfully')
  } catch {
    toast.error({
      title: 'Nueva tarjeta',
      message: 'Error al agregar'
    })
  }
  emit('close')
}

function onError(err: Error) { }
</script>
