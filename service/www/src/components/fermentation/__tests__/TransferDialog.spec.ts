import type { Batch, Occupancy, Vessel, Volume } from '@/types'
import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick } from 'vue'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import TransferDialog from '../TransferDialog.vue'

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

// Mock vue-router
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
  }),
  useRoute: () => ({
    path: '/fermentation',
  }),
}))

function makeOccupancy (overrides: Partial<Occupancy> = {}): Occupancy {
  return {
    uuid: 'occ-1',
    vessel_uuid: 'vessel-1',
    volume_uuid: 'vol-1',
    batch_uuid: 'batch-1',
    status: 'fermenting',
    in_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(),
    out_at: null,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    ...overrides,
  }
}

function makeVessel (overrides: Partial<Vessel> = {}): Vessel {
  return {
    uuid: 'vessel-1',
    type: 'fermenter',
    name: 'FV-1',
    capacity: 10,
    capacity_unit: 'bbl',
    make: null,
    model: null,
    status: 'active',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    deleted_at: null,
    ...overrides,
  }
}

function makeBatch (overrides: Partial<Batch> = {}): Batch {
  return {
    uuid: 'batch-1',
    short_name: 'IPA #42',
    brew_date: '2025-01-01',
    recipe_uuid: null,
    recipe_name: null,
    current_phase: 'fermenting',
    notes: null,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    ...overrides,
  }
}

function makeVolume (overrides: Partial<Volume> = {}): Volume {
  return {
    uuid: 'vol-1',
    name: 'IPA Wort',
    description: null,
    amount: 7,
    amount_unit: 'bbl',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    deleted_at: null,
    ...overrides,
  }
}

/** Helper to flush all pending promises and Vue reactivity */
async function flushAll () {
  await vi.dynamicImportSettled()
  await nextTick()
  // Flush microtask queue for async operations in resetAndLoad
  await new Promise(resolve => setTimeout(resolve, 10))
  await nextTick()
}

type TransferMode = 'transfer' | 'split' | 'blend'

function mountDialog (props: {
  modelValue?: boolean
  mode?: TransferMode
  sourceOccupancy?: Occupancy | null
  sourceVessel?: Vessel | null
  sourceBatch?: Batch | null
  sourceVolume?: Volume | null
} = {}) {
  const div = document.createElement('div')
  div.id = 'app'
  document.body.append(div)

  // Mock API responses for loading reference data
  mockRequest.mockImplementation((url: string) => {
    if (url.includes('/occupancies')) {
      return Promise.resolve([
        makeOccupancy(),
        makeOccupancy({ uuid: 'occ-2', vessel_uuid: 'vessel-3', volume_uuid: 'vol-2', batch_uuid: 'batch-2', status: 'conditioning' }),
      ])
    }
    if (url.includes('/vessels')) {
      return Promise.resolve([
        makeVessel(),
        makeVessel({ uuid: 'vessel-2', name: 'BT-1', type: 'brite_tank' }),
        makeVessel({ uuid: 'vessel-3', name: 'FV-2', type: 'fermenter' }),
        makeVessel({ uuid: 'vessel-4', name: 'BT-2', type: 'brite_tank' }),
      ])
    }
    return Promise.resolve([])
  })

  return mount(TransferDialog, {
    attachTo: div,
    global: {
      plugins: [vuetify],
    },
    props: {
      modelValue: props.modelValue ?? false,
      mode: props.mode,
      sourceOccupancy: props.sourceOccupancy ?? makeOccupancy(),
      sourceVessel: props.sourceVessel ?? makeVessel(),
      sourceBatch: props.sourceBatch ?? makeBatch(),
      sourceVolume: props.sourceVolume ?? makeVolume(),
    },
  })
}

/** Mount dialog closed, then open it (triggers the watch) */
async function mountAndOpen (props: {
  mode?: TransferMode
  sourceOccupancy?: Occupancy | null
  sourceVessel?: Vessel | null
  sourceBatch?: Batch | null
  sourceVolume?: Volume | null
} = {}) {
  const wrapper = mountDialog({ modelValue: false, ...props })
  await wrapper.setProps({ modelValue: true })
  await flushAll()
  return wrapper
}

describe('TransferDialog', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })

  describe('rendering', () => {
    it('renders the dialog title', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Transfer Beer')
    })

    it('renders source vessel info when sourceOccupancy is provided', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('FV-1')
    })

    it('renders source batch name when sourceBatch is provided', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('IPA #42')
    })

    it('renders the occupancy status chip', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Fermenting')
    })

    it('renders Transfer Details step label', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Transfer Details')
    })

    it('renders Review & Confirm step label', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Review & Confirm')
    })

    it('renders Cancel button', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Cancel')
    })

    it('renders Next button on step 1', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Next →')
    })

    it('renders volume fields', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Transfer volume')
      expect(document.body.textContent).toContain('Transfer loss')
    })

    it('renders destination status select', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Destination status')
    })

    it('renders close source checkbox', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Close source vessel after transfer')
    })

    it('renders Change time button', async () => {
      await mountAndOpen()
      expect(document.body.textContent).toContain('Change time')
    })
  })

  describe('smart defaults', () => {
    it('defaults dest status to conditioning when source is fermenting', async () => {
      const wrapper = await mountAndOpen({
        sourceOccupancy: makeOccupancy({ status: 'fermenting' }),
      })
      const vm = wrapper.vm as unknown as { form: { destStatus: string } }
      expect(vm.form.destStatus).toBe('conditioning')
    })

    it('defaults dest status to carbonating when source is conditioning', async () => {
      const wrapper = await mountAndOpen({
        sourceOccupancy: makeOccupancy({ status: 'conditioning' }),
      })
      const vm = wrapper.vm as unknown as { form: { destStatus: string } }
      expect(vm.form.destStatus).toBe('carbonating')
    })

    it('defaults dest status to conditioning when source is cold_crashing', async () => {
      const wrapper = await mountAndOpen({
        sourceOccupancy: makeOccupancy({ status: 'cold_crashing' }),
      })
      const vm = wrapper.vm as unknown as { form: { destStatus: string } }
      expect(vm.form.destStatus).toBe('conditioning')
    })

    it('defaults dest status to holding when source has no status', async () => {
      const wrapper = await mountAndOpen({
        sourceOccupancy: makeOccupancy({ status: null }),
      })
      const vm = wrapper.vm as unknown as { form: { destStatus: string } }
      expect(vm.form.destStatus).toBe('holding')
    })

    it('pre-fills transfer amount from source volume', async () => {
      const wrapper = await mountAndOpen({
        sourceVolume: makeVolume({ amount: 7, amount_unit: 'bbl' }),
      })
      const vm = wrapper.vm as unknown as { form: { transferAmount: string } }
      expect(vm.form.transferAmount).toBe('7.00')
    })

    it('defaults close source to true', async () => {
      const wrapper = await mountAndOpen()
      const vm = wrapper.vm as unknown as { form: { closeSource: boolean } }
      expect(vm.form.closeSource).toBe(true)
    })
  })

  describe('events', () => {
    it('loads reference data when dialog opens', async () => {
      await mountAndOpen()
      // Verify API calls were made to load occupancies and vessels
      const occupancyCall = mockRequest.mock.calls.find(
        (call: string[]) => call[0].includes('/occupancies'),
      )
      const vesselCall = mockRequest.mock.calls.find(
        (call: string[]) => call[0].includes('/vessels'),
      )
      expect(occupancyCall).toBeDefined()
      expect(vesselCall).toBeDefined()
    })

    it('does not render content when modelValue is false', () => {
      mountDialog({ modelValue: false })
      // Dialog should not show its content when closed
      expect(document.body.textContent).not.toContain('Transfer Details')
    })
  })

  describe('mode selector', () => {
    it('renders mode toggle buttons', async () => {
      await mountAndOpen()
      const text = document.body.textContent ?? ''
      expect(text).toContain('Transfer')
      expect(text).toContain('Split')
      expect(text).toContain('Blend')
    })

    it('defaults to transfer mode', async () => {
      const wrapper = await mountAndOpen()
      const vm = wrapper.vm as unknown as { activeMode: string }
      expect(vm.activeMode).toBe('transfer')
    })

    it('respects mode prop for split', async () => {
      const wrapper = await mountAndOpen({ mode: 'split' })
      const vm = wrapper.vm as unknown as { activeMode: string }
      expect(vm.activeMode).toBe('split')
    })

    it('respects mode prop for blend', async () => {
      const wrapper = await mountAndOpen({ mode: 'blend' })
      const vm = wrapper.vm as unknown as { activeMode: string }
      expect(vm.activeMode).toBe('blend')
    })
  })

  describe('split mode', () => {
    it('renders Split Beer title when in split mode', async () => {
      await mountAndOpen({ mode: 'split' })
      expect(document.body.textContent).toContain('Split Beer')
    })

    it('renders Split Into section', async () => {
      await mountAndOpen({ mode: 'split' })
      expect(document.body.textContent).toContain('Split Into')
    })

    it('renders 2 destination rows by default', async () => {
      const wrapper = await mountAndOpen({ mode: 'split' })
      const vm = wrapper.vm as unknown as { splitDestinations: Array<{ vesselUuid: string }> }
      expect(vm.splitDestinations).toHaveLength(2)
    })

    it('renders Destination 1 and Destination 2 labels', async () => {
      await mountAndOpen({ mode: 'split' })
      expect(document.body.textContent).toContain('Destination 1')
      expect(document.body.textContent).toContain('Destination 2')
    })

    it('renders Add Destination button', async () => {
      await mountAndOpen({ mode: 'split' })
      expect(document.body.textContent).toContain('Add Destination')
    })

    it('renders close source checkbox for split', async () => {
      await mountAndOpen({ mode: 'split' })
      expect(document.body.textContent).toContain('Close source vessel after split')
    })

    it('renders Confirm Split button on step 2', async () => {
      await mountAndOpen({ mode: 'split' })
      // The confirm button text should be "Confirm Split" (visible on step 2)
      const vm = (await mountAndOpen({ mode: 'split' })).vm as unknown as { confirmButtonLabel: string }
      expect(vm.confirmButtonLabel).toBe('Confirm Split')
    })

    it('adds a destination row when Add Destination is clicked', async () => {
      const wrapper = await mountAndOpen({ mode: 'split' })
      const vm = wrapper.vm as unknown as {
        splitDestinations: Array<{ vesselUuid: string }>
        addSplitDestination: () => void
      }
      // Simulate adding a destination
      vm.addSplitDestination()
      await nextTick()
      expect(vm.splitDestinations).toHaveLength(3)
    })

    it('limits destinations to 4 maximum', async () => {
      const wrapper = await mountAndOpen({ mode: 'split' })
      const vm = wrapper.vm as unknown as {
        splitDestinations: Array<{ vesselUuid: string }>
        addSplitDestination: () => void
      }
      vm.addSplitDestination() // 3
      vm.addSplitDestination() // 4
      vm.addSplitDestination() // should not add (max 4)
      await nextTick()
      expect(vm.splitDestinations).toHaveLength(4)
    })

    it('does not remove below 2 destinations', async () => {
      const wrapper = await mountAndOpen({ mode: 'split' })
      const vm = wrapper.vm as unknown as {
        splitDestinations: Array<{ vesselUuid: string }>
        removeSplitDestination: (index: number) => void
      }
      vm.removeSplitDestination(0) // should not remove (min 2)
      await nextTick()
      expect(vm.splitDestinations).toHaveLength(2)
    })

    it('computes split volume match correctly', async () => {
      const wrapper = await mountAndOpen({
        mode: 'split',
        sourceVolume: makeVolume({ amount: 10, amount_unit: 'bbl' }),
      })
      const vm = wrapper.vm as unknown as {
        splitDestinations: Array<{ vesselUuid: string, amount: string, status: string }>
        splitVolumeMatch: string
        form: { lossAmount: string }
      }
      // Set destinations that add up to source
      vm.splitDestinations[0]!.amount = '5'
      vm.splitDestinations[1]!.amount = '4.5'
      vm.form.lossAmount = '0.5'
      await nextTick()
      expect(vm.splitVolumeMatch).toBe('match')
    })

    it('detects split volume mismatch', async () => {
      const wrapper = await mountAndOpen({
        mode: 'split',
        sourceVolume: makeVolume({ amount: 10, amount_unit: 'bbl' }),
      })
      const vm = wrapper.vm as unknown as {
        splitDestinations: Array<{ vesselUuid: string, amount: string, status: string }>
        splitVolumeMatch: string
        form: { lossAmount: string }
      }
      vm.splitDestinations[0]!.amount = '3'
      vm.splitDestinations[1]!.amount = '3'
      vm.form.lossAmount = '0'
      await nextTick()
      expect(vm.splitVolumeMatch).toBe('mismatch')
    })
  })

  describe('blend mode', () => {
    it('renders Blend Beer title when in blend mode', async () => {
      await mountAndOpen({ mode: 'blend' })
      expect(document.body.textContent).toContain('Blend Beer')
    })

    it('renders Blend From section', async () => {
      await mountAndOpen({ mode: 'blend' })
      expect(document.body.textContent).toContain('Blend From')
    })

    it('renders 2 source rows by default', async () => {
      const wrapper = await mountAndOpen({ mode: 'blend' })
      const vm = wrapper.vm as unknown as { blendSources: Array<{ occupancyUuid: string }> }
      expect(vm.blendSources).toHaveLength(2)
    })

    it('renders Source 1 and Source 2 labels', async () => {
      await mountAndOpen({ mode: 'blend' })
      expect(document.body.textContent).toContain('Source 1')
      expect(document.body.textContent).toContain('Source 2')
    })

    it('renders Add Source button', async () => {
      await mountAndOpen({ mode: 'blend' })
      expect(document.body.textContent).toContain('Add Source')
    })

    it('renders Into section', async () => {
      await mountAndOpen({ mode: 'blend' })
      expect(document.body.textContent).toContain('Into')
    })

    it('renders Confirm Blend button label', async () => {
      const wrapper = await mountAndOpen({ mode: 'blend' })
      const vm = wrapper.vm as unknown as { confirmButtonLabel: string }
      expect(vm.confirmButtonLabel).toBe('Confirm Blend')
    })

    it('pre-populates first source from sourceOccupancy', async () => {
      const wrapper = await mountAndOpen({
        mode: 'blend',
        sourceOccupancy: makeOccupancy({ uuid: 'occ-1' }),
      })
      const vm = wrapper.vm as unknown as { blendSources: Array<{ occupancyUuid: string }> }
      expect(vm.blendSources[0]!.occupancyUuid).toBe('occ-1')
    })

    it('adds a source row when Add Source is called', async () => {
      const wrapper = await mountAndOpen({ mode: 'blend' })
      const vm = wrapper.vm as unknown as {
        blendSources: Array<{ occupancyUuid: string }>
        addBlendSource: () => void
      }
      vm.addBlendSource()
      await nextTick()
      expect(vm.blendSources).toHaveLength(3)
    })

    it('limits sources to 4 maximum', async () => {
      const wrapper = await mountAndOpen({ mode: 'blend' })
      const vm = wrapper.vm as unknown as {
        blendSources: Array<{ occupancyUuid: string }>
        addBlendSource: () => void
      }
      vm.addBlendSource() // 3
      vm.addBlendSource() // 4
      vm.addBlendSource() // should not add (max 4)
      await nextTick()
      expect(vm.blendSources).toHaveLength(4)
    })

    it('does not remove below 2 sources', async () => {
      const wrapper = await mountAndOpen({ mode: 'blend' })
      const vm = wrapper.vm as unknown as {
        blendSources: Array<{ occupancyUuid: string }>
        removeBlendSource: (index: number) => void
      }
      vm.removeBlendSource(0) // should not remove (min 2)
      await nextTick()
      expect(vm.blendSources).toHaveLength(2)
    })

    it('computes blend destination amount correctly', async () => {
      const wrapper = await mountAndOpen({ mode: 'blend' })
      const vm = wrapper.vm as unknown as {
        blendSources: Array<{ occupancyUuid: string, amount: string }>
        blendDestAmount: number
        form: { lossAmount: string }
      }
      vm.blendSources[0]!.amount = '5'
      vm.blendSources[1]!.amount = '5'
      vm.form.lossAmount = '0.3'
      await nextTick()
      expect(vm.blendDestAmount).toBeCloseTo(9.7, 1)
    })

    it('defaults closeAfterBlend to true for each source', async () => {
      const wrapper = await mountAndOpen({ mode: 'blend' })
      const vm = wrapper.vm as unknown as { blendSources: Array<{ closeAfterBlend: boolean }> }
      expect(vm.blendSources[0]!.closeAfterBlend).toBe(true)
      expect(vm.blendSources[1]!.closeAfterBlend).toBe(true)
    })
  })

  describe('validation', () => {
    it('cannot proceed in split mode without valid destinations', async () => {
      const wrapper = await mountAndOpen({ mode: 'split' })
      const vm = wrapper.vm as unknown as { canProceedToReview: boolean }
      // Destinations have no vessel or amount set
      expect(vm.canProceedToReview).toBe(false)
    })

    it('cannot proceed in blend mode without valid sources and destination', async () => {
      const wrapper = await mountAndOpen({ mode: 'blend' })
      const vm = wrapper.vm as unknown as { canProceedToReview: boolean }
      // Sources have no occupancy or amount set, no dest vessel
      expect(vm.canProceedToReview).toBe(false)
    })
  })
})
