<template>
  <v-container class="pa-4" fluid>
    <v-row>
      <v-col cols="12">
        <p class="text-body-2 text-medium-emphasis mb-4">
          Enter receiving details for each selected line item.
        </p>
      </v-col>
    </v-row>

    <!-- Mobile: Expandable panels -->
    <template v-if="$vuetify.display.xs">
      <v-expansion-panels v-model="expandedPanel" variant="accordion">
        <v-expansion-panel
          v-for="(orderLine, index) in lines"
          :key="orderLine.uuid"
          :value="index"
        >
          <v-expansion-panel-title>
            <div class="d-flex flex-column">
              <span class="font-weight-medium">{{ orderLine.item_name }}</span>
              <span class="text-body-2 text-medium-emphasis">
                {{ getDetail(orderLine.uuid)?.quantity ?? 0 }} {{ orderLine.quantity_unit }}
                <span v-if="getDetail(orderLine.uuid)?.locationUuid">
                  @ {{ getLocationName(getDetail(orderLine.uuid)?.locationUuid ?? '') }}
                </span>
              </span>
            </div>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <LineDetailForm
              :detail="getDetail(orderLine.uuid)"
              :ingredient-name="getIngredientName(orderLine.inventory_item_uuid)"
              :line="orderLine"
              :previously-received="previouslyReceived.get(orderLine.uuid) ?? 0"
              :stock-locations="stockLocations"
              @update:detail="updateDetail(orderLine.uuid, $event)"
            />
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </template>

    <!-- Desktop: Cards -->
    <template v-else>
      <v-row v-for="orderLine in lines" :key="orderLine.uuid">
        <v-col cols="12">
          <v-card variant="outlined">
            <v-card-title class="text-subtitle-1 pb-0">
              {{ orderLine.item_name }}
              <span class="text-body-2 text-medium-emphasis ml-2">
                ({{ getIngredientName(orderLine.inventory_item_uuid) }})
              </span>
            </v-card-title>
            <v-card-text>
              <LineDetailForm
                :detail="getDetail(orderLine.uuid)"
                :ingredient-name="getIngredientName(orderLine.inventory_item_uuid)"
                :line="orderLine"
                :previously-received="previouslyReceived.get(orderLine.uuid) ?? 0"
                :stock-locations="stockLocations"
                @update:detail="updateDetail(orderLine.uuid, $event)"
              />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </v-container>
</template>

<script lang="ts" setup>
  import type { Ingredient, LineReceivingDetails, PurchaseOrderLine, StockLocation } from '@/types'
  import { ref } from 'vue'
  import LineDetailForm from './LineDetailForm.vue'

  const props = defineProps<{
    lines: PurchaseOrderLine[]
    lineDetails: LineReceivingDetails[]
    stockLocations: StockLocation[]
    ingredients: Ingredient[]
    previouslyReceived: Map<string, number>
  }>()

  const emit = defineEmits<{
    'update:lineDetails': [value: LineReceivingDetails[]]
  }>()

  const expandedPanel = ref<number | undefined>(0)

  function getDetail (lineUuid: string): LineReceivingDetails | undefined {
    return props.lineDetails.find(d => d.lineUuid === lineUuid)
  }

  function updateDetail (lineUuid: string, updates: Partial<LineReceivingDetails>) {
    const newDetails = props.lineDetails.map(d => {
      if (d.lineUuid === lineUuid) {
        return { ...d, ...updates }
      }
      return d
    })
    emit('update:lineDetails', newDetails)
  }

  function getIngredientName (ingredientUuid: string | null): string {
    if (!ingredientUuid) return 'Unknown'
    return props.ingredients.find(i => i.uuid === ingredientUuid)?.name ?? 'Unknown'
  }

  function getLocationName (locationUuid: string): string {
    return props.stockLocations.find(l => l.uuid === locationUuid)?.name ?? ''
  }
</script>
