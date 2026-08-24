<template>
  <AppSheet :is-open="isOpen" @close="$emit('close')">
    <div class="flex items-center justify-between pt-3">
      <div>
        <h3 class="text-lg font-bold tracking-tight text-ink">Unirse a un álbum</h3>
        <p class="text-sm text-mist">Escanea el QR del álbum o escribe su código.</p>
      </div>
      <div class="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-accent-soft text-accent">
        <PhTray :size="24" weight="duotone" />
      </div>
    </div>

    <div class="mt-5 grid grid-cols-2 gap-1 rounded-xl bg-raise p-1 ring-1 ring-edge">
      <button type="button" @click="mode = 'scan'"
        class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition-colors"
        :class="mode === 'scan' ? 'bg-accent text-accent-ink' : 'text-mist'">
        <PhScan :size="15" weight="bold" />
        Escanear QR
      </button>
      <button type="button" @click="mode = 'type'"
        class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition-colors"
        :class="mode === 'type' ? 'bg-accent text-accent-ink' : 'text-mist'">
        <PhKeyboard :size="15" weight="bold" />
        Escribir código
      </button>
    </div>

    <!-- Modo escanear -->
    <div v-if="mode === 'scan'" class="mt-5">
      <div class="relative aspect-square w-full overflow-hidden rounded-2xl border border-edge bg-black shadow-inner">
        <QrcodeStream @error="onError" @detect="onDetect" />
        <div class="pointer-events-none absolute inset-6 rounded-xl border-2 border-dashed border-accent/40"></div>
        <div class="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/40 to-transparent"></div>
      </div>
      <p v-if="joining" class="mt-4 text-center text-xs text-mist">Uniéndote al álbum…</p>
    </div>

    <!-- Modo escribir -->
    <form v-else class="mt-5 flex flex-col gap-4" @submit.prevent="joinWithCode">
      <div class="flex flex-col gap-1.5">
        <label for="album-code" class="text-[13px] font-semibold text-mist">Código del álbum</label>
        <input id="album-code" v-model="code" type="text" required autofocus placeholder="Pídelo al organizador"
          class="rounded-xl bg-panel px-4 py-3 font-mono text-sm tracking-wide text-ink ring-1 ring-edge placeholder:font-sans placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
      </div>

      <button type="submit" :disabled="joining"
        class="w-full rounded-2xl bg-accent px-4 py-3.5 text-sm font-semibold text-accent-ink shadow-glow transition-transform active:scale-[0.97] disabled:opacity-50">
        {{ joining ? 'Uniéndote…' : 'Unirme al álbum' }}
      </button>
    </form>
  </AppSheet>
</template>

<script setup lang="ts">
import { PhKeyboard, PhScan, PhTray } from '@phosphor-icons/vue'
import type { DetectedBarcode } from 'nuxt-qrcode'
import { postApiV1AlbumParticipant } from '~/services/album-participant/album-participant'

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'joined', albumId: string): void
}>()

const toast = useToast()
const mode = ref<'scan' | 'type'>('scan')
const code = ref('')
const joining = ref(false)

async function joinWithCode() {
  if (!code.value.trim()) {
    toast.error({ title: 'Falta el código', message: 'Ingresa el código del álbum.' })
    return
  }
  await join(code.value.trim())
}

async function onDetect(detectedCodes: DetectedBarcode[]) {
  if (joining.value) return
  const value = detectedCodes.map(payload => payload.rawValue).join('')
  await join(value)
}

async function join(albumId: string) {
  joining.value = true
  try {
    const response = await postApiV1AlbumParticipant({ album_id: albumId })
    toast.success({ title: '¡Ya estás dentro!', message: 'Revisa tu tarjeta asignada.' })
    code.value = ''
    emit('joined', response.data.album_id || albumId)
    emit('close')
  } catch {
    toast.error({ title: 'No se pudo unir', message: 'Verifica el código e inténtalo de nuevo.' })
  } finally {
    joining.value = false
  }
}

function onError(_err: Error) { }
</script>