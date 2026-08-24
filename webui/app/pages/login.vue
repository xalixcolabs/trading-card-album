<template>
  <div class="grid place-content-center h-screen font-sans">
    <div class="card bg-base-100 shadow-xl border border-base-300 w-80">
      <div class="card-body p-6 flex flex-col gap-4 text-center">
        <div class="flex flex-col gap-1 items-center">
          <h2 class="text-lg font-black uppercase tracking-wider text-base-content">Trading Card Album</h2>
          <p class="text-xs text-base-content/60 font-medium">Inicia sesión para comenzar tu colección</p>
        </div>

        <a href="/api/v1/auth"
          class="btn btn-primary w-full font-bold uppercase tracking-wider cursor-pointer shadow-md mt-2 flex items-center justify-center gap-2">
          Ingresar con Google
        </a>
      </div>
    </div>
  </div>
</template>