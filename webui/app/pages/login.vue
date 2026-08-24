<template>
  <div class="grain relative flex min-h-[100dvh] flex-col items-center justify-between px-8 py-10">
    <div class="flex-1"></div>

    <div class="flex flex-col items-center text-center">
      <div class="relative">
        <div
          class="flex h-20 w-20 items-center justify-center rounded-[1.75rem] bg-gradient-to-br from-accent to-teal-500 text-accent-ink shadow-glow ring-1 ring-inset ring-white/20">
          <PhCardholder :size="40" weight="duotone" />
        </div>
        <div
          class="absolute -right-2 -top-2 flex h-8 w-8 items-center justify-center rounded-full bg-accent text-accent-ink shadow-glow">
          <PhSparkle :size="17" weight="fill" />
        </div>
      </div>

      <h1 class="mt-7 text-3xl font-bold tracking-tight text-ink text-balance">
        Trading Card<br />Album
      </h1>
      <p class="mt-3 max-w-[32ch] text-[15px] leading-relaxed text-mist">
        Colecciona las tarjetas del DevFest, intercambia con otros y completa tu álbum escaneando sus códigos.
      </p>
    </div>

    <div class="mt-10 w-full max-w-[340px]">
      <a href="/api/v1/auth"
        class="flex w-full items-center justify-center gap-3 rounded-2xl bg-white px-5 py-4 text-[15px] font-semibold text-zinc-900 shadow-glow ring-1 ring-inset ring-black/5 transition-transform active:scale-[0.97]">
        <GoogleMark class="h-5 w-5" />
        Continuar con Google
      </a>
      <p class="mt-4 text-center text-xs leading-relaxed text-faint">
        Al continuar aceptas participar como coleccionista del evento.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { PhCardholder, PhSparkle } from '@phosphor-icons/vue'
import GoogleMark from '~/components/GoogleMark.vue'
</script>