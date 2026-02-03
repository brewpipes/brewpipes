/**
 * Shared test utilities for Vue component tests with Vuetify.
 */
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

/**
 * Creates a Vuetify instance configured for testing.
 * Includes all components and directives.
 */
export function createTestVuetify () {
  return createVuetify({
    components,
    directives,
  })
}

/**
 * Creates a wrapper component that provides VApp context.
 * Required for layout components like VFooter, VAppBar, etc.
 */
export function createVAppWrapper (componentName: string, component: object) {
  return {
    components: { [componentName]: component },
    template: `<v-app><${componentName} /></v-app>`,
  }
}
