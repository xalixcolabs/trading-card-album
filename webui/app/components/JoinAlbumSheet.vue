<template>
  <AppSheet :is-open="isOpen" @close="$emit('close')">
    <div class="flex items-center justify-between pt-3">
      <div>
        <h3 class="text-lg font-bold tracking-tight text-ink">Unirse a un álbum</h3>
        <p class="text-sm text-mist">Empieza a coleccionar las tarjetas del evento.</p>
      </div>
      <div class="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-accent-soft text-accent">
        <PhTray :size="24" weight="duotone" />
      </div>
    </div>

    <form class="mt-5 flex flex-col gap-4" @submit.prevent="handleSubmit">
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
import { PhTray } from '@phosphor-icons/vue'
import { postApiV1AlbumParticipant } from '~/services/album-participant/album-participant'

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'joined', albumId: string): void
}>()

const toast = useToast()
const code = ref('')
const joining = ref(false)

async function handleSubmit() {
  if (!code.value.trim()) {
    toast.error({ title: 'Falta el código', message: 'Ingresa el código del álbum.' })
    return
  }
  joining.value = true
  try {
    const response = await postApiV1AlbumParticipant({ album_id: code.value.trim() })
    const albumId = response.data.album_id
    toast.success({ title: '¡Ya estás dentro!', message: 'Revisa tu tarjeta asignada.' })
    code.value = ''
    emit('joined', albumId)
    emit('close')
  } catch {
    toast.error({ title: 'No se pudo unir', message: 'Verifica el código e inténtalo de nuevo.' })
  } finally {
    joining.value = false
  }
}
</script>