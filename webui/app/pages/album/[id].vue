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

                <div class="mt-4 pt-3 border-t border-base-300 flex flex-col gap-1">
                    <p class="text-[10px] font-bold text-base-content/50 uppercase tracking-wider">Portador de:</p>
                    <div v-if="assignedCard" class="w-1/2 mx-auto">
                        <Card :card="assignedCard" />
                    </div>
                </div>
            </div>
        </div>

        <div class="card bg-base-100 shadow mt-6 border border-base-300">
            <div class="card-body p-5">
                <h1 class="card-title text-base font-black uppercase tracking-wider text-base-content mb-3">
                    Tarjetas
                </h1>
                <div class="grid grid-cols-2 gap-4">
                    <Card v-for="card in cards" :card="card" :key="card.id" />
                </div>
            </div>
        </div>

        <FloatingMenu :albumId="albumId" @card-added-successfully="cardAddedSuccessfully" />
    </div>
</template>

<script setup>
import { getApiV1AlbumIdAssignedCard, getApiV1AlbumIdCard } from '~/services/album/album';
import { getApiV1AuthMe } from '~/services/auth/auth';

const route = useRoute()
const albumId = route.params.id

const { data: profile } = useApiData(() => getApiV1AuthMe(), 'profile')
const { data: assignedCard } = useApiData(() => getApiV1AlbumIdAssignedCard(albumId))
const { data: cards, refresh: refreshCards } = useApiData(() => getApiV1AlbumIdCard(albumId))

watch(profile, () => useProfile(profile.value))

const cardAddedSuccessfully = () => {
    refreshCards()
}
</script>
