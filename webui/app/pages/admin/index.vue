<template>
  <AdminShell active="overview">
    <h2 class="text-lg font-bold tracking-tight text-ink">Resumen</h2>
    <p class="mt-1 text-sm text-mist">Estado general de la plataforma.</p>

    <div v-if="pending" class="mt-5 grid grid-cols-2 gap-3">
      <div v-for="i in 6" :key="i" class="h-24 animate-pulse rounded-2xl bg-raise ring-1 ring-edge"></div>
    </div>

    <div v-else class="mt-5 grid grid-cols-2 gap-3">
      <div v-for="metric in metrics" :key="metric.label"
        class="rounded-2xl bg-raise p-4 ring-1 ring-edge">
        <p class="font-mono text-2xl font-bold text-accent tabular">{{ metric.value }}</p>
        <p class="mt-1 text-xs text-mist">{{ metric.label }}</p>
      </div>
    </div>
  </AdminShell>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'admin' })
import { getApiV1AdminOverview } from '~/services/admin/admin'

const { data: overview, pending } = useApiData(() => getApiV1AdminOverview(), 'admin-overview')

const metrics = computed(() => [
  { label: 'Álbumes', value: overview.value?.albums ?? 0 },
  { label: 'Usuarios', value: overview.value?.users ?? 0 },
  { label: 'Tarjetas', value: overview.value?.cards ?? 0 },
  { label: 'Participantes', value: overview.value?.participants ?? 0 },
  { label: 'Contactos', value: overview.value?.contacts ?? 0 },
  { label: 'Tarjetas recolectadas', value: overview.value?.collected ?? 0 },
])
</script>