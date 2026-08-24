import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: false,
  css: ['./app/assets/css/main.css'],
  devServer: {
    host: 'localhost',
    port: 3000,
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.NODE_ENV === 'development' 
        ? 'http://localhost:8080' 
        : ''
    }
  },

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  modules: ['nuxt-toast', 'nuxt-qrcode'],
  toast: {
    settings: {
      position: 'topRight',
      timeout: 2000,
    }
  }
})