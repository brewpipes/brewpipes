import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick } from 'vue'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import BeerLotCreateDialog from '../BeerLotCreateDialog.vue'

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

// Mock useApiClient (used by useProductionApi and useInventoryApi)
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

// Mock vue-router
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
  useRoute: () => ({
    path: '/inventory/product',
  }),
}))

const mockBatches = [
  {
    uuid: 'batch-uuid-1',
    short_name: 'IPA-2024-001',
    brew_date: '2024-01-15T00:00:00Z',
    recipe_uuid: 'recipe-1',
    recipe_name: 'West Coast IPA',
    current_phase: 'fermenting',
    notes: null,
    created_at: '2024-01-15T00:00:00Z',
    updated_at: '2024-01-15T00:00:00Z',
  },
  {
    uuid: 'batch-uuid-2',
    short_name: 'STOUT-2024-002',
    brew_date: '2024-02-01T00:00:00Z',
    recipe_uuid: 'recipe-2',
    recipe_name: 'Milk Stout',
    current_phase: 'conditioning',
    notes: null,
    created_at: '2024-02-01T00:00:00Z',
    updated_at: '2024-02-01T00:00:00Z',
  },
]

const mockStockLocations = [
  {
    uuid: 'loc-uuid-1',
    name: 'Cold Room A',
    location_type: 'cold_storage',
    description: 'Main cold room',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
]

/** Helper to flush all pending promises and Vue reactivity */
async function flushAll () {
  await vi.dynamicImportSettled()
  await nextTick()
  await new Promise(resolve => setTimeout(resolve, 10))
  await nextTick()
}

function mountDialog () {
  const div = document.createElement('div')
  div.id = 'app'
  document.body.append(div)

  return mount(BeerLotCreateDialog, {
    attachTo: div,
    global: {
      plugins: [vuetify],
    },
    props: {
      modelValue: false,
    },
  })
}

/** Mount dialog closed, then open it (triggers the watch) */
async function mountAndOpen () {
  const wrapper = mountDialog()
  await wrapper.setProps({ modelValue: true })
  await flushAll()
  return wrapper
}

describe('BeerLotCreateDialog', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockRequest.mockImplementation((url: string) => {
      if (url === '/batches') return Promise.resolve(mockBatches)
      if (url === '/stock-locations') return Promise.resolve(mockStockLocations)
      if (url === '/beer-lots') return Promise.resolve({ uuid: 'new-lot-uuid' })
      return Promise.resolve([])
    })
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })

  describe('rendering', () => {
    it('renders the dialog title', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Create Beer Lot')
    })

    it('renders batch autocomplete field', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Batch')
    })

    it('renders all form fields', async () => {
      await mountAndOpen()
      const text = document.body.textContent ?? ''
      expect(text).toContain('Batch')
      expect(text).toContain('Lot code')
      expect(text).toContain('Container type')
      expect(text).toContain('Package format name')
      expect(text).toContain('Volume per unit')
      expect(text).toContain('Quantity')
      expect(text).toContain('Stock location')
      expect(text).toContain('Packaged at')
      expect(text).toContain('Best by')
      expect(text).toContain('Notes')
    })

    it('renders Cancel and Create lot buttons', async () => {
      await mountAndOpen()
      const text = document.body.textContent ?? ''
      expect(text).toContain('Cancel')
      expect(text).toContain('Create lot')
    })
  })

  describe('data loading', () => {
    it('loads batches when dialog opens', async () => {
      await mountAndOpen()
      expect(mockRequest).toHaveBeenCalledWith('/batches')
    })

    it('loads stock locations when dialog opens', async () => {
      await mountAndOpen()
      expect(mockRequest).toHaveBeenCalledWith('/stock-locations')
    })
  })

  describe('form validation', () => {
    it('disables submit button when no batch is selected', async () => {
      const wrapper = await mountAndOpen()
      const vm = wrapper.vm as unknown as { isFormValid: boolean }
      expect(vm.isFormValid).toBe(false)
    })

    it('enables submit when batch is selected', async () => {
      const wrapper = await mountAndOpen()
      const vm = wrapper.vm as unknown as {
        form: { production_batch_uuid: string | null }
        isFormValid: boolean
      }
      vm.form.production_batch_uuid = 'batch-uuid-1'
      await nextTick()
      expect(vm.isFormValid).toBe(true)
    })
  })

  describe('cancel behavior', () => {
    it('emits update:modelValue false when cancel is clicked', async () => {
      const wrapper = await mountAndOpen()

      // Find cancel button in the teleported content
      const buttons = document.body.querySelectorAll('.v-btn')
      let cancelBtn: HTMLElement | null = null
      buttons.forEach(btn => {
        if (btn.textContent?.trim() === 'Cancel') {
          cancelBtn = btn as HTMLElement
        }
      })
      expect(cancelBtn).toBeTruthy()
      cancelBtn!.click()
      await nextTick()

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')![0]).toEqual([false])
    })
  })

  describe('submission', () => {
    it('submits with correct payload when form is valid', async () => {
      mockRequest.mockImplementation((url: string, options?: RequestInit) => {
        if (url === '/batches') return Promise.resolve(mockBatches)
        if (url === '/stock-locations') return Promise.resolve(mockStockLocations)
        if (url === '/beer-lots' && options?.method === 'POST') {
          return Promise.resolve({ uuid: 'new-lot-uuid' })
        }
        return Promise.resolve([])
      })

      const wrapper = await mountAndOpen()

      // Set batch UUID directly on the form
      const vm = wrapper.vm as unknown as {
        form: { production_batch_uuid: string | null }
      }
      vm.form.production_batch_uuid = 'batch-uuid-1'
      await nextTick()

      // Find and click the Create lot button
      const buttons = document.body.querySelectorAll('.v-btn')
      let createBtn: HTMLElement | null = null
      buttons.forEach(btn => {
        if (btn.textContent?.trim() === 'Create lot') {
          createBtn = btn as HTMLElement
        }
      })
      expect(createBtn).toBeTruthy()
      createBtn!.click()
      await flushAll()

      // Verify the API was called with the correct payload
      const postCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => call[0] === '/beer-lots' && (call[1] as RequestInit)?.method === 'POST',
      )
      expect(postCalls.length).toBe(1)
      const body = JSON.parse((postCalls[0][1] as RequestInit).body as string)
      expect(body.production_batch_uuid).toBe('batch-uuid-1')

      // Verify snackbar was shown
      expect(mockShowNotice).toHaveBeenCalledWith('Beer lot created')

      // Verify events were emitted
      expect(wrapper.emitted('created')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    })

    it('shows error via snackbar when submission fails', async () => {
      mockRequest.mockImplementation((url: string, options?: RequestInit) => {
        if (url === '/batches') return Promise.resolve(mockBatches)
        if (url === '/stock-locations') return Promise.resolve(mockStockLocations)
        if (url === '/beer-lots' && options?.method === 'POST') {
          return Promise.reject(new Error('Server error'))
        }
        return Promise.resolve([])
      })

      const wrapper = await mountAndOpen()

      // Set batch UUID directly
      const vm = wrapper.vm as unknown as {
        form: { production_batch_uuid: string | null }
      }
      vm.form.production_batch_uuid = 'batch-uuid-1'
      await nextTick()

      // Click submit
      const buttons = document.body.querySelectorAll('.v-btn')
      let createBtn: HTMLElement | null = null
      buttons.forEach(btn => {
        if (btn.textContent?.trim() === 'Create lot') {
          createBtn = btn as HTMLElement
        }
      })
      createBtn!.click()
      await flushAll()

      // Verify error snackbar was shown
      expect(mockShowNotice).toHaveBeenCalledWith('Server error', 'error')

      // Should NOT emit created on failure
      expect(wrapper.emitted('created')).toBeFalsy()
    })
  })
})
