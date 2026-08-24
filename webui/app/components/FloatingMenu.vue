<template>
  <div
    class="fixed bottom-6 left-1/2 -translate-x-1/2 w-full max-w-120 z-40 flex justify-center px-6 pointer-events-none">
    <div
      class="flex items-center gap-3 bg-base-100/90 backdrop-blur-md p-2 rounded-full shadow-2xl border border-base-300 pointer-events-auto">

      <button @click="openScan"
        class="btn btn-sm rounded-full bg-primary text-primary-content hover:scale-102 active:scale-98 flex items-center gap-2 px-4 font-sans font-bold cursor-pointer border-none shadow-md transition-all">
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#fff">
          <path
            d="M80-680v-200h200v80H160v120H80Zm0 600v-200h80v120h120v80H80Zm600 0v-80h120v-120h80v200H680Zm120-600v-120H680v-80h200v200h-80ZM700-260h60v60h-60v-60Zm0-120h60v60h-60v-60Zm-60 60h60v60h-60v-60Zm-60 60h60v60h-60v-60Zm-60-60h60v60h-60v-60Zm120-120h60v60h-60v-60Zm-60 60h60v60h-60v-60Zm-60-60h60v60h-60v-60Zm240-320v240H520v-240h240ZM440-440v240H200v-240h240Zm0-320v240H200v-240h240Zm-60 500v-120H260v120h120Zm0-320v-120H260v120h120Zm320 0v-120H580v120h120Z" />
        </svg>
        <span class="text-xs">Escanear</span>
      </button>

      <div class="h-5 w-px bg-base-300"></div>

      <button @click="openQr"
        class="btn btn-sm rounded-full bg-secondary text-secondary-content hover:scale-102 active:scale-98 flex items-center gap-2 px-4 font-sans font-bold cursor-pointer border-none shadow-md transition-all">
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#fff">
          <path
            d="M120-520v-320h320v320H120Zm80-80h160v-160H200v160Zm-80 480v-320h320v320H120Zm80-80h160v-160H200v160Zm320-320v-320h320v320H520Zm80-80h160v-160H600v160Zm160 480v-80h80v80h-80ZM520-360v-80h80v80h-80Zm80 80v-80h80v80h-80Zm-80 80v-80h80v80h-80Zm80 80v-80h80v80h-80Zm80-80v-80h80v80h-80Zm0-160v-80h80v80h-80Zm80 80v-80h80v80h-80Z" />
        </svg>
        <span class="text-xs">Mi QR</span>
      </button>
    </div>

    <MyQrModal :is-open="isQrOpen" :albumId="albumId" @close="closeQr" />
    <ScanQrModal :is-open="isScanOpen" @close="closeScan" @card-added-successfully="cardAddedSuccessfully" />
  </div>
</template>

<script setup lang="ts">

const emit = defineEmits<{
  (e: 'card-added-successfully'): void,
}>()

const { albumId } = defineProps<{
  albumId: string
}>()

const isQrOpen = useState(() => false)
const isScanOpen = useState(() => false)

const openQr = () => {
  isQrOpen.value = true
}

const closeQr = () => {
  isQrOpen.value = false
}

const openScan = () => {
  isScanOpen.value = true
}

const closeScan = () => {
  isScanOpen.value = false
}

const cardAddedSuccessfully = () => {
  emit('card-added-successfully')
}
</script>
