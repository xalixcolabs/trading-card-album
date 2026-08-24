<template>
  <AppSheet :is-open="isOpen" @close="$emit('close')">
    <div v-if="loading" class="flex flex-col gap-3 pt-4">
      <div v-for="i in 3" :key="i" class="h-16 animate-pulse rounded-2xl bg-raise ring-1 ring-edge"></div>
    </div>

    <div v-else-if="detail" class="pt-3">
      <div class="flex items-center gap-3">
        <div
          class="flex h-12 w-12 items-center justify-center rounded-2xl bg-raise text-lg font-bold text-accent ring-1 ring-edge">
          {{ initialOf(detail.user?.name) }}
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-[15px] font-semibold text-ink">{{ detail.user?.name || 'Sin nombre' }}</p>
          <p class="truncate text-[13px] text-mist">{{ detail.user?.email }}</p>
        </div>
        <span v-if="detail.user?.is_admin"
          class="shrink-0 rounded-full bg-accent-soft px-2.5 py-1 text-[11px] font-semibold text-accent">Admin</span>
      </div>

      <div class="mt-6">
        <h4 class="text-[11px] font-semibold uppercase tracking-[0.18em] text-faint">Álbumes ({{ detail.albums?.length ?? 0 }})</h4>
        <div v-if="detail.albums?.length" class="mt-3 flex flex-col gap-2">
          <div v-for="album in detail.albums" :key="album.id"
            class="flex items-center justify-between rounded-xl bg-raise px-4 py-3 ring-1 ring-edge-soft">
            <span class="truncate text-sm font-semibold text-ink">{{ album.title }}</span>
            <span class="font-mono text-xs text-mist tabular">{{ album.total_cards }}</span>
          </div>
        </div>
        <p v-else class="mt-3 text-sm text-faint">Sin álbumes.</p>
      </div>

      <div class="mt-6">
        <h4 class="text-[11px] font-semibold uppercase tracking-[0.18em] text-faint">Tarjetas recolectadas ({{ detail.cards?.length ?? 0 }})</h4>
        <div v-if="detail.cards?.length" class="mt-3 grid grid-cols-3 gap-3">
          <div v-for="card in detail.cards" :key="card.id" class="text-center">
            <div class="mx-auto aspect-2/3 w-full overflow-hidden rounded-lg bg-raise ring-1 ring-edge">
              <img :src="card.image_url" :alt="card.name" loading="lazy" class="h-full w-full object-cover" />
            </div>
            <p class="mt-1.5 truncate text-[11px] font-semibold text-mist">{{ card.name }}</p>
          </div>
        </div>
        <p v-else class="mt-3 text-sm text-faint">Aún no recolecta tarjetas.</p>
      </div>
    </div>
  </AppSheet>
</template>

<script setup lang="ts">
import { getApiV1AdminUsersId } from '~/services/admin/admin'
import type { AdminDtoUserDetail } from '~/models'

const props = defineProps<{
  isOpen: boolean
  userId: string
}>()

defineEmits<{
  (e: 'close'): void
}>()

const detail = ref<AdminDtoUserDetail | null>(null)
const loading = ref(false)

watch(() => props.userId, async (id) => {
  if (!id) {
    detail.value = null
    return
  }
  loading.value = true
  try {
    const response = await getApiV1AdminUsersId(id)
    detail.value = response.data
  } catch {
    detail.value = null
  } finally {
    loading.value = false
  }
}, { immediate: true })

function initialOf(name?: string) {
  return (name || '?').charAt(0).toUpperCase()
}
</script>