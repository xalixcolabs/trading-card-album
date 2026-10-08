<template>
  <div class="px-5 pt-4 pb-6">
    <!-- Encabezado -->
    <div class="flex items-center gap-2">
      <button type="button" @click="router.back()"
        class="flex h-10 w-10 items-center justify-center rounded-xl bg-raise text-mist ring-1 ring-edge transition-transform active:scale-90"
        aria-label="Volver">
        <PhArrowLeft :size="20" />
      </button>
      <h1 class="text-xl font-bold tracking-tight text-ink">Editar perfil</h1>
    </div>

    <!-- Avatar -->
    <div class="mt-6 flex flex-col items-center">
      <div class="h-20 w-20 overflow-hidden rounded-[1.75rem] ring-1 ring-edge">
        <UserAvatar :name="name" :picture="profile?.picture" size-class="text-3xl" />
      </div>
      <p class="mt-3 text-sm text-mist">Estos datos se muestran a quien escanea tu tarjeta.</p>
    </div>

    <!-- Formulario -->
    <form class="mt-7 flex flex-col gap-6" @submit.prevent="handleSubmit">
      <section class="flex flex-col gap-3">
        <h2 class="text-[11px] font-semibold uppercase tracking-[0.18em] text-faint">Información</h2>
        <div class="flex flex-col gap-3">
          <div class="flex flex-col gap-1.5">
            <label for="name" class="text-[13px] font-semibold text-mist">Nombre completo</label>
            <input id="name" v-model="name" type="text" required placeholder="Tu nombre"
              class="rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
          </div>
          <div class="flex flex-col gap-1.5">
            <label for="email" class="text-[13px] font-semibold text-mist">Correo (Google)</label>
            <input id="email" v-model="email" type="email" disabled readonly
              class="rounded-xl bg-raise/60 px-4 py-3 text-sm text-faint ring-1 ring-edge" />
          </div>
          <div class="flex items-center justify-between gap-3 rounded-xl bg-panel px-4 py-3 ring-1 ring-edge">
            <div class="min-w-0">
              <p class="text-[13px] font-semibold text-mist">Compartir mi correo</p>
              <p class="mt-0.5 text-[11px] text-faint">
                {{ publicEmail ? 'Tus contactos verán tu correo.' : 'Privado: verán lo que pongas abajo.' }}
              </p>
            </div>
            <button type="button" role="switch" :aria-checked="publicEmail" aria-label="Compartir mi correo"
              class="relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors"
              :class="publicEmail ? 'bg-accent' : 'bg-raise ring-1 ring-edge'"
              @click="publicEmail = !publicEmail">
              <span class="inline-block h-5 w-5 rounded-full bg-white shadow transition-transform"
                :class="publicEmail ? 'translate-x-[22px]' : 'translate-x-0.5'" />
            </button>
          </div>
          <div v-if="!publicEmail" class="flex flex-col gap-1.5">
            <label for="public_contact" class="text-[13px] font-semibold text-mist">Contacto público</label>
            <input id="public_contact" v-model="publicContact" type="text"
              placeholder="Otro correo, teléfono, @telegram…"
              class="rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
            <p class="text-[11px] text-faint">Lo verán tus contactos en lugar de tu correo.</p>
          </div>
        </div>
      </section>

      <section class="flex flex-col gap-3">
        <h2 class="text-[11px] font-semibold uppercase tracking-[0.18em] text-faint">Redes</h2>
        <div class="flex flex-col gap-3">
          <div class="flex flex-col gap-1.5">
            <label for="github" class="text-[13px] font-semibold text-mist">GitHub</label>
            <div class="relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 font-mono text-sm text-faint">@</span>
              <input id="github" v-model="github" type="text" placeholder="usuario"
                class="w-full rounded-xl bg-panel py-3 pl-8 pr-4 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
            </div>
          </div>
          <div class="flex flex-col gap-1.5">
            <label for="linkedin" class="text-[13px] font-semibold text-mist">LinkedIn</label>
            <div class="relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 font-mono text-sm text-faint">in/</span>
              <input id="linkedin" v-model="linkedin" type="text" placeholder="usuario-perfil"
                class="w-full rounded-xl bg-panel py-3 pl-10 pr-4 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
            </div>
          </div>
          <div class="flex flex-col gap-1.5">
            <label for="web" class="text-[13px] font-semibold text-mist">Sitio web</label>
            <input id="web" v-model="web" type="url" placeholder="https://tusitio.dev"
              class="rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent" />
          </div>
        </div>
      </section>

      <section class="flex flex-col gap-3">
        <h2 class="text-[11px] font-semibold uppercase tracking-[0.18em] text-faint">Sobre ti</h2>
        <div class="flex flex-col gap-1.5">
          <label for="description" class="text-[13px] font-semibold text-mist">Descripción</label>
          <textarea id="description" v-model="description" rows="4" placeholder="Un dato curioso o tu stack favorito…"
            class="resize-none rounded-xl bg-panel px-4 py-3 text-sm text-ink ring-1 ring-edge placeholder:text-faint focus:outline-none focus:ring-2 focus:ring-accent"></textarea>
        </div>
      </section>

      <!-- Guardar -->
      <button type="submit" :disabled="saving"
        class="mt-2 w-full rounded-2xl bg-accent px-4 py-3.5 text-sm font-semibold text-accent-ink shadow-glow transition-transform active:scale-[0.97] disabled:opacity-50">
        {{ saving ? 'Guardando…' : 'Guardar perfil' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { PhArrowLeft } from '@phosphor-icons/vue'
import { getApiV1AuthMe } from '~/services/auth/auth'
import { putApiV1UserId } from '~/services/user/user'

const router = useRouter()
const toast = useToast()
const saving = ref(false)

const { data: profile } = useApiData(() => getApiV1AuthMe(), 'profile')

const name = ref('')
const email = ref('')
const github = ref('')
const linkedin = ref('')
const web = ref('')
const description = ref('')
const publicEmail = ref(false)
const publicContact = ref('')

watch(profile, (value) => {
  if (!value) return
  name.value = value.name ?? ''
  email.value = value.email ?? ''
  github.value = value.github ?? ''
  linkedin.value = value.linkedin ?? ''
  web.value = value.web ?? ''
  description.value = value.description ?? ''
  publicEmail.value = (value.public_email ?? 0) === 1
  publicContact.value = value.public_contact ?? ''
}, { immediate: true })

async function handleSubmit() {
  if (!name.value.trim()) {
    toast.error({ title: 'Falta el nombre', message: 'Escribe tu nombre para continuar.' })
    return
  }
  saving.value = true
  try {
    await putApiV1UserId(profile.value!.id!, {
      name: name.value.trim(),
      email: email.value,
      github: github.value.trim(),
      linkedin: linkedin.value.trim(),
      web: web.value.trim(),
      description: description.value.trim(),
      public_email: publicEmail.value ? 1 : 0,
      public_contact: publicContact.value.trim(),
    })
    toast.success({ title: 'Perfil actualizado', message: 'Tus datos ya están guardados.' })
    await refreshNuxtData('profile')
    router.push('/')
  } catch {
    toast.error({ title: 'No se pudo guardar', message: 'Revisa tu conexión e inténtalo de nuevo.' })
  } finally {
    saving.value = false
  }
}
</script>