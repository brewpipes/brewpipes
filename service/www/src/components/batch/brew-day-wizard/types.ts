/**
 * Wizard-specific types for the Brew Day Wizard.
 */

import type { BrewSession, Occupancy } from '@/types'

/** Identifies each step in the wizard */
export type WizardStepId = 'pick' | 'session' | 'fermenter'

/** Status of a wizard step */
export type WizardStepStatus = 'not_started' | 'in_progress' | 'complete' | 'skipped'

/** A single step in the wizard */
export interface WizardStep {
  id: WizardStepId
  label: string
  icon: string
  status: WizardStepStatus
}

/** Summary data collected across wizard steps */
export interface WizardCompletionData {
  ingredientCount: number
  lotCount: number
  session: BrewSession | null
  mashTemp: string | null
  originalGravity: string | null
  occupancy: Occupancy | null
  fermenterName: string | null
}
