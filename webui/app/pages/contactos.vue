<template>
  <div class="px-5 pt-6">
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-ink">Contactos</h1>
      <p class="mt-1 text-sm text-mist">Los desarrolladores con los que has intercambiado tarjetas.</p>
    </div>

    <div class="mt-6">
      <!-- Loading -->
      <div v-if="pending" class="flex flex-col gap-3">
        <div v-for="i in 4" :key="i" class="flex animate-pulse items-center gap-3 rounded-2xl bg-raise p-4 ring-1 ring-edge">
          <div class="h-12 w-12 rounded-2xl bg-panel"></div>
          <div class="flex-1 space-y-2">
            <div class="h-3 w-2/3 rounded-full bg-panel"></div>
            <div class="h-2.5 w-1/2 rounded-full bg-panel"></div>
          </div>
        </div>
      </div>

      <!-- Lista -->
      <div v-else-if="contacts && contacts.length" class="flex flex-col gap-3">
        <div v-for="contact in contacts" :key="contact.user_id"
          class="overflow-hidden rounded-2xl bg-raise ring-1 ring-edge transition-transform active:scale-[0.99]">
          <div class="flex items-center gap-3 p-4">
            <div class="h-12 w-12 shrink-0 overflow-hidden rounded-2xl ring-1 ring-edge">
              <UserAvatar :name="contact.name" :picture="contact.picture" size-class="text-lg" />
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-[15px] font-semibold text-ink">{{ contact.name || 'Desarrollador' }}</p>
              <p class="truncate text-[13px] text-mist">{{ contact.email }}</p>
              <p class="mt-0.5 text-[11px] text-faint">Conocido el {{ dateOf(contact.scanned_at) }}</p>
            </div>
          </div>

          <div class="flex flex-wrap gap-2 border-t border-edge-soft px-4 py-3">
            <a v-if="contact.github" :href="`https://github.com/${contact.github}`" target="_blank" rel="noopener"
              class="inline-flex items-center gap-1.5 rounded-lg bg-panel px-2.5 py-1.5 text-xs font-semibold text-mist ring-1 ring-edge transition-transform active:scale-95">
              <PhGithubLogo :size="14" /> {{ contact.github }}
            </a>
            <a v-if="contact.linkedin" :href="`https://linkedin.com/in/${contact.linkedin}`" target="_blank" rel="noopener"
              class="inline-flex items-center gap-1.5 rounded-lg bg-panel px-2.5 py-1.5 text-xs font-semibold text-mist ring-1 ring-edge transition-transform active:scale-95">
              <PhLinkedinLogo :size="14" /> {{ contact.linkedin }}
            </a>
            <a v-if="contact.web" :href="contact.web" target="_blank" rel="noopener"
              class="inline-flex items-center gap-1.5 rounded-lg bg-panel px-2.5 py-1.5 text-xs font-semibold text-mist ring-1 ring-edge transition-transform active:scale-95">
              <PhGlobe :size="14" /> Sitio web
            </a>
            <span v-if="!contact.github && !contact.linkedin && !contact.web"
              class="text-xs italic text-faint">Aún sin redes públicas</span>
          </div>
        </div>
      </div>

      <!-- Vacío -->
      <div v-else class="overflow-hidden rounded-3xl border border-dashed border-edge bg-raise/60 p-6 text-center">
        <div class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-accent-soft text-accent">
          <PhUsers :size="26" weight="duotone" />
        </div>
        <h2 class="mt-4 text-lg font-bold text-ink">Aún no tienes contactos</h2>
        <p class="mx-auto mt-1 max-w-[34ch] text-sm leading-relaxed text-mist">
          Escanea el QR de otro desarrollador y su perfil aparecerá aquí para que sigas en contacto.
        </p>
        <button type="button" @click="openScanner"
          class="mt-5 inline-flex items-center gap-2 rounded-xl bg-accent px-5 py-3 text-sm font-semibold text-accent-ink shadow-glow transition-transform active:scale-95">
          <PhScan :size="18" weight="fill" />
          Escanear ahora
        </button>
      </div>
    </div>

    <ScanQrModal :is-open="scannerOpen" @close="scannerOpen = false" @card-added-successfully="onScanned" />
  </div>
</template>

<script setup lang="ts">
import { PhGlobe, PhGithubLogo, PhLinkedinLogo, PhScan, PhUsers } from '@phosphor-icons/vue'
import { getApiV1Contact } from '~/services/contact/contact'

const scannerOpen = ref(false)
const { data: contacts, pending, refresh: refreshContacts } = useApiData(() => getApiV1Contact(), 'contacts')

function dateOf(unix?: number) {
  if (!unix) return '—'
  return new Date(unix * 1000).toLocaleDateString('es-MX', { day: 'numeric', month: 'short', year: 'numeric' })
}

const openScanner = () => { scannerOpen.value = true }

const onScanned = async () => {
  await refreshContacts()
}
</script>