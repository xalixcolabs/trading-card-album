<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-xs"
        @click.self="$emit('close')">
        <div
          class="relative w-full max-w-85 bg-base-100 text-base-content rounded-2xl p-6 shadow-2xl border border-base-300 flex flex-col items-center">

          <button @click="$emit('close')"
            class="absolute top-3 right-3 btn btn-xs btn-circle btn-ghost cursor-pointer text-base-content/70">✕</button>

          <div class="text-center">
            <h3 class="font-extrabold text-base uppercase tracking-wider text-primary">Mi Código QR</h3>
            <p class="text-xs text-base-content/60 font-medium mb-4 ">{{ profile ? profile.name : '' }}</p>
            <p class="text-xs text-base-content/60 font-bold">Codigo de un solo uso</p>
          </div>

          <div
            class="p-3 bg-white rounded-xl border border-slate-200/50 shadow-inner flex items-center justify-center mt-1 mb-3 relative overflow-hidden">
            <div class="absolute inset-2 border border-slate-100 pointer-events-none border-dashed"></div>
            <img v-if="imageUrl" :src="imageUrl" alt="Mi Código QR" class="w-50 h-50 select-none" />
          </div>
          <button @click="refreshQr()" class="btn btn-primary">Refrescar QR</button>
          <p class="text-center text-xs text-base-content/70 mt-3 font-medium leading-relaxed px-2">
            Muestra este código a otros desarrolladores del GDG para registrar su interacción.
          </p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { getApiV1AlbumIdShareAssignedCard } from '~/services/album/album';

const { albumId } = defineProps<{
  albumId: string,
  isOpen: boolean,
}>()

defineEmits<{
  (e: 'close'): void
}>()

const imageUrl = useState<string | null>()
const profile = useProfile()

const { data: qrImage, refresh: refreshQr } = useApiData(
  () => getApiV1AlbumIdShareAssignedCard(albumId, { qr: 't' }) as any,
  'getApiV1AlbumIdShareAssignedCard'
)

watch(qrImage, () => {
  imageUrl.value = URL.createObjectURL(qrImage.value as any)
})

</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.35s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
