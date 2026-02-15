import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import QuickReadingSheet from '../QuickReadingSheet.vue'

// Polyfill visualViewport for happy-dom (used by Vuetify overlay positioning)
if (globalThis.visualViewport === undefined) {
  Object.defineProperty(globalThis, 'visualViewport', {
    value: {
      width: 1024,
      height: 768,
      offsetLeft: 0,
      offsetTop: 0,
      pageLeft: 0,
      pageTop: 0,
      scale: 1,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    },
    writable: true,
  })
}

const vuetify = createVuetify({
  components,
  directives,
})

// Mock useApiClient (used by useProductionApi)
const mockRequest = vi.fn()

vi.mock('@/composables/useApiClient', () => ({
  useApiClient: () => ({
    request: mockRequest,
  }),
}))

// Mock useSnackbar
const mockShowNotice = vi.fn()

vi.mock('@/composables/useSnackbar', () => ({
  useSnackbar: () => ({
    snackbar: { show: false, text: '', color: 'success' },
    showNotice: mockShowNotice,
  }),
}))

function mountSheet (props: {
  modelValue?: boolean
  batchUuid?: string
  occupancyUuid?: string
  vesselName?: string
  batchName?: string
} = {}) {
  const div = document.createElement('div')
  div.id = 'app'
  document.body.append(div)

  return mount(QuickReadingSheet, {
    attachTo: div,
    global: {
      plugins: [vuetify],
    },
    props: {
      modelValue: props.modelValue ?? true,
      batchUuid: props.batchUuid ?? 'test-batch-uuid',
      occupancyUuid: props.occupancyUuid,
      vesselName: props.vesselName,
      batchName: props.batchName,
    },
  })
}

describe('QuickReadingSheet', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })

  describe('rendering', () => {
    it('renders the component title', () => {
      mountSheet()
      expect(document.body.textContent).toContain('Log Reading')
    })

    it('renders subtitle with vessel and batch names', () => {
      mountSheet({
        vesselName: 'FV-1',
        batchName: 'IPA #42',
      })
      expect(document.body.textContent).toContain('FV-1')
      expect(document.body.textContent).toContain('IPA #42')
    })

    it('renders gravity field', () => {
      mountSheet()
      expect(document.body.textContent).toContain('Gravity')
    })

    it('renders temperature field', () => {
      mountSheet()
      expect(document.body.textContent).toContain('Temperature')
    })

    it('renders Add pH button initially (pH field collapsed)', () => {
      mountSheet()
      expect(document.body.textContent).toContain('Add pH')
    })

    it('renders notes field', () => {
      mountSheet()
      expect(document.body.textContent).toContain('Notes')
    })

    it('renders Change time button initially (time picker collapsed)', () => {
      mountSheet()
      expect(document.body.textContent).toContain('Change time')
    })

    it('renders Save Reading button', () => {
      mountSheet()
      expect(document.body.textContent).toContain('Save Reading')
    })
  })

  describe('save behavior', () => {
    it('calls createMeasurement with batch_uuid when no occupancy provided', async () => {
      mockRequest.mockResolvedValue({})

      mountSheet({
        batchUuid: 'my-batch-uuid',
      })

      // Set gravity value via the input in the teleported content
      const inputs = document.querySelectorAll('input')
      const gravityInput = inputs[0] as HTMLInputElement
      gravityInput.value = '1.050'
      gravityInput.dispatchEvent(new Event('input', { bubbles: true }))
      await vi.dynamicImportSettled()

      // Submit form
      const form = document.querySelector('form') as HTMLFormElement
      form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }))
      await vi.dynamicImportSettled()

      // Wait for async operations
      await vi.waitFor(() => {
        expect(mockRequest).toHaveBeenCalled()
      })

      const callArgs = mockRequest.mock.calls[0]
      expect(callArgs[0]).toBe('/measurements')
      const body = JSON.parse(callArgs[1].body)
      expect(body.batch_uuid).toBe('my-batch-uuid')
      expect(body.kind).toBe('gravity')
      expect(body.unit).toBe('sg')
    })

    it('calls createMeasurement with occupancy_uuid when provided', async () => {
      mockRequest.mockResolvedValue({})

      mountSheet({
        batchUuid: 'my-batch-uuid',
        occupancyUuid: 'my-occupancy-uuid',
      })

      // Set gravity value
      const inputs = document.querySelectorAll('input')
      const gravityInput = inputs[0] as HTMLInputElement
      gravityInput.value = '1.050'
      gravityInput.dispatchEvent(new Event('input', { bubbles: true }))
      await vi.dynamicImportSettled()

      // Submit form
      const form = document.querySelector('form') as HTMLFormElement
      form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }))
      await vi.dynamicImportSettled()

      await vi.waitFor(() => {
        expect(mockRequest).toHaveBeenCalled()
      })

      const callArgs = mockRequest.mock.calls[0]
      const body = JSON.parse(callArgs[1].body)
      expect(body.occupancy_uuid).toBe('my-occupancy-uuid')
      expect(body.batch_uuid).toBeUndefined()
    })

    it('emits saved and closes sheet on success', async () => {
      mockRequest.mockResolvedValue({})

      const wrapper = mountSheet()

      // Set gravity value
      const inputs = document.querySelectorAll('input')
      const gravityInput = inputs[0] as HTMLInputElement
      gravityInput.value = '1.050'
      gravityInput.dispatchEvent(new Event('input', { bubbles: true }))
      await vi.dynamicImportSettled()

      // Submit form
      const form = document.querySelector('form') as HTMLFormElement
      form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }))
      await vi.dynamicImportSettled()

      await vi.waitFor(() => {
        expect(wrapper.emitted('saved')).toBeTruthy()
      })

      expect(mockShowNotice).toHaveBeenCalledWith('Reading recorded')
    })

    it('shows error snackbar on failure', async () => {
      mockRequest.mockRejectedValue(new Error('Network error'))

      mountSheet()

      // Set gravity value
      const inputs = document.querySelectorAll('input')
      const gravityInput = inputs[0] as HTMLInputElement
      gravityInput.value = '1.050'
      gravityInput.dispatchEvent(new Event('input', { bubbles: true }))
      await vi.dynamicImportSettled()

      // Submit form
      const form = document.querySelector('form') as HTMLFormElement
      form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }))
      await vi.dynamicImportSettled()

      await vi.waitFor(() => {
        expect(mockShowNotice).toHaveBeenCalledWith('Network error', 'error')
      })
    })
  })
})
