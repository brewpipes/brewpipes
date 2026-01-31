/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: 'brewDark',
    themes: {
      brewLight: {
        dark: false,
        colors: {
          background: '#f5f2ea',
          surface: '#fdfbf7',
          primary: '#2f5d50',
          secondary: '#c4753c',
          error: '#b23c2d',
          info: '#2d5f7a',
          success: '#2f7a4a',
          warning: '#c08322',
        },
      },
      brewDark: {
        dark: true,
        colors: {
          background: '#12100d',
          surface: '#1a1714',
          primary: '#6db28f',
          secondary: '#d39a5f',
          error: '#e06a5a',
          info: '#6ea0c1',
          success: '#6fbf86',
          warning: '#d6a34a',
        },
      },
    },
  },
})
