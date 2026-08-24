<template>
  <div class="min-h-screen w-full flex justify-center items-start overflow-x-hidden">
    <div class="w-full max-w-120 min-h-screen shadow-2xl px-2 flex flex-col relative border-x border-base-200">
      <slot />
    </div>
  </div>
</template>

<style scoped>
@supports (min-height: 100dvh) {
  .min-h-screen {
    min-height: 100dvh;
  }
}
</style>
