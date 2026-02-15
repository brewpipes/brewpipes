import type { BatchSummary, Measurement, Occupancy, Vessel } from '@/types'
import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import FermentationCard from '../FermentationCard.vue'

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
  RouterLink: {
    template: '<a><slot /></a>',
    props: ['to'],
  },
}))

function makeOccupancy (overrides: Partial<Occupancy> = {}): Occupancy {
  return {
    uuid: 'occ-1',
    vessel_uuid: 'vessel-1',
    volume_uuid: 'vol-1',
    batch_uuid: 'batch-1',
    status: 'fermenting',
    in_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(), // 5 days ago
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

function makeBatchSummary (overrides: Partial<BatchSummary> = {}): BatchSummary {
  return {
    uuid: 'batch-1',
    short_name: 'IPA #42',
    brew_date: '2025-01-01',
    notes: null,
    recipe_name: 'West Coast IPA',
    style_name: 'IPA',
    brew_sessions: [],
    current_phase: 'fermenting',
    current_vessel: 'FV-1',
    current_occupancy_status: 'fermenting',
    current_occupancy_uuid: 'occ-1',
    original_gravity: 1.065,
    final_gravity: null,
    abv: null,
    abv_calculated: false,
    ibu: 65,
    days_in_fermenter: 5,
    days_in_brite: null,
    days_grain_to_glass: null,
    starting_volume_bbl: 10,
    current_volume_bbl: 9.5,
    total_loss_bbl: 0.5,
    loss_percentage: 5,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    ...overrides,
  }
}

function makeMeasurement (overrides: Partial<Measurement> = {}): Measurement {
  return {
    uuid: 'meas-1',
    batch_uuid: 'batch-1',
    occupancy_uuid: 'occ-1',
    volume_uuid: null,
    kind: 'gravity',
    value: 1.05,
    unit: 'sg',
    observed_at: new Date().toISOString(),
    notes: null,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    deleted_at: null,
    ...overrides,
  }
}

function mountCard (props: {
  occupancy?: Occupancy
  vessel?: Vessel
  batchSummary?: BatchSummary | null
  measurements?: Measurement[]
} = {}) {
  const div = document.createElement('div')
  div.id = 'app'
  document.body.append(div)

  return mount(FermentationCard, {
    attachTo: div,
    global: {
      plugins: [vuetify],
      stubs: {
        RouterLink: {
          template: '<a><slot /></a>',
          props: ['to'],
        },
      },
    },
    props: {
      occupancy: props.occupancy ?? makeOccupancy(),
      vessel: props.vessel ?? makeVessel(),
      batchSummary: props.batchSummary === undefined ? makeBatchSummary() : props.batchSummary,
      measurements: props.measurements ?? [],
    },
  })
}

describe('FermentationCard', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })

  describe('rendering', () => {
    it('renders the vessel name', () => {
      mountCard()
      expect(document.body.textContent).toContain('FV-1')
    })

    it('renders the batch name', () => {
      mountCard()
      expect(document.body.textContent).toContain('IPA #42')
    })

    it('renders the recipe name', () => {
      mountCard()
      expect(document.body.textContent).toContain('West Coast IPA')
    })

    it('renders the day count badge', () => {
      mountCard()
      expect(document.body.textContent).toContain('Day 5')
    })

    it('renders the occupancy status', () => {
      mountCard()
      expect(document.body.textContent).toContain('Fermenting')
    })

    it('renders Log Reading button', () => {
      mountCard()
      expect(document.body.textContent).toContain('Log Reading')
    })

    it('renders gravity sparkline label', () => {
      mountCard()
      expect(document.body.textContent).toContain('Gravity')
    })

    it('renders temperature sparkline label', () => {
      mountCard()
      expect(document.body.textContent).toContain('Temp')
    })

    it('renders attenuation label', () => {
      mountCard()
      expect(document.body.textContent).toContain('Attenuation')
    })

    it('renders ABV label', () => {
      mountCard()
      expect(document.body.textContent).toContain('Est. ABV')
    })
  })

  describe('metrics', () => {
    it('shows latest gravity value', () => {
      const measurements = [
        makeMeasurement({ kind: 'gravity', value: 1.05, observed_at: '2025-01-02T00:00:00Z' }),
        makeMeasurement({ uuid: 'meas-2', kind: 'gravity', value: 1.03, observed_at: '2025-01-03T00:00:00Z' }),
      ]
      mountCard({ measurements })
      // Should show the latest gravity (1.030) formatted as SG
      expect(document.body.textContent).toContain('1.030')
    })

    it('shows latest temperature value', () => {
      const measurements = [
        makeMeasurement({ kind: 'temperature', value: 20, unit: 'c', observed_at: '2025-01-02T00:00:00Z' }),
      ]
      mountCard({ measurements })
      // Should show temperature formatted (20°C → 68.0°F in default US prefs)
      expect(document.body.textContent).toContain('68.0')
    })

    it('computes attenuation correctly', () => {
      const summary = makeBatchSummary({ original_gravity: 1.06 })
      const measurements = [
        makeMeasurement({ kind: 'gravity', value: 1.015, observed_at: '2025-01-03T00:00:00Z' }),
      ]
      mountCard({ batchSummary: summary, measurements })
      // Attenuation = (1.060 - 1.015) / (1.060 - 1.0) * 100 = 75.0%
      expect(document.body.textContent).toContain('75.0%')
    })

    it('computes ABV correctly when no summary ABV', () => {
      const summary = makeBatchSummary({ original_gravity: 1.06, abv: null })
      const measurements = [
        makeMeasurement({ kind: 'gravity', value: 1.015, observed_at: '2025-01-03T00:00:00Z' }),
      ]
      mountCard({ batchSummary: summary, measurements })
      // ABV = (1.060 - 1.015) * 131.25 = 5.9%
      expect(document.body.textContent).toContain('5.9%')
    })

    it('shows dash when no measurements', () => {
      mountCard({ measurements: [] })
      // Should show dashes for gravity and temp
      const text = document.body.textContent ?? ''
      // Count dashes — there should be at least 2 (gravity and temp latest labels)
      const dashCount = (text.match(/—/g) ?? []).length
      expect(dashCount).toBeGreaterThanOrEqual(2)
    })
  })

  describe('attention indicators', () => {
    it('applies warning class when no gravity readings', () => {
      const wrapper = mountCard({ measurements: [] })
      expect(wrapper.find('.attention-warning').exists()).toBe(true)
    })

    it('applies warning class when gravity reading is stale (24+ hours)', () => {
      const staleTime = new Date(Date.now() - 25 * 60 * 60 * 1000).toISOString()
      const measurements = [
        makeMeasurement({ kind: 'gravity', value: 1.05, observed_at: staleTime }),
      ]
      const wrapper = mountCard({ measurements })
      expect(wrapper.find('.attention-warning').exists()).toBe(true)
    })

    it('applies info class when gravity is stable for 3+ readings', () => {
      const now = Date.now()
      const measurements = [
        makeMeasurement({ uuid: 'm1', kind: 'gravity', value: 1.012, observed_at: new Date(now - 3 * 60 * 60 * 1000).toISOString() }),
        makeMeasurement({ uuid: 'm2', kind: 'gravity', value: 1.012, observed_at: new Date(now - 2 * 60 * 60 * 1000).toISOString() }),
        makeMeasurement({ uuid: 'm3', kind: 'gravity', value: 1.012, observed_at: new Date(now - 1 * 60 * 60 * 1000).toISOString() }),
      ]
      const wrapper = mountCard({ measurements })
      expect(wrapper.find('.attention-info').exists()).toBe(true)
    })

    it('has no attention class when readings are normal', () => {
      const now = Date.now()
      const measurements = [
        makeMeasurement({ uuid: 'm1', kind: 'gravity', value: 1.05, observed_at: new Date(now - 2 * 60 * 60 * 1000).toISOString() }),
        makeMeasurement({ uuid: 'm2', kind: 'gravity', value: 1.04, observed_at: new Date(now - 1 * 60 * 60 * 1000).toISOString() }),
      ]
      const wrapper = mountCard({ measurements })
      expect(wrapper.find('.attention-warning').exists()).toBe(false)
      expect(wrapper.find('.attention-info').exists()).toBe(false)
    })
  })

  describe('events', () => {
    it('emits logReading when Log Reading button is clicked', async () => {
      const wrapper = mountCard()
      const buttons = wrapper.findAll('.v-btn')
      const logButton = buttons.find(b => b.text().includes('Log Reading'))
      expect(logButton).toBeDefined()
      await logButton!.trigger('click')
      expect(wrapper.emitted('logReading')).toBeTruthy()
    })
  })
})
