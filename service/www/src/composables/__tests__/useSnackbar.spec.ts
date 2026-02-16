import { describe, expect, it } from 'vitest'
import { useSnackbar } from '../useSnackbar'

describe('useSnackbar', () => {
  describe('initial state', () => {
    it('starts with show false', () => {
      const { snackbar } = useSnackbar()
      expect(snackbar.show).toBe(false)
    })

    it('starts with empty text', () => {
      const { snackbar } = useSnackbar()
      expect(snackbar.text).toBe('')
    })

    it('starts with success color', () => {
      const { snackbar } = useSnackbar()
      expect(snackbar.color).toBe('success')
    })

    it('starts with 3000ms timeout', () => {
      const { snackbar } = useSnackbar()
      expect(snackbar.timeout).toBe(3000)
    })
  })

  describe('showNotice', () => {
    it('sets show to true', () => {
      const { snackbar, showNotice } = useSnackbar()
      snackbar.show = false
      showNotice('Test message')
      expect(snackbar.show).toBe(true)
    })

    it('sets the text', () => {
      const { snackbar, showNotice } = useSnackbar()
      showNotice('Vessel updated successfully')
      expect(snackbar.text).toBe('Vessel updated successfully')
    })

    it('defaults color to success', () => {
      const { snackbar, showNotice } = useSnackbar()
      showNotice('Saved')
      expect(snackbar.color).toBe('success')
    })

    it('defaults timeout to 3000', () => {
      const { snackbar, showNotice } = useSnackbar()
      showNotice('Saved')
      expect(snackbar.timeout).toBe(3000)
    })

    it('accepts a custom color', () => {
      const { snackbar, showNotice } = useSnackbar()
      showNotice('Error occurred', 'error')
      expect(snackbar.color).toBe('error')
    })

    it('accepts a custom timeout', () => {
      const { snackbar, showNotice } = useSnackbar()
      showNotice('Quick message', 'info', 1000)
      expect(snackbar.timeout).toBe(1000)
    })

    it('accepts all custom parameters', () => {
      const { snackbar, showNotice } = useSnackbar()
      showNotice('Warning!', 'warning', 5000)
      expect(snackbar.text).toBe('Warning!')
      expect(snackbar.color).toBe('warning')
      expect(snackbar.timeout).toBe(5000)
      expect(snackbar.show).toBe(true)
    })

    it('overwrites previous message', () => {
      const { snackbar, showNotice } = useSnackbar()
      showNotice('First message')
      showNotice('Second message')
      expect(snackbar.text).toBe('Second message')
    })

    it('overwrites previous color', () => {
      const { snackbar, showNotice } = useSnackbar()
      showNotice('Error', 'error')
      showNotice('Success')
      expect(snackbar.color).toBe('success')
    })
  })

  describe('singleton behavior', () => {
    it('shares state across multiple calls', () => {
      const instance1 = useSnackbar()
      const instance2 = useSnackbar()

      instance1.showNotice('Shared message', 'info', 2000)

      expect(instance2.snackbar.show).toBe(true)
      expect(instance2.snackbar.text).toBe('Shared message')
      expect(instance2.snackbar.color).toBe('info')
      expect(instance2.snackbar.timeout).toBe(2000)
    })

    it('returns the same snackbar reactive object', () => {
      const instance1 = useSnackbar()
      const instance2 = useSnackbar()
      expect(instance1.snackbar).toBe(instance2.snackbar)
    })
  })

  describe('manual hide', () => {
    it('can be hidden by setting show to false', () => {
      const { snackbar, showNotice } = useSnackbar()
      showNotice('Visible')
      expect(snackbar.show).toBe(true)

      snackbar.show = false
      expect(snackbar.show).toBe(false)
    })
  })
})
