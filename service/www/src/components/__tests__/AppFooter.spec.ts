import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import AppFooter from '../AppFooter.vue'

const vuetify = createVuetify({
  components,
  directives,
})

// Wrapper component that provides VApp context for layout components
const WrapperComponent = {
  components: { AppFooter },
  template: '<v-app><AppFooter /></v-app>',
}

function mountAppFooter () {
  return mount(WrapperComponent, {
    global: {
      plugins: [vuetify],
    },
  })
}

describe('AppFooter', () => {
  describe('rendering', () => {
    it('renders the footer component', () => {
      const wrapper = mountAppFooter()
      expect(wrapper.find('.v-footer').exists()).toBe(true)
    })

    it('does not display the application name text', () => {
      const wrapper = mountAppFooter()
      expect(wrapper.text()).not.toContain('BrewPipes Production UI')
    })

    it('renders the Repo link', () => {
      const wrapper = mountAppFooter()
      const repoLink = wrapper.find('a[href="https://github.com/brewpipes/brewpipes"]')
      expect(repoLink.exists()).toBe(true)
      expect(wrapper.text()).toContain('Repo')
    })

    it('renders the API link', () => {
      const wrapper = mountAppFooter()
      // API link uses window.location.origin dynamically
      const apiLink = wrapper.find('a[href$="/api"]')
      expect(apiLink.exists()).toBe(true)
      expect(wrapper.text()).toContain('API')
    })

    it('opens links in new tab', () => {
      const wrapper = mountAppFooter()
      const links = wrapper.findAll('a[target="_blank"]')
      expect(links.length).toBe(2)
    })

    it('has noopener noreferrer on external links', () => {
      const wrapper = mountAppFooter()
      const links = wrapper.findAll('a[rel="noopener noreferrer"]')
      expect(links.length).toBe(2)
    })
  })
})
