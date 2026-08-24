<template>
  <AppSheet :is-open="isOpen" @close="onClose">
    <div class="flex items-center justify-between pt-3">
      <div>
        <h3 class="text-lg font-bold tracking-tight text-ink">Nuevo álbum</h3>
        <p class="text-sm text-mist">Crea la baraja de tarjetas del evento.</p>
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

      <div class="flex flex-col gap-3">
        <div class="flex items-center justify-between">
          <span class="text-[13px] font-semibold text-mist">Tarjetas</span>
          <button type="button" @click="addCard"
            class="inline-flex items-center gap-1 rounded-lg bg-accent-soft px-2.5 py-1.5 text-xs font-semibold text-accent transition-transform active:scale-95">
            <PhPlus :size="14" weight="bold" />
            Agregar
          </button>
        </div>

        <div v-if="cards.length === 0" class="rounded-2xl border border-dashed border-edge p-4 text-center text-sm text-faint">
          Agrega al menos una tarjeta a la baraja.
        </div>

        <div v-for="(card, index) in cards" :key="index"
          class="rounded-2xl bg-raise/60 p-3 ring-1 ring-edge-soft">
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

      <button type="submit" :disabled="submitting"
        class="w-full rounded-2xl bg-accent px-4 py-3.5 text-sm font-semibold text-accent-ink shadow-glow transition-transform active:scale-[0.97] disabled:opacity-50">
        {{ submitting ? 'Creando…' : 'Crear álbum' }}
      </button>
    </form>
  </AppSheet>
</template>

<script setup lang="ts">
import { PhPlus, PhStack, PhTrash } from '@phosphor-icons/vue'
import { postApiV1Album } from '~/services/album/album'

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created'): void
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

const onClose = () => {
  emit('close')
}

async function handleSubmit() {
  if (!title.value.trim()) {
    toast.error({ title: 'Falta el título', message: 'Escribe el nombre del álbum.' })
    return
  }
  if (cards.value.length === 0) {
    toast.error({ title: 'Sin tarjetas', message: 'Agrega al menos una tarjeta.' })
    return
  }
  submitting.value = true
  try {
    await postApiV1Album({
      title: title.value.trim(),
      cards: cards.value.map(card => ({
        number: card.number,
        name: card.name.trim(),
        description: card.description.trim(),
        image_url: card.image_url.trim(),
      })),
    })
    toast.success({ title: 'Álbum creado', message: 'Ya puedes compartir el código.' })
    title.value = ''
    cards.value = []
    emit('created')
    emit('close')
  } catch {
    toast.error({ title: 'No se pudo crear', message: 'Solo los organizadores pueden crear álbumes.' })
  } finally {
    submitting.value = false
  }
}
</script>