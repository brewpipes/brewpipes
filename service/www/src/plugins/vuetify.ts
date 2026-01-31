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
          background: '#f7f3ee',
          surface: '#fefcf9',
          primary: '#b66a3c',
          secondary: '#d49a6a',
          error: '#b94d3a',
          info: '#a8643d',
          success: '#c1774b',
          warning: '#d7a06c',
        },
      },
      brewDark: {
        dark: true,
        colors: {
          background: '#12100f',
          surface: '#1a1715',
          primary: '#d3a077',
          secondary: '#b8734a',
          error: '#e08a6f',
          info: '#c4835c',
          success: '#d09a72',
          warning: '#e3b083',
        },
      },
    },
  },
})
