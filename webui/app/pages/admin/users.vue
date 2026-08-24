<template>
  <AdminShell active="users">
    <h2 class="text-lg font-bold tracking-tight text-ink">Usuarios</h2>
    <p class="mt-1 text-sm text-mist">{{ users?.length ?? 0 }} desarrolladores registrados.</p>

    <div v-if="pending" class="mt-5 flex flex-col gap-3">
      <div v-for="i in 4" :key="i" class="h-16 animate-pulse rounded-2xl bg-raise ring-1 ring-edge"></div>
    </div>

    <div v-else class="mt-5 flex flex-col gap-3">
      <div v-for="user in users" :key="user.id"
        class="flex items-center gap-3 rounded-2xl bg-raise p-4 ring-1 ring-edge">
        <div
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-panel text-base font-bold text-accent ring-1 ring-edge">
          {{ initialOf(user.name) }}
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-[15px] font-semibold text-ink">{{ user.name || 'Sin nombre' }}</p>
          <p class="truncate text-xs text-mist">{{ user.email }}</p>
        </div>
        <button type="button" @click="toggleRole(user)" :disabled="isSelf(user) || toggling"
          class="shrink-0 rounded-lg px-3 py-2 text-xs font-semibold transition-colors disabled:opacity-50"
          :class="user.is_admin ? 'bg-accent-soft text-accent' : 'bg-panel text-mist ring-1 ring-edge'"
          :title="isSelf(user) ? 'No puedes quitarte tu propio rol' : ''">
          {{ user.is_admin ? 'Admin' : 'Usuario' }}
        </button>
      </div>

      <p v-if="users && users.length === 0" class="mt-4 text-center text-sm text-faint">Sin usuarios.</p>
    </div>
  </AdminShell>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'admin' })
import { getApiV1AdminUsers, putApiV1AdminUsersIdRole } from '~/services/admin/admin'
import type { AdminDtoUser } from '~/models'

const toast = useToast()
const toggling = ref(false)

const { data: users, pending, refresh } = useApiData(() => getApiV1AdminUsers(), 'admin-users')
const profile = useProfile()

function initialOf(name?: string) {
  return (name || '?').charAt(0).toUpperCase()
}

function isSelf(user: AdminDtoUser) {
  return user.id === profile.value?.id
}

async function toggleRole(user: AdminDtoUser) {
  toggling.value = true
  try {
    await putApiV1AdminUsersIdRole(user.id!, { is_admin: !user.is_admin })
    toast.success({ title: 'Rol actualizado', message: user.email || '' })
    await refresh()
  } catch {
    toast.error({ title: 'No se pudo actualizar', message: 'Inténtalo de nuevo.' })
  } finally {
    toggling.value = false
  }
}
</script>