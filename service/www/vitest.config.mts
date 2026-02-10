import { fileURLToPath, URL } from 'node:url'
import Vue from '@vitejs/plugin-vue'
import Vuetify from 'vite-plugin-vuetify'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [
    Vue(),
    Vuetify({ autoImport: true }),
  ],
  test: {
    environment: 'happy-dom',
    globals: true,
    include: ['src/**/__tests__/**/*.spec.ts'],
    server: {
      deps: {
        inline: ['vuetify'],
      },
    },
    css: true,
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('src', import.meta.url)),
    },
  },
})
