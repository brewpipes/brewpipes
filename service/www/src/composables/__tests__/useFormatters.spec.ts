import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  useAdditionTypeFormatters,
  useFeeTypeFormatters,
  useFormatters,
  useLineItemTypeFormatters,
  useOccupancyStatusFormatters,
  usePurchaseOrderStatusFormatters,
  useVesselStatusFormatters,
  type VesselStatus,
} from '../useFormatters'

describe('useFormatters', () => {
  const { formatDateTime, formatDate, formatRelativeTime } = useFormatters()

  describe('formatDateTime', () => {
    it('formats a valid date/time string', () => {
      const result = formatDateTime('2024-06-15T14:30:00Z')
      // The exact format depends on locale, but should contain date and time parts
      expect(result).toMatch(/Jun/)
      expect(result).toMatch(/15/)
      expect(result).toMatch(/2024/)
    })

    it('returns "Unknown" for null input', () => {
      expect(formatDateTime(null)).toBe('Unknown')
    })

    it('returns "Unknown" for undefined input', () => {
      expect(formatDateTime(undefined)).toBe('Unknown')
    })

    it('returns "Unknown" for empty string', () => {
      expect(formatDateTime('')).toBe('Unknown')
    })
  })

  describe('formatDate', () => {
    it('formats a valid date string', () => {
      const result = formatDate('2024-06-15T14:30:00Z')
      expect(result).toMatch(/Jun/)
      expect(result).toMatch(/15/)
      expect(result).toMatch(/2024/)
    })

    it('returns "Unknown" for null input', () => {
      expect(formatDate(null)).toBe('Unknown')
    })

    it('returns "Unknown" for undefined input', () => {
      expect(formatDate(undefined)).toBe('Unknown')
    })

    it('returns "Unknown" for empty string', () => {
      expect(formatDate('')).toBe('Unknown')
    })
  })

  describe('formatRelativeTime', () => {
    beforeEach(() => {
      // Mock Date.now() to a fixed time for consistent testing
      vi.useFakeTimers()
      vi.setSystemTime(new Date('2024-06-15T12:00:00Z'))
    })

    afterEach(() => {
      vi.useRealTimers()
    })

    it('returns "just now" for times less than 60 seconds ago', () => {
      const result = formatRelativeTime('2024-06-15T11:59:30Z')
      expect(result).toBe('just now')
    })

    it('returns minutes ago for times less than 60 minutes ago', () => {
      const result = formatRelativeTime('2024-06-15T11:45:00Z')
      expect(result).toBe('15m ago')
    })

    it('returns hours ago for times less than 24 hours ago', () => {
      const result = formatRelativeTime('2024-06-15T09:00:00Z')
      expect(result).toBe('3h ago')
    })

    it('returns days ago for times less than 7 days ago', () => {
      const result = formatRelativeTime('2024-06-13T12:00:00Z')
      expect(result).toBe('2d ago')
    })

    it('returns formatted date for times 7 or more days ago', () => {
      const result = formatRelativeTime('2024-06-01T12:00:00Z')
      expect(result).toMatch(/Jun/)
      expect(result).toMatch(/1/)
      expect(result).toMatch(/2024/)
    })

    it('returns "Unknown" for null input', () => {
      expect(formatRelativeTime(null)).toBe('Unknown')
    })

    it('returns "Unknown" for undefined input', () => {
      expect(formatRelativeTime(undefined)).toBe('Unknown')
    })

    it('returns "Unknown" for empty string', () => {
      expect(formatRelativeTime('')).toBe('Unknown')
    })
  })
})

describe('useVesselStatusFormatters', () => {
  const { formatVesselStatus, getVesselStatusColor } = useVesselStatusFormatters()

  describe('formatVesselStatus', () => {
    it('formats "active" status', () => {
      expect(formatVesselStatus('active')).toBe('Active')
    })

    it('formats "inactive" status', () => {
      expect(formatVesselStatus('inactive')).toBe('Inactive')
    })

    it('formats "retired" status', () => {
      expect(formatVesselStatus('retired')).toBe('Retired')
    })

    it('returns the original value for unknown status', () => {
      expect(formatVesselStatus('unknown' as VesselStatus)).toBe('unknown')
    })
  })

  describe('getVesselStatusColor', () => {
    it('returns "success" for active status', () => {
      expect(getVesselStatusColor('active')).toBe('success')
    })

    it('returns "grey" for inactive status', () => {
      expect(getVesselStatusColor('inactive')).toBe('grey')
    })

    it('returns "error" for retired status', () => {
      expect(getVesselStatusColor('retired')).toBe('error')
    })

    it('returns "secondary" for unknown status', () => {
      expect(getVesselStatusColor('unknown' as VesselStatus)).toBe('secondary')
    })
  })
})

describe('useAdditionTypeFormatters', () => {
  const { formatAdditionType } = useAdditionTypeFormatters()

  describe('formatAdditionType', () => {
    it('formats known addition types correctly', () => {
      expect(formatAdditionType('hop')).toBe('Hop')
      expect(formatAdditionType('malt')).toBe('Malt')
      expect(formatAdditionType('yeast')).toBe('Yeast')
      expect(formatAdditionType('adjunct')).toBe('Adjunct')
      expect(formatAdditionType('water_chem')).toBe('Water Chemistry')
      expect(formatAdditionType('fining')).toBe('Fining')
      expect(formatAdditionType('gas')).toBe('Gas')
      expect(formatAdditionType('other')).toBe('Other')
    })

    it('capitalizes and formats unknown types with underscores', () => {
      expect(formatAdditionType('dry_hop')).toBe('Dry hop')
    })

    it('capitalizes simple unknown types', () => {
      expect(formatAdditionType('enzyme')).toBe('Enzyme')
    })
  })
})

describe('useOccupancyStatusFormatters', () => {
  const { formatOccupancyStatus, getOccupancyStatusColor, getOccupancyStatusIcon }
    = useOccupancyStatusFormatters()

  describe('formatOccupancyStatus', () => {
    it('formats known statuses correctly', () => {
      expect(formatOccupancyStatus('fermenting')).toBe('Fermenting')
      expect(formatOccupancyStatus('conditioning')).toBe('Conditioning')
      expect(formatOccupancyStatus('cold_crashing')).toBe('Cold Crashing')
      expect(formatOccupancyStatus('dry_hopping')).toBe('Dry Hopping')
      expect(formatOccupancyStatus('carbonating')).toBe('Carbonating')
      expect(formatOccupancyStatus('holding')).toBe('Holding')
      expect(formatOccupancyStatus('packaging')).toBe('Packaging')
    })

    it('returns "No status" for null input', () => {
      expect(formatOccupancyStatus(null)).toBe('No status')
    })

    it('returns "No status" for undefined input', () => {
      expect(formatOccupancyStatus(undefined)).toBe('No status')
    })

    it('capitalizes and formats unknown status with underscores', () => {
      expect(formatOccupancyStatus('some_unknown_status')).toBe('Some unknown status')
    })

    it('capitalizes simple unknown status', () => {
      expect(formatOccupancyStatus('testing')).toBe('Testing')
    })
  })

  describe('getOccupancyStatusColor', () => {
    it('returns correct colors for known statuses', () => {
      expect(getOccupancyStatusColor('fermenting')).toBe('orange')
      expect(getOccupancyStatusColor('conditioning')).toBe('blue')
      expect(getOccupancyStatusColor('cold_crashing')).toBe('cyan')
      expect(getOccupancyStatusColor('dry_hopping')).toBe('green')
      expect(getOccupancyStatusColor('carbonating')).toBe('purple')
      expect(getOccupancyStatusColor('holding')).toBe('grey')
      expect(getOccupancyStatusColor('packaging')).toBe('teal')
    })

    it('returns "grey" for null input', () => {
      expect(getOccupancyStatusColor(null)).toBe('grey')
    })

    it('returns "grey" for undefined input', () => {
      expect(getOccupancyStatusColor(undefined)).toBe('grey')
    })

    it('returns "secondary" for unknown status', () => {
      expect(getOccupancyStatusColor('unknown_status')).toBe('secondary')
    })
  })

  describe('getOccupancyStatusIcon', () => {
    it('returns correct icons for known statuses', () => {
      expect(getOccupancyStatusIcon('fermenting')).toBe('mdi-molecule')
      expect(getOccupancyStatusIcon('conditioning')).toBe('mdi-clock-outline')
      expect(getOccupancyStatusIcon('cold_crashing')).toBe('mdi-snowflake')
      expect(getOccupancyStatusIcon('dry_hopping')).toBe('mdi-leaf')
      expect(getOccupancyStatusIcon('carbonating')).toBe('mdi-shimmer')
      expect(getOccupancyStatusIcon('holding')).toBe('mdi-pause-circle-outline')
      expect(getOccupancyStatusIcon('packaging')).toBe('mdi-package-variant')
    })

    it('returns help icon for null input', () => {
      expect(getOccupancyStatusIcon(null)).toBe('mdi-help-circle-outline')
    })

    it('returns help icon for undefined input', () => {
      expect(getOccupancyStatusIcon(undefined)).toBe('mdi-help-circle-outline')
    })

    it('returns circle icon for unknown status', () => {
      expect(getOccupancyStatusIcon('unknown_status')).toBe('mdi-circle')
    })
  })
})

describe('usePurchaseOrderStatusFormatters', () => {
  const {
    formatPurchaseOrderStatus,
    getPurchaseOrderStatusColor,
    purchaseOrderStatusOptions,
  } = usePurchaseOrderStatusFormatters()

  describe('formatPurchaseOrderStatus', () => {
    it('formats known statuses correctly', () => {
      expect(formatPurchaseOrderStatus('draft')).toBe('Draft')
      expect(formatPurchaseOrderStatus('submitted')).toBe('Submitted')
      expect(formatPurchaseOrderStatus('confirmed')).toBe('Confirmed')
      expect(formatPurchaseOrderStatus('partially_received')).toBe('Partially Received')
      expect(formatPurchaseOrderStatus('received')).toBe('Received')
      expect(formatPurchaseOrderStatus('cancelled')).toBe('Cancelled')
    })

    it('capitalizes and formats unknown status with underscores', () => {
      expect(formatPurchaseOrderStatus('on_hold')).toBe('On hold')
    })

    it('capitalizes simple unknown status', () => {
      expect(formatPurchaseOrderStatus('pending')).toBe('Pending')
    })
  })

  describe('getPurchaseOrderStatusColor', () => {
    it('returns correct colors for known statuses', () => {
      expect(getPurchaseOrderStatusColor('draft')).toBe('grey')
      expect(getPurchaseOrderStatusColor('submitted')).toBe('blue')
      expect(getPurchaseOrderStatusColor('confirmed')).toBe('green')
      expect(getPurchaseOrderStatusColor('partially_received')).toBe('orange')
      expect(getPurchaseOrderStatusColor('received')).toBe('green')
      expect(getPurchaseOrderStatusColor('cancelled')).toBe('red')
    })

    it('returns "grey" for unknown status', () => {
      expect(getPurchaseOrderStatusColor('unknown')).toBe('grey')
    })
  })

  describe('purchaseOrderStatusOptions', () => {
    it('returns options with title and value for each status', () => {
      expect(purchaseOrderStatusOptions).toEqual([
        { title: 'Draft', value: 'draft' },
        { title: 'Submitted', value: 'submitted' },
        { title: 'Confirmed', value: 'confirmed' },
        { title: 'Partially Received', value: 'partially_received' },
        { title: 'Received', value: 'received' },
        { title: 'Cancelled', value: 'cancelled' },
      ])
    })
  })
})

describe('useLineItemTypeFormatters', () => {
  const { formatLineItemType, lineItemTypeOptions } = useLineItemTypeFormatters()

  describe('formatLineItemType', () => {
    it('formats known types correctly', () => {
      expect(formatLineItemType('ingredient')).toBe('Ingredient')
      expect(formatLineItemType('packaging')).toBe('Packaging')
      expect(formatLineItemType('service')).toBe('Service')
      expect(formatLineItemType('equipment')).toBe('Equipment')
      expect(formatLineItemType('other')).toBe('Other')
    })

    it('capitalizes and formats unknown types with underscores', () => {
      expect(formatLineItemType('raw_material')).toBe('Raw material')
    })

    it('capitalizes simple unknown types', () => {
      expect(formatLineItemType('chemical')).toBe('Chemical')
    })
  })

  describe('lineItemTypeOptions', () => {
    it('returns options with title and value for each type', () => {
      expect(lineItemTypeOptions).toEqual([
        { title: 'Ingredient', value: 'ingredient' },
        { title: 'Packaging', value: 'packaging' },
        { title: 'Service', value: 'service' },
        { title: 'Equipment', value: 'equipment' },
        { title: 'Other', value: 'other' },
      ])
    })
  })
})

describe('useFeeTypeFormatters', () => {
  const { formatFeeType, feeTypeOptions } = useFeeTypeFormatters()

  describe('formatFeeType', () => {
    it('formats known types correctly', () => {
      expect(formatFeeType('shipping')).toBe('Shipping')
      expect(formatFeeType('handling')).toBe('Handling')
      expect(formatFeeType('tax')).toBe('Tax')
      expect(formatFeeType('insurance')).toBe('Insurance')
      expect(formatFeeType('customs')).toBe('Customs')
      expect(formatFeeType('freight')).toBe('Freight')
      expect(formatFeeType('hazmat')).toBe('Hazmat')
      expect(formatFeeType('other')).toBe('Other')
    })

    it('capitalizes and formats unknown types with underscores', () => {
      expect(formatFeeType('cold_storage')).toBe('Cold storage')
    })

    it('capitalizes simple unknown types', () => {
      expect(formatFeeType('surcharge')).toBe('Surcharge')
    })
  })

  describe('feeTypeOptions', () => {
    it('returns options with title and value for each type', () => {
      expect(feeTypeOptions).toEqual([
        { title: 'Shipping', value: 'shipping' },
        { title: 'Handling', value: 'handling' },
        { title: 'Tax', value: 'tax' },
        { title: 'Insurance', value: 'insurance' },
        { title: 'Customs', value: 'customs' },
        { title: 'Freight', value: 'freight' },
        { title: 'Hazmat', value: 'hazmat' },
        { title: 'Other', value: 'other' },
      ])
    })
  })
})
