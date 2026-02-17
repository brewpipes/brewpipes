import { mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { nextTick } from 'vue'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import LotDetailDialog from '../LotDetailDialog.vue'
import type { Ingredient, IngredientLot } from '@/types'

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

// Mock useApiClient (used by useInventoryApi)
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
    currentRoute: { value: { path: '/inventory/ingredients' } },
  }),
  useRoute: () => ({
    path: '/inventory/ingredients',
  }),
}))

// --- Test data factories ---

function makeIngredient (overrides: Partial<Ingredient> = {}): Ingredient {
  return {
    uuid: 'ingredient-uuid-1',
    name: 'Pale Malt',
    category: 'fermentable',
    default_unit: 'lb',
    description: 'Base malt',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
    ...overrides,
  }
}

function makeLot (overrides: Partial<IngredientLot> = {}): IngredientLot {
  return {
    uuid: 'lot-uuid-1',
    ingredient_uuid: 'ingredient-uuid-1',
    receipt_uuid: null,
    supplier_uuid: null,
    purchase_order_line_uuid: null,
    brewery_lot_code: 'LOT-001',
    originator_lot_code: null,
    originator_name: null,
    originator_type: null,
    received_at: '2024-01-15T00:00:00Z',
    received_amount: 50,
    received_unit: 'lb',
    current_amount: 45,
    current_unit: 'lb',
    best_by_at: null,
    expires_at: null,
    notes: null,
    created_at: '2024-01-15T00:00:00Z',
    updated_at: '2024-01-15T00:00:00Z',
    ...overrides,
  }
}

/** Helper to flush all pending promises and Vue reactivity */
async function flushAll () {
  await vi.dynamicImportSettled()
  await nextTick()
  await new Promise(resolve => setTimeout(resolve, 10))
  await nextTick()
}

function mountDialog (props: {
  modelValue: boolean
  lot: IngredientLot | null
  ingredients: Ingredient[]
}) {
  const div = document.createElement('div')
  div.id = 'app'
  document.body.append(div)

  return mount(LotDetailDialog, {
    attachTo: div,
    global: {
      plugins: [vuetify],
    },
    props,
  })
}

/** Mount dialog closed, then open it (triggers the watch) */
async function mountAndOpen (
  lot: IngredientLot,
  ingredients: Ingredient[],
) {
  const wrapper = mountDialog({
    modelValue: false,
    lot,
    ingredients,
  })
  await wrapper.setProps({ modelValue: true })
  await flushAll()
  return wrapper
}

describe('LotDetailDialog', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Default: all detail fetches return 404 (no existing detail)
    mockRequest.mockRejectedValue(new Error('Not found'))
  })

  afterEach(() => {
    document.body.innerHTML = ''
  })

  describe('conditional rendering by category', () => {
    it('shows only malt detail section for fermentable category', async () => {
      const ingredient = makeIngredient({ category: 'fermentable', name: 'Pale Malt' })
      const lot = makeLot({ ingredient_uuid: ingredient.uuid })

      await mountAndOpen(lot, [ingredient])

      const text = document.body.textContent ?? ''
      expect(text).toContain('Malt lot detail')
      expect(text).not.toContain('Hop lot detail')
      expect(text).not.toContain('Yeast lot detail')
      expect(text).not.toContain('No specialized details')
    })

    it('shows only hop detail section for hop category', async () => {
      const ingredient = makeIngredient({ uuid: 'hop-ing', category: 'hop', name: 'Cascade' })
      const lot = makeLot({ uuid: 'hop-lot', ingredient_uuid: 'hop-ing' })

      await mountAndOpen(lot, [ingredient])

      const text = document.body.textContent ?? ''
      expect(text).toContain('Hop lot detail')
      expect(text).not.toContain('Malt lot detail')
      expect(text).not.toContain('Yeast lot detail')
      expect(text).not.toContain('No specialized details')
    })

    it('shows only yeast detail section for yeast category', async () => {
      const ingredient = makeIngredient({ uuid: 'yeast-ing', category: 'yeast', name: 'US-05' })
      const lot = makeLot({ uuid: 'yeast-lot', ingredient_uuid: 'yeast-ing' })

      await mountAndOpen(lot, [ingredient])

      const text = document.body.textContent ?? ''
      expect(text).toContain('Yeast lot detail')
      expect(text).not.toContain('Malt lot detail')
      expect(text).not.toContain('Hop lot detail')
      expect(text).not.toContain('No specialized details')
    })

    it('shows "no specialized details" for adjunct category', async () => {
      const ingredient = makeIngredient({ uuid: 'adj-ing', category: 'adjunct', name: 'Honey' })
      const lot = makeLot({ uuid: 'adj-lot', ingredient_uuid: 'adj-ing' })

      await mountAndOpen(lot, [ingredient])

      const text = document.body.textContent ?? ''
      expect(text).toContain('No specialized details')
      expect(text).not.toContain('Malt lot detail')
      expect(text).not.toContain('Hop lot detail')
      expect(text).not.toContain('Yeast lot detail')
    })

    it('shows "no specialized details" for salt category', async () => {
      const ingredient = makeIngredient({ uuid: 'salt-ing', category: 'salt', name: 'Gypsum' })
      const lot = makeLot({ uuid: 'salt-lot', ingredient_uuid: 'salt-ing' })

      await mountAndOpen(lot, [ingredient])

      const text = document.body.textContent ?? ''
      expect(text).toContain('No specialized details')
      expect(text).not.toContain('Malt lot detail')
      expect(text).not.toContain('Hop lot detail')
      expect(text).not.toContain('Yeast lot detail')
    })

    it('shows "no specialized details" for other category', async () => {
      const ingredient = makeIngredient({ uuid: 'other-ing', category: 'other', name: 'Misc' })
      const lot = makeLot({ uuid: 'other-lot', ingredient_uuid: 'other-ing' })

      await mountAndOpen(lot, [ingredient])

      const text = document.body.textContent ?? ''
      expect(text).toContain('No specialized details')
      expect(text).not.toContain('Malt lot detail')
      expect(text).not.toContain('Hop lot detail')
      expect(text).not.toContain('Yeast lot detail')
    })
  })

  describe('API calls by category', () => {
    it('fetches malt detail for fermentable lots', async () => {
      const ingredient = makeIngredient({ category: 'fermentable' })
      const lot = makeLot({ ingredient_uuid: ingredient.uuid })

      await mountAndOpen(lot, [ingredient])

      const maltCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('malt-detail'),
      )
      expect(maltCalls.length).toBe(1)

      // Should NOT call hop or yeast detail endpoints
      const hopCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('hop-detail'),
      )
      const yeastCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('yeast-detail'),
      )
      expect(hopCalls.length).toBe(0)
      expect(yeastCalls.length).toBe(0)
    })

    it('fetches hop detail for hop lots', async () => {
      const ingredient = makeIngredient({ uuid: 'hop-ing', category: 'hop', name: 'Cascade' })
      const lot = makeLot({ uuid: 'hop-lot', ingredient_uuid: 'hop-ing' })

      await mountAndOpen(lot, [ingredient])

      const hopCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('hop-detail'),
      )
      expect(hopCalls.length).toBe(1)

      // Should NOT call malt or yeast detail endpoints
      const maltCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('malt-detail'),
      )
      const yeastCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('yeast-detail'),
      )
      expect(maltCalls.length).toBe(0)
      expect(yeastCalls.length).toBe(0)
    })

    it('fetches yeast detail for yeast lots', async () => {
      const ingredient = makeIngredient({ uuid: 'yeast-ing', category: 'yeast', name: 'US-05' })
      const lot = makeLot({ uuid: 'yeast-lot', ingredient_uuid: 'yeast-ing' })

      await mountAndOpen(lot, [ingredient])

      const yeastCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('yeast-detail'),
      )
      expect(yeastCalls.length).toBe(1)

      // Should NOT call malt or hop detail endpoints
      const maltCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('malt-detail'),
      )
      const hopCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('hop-detail'),
      )
      expect(maltCalls.length).toBe(0)
      expect(hopCalls.length).toBe(0)
    })

    it('makes no detail API calls for adjunct lots', async () => {
      const ingredient = makeIngredient({ uuid: 'adj-ing', category: 'adjunct', name: 'Honey' })
      const lot = makeLot({ uuid: 'adj-lot', ingredient_uuid: 'adj-ing' })

      await mountAndOpen(lot, [ingredient])

      // Should NOT call any detail endpoints
      const detailCalls = mockRequest.mock.calls.filter(
        (call: unknown[]) => typeof call[0] === 'string' && (call[0] as string).includes('detail'),
      )
      expect(detailCalls.length).toBe(0)
    })
  })

  describe('ingredient name display', () => {
    it('shows the ingredient name in the summary', async () => {
      const ingredient = makeIngredient({ name: 'Cascade Hops', category: 'hop', uuid: 'hop-ing' })
      const lot = makeLot({ ingredient_uuid: 'hop-ing', brewery_lot_code: 'HOP-42' })

      await mountAndOpen(lot, [ingredient])

      const text = document.body.textContent ?? ''
      expect(text).toContain('Cascade Hops')
      expect(text).toContain('Lot HOP-42')
    })

    it('shows "Unknown Ingredient" when ingredient is not found', async () => {
      const lot = makeLot({ ingredient_uuid: 'nonexistent-uuid' })

      await mountAndOpen(lot, [])

      const text = document.body.textContent ?? ''
      expect(text).toContain('Unknown Ingredient')
    })
  })

  describe('close behavior', () => {
    it('emits update:modelValue false when close button is clicked', async () => {
      const ingredient = makeIngredient({ category: 'fermentable' })
      const lot = makeLot({ ingredient_uuid: ingredient.uuid })

      const wrapper = await mountAndOpen(lot, [ingredient])

      // Find the Close button in the card actions
      const buttons = document.body.querySelectorAll('.v-btn')
      let closeBtn: HTMLElement | null = null
      buttons.forEach(btn => {
        if (btn.textContent?.trim() === 'Close') {
          closeBtn = btn as HTMLElement
        }
      })
      expect(closeBtn).toBeTruthy()
      closeBtn!.click()
      await nextTick()

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')![0]).toEqual([false])
    })
  })
})
