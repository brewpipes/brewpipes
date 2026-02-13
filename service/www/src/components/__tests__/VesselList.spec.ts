import type { Occupancy, Vessel } from '@/types'
import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import VesselList from '../VesselList.vue'

// Mock the useUnitPreferences composable
vi.mock('@/composables/useUnitPreferences', () => ({
  useUnitPreferences: () => ({
    formatVolumePreferred: (value: number, unit: string) => `${value} ${unit}`,
  }),
}))

const vuetify = createVuetify({
  components,
  directives,
})

function createVessel (overrides: Partial<Vessel> = {}): Vessel {
  return {
    uuid: 'vessel-uuid-1',
    type: 'fermenter',
    name: 'FV-01',
    capacity: 10,
    capacity_unit: 'bbl',
    make: null,
    model: null,
    status: 'active',
    created_at: '2024-06-01T10:00:00Z',
    updated_at: '2024-06-01T10:00:00Z',
    deleted_at: null,
    ...overrides,
  }
}

function createOccupancy (overrides: Partial<Occupancy> = {}): Occupancy {
  return {
    uuid: 'occupancy-uuid-1',
    vessel_uuid: 'vessel-uuid-1',
    volume_uuid: 'volume-uuid-1',
    batch_uuid: 'batch-uuid-1',
    status: 'fermenting',
    in_at: '2024-06-10T10:00:00Z',
    out_at: null,
    created_at: '2024-06-10T10:00:00Z',
    updated_at: '2024-06-10T10:00:00Z',
    ...overrides,
  }
}

function mountVesselList (props: {
  vessels?: Vessel[]
  occupancies?: Occupancy[]
  selectedVesselUuid?: string | null
  loading?: boolean
} = {}) {
  return mount(VesselList, {
    global: {
      plugins: [vuetify],
    },
    props: {
      vessels: props.vessels ?? [],
      occupancies: props.occupancies ?? [],
      selectedVesselUuid: props.selectedVesselUuid ?? null,
      loading: props.loading ?? false,
    },
  })
}

describe('VesselList', () => {
  describe('rendering', () => {
    it('renders the component title', () => {
      const wrapper = mountVesselList()
      expect(wrapper.text()).toContain('Active Vessels')
    })

    it('renders vessel list items', () => {
      const vessels = [
        createVessel({ uuid: 'vessel-uuid-1', name: 'FV-01' }),
        createVessel({ uuid: 'vessel-uuid-2', name: 'FV-02' }),
      ]
      const wrapper = mountVesselList({ vessels })

      expect(wrapper.text()).toContain('FV-01')
      expect(wrapper.text()).toContain('FV-02')
    })

    it('renders vessel type and capacity', () => {
      const vessels = [
        createVessel({ type: 'fermenter', capacity: 15, capacity_unit: 'bbl' }),
      ]
      const wrapper = mountVesselList({ vessels })

      expect(wrapper.text()).toContain('fermenter')
      expect(wrapper.text()).toContain('15 bbl')
    })

    it('accepts loading prop', () => {
      // Component should accept loading prop without errors
      const loadingWrapper = mountVesselList({ loading: true })
      expect(loadingWrapper.exists()).toBe(true)

      const notLoadingWrapper = mountVesselList({ loading: false })
      expect(notLoadingWrapper.exists()).toBe(true)
    })

    it('shows refresh button', () => {
      const wrapper = mountVesselList()
      expect(wrapper.text()).toContain('Refresh')
    })
  })

  describe('empty state', () => {
    it('shows empty state message when no vessels and not loading', () => {
      const wrapper = mountVesselList({ vessels: [], loading: false })
      expect(wrapper.text()).toContain('No active vessels')
      expect(wrapper.text()).toContain('Register vessels in All Vessels')
    })
  })

  describe('occupancy status', () => {
    it('shows "Occupied" chip for vessels with occupancy', () => {
      const vessels = [createVessel({ uuid: 'vessel-uuid-1' })]
      const occupancies = [createOccupancy({ vessel_uuid: 'vessel-uuid-1' })]
      const wrapper = mountVesselList({ vessels, occupancies })

      expect(wrapper.text()).toContain('Occupied')
    })

    it('shows "Available" chip for vessels without occupancy', () => {
      const vessels = [createVessel({ uuid: 'vessel-uuid-1' })]
      const wrapper = mountVesselList({ vessels, occupancies: [] })

      expect(wrapper.text()).toContain('Available')
    })

    it('correctly identifies occupied vs available vessels', () => {
      const vessels = [
        createVessel({ uuid: 'vessel-uuid-1', name: 'FV-01' }),
        createVessel({ uuid: 'vessel-uuid-2', name: 'FV-02' }),
      ]
      const occupancies = [createOccupancy({ vessel_uuid: 'vessel-uuid-1' })]
      const wrapper = mountVesselList({ vessels, occupancies })

      const chips = wrapper.findAll('.v-chip')
      // Should have one Occupied and one Available
      const chipTexts = chips.map(c => c.text())
      expect(chipTexts).toContain('Occupied')
      expect(chipTexts).toContain('Available')
    })
  })

  describe('sorting', () => {
    it('sorts occupied vessels before available vessels', () => {
      const vessels = [
        createVessel({ uuid: 'vessel-uuid-1', name: 'AAA-Available' }),
        createVessel({ uuid: 'vessel-uuid-2', name: 'BBB-Occupied' }),
        createVessel({ uuid: 'vessel-uuid-3', name: 'CCC-Available' }),
      ]
      const occupancies = [createOccupancy({ vessel_uuid: 'vessel-uuid-2' })]
      const wrapper = mountVesselList({ vessels, occupancies })

      const listItems = wrapper.findAll('.v-list-item')
      // First item should be the occupied vessel (BBB-Occupied)
      expect(listItems[0].text()).toContain('BBB-Occupied')
    })

    it('sorts alphabetically within same occupancy group', () => {
      const vessels = [
        createVessel({ uuid: 'vessel-uuid-1', name: 'Zebra' }),
        createVessel({ uuid: 'vessel-uuid-2', name: 'Alpha' }),
        createVessel({ uuid: 'vessel-uuid-3', name: 'Beta' }),
      ]
      const wrapper = mountVesselList({ vessels, occupancies: [] })

      const listItems = wrapper.findAll('.v-list-item')
      expect(listItems[0].text()).toContain('Alpha')
      expect(listItems[1].text()).toContain('Beta')
      expect(listItems[2].text()).toContain('Zebra')
    })
  })

  describe('selection', () => {
    it('marks selected vessel as active', () => {
      const vessels = [
        createVessel({ uuid: 'vessel-uuid-1', name: 'FV-01' }),
        createVessel({ uuid: 'vessel-uuid-2', name: 'FV-02' }),
      ]
      const wrapper = mountVesselList({ vessels, selectedVesselUuid: 'vessel-uuid-2' })

      const listItems = wrapper.findAll('.v-list-item')
      // Find the item with FV-02 and check if it's active
      const fv02Item = listItems.find(item => item.text().includes('FV-02'))
      expect(fv02Item?.classes()).toContain('v-list-item--active')
    })

    it('emits select event when vessel is clicked', async () => {
      const vessels = [createVessel({ uuid: 'vessel-uuid-42', name: 'FV-01' })]
      const wrapper = mountVesselList({ vessels })

      const listItem = wrapper.find('.v-list-item')
      await listItem.trigger('click')

      expect(wrapper.emitted('select')).toBeTruthy()
      expect(wrapper.emitted('select')![0]).toEqual(['vessel-uuid-42'])
    })
  })

  describe('events', () => {
    it('emits refresh event when refresh button is clicked', async () => {
      const wrapper = mountVesselList()
      const refreshButton = wrapper.find('.v-card-title .v-btn')
      await refreshButton.trigger('click')

      expect(wrapper.emitted('refresh')).toBeTruthy()
      expect(wrapper.emitted('refresh')!.length).toBe(1)
    })
  })
})
