<template>
  <AppSheet :is-open="isOpen" @close="$emit('close')">
    <div class="flex items-center justify-between pt-3">
      <div>
        <h3 class="text-lg font-bold tracking-tight text-ink">{{ mode === 'create' ? 'Nueva tarjeta' : 'Editar tarjeta' }}</h3>
        <p class="text-sm text-mist">{{ mode === 'create' ? 'Agrega una carta a un álbum.' : 'Actualiza los datos de la carta.' }}</p>
      </div>
      <div class="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-accent-soft text-accent">
        <PhCardholder :size="24" weight="duotone" />
      </div>
    </div>

    <form class="mt-5 flex flex-col gap-4" @submit.prevent="handleSubmit">
      <div v-if="mode === 'create' && !fixedAlbumId" class="flex flex-col gap-1.5">
        <label for="card-album" class="text-[13px] font-semibold text-mist">Álbum</label>
        <select id="card-album" v-model="albumId" required
          class="rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge focus:outline-none focus:ring-2 focus:ring-accent">
          <option value="" disabled>Selecciona un álbum</option>
          <option v-for="album in albums" :key="album.id" :value="album.id">{{ album.title }}</option>
        </select>
      </div>

      <div class="grid grid-cols-2 gap-3">
        <div class="flex flex-col gap-1.5">
          <label for="card-number" class="text-[13px] font-semibold text-mist">Número</label>
          <input id="card-number" v-model="number" type="text" required placeholder="01"
            class="rounded-xl bg-panel px-4 py-3 font-mono text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
        </div>
        <div class="flex flex-col gap-1.5">
          <label for="card-name" class="text-[13px] font-semibold text-mist">Nombre</label>
          <input id="card-name" v-model="name" type="text" required placeholder="Nombre de la carta"
            class="rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
        </div>
      </div>

      <div class="flex flex-col gap-1.5">
        <label for="card-image" class="text-[13px] font-semibold text-mist">URL de la imagen (.webp)</label>
        <input id="card-image" v-model="imageUrl" type="url" required placeholder="https://…/carta.webp"
          class="rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
      </div>

      <div class="flex flex-col gap-1.5">
        <label for="card-desc" class="text-[13px] font-semibold text-mist">Descripción</label>
        <textarea id="card-desc" v-model="description" rows="3" placeholder="Descripción de la carta"
          class="resize-none rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent"></textarea>
      </div>

      <button type="submit" :disabled="submitting"
        class="w-full rounded-2xl bg-accent px-4 py-3.5 text-sm font-semibold text-accent-ink shadow-glow transition-transform active:scale-[0.97] disabled:opacity-50">
        {{ submitting ? 'Guardando…' : (mode === 'create' ? 'Crear tarjeta' : 'Guardar cambios') }}
      </button>
    </form>
  </AppSheet>
</template>

<script setup lang="ts">
import { PhCardholder } from '@phosphor-icons/vue'
import { postApiV1AdminCards, putApiV1AdminCardsId } from '~/services/admin/admin'
import type { AdminDtoAlbum, CardModelCard } from '~/models'

const { isOpen, mode, card, albums, fixedAlbumId } = defineProps<{
  isOpen: boolean
  mode: 'create' | 'edit'
  card?: CardModelCard
  albums: AdminDtoAlbum[]
  fixedAlbumId?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const toast = useToast()
const albumId = ref('')
const number = ref('')
const name = ref('')
const description = ref('')
const imageUrl = ref('')
const submitting = ref(false)

watch(() => [card, fixedAlbumId] as const, ([value, fixedId]) => {
  number.value = value?.number ?? ''
  name.value = value?.name ?? ''
  description.value = value?.description ?? ''
  imageUrl.value = value?.image_url ?? ''
  albumId.value = value?.album_id ?? fixedId ?? ''
}, { immediate: true })

async function handleSubmit() {
  if (!name.value.trim() || !number.value.trim() || !imageUrl.value.trim()) {
    toast.error({ title: 'Datos incompletos', message: 'Completa número, nombre e imagen.' })
    return
  }
  submitting.value = true
  try {
    if (mode === 'create') {
      if (!albumId.value) {
        toast.error({ title: 'Falta el álbum', message: 'Selecciona un álbum.' })
        return
      }
      await postApiV1AdminCards({
        album_id: albumId.value,
        number: number.value.trim(),
        name: name.value.trim(),
        description: description.value.trim(),
        image_url: imageUrl.value.trim(),
      })
      toast.success({ title: 'Tarjeta creada', message: name.value.trim() })
    } else {
      await putApiV1AdminCardsId(card!.id!, {
        number: number.value.trim(),
        name: name.value.trim(),
        description: description.value.trim(),
        image_url: imageUrl.value.trim(),
      })
      toast.success({ title: 'Tarjeta actualizada', message: name.value.trim() })
    }
    emit('saved')
    emit('close')
  } catch {
    toast.error({ title: 'No se pudo guardar', message: 'Inténtalo de nuevo.' })
  } finally {
    submitting.value = false
  }
}
</script>