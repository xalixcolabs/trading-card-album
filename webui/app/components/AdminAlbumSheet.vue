<template>
  <AppSheet :is-open="isOpen" @close="$emit('close')">
    <div class="flex items-center justify-between pt-3">
      <div>
        <h3 class="text-lg font-bold tracking-tight text-ink">{{ mode === 'create' ? 'Nuevo álbum' : 'Editar álbum' }}</h3>
        <p class="text-sm text-mist">{{ mode === 'create' ? 'Crea la baraja del evento.' : 'Cambia el nombre del álbum.' }}</p>
      </div>
      <div class="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-accent-soft text-accent">
        <PhStack :size="24" weight="duotone" />
      </div>
    </div>

    <form class="mt-5 flex flex-col gap-5" @submit.prevent="handleSubmit">
      <div class="flex flex-col gap-1.5">
        <label for="album-title" class="text-[13px] font-semibold text-mist">Título del álbum</label>
        <input id="album-title" v-model="title" type="text" required placeholder="Ej. DevFest 2026"
          class="rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
      </div>

      <template v-if="mode === 'create'">
        <div class="flex flex-col gap-3">
          <div class="flex items-center justify-between">
            <span class="text-[13px] font-semibold text-mist">Tarjetas</span>
            <button type="button" @click="addCard"
              class="inline-flex items-center gap-1 rounded-lg bg-accent-soft px-2.5 py-1.5 text-xs font-semibold text-accent transition-transform active:scale-95">
              <PhPlus :size="14" weight="bold" />
              Agregar
            </button>
          </div>

          <p v-if="cards.length === 0" class="rounded-2xl border border-dashed border-edge p-4 text-center text-sm text-faint">
            Agrega al menos una tarjeta a la baraja.
          </p>

          <div v-for="(card, index) in cards" :key="index" class="rounded-2xl bg-raise/60 p-3 ring-1 ring-edge-soft">
            <div class="flex items-center justify-between gap-2">
              <span class="font-mono text-[11px] font-bold text-accent tabular">#{{ card.number }}</span>
              <button type="button" @click="removeCard(index)" aria-label="Quitar tarjeta"
                class="flex h-7 w-7 items-center justify-center rounded-lg bg-panel text-faint ring-1 ring-edge transition-transform active:scale-90">
                <PhTrash :size="14" />
              </button>
            </div>
            <div class="mt-2 flex flex-col gap-2">
              <input v-model="card.name" type="text" required placeholder="Nombre de la tarjeta"
                class="rounded-xl bg-panel px-3.5 py-2.5 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
              <input v-model="card.image_url" type="url" required placeholder="URL de la imagen (.webp)"
                class="rounded-xl bg-panel px-3.5 py-2.5 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
              <textarea v-model="card.description" rows="2" placeholder="Descripción (opcional)"
                class="resize-none rounded-xl bg-panel px-3.5 py-2.5 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent"></textarea>
            </div>
          </div>
        </div>
      </template>

      <button type="submit" :disabled="submitting"
        class="w-full rounded-2xl bg-accent px-4 py-3.5 text-sm font-semibold text-accent-ink shadow-glow transition-transform active:scale-[0.97] disabled:opacity-50">
        {{ submitting ? 'Guardando…' : (mode === 'create' ? 'Crear álbum' : 'Guardar cambios') }}
      </button>
    </form>
  </AppSheet>
</template>

<script setup lang="ts">
import { PhPlus, PhStack, PhTrash } from '@phosphor-icons/vue'
import { postApiV1AdminAlbums, putApiV1AdminAlbumsId } from '~/services/admin/admin'
import type { AdminDtoAlbum } from '~/models'

const { isOpen, mode, album } = defineProps<{
  isOpen: boolean
  mode: 'create' | 'edit'
  album?: AdminDtoAlbum
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const toast = useToast()
const title = ref('')
const submitting = ref(false)

interface NewCard {
  number: string
  name: string
  description: string
  image_url: string
}

const cards = ref<NewCard[]>([])

watch(() => album, (value) => {
  title.value = value?.title ?? ''
}, { immediate: true })

const addCard = () => {
  cards.value.push({
    number: String(cards.value.length + 1).padStart(2, '0'),
    name: '',
    description: '',
    image_url: '',
  })
}

const removeCard = (index: number) => {
  cards.value.splice(index, 1)
  cards.value.forEach((card, i) => {
    card.number = String(i + 1).padStart(2, '0')
  })
}

async function handleSubmit() {
  if (!title.value.trim()) {
    toast.error({ title: 'Falta el título', message: 'Escribe el nombre del álbum.' })
    return
  }
  if (mode === 'create' && cards.value.length === 0) {
    toast.error({ title: 'Sin tarjetas', message: 'Agrega al menos una tarjeta.' })
    return
  }
  submitting.value = true
  try {
    if (mode === 'create') {
      await postApiV1AdminAlbums({
        title: title.value.trim(),
        cards: cards.value.map(card => ({
          number: card.number,
          name: card.name.trim(),
          description: card.description.trim(),
          image_url: card.image_url.trim(),
        })),
      })
      toast.success({ title: 'Álbum creado', message: 'Ya puedes invitar participantes.' })
    } else {
      await putApiV1AdminAlbumsId(album!.id!, { title: title.value.trim() })
      toast.success({ title: 'Álbum actualizado', message: title.value.trim() })
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