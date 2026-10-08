import tailwindcss from "@tailwindcss/vite";

// URL pública usada para las meta tags absolutas (OpenGraph/Twitter). Se
// resuelve en build time; en Docker se reemplaza con el build-arg
// NUXT_PUBLIC_SITE_URL (ver Dockerfile.full).
const siteUrl = (process.env.NUXT_PUBLIC_SITE_URL || 'http://localhost:8080').replace(/\/+$/, '')
const siteTitle = 'Trading Card Album'
const siteDescription = 'Colecciona, intercambia y comparte tarjetas. Recibe una tarjeta al azar, compártela por QR y desbloquea las de otros participantes en tu colección.'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: false,
  css: ['./app/assets/css/main.css'],
  app: {
    head: {
      title: siteTitle,
      htmlAttrs: { lang: 'es' },
      meta: [
        { name: 'description', content: siteDescription },
        { name: 'theme-color', content: '#0b0d11' },
        { name: 'apple-mobile-web-app-capable', content: 'yes' },
        { name: 'apple-mobile-web-app-status-bar-style', content: 'black-translucent' },
        // OpenGraph
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: siteTitle },
        { property: 'og:title', content: siteTitle },
        { property: 'og:description', content: siteDescription },
        { property: 'og:url', content: siteUrl },
        { property: 'og:locale', content: 'es_MX' },
        { property: 'og:image', content: `${siteUrl}/opengraph.png` },
        { property: 'og:image:secure_url', content: `${siteUrl}/opengraph.png` },
        { property: 'og:image:width', content: '1200' },
        { property: 'og:image:height', content: '630' },
        { property: 'og:image:type', content: 'image/png' },
        { property: 'og:image:alt', content: siteTitle },
        // Twitter
        { name: 'twitter:card', content: 'summary_large_image' },
        { name: 'twitter:title', content: siteTitle },
        { name: 'twitter:description', content: siteDescription },
        { name: 'twitter:image', content: `${siteUrl}/opengraph.png` },
        { name: 'twitter:image:alt', content: siteTitle },
      ],
      link: [
        { rel: 'icon', type: 'image/png', href: '/favicon.png' }
      ]
    },
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
      apiBase: '',
      // Dominio público de la app (para meta tags y enlaces absolutos).
      siteUrl,
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