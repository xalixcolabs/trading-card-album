import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: false,
  css: ['./app/assets/css/main.css'],
  app: {
    pageTransition: { name: 'page', mode: 'out-in' },
    // Los assets salen en /assets en vez de /_nuxt porque go:embed
    // excluye directorios que empiezan con "_".
    buildAssetsDir: 'assets',
  },
  devServer: {
    host: 'localhost',
    port: 3000,
  },
  runtimeConfig: {
    public: {
      // Mismo origen: en dev Nuxt proxya /api hacia el backend y en
      // producción Fiber sirve la SPA y el API juntos.
      apiBase: ''
    }
  },

  vite: {
    plugins: [
      tailwindcss(),
    ],
    server: {
      proxy: {
        '/api': 'http://localhost:8080',
      },
    },
  },

  modules: ['nuxt-toast', 'nuxt-qrcode'],
  toast: {
    settings: {
      position: 'bottomRight',
      timeout: 2500,
      progressBar: true,
      progressBarColor: '#2dd4bf',
      theme: 'dark',
      layout: 1,
      close: false,
    }
  }
})