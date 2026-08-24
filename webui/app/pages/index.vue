<template>
  <div class="px-5 pt-6">
    <!-- Perfil -->
    <div class="flex items-center gap-3">
      <div class="h-12 w-12 shrink-0 overflow-hidden rounded-2xl ring-1 ring-edge">
        <UserAvatar :name="profile?.name" :picture="profile?.picture" size-class="text-lg" />
      </div>
      <div class="min-w-0 flex-1">
        <p class="truncate text-[15px] font-semibold text-ink">{{ profile?.name || 'Coleccionista' }}</p>
        <p class="truncate text-[13px] text-mist">{{ profile?.email }}</p>
      </div>
      <NuxtLink to="/profile"
        class="flex h-9 w-9 items-center justify-center rounded-xl bg-raise text-mist ring-1 ring-edge transition-transform active:scale-90">
        <PhSlidersHorizontal :size="18" />
      </NuxtLink>
      <NuxtLink v-if="profile?.is_admin" to="/admin"
        class="flex h-9 w-9 items-center justify-center rounded-xl bg-accent-soft text-accent ring-1 ring-accent/30 transition-transform active:scale-90"
        aria-label="Administración">
        <PhShieldStar :size="18" weight="fill" />
      </NuxtLink>
    </div>

    <!-- Álbumes -->
    <div class="mt-8">
      <div class="flex items-end justify-between">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-ink">Tus álbumes</h1>
          <p class="mt-1 text-sm text-mist">El DevFest trae una colección nueva por evento.</p>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="pending" class="mt-5 flex gap-4 overflow-hidden">
        <div v-for="i in 2" :key="i" class="w-36 shrink-0 animate-pulse">
          <div class="aspect-2/3 w-full rounded-card bg-raise ring-1 ring-edge"></div>
          <div class="mt-2 h-3 w-24 rounded-full bg-raise"></div>
        </div>
      </div>

      <!-- Shelf -->
      <div v-else-if="albums && albums.length" class="no-scrollbar -mx-5 mt-5 flex snap-x snap-mandatory gap-4 overflow-x-auto px-5 pb-2">
        <NuxtLink v-for="album in albums" :key="album.id" :to="`/album/${album.id}`"
          class="w-36 shrink-0 snap-start">
          <div
            class="group relative aspect-2/3 w-full overflow-hidden rounded-card bg-raise ring-1 ring-edge transition-all duration-300 hover:ring-accent/50 active:scale-[0.97]">
            <img :src="'/album.webp'" :alt="album.title"
              class="h-full w-full select-none object-cover transition-transform duration-500 group-hover:scale-105" />
            <div class="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/70 via-black/10 to-transparent"></div>
            <div class="absolute bottom-2 left-2 right-2">
              <span class="font-mono text-[10px] font-bold text-accent tabular">{{ album.totalCards ?? 0 }} tarjetas</span>
              <p class="truncate text-[13px] font-semibold text-white">{{ album.title }}</p>
            </div>
          </div>
        </NuxtLink>
      </div>

      <!-- Empty -->
      <div v-else class="mt-6 overflow-hidden rounded-3xl border border-dashed border-edge bg-raise/60 p-6 text-center">
        <div class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-accent-soft text-accent">
          <PhStack :size="26" weight="duotone" />
        </div>
        <h2 class="mt-4 text-lg font-bold text-ink">Aún no estás en un álbum</h2>
        <p class="mx-auto mt-1 max-w-[34ch] text-sm leading-relaxed text-mist">
          Pide el código del álbum al organizador y únete para recibir tu tarjeta asignada.
        </p>

        <form class="mt-5 flex gap-2" @submit.prevent="joinAlbum">
          <input v-model="joinCode" type="text" inputmode="text" placeholder="Código del álbum"
            class="min-w-0 flex-1 rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
          <button type="submit" :disabled="joining"
            class="shrink-0 rounded-xl bg-accent px-4 py-3 text-sm font-semibold text-accent-ink shadow-glow transition-transform active:scale-95 disabled:opacity-50">
            {{ joining ? '…' : 'Unirme' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { PhShieldStar, PhStack, PhSlidersHorizontal } from '@phosphor-icons/vue'
import { getApiV1AuthMe } from '~/services/auth/auth'
import { getApiV1Album } from '~/services/album/album'
import { postApiV1AlbumParticipant } from '~/services/album-participant/album-participant'

const toast = useToast()
const joinCode = ref('')
const joining = ref(false)

const { data: profile } = useApiData(() => getApiV1AuthMe(), 'profile')
const { data: albums, pending, refresh: refreshAlbums } = useApiData(() => getApiV1Album(), 'albums')

watch(profile, (value) => {
  if (value) useProfile(value)
}, { immediate: true })

async function joinAlbum() {
  if (!joinCode.value.trim()) {
    toast.error({ title: 'Falta el código', message: 'Ingresa el código del álbum.' })
    return
  }
  joining.value = true
  try {
    await postApiV1AlbumParticipant({ album_id: joinCode.value.trim() })
    toast.success({ title: '¡Ya estás dentro!', message: 'Revisa tu tarjeta asignada.' })
    joinCode.value = ''
    await refreshAlbums()
  } catch {
    toast.error({ title: 'No se pudo unir', message: 'Verifica el código e inténtalo de nuevo.' })
  } finally {
    joining.value = false
  }
}
</script>