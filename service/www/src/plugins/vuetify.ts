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
    defaultTheme: 'brew',
    themes: {
      brew: {
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
    },
  },
})
