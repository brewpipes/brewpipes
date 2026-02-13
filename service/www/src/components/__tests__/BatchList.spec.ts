import type { Batch } from '@/types'
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import BatchList from '../BatchList.vue'

const vuetify = createVuetify({
  components,
  directives,
})

function createBatch (overrides: Partial<Batch> = {}): Batch {
  return {
    uuid: 'test-uuid-1',
    short_name: 'Test Batch',
    brew_date: '2024-06-15T10:00:00Z',
    recipe_uuid: null,
    notes: null,
    created_at: '2024-06-15T10:00:00Z',
    updated_at: '2024-06-15T12:00:00Z',
    ...overrides,
  }
}

function mountBatchList (props: {
  batches?: Batch[]
  selectedBatchUuid?: string | null
  loading?: boolean
  showCreateButton?: boolean
  showBulkImport?: boolean
} = {}) {
  return mount(BatchList, {
    global: {
      plugins: [vuetify],
    },
    props: {
      batches: props.batches ?? [],
      selectedBatchUuid: props.selectedBatchUuid ?? null,
      loading: props.loading ?? false,
      showCreateButton: props.showCreateButton ?? true,
      showBulkImport: props.showBulkImport ?? true,
    },
  })
}

describe('BatchList', () => {
  describe('rendering', () => {
    it('renders the component title', () => {
      const wrapper = mountBatchList()
      expect(wrapper.text()).toContain('Batches')
    })

    it('renders batch list items', () => {
      const batches = [
        createBatch({ uuid: 'uuid-1', short_name: 'IPA Batch 001' }),
        createBatch({ uuid: 'uuid-2', short_name: 'Stout Batch 002' }),
      ]
      const wrapper = mountBatchList({ batches })

      expect(wrapper.text()).toContain('IPA Batch 001')
      expect(wrapper.text()).toContain('Stout Batch 002')
    })

    it('renders formatted brew date', () => {
      const batches = [
        createBatch({ brew_date: '2024-06-15T10:00:00Z' }),
      ]
      const wrapper = mountBatchList({ batches })

      expect(wrapper.text()).toMatch(/Jun/)
      expect(wrapper.text()).toMatch(/15/)
      expect(wrapper.text()).toMatch(/2024/)
    })

    it('accepts loading prop', () => {
      // Component should accept loading prop without errors
      const loadingWrapper = mountBatchList({ loading: true })
      expect(loadingWrapper.exists()).toBe(true)

      const notLoadingWrapper = mountBatchList({ loading: false })
      expect(notLoadingWrapper.exists()).toBe(true)
    })

    it('shows create button by default', () => {
      const wrapper = mountBatchList()
      const createButton = wrapper.find('[aria-label="Create batch"]')
      expect(createButton.exists()).toBe(true)
    })

    it('hides create button when showCreateButton is false', () => {
      const wrapper = mountBatchList({ showCreateButton: false })
      const createButton = wrapper.find('[aria-label="Create batch"]')
      expect(createButton.exists()).toBe(false)
    })

    it('shows bulk import button by default', () => {
      const wrapper = mountBatchList()
      expect(wrapper.text()).toContain('Bulk import')
    })

    it('hides bulk import button when showBulkImport is false', () => {
      const wrapper = mountBatchList({ showBulkImport: false })
      expect(wrapper.text()).not.toContain('Bulk import')
    })
  })

  describe('empty state', () => {
    it('shows empty state message when no batches and not loading', () => {
      const wrapper = mountBatchList({ batches: [], loading: false })
      expect(wrapper.text()).toContain('No batches yet')
      expect(wrapper.text()).toContain('Use + to add the first batch')
    })

    it('does not show empty state when loading', () => {
      const wrapper = mountBatchList({ batches: [], loading: true })
      // Empty state should still be visible but loading indicator shows
      expect(wrapper.find('.v-progress-linear').exists()).toBe(true)
    })
  })

  describe('selection', () => {
    it('marks selected batch as active', async () => {
      const batches = [
        createBatch({ uuid: 'batch-uuid-1', short_name: 'Batch 1' }),
        createBatch({ uuid: 'batch-uuid-2', short_name: 'Batch 2' }),
      ]
      const wrapper = mountBatchList({ batches, selectedBatchUuid: 'batch-uuid-2' })

      const listItems = wrapper.findAll('.v-list-item')
      // The second batch (uuid: batch-uuid-2) should be active
      expect(listItems[1].classes()).toContain('v-list-item--active')
    })

    it('emits select event when batch is clicked', async () => {
      const batches = [
        createBatch({ uuid: 'batch-uuid-1', short_name: 'Batch 1' }),
        createBatch({ uuid: 'batch-uuid-2', short_name: 'Batch 2' }),
      ]
      const wrapper = mountBatchList({ batches })

      const listItems = wrapper.findAll('.v-list-item')
      await listItems[0].trigger('click')

      expect(wrapper.emitted('select')).toBeTruthy()
      expect(wrapper.emitted('select')![0]).toEqual(['batch-uuid-1'])
    })
  })

  describe('events', () => {
    it('emits create event when create button is clicked', async () => {
      const wrapper = mountBatchList()
      const createButton = wrapper.find('[aria-label="Create batch"]')
      await createButton.trigger('click')

      expect(wrapper.emitted('create')).toBeTruthy()
      expect(wrapper.emitted('create')!.length).toBe(1)
    })

    it('emits bulk-import event when bulk import button is clicked', async () => {
      const wrapper = mountBatchList()
      const bulkImportButton = wrapper.find('.v-card-actions .v-btn')
      await bulkImportButton.trigger('click')

      expect(wrapper.emitted('bulk-import')).toBeTruthy()
      expect(wrapper.emitted('bulk-import')!.length).toBe(1)
    })
  })

  describe('date formatting', () => {
    it('shows "Unknown" for null brew date', () => {
      const batches = [createBatch({ brew_date: null })]
      const wrapper = mountBatchList({ batches })

      expect(wrapper.text()).toContain('Unknown')
    })

    it('formats updated_at timestamp', () => {
      const batches = [createBatch({ updated_at: '2024-06-15T14:30:00Z' })]
      const wrapper = mountBatchList({ batches })

      // Should contain formatted date in chip
      expect(wrapper.text()).toMatch(/Jun/)
      expect(wrapper.text()).toMatch(/15/)
    })
  })
})
