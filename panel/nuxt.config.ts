// https://nuxt.com/docs/api/configuration/nuxt-config

export default defineNuxtConfig({
  css: [
    'vuetify/lib/styles/main.sass',
    "@mdi/font/css/materialdesignicons.css",
  ],
  build: {
    transpile: ['vuetify']
  },
  ssr: false,
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true }
})
