<template>
    <div class="flex flex-col min-h-full pb-20 font-sans">
        <div class="card bg-base-100 shadow mt-10 border border-base-300">
            <div class="card-body p-5">
                <div class="flex justify-between items-start">
                    <div>
                        <h1 class="card-title text-xl font-black text-base-content">
                            {{ profile ? profile.name : '' }}
                        </h1>
                        <p class="text-xs text-base-content/60 font-medium mt-0.5">
                            {{ profile ? profile.email : '' }}
                        </p>
                    </div>
                    <NuxtLink to="/profile" class="btn btn-xs btn-ghost text-primary font-bold">Editar Perfil</NuxtLink>
                </div>
            </div>
        </div>

        <div class="card bg-base-100 shadow mt-6 border border-base-300">
            <div class="card-body p-5">
                <h1 class="card-title text-base font-black uppercase tracking-wider text-base-content mb-3">
                    Álbums
                </h1>
                <div class="grid grid-cols-2 gap-4">
                    <NuxtLink v-for="album in albums" :key="album.id" :to="`/album/${album.id}`"
                        class="group relative aspect-2/3 w-full rounded-2xl overflow-hidden shadow-md cursor-pointer transition-all duration-300 hover:scale-[1.03] hover:shadow-xl active:scale-[0.98] border border-base-200/80 bg-base-200 text-base-content">
                        <img :src="`/album.webp`" :alt="album.title"
                            class="w-full h-full object-cover select-none transition-transform duration-500 group-hover:scale-105" />
                        <div
                            class="absolute inset-0 bg-linear-to-t from-black/80 via-black/20 to-transparent opacity-60 group-hover:opacity-40 transition-opacity">
                        </div>

                        <div
                            class="absolute bottom-2 left-2 right-2 flex justify-between items-center bg-black/60 backdrop-blur-xs py-1 px-2.5 rounded-lg border border-white/10">
                            <span class="text-xs font-semibold text-primary-content truncate max-w-[80%]">
                                {{ album.title }}
                            </span>
                        </div>
                    </NuxtLink>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { getApiV1Album } from '~/services/album/album';
import { getApiV1AuthMe } from '~/services/auth/auth';

const { data: profile } = useApiData(() => getApiV1AuthMe(), 'profile')
const { data: albums } = useApiData(() => getApiV1Album())

watch(profile, () => useProfile(profile.value))
</script>
