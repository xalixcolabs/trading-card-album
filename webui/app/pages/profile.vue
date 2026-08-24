<template>
  <div class="p-4 flex flex-col gap-6 font-sans">
    <!-- Header Avatar Section -->
    <div class="mt-8 flex flex-col items-center gap-3">
      <p class="text-xl text-base-content/60 font-semibold uppercase tracking-wider text-center">
        Perfil de Coleccionista
      </p>
      <p class="text-xs text-center text-base-content/70 px-6 max-w-sm">
        Completa tu perfil para que otros desarrolladores puedan ver tus redes al interactuar contigo.
      </p>
    </div>

    <!-- Form Sections -->
    <div class="flex flex-col gap-4">

      <!-- Personal Information Card -->
      <div class="card bg-base-100 shadow-md border border-base-300">
        <div class="card-body p-5 flex flex-col gap-3.5">
          <h3 class="font-extrabold text-sm uppercase tracking-wider text-primary">Información Personal</h3>

          <!-- Name Input -->
          <div class="flex flex-col gap-1">
            <label class="text-[10px] font-bold text-base-content/60 uppercase tracking-wider">Nombre Completo</label>
            <input type="text" v-model="name" placeholder="Ingresa tu nombre"
              class="input input-bordered w-full bg-base-200" required />
          </div>

          <!-- Email Input (Prefilled from Google Login - Readonly) -->
          <div class="flex flex-col gap-1">
            <label class="text-[10px] font-bold text-base-content/60 uppercase tracking-wider">Correo Electrónico
              (Google)</label>
            <input type="email" v-model="email"
              class="input input-bordered w-full bg-base-300 text-base-content/50 cursor-not-allowed" readonly
              disabled />
          </div>
        </div>
      </div>

      <!-- Social Networks Card -->
      <div class="card bg-base-100 shadow-md border border-base-300">
        <div class="card-body p-5 flex flex-col gap-3.5">
          <h3 class="font-extrabold text-sm uppercase tracking-wider text-secondary">Redes Sociales</h3>

          <!-- Github Input -->
          <div class="flex flex-col gap-1">
            <label class="text-[10px] font-bold text-base-content/60 uppercase tracking-wider">Usuario GitHub</label>
            <div class="relative flex items-center">
              <span class="absolute left-3 text-lg text-base-content/50 font-mono">@</span>
              <input type="text" v-model="github" placeholder="usuario"
                class="input input-bordered w-full bg-base-200" />
            </div>
          </div>

          <!-- Linkedin Input -->
          <div class="flex flex-col gap-1">
            <label class="text-[10px] font-bold text-base-content/60 uppercase tracking-wider">Usuario LinkedIn</label>
            <div class="relative flex items-center">
              <span class="absolute left-3 text-[10px] text-base-content/50 font-mono">in/</span>
              <input type="text" v-model="linkedin" placeholder="usuario-perfil"
                class="input input-bordered w-full bg-base-200" />
            </div>
          </div>

          <!-- Web URL Input -->
          <div class="flex flex-col gap-1">
            <label class="text-[10px] font-bold text-base-content/60 uppercase tracking-wider">Sitio Web /
              Portafolio</label>
            <input type="url" v-model="web" placeholder="https://tuweb.com"
              class="input input-bordered w-full bg-base-200" />
          </div>
        </div>
      </div>

      <div class="card bg-base-100 shadow-md border border-base-300">
        <div class="card-body p-5 flex flex-col gap-3.5">
          <h3 class="font-extrabold text-sm uppercase tracking-wider text-secondary">Sobre Mí</h3>
          <div class="flex flex-col gap-1">
            <label class="text-[10px] font-bold text-base-content/60 uppercase tracking-wider">Descripción o Dato
              Curioso</label>
            <textarea v-model="description" placeholder="Cuéntanos algo interesante sobre ti o tu experiencia dev..."
              class="textarea textarea-bordered w-full h-24 bg-base-200"></textarea>
          </div>
        </div>
      </div>

      <!-- Action Button -->
      <button @click="handleSubmit"
        class="btn btn-primary mb-10 w-full shadow-lg font-bold uppercase tracking-wider cursor-pointer">
        Guardar Perfil
      </button>
    </div>
  </div>
</template>

<script setup>
import { getApiV1AuthMe } from '~/services/auth/auth';
import { putApiV1UserId } from '~/services/user/user';

const { data: profile } = useApiData(() => getApiV1AuthMe(), 'profile')

const router = useRouter()
const toast = useToast()

const name = useState(() => profile.value.name ?? '')
const email = useState(() => profile.value.email ?? 'usuario@gmail.com')
const github = useState(() => profile.value.github ?? '')
const linkedin = useState(() => profile.value.linkedin ?? '')
const web = useState(() => profile.value.web ?? '')
const description = useState(() => profile.value.description ?? '')

watch(profile, () => {
  name.value = profile.value.name
  email.value = profile.value.email
  github.value = profile.value.github
  linkedin.value = profile.value.linkedin
  web.value = profile.value.web
  description.value = profile.value.description
})

const handleSubmit = async () => {
  if (!name.value.trim()) {
    toast.error({
      title: 'Error',
      message: 'Por favor, ingresa tu nombre.'
    })
    return
  }

  try {
    const newData = {
      name: name.value,
      email: email.value,
      github: github.value,
      linkedin: linkedin.value,
      web: web.value,
      description: description.value,
    }
    await putApiV1UserId(profile.value.id, newData)
    toast.success({
      title: "OK",
      message: "Perfil actualizado con exito  "
    })
    await refreshNuxtData('profile')
    navigateTo('/')
  } catch (err) {
    console.error('Error saving profile changes:', err)
    toast.error({
      title: 'Error',
      message: 'Error al guardar el perfil. Intenta de nuevo.'
    })
  }
}
</script>