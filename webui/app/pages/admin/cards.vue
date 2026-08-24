<template>
  <AdminShell active="cards">
    <h2 class="text-lg font-bold tracking-tight text-ink">Tarjetas</h2>
    <p class="mt-1 text-sm text-mist">{{ cards?.length ?? 0 }} tarjetas en circulación.</p>

    <div v-if="pending" class="mt-5 flex flex-col gap-3">
      <div v-for="i in 4" :key="i" class="h-16 animate-pulse rounded-2xl bg-raise ring-1 ring-edge"></div>
    </div>

    <div v-else class="mt-5 flex flex-col gap-3">
      <div v-for="card in cards" :key="card.id"
        class="flex items-center gap-3 rounded-2xl bg-raise p-3 ring-1 ring-edge">
        <div class="h-14 w-10 shrink-0 overflow-hidden rounded-lg bg-panel ring-1 ring-edge">
          <img :src="card.image_url" :alt="card.name" loading="lazy" class="h-full w-full object-cover" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-[15px] font-semibold text-ink">{{ card.name }}</p>
          <p class="font-mono text-xs text-mist tabular">#{{ card.number }}</p>
        </div>
        <button type="button" @click="removeCard(card)"
          class="shrink-0 rounded-lg px-3 py-2 text-xs font-semibold transition-colors"
          :class="confirmingId === card.id ? 'bg-danger text-white' : 'bg-panel text-faint ring-1 ring-edge'"
          :disabled="deleting">
          {{ confirmingId === card.id ? '¿Eliminar?' : 'Eliminar' }}
        </button>
      </div>

      <p v-if="cards && cards.length === 0" class="mt-4 text-center text-sm text-faint">Aún no hay tarjetas.</p>
    </div>
  </AdminShell>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'admin' })
import { getApiV1AdminCards, deleteApiV1AdminCardsId } from '~/services/admin/admin'
import type { CardModelCard } from '~/models'

const toast = useToast()
const confirmingId = ref<string | null>(null)
const deleting = ref(false)

const { data: cards, pending, refresh } = useApiData(() => getApiV1AdminCards(), 'admin-cards')

function removeCard(card: CardModelCard) {
  if (confirmingId.value !== card.id) {
    confirmingId.value = card.id
    setTimeout(() => {
      if (confirmingId.value === card.id) confirmingId.value = null
    }, 2500)
    return
  }
  performDelete(card)
}

async function performDelete(card: CardModelCard) {
  deleting.value = true
  try {
    await deleteApiV1AdminCardsId(card.id!)
    toast.success({ title: 'Tarjeta eliminada', message: card.name || '' })
    confirmingId.value = null
    await refresh()
  } catch {
    toast.error({ title: 'No se pudo eliminar', message: 'Inténtalo de nuevo.' })
  } finally {
    deleting.value = false
  }
}
</script>