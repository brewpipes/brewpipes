<template>
  <v-container class="pa-4" fluid>
    <v-row>
      <v-col cols="12">
        <p class="text-body-2 text-medium-emphasis mb-4">
          Review the receiving details below and confirm.
        </p>
      </v-col>
    </v-row>

    <!-- Error alert -->
    <v-row v-if="error">
      <v-col cols="12">
        <v-alert
          closable
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ error }}
        </v-alert>
      </v-col>
    </v-row>

    <!-- Summary card -->
    <v-row>
      <v-col cols="12">
        <v-card variant="outlined">
          <v-card-title class="text-subtitle-1">Receipt Summary</v-card-title>
          <v-card-text>
            <v-list density="compact" lines="two">
              <v-list-item>
                <v-list-item-title>Supplier</v-list-item-title>
                <v-list-item-subtitle>{{ supplierName }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>Items to receive</v-list-item-title>
                <v-list-item-subtitle>{{ lines.length }} line item(s)</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>New PO status</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip
                    :color="newStatus === 'received' ? 'success' : 'warning'"
                    density="compact"
                    size="small"
                    variant="tonal"
                  >
                    {{ newStatus === 'received' ? 'Fully Received' : 'Partially Received' }}
                  </v-chip>
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Line items table -->
    <v-row>
      <v-col cols="12">
        <v-card variant="outlined">
          <v-card-title class="text-subtitle-1">Items Being Received</v-card-title>
          <v-card-text class="pa-0">
            <v-table density="compact">
              <thead>
                <tr>
                  <th>Item</th>
                  <th class="text-right">Quantity</th>
                  <th>Location</th>
                  <th class="d-none d-sm-table-cell">Brewery Lot</th>
                  <th class="d-none d-sm-table-cell">Supplier Lot</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="line in lines" :key="line.uuid">
                  <td>
                    <div class="font-weight-medium">{{ line.item_name }}</div>
                    <div class="text-body-2 text-medium-emphasis">
                      {{ getIngredientName(line.inventory_item_uuid) }}
                    </div>
                  </td>
                  <td class="text-right">
                    {{ getDetail(line.uuid)?.quantity ?? 0 }}
                    {{ getDetail(line.uuid)?.unit ?? line.quantity_unit }}
                  </td>
                  <td>{{ getLocationName(getDetail(line.uuid)?.locationUuid ?? '') }}</td>
                  <td class="d-none d-sm-table-cell">
                    {{ getDetail(line.uuid)?.breweryLotCode || '-' }}
                  </td>
                  <td class="d-none d-sm-table-cell">
                    {{ getDetail(line.uuid)?.supplierLotCode || '-' }}
                  </td>
                </tr>
              </tbody>
            </v-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Mobile: Show lot codes in expandable section -->
    <v-row v-if="$vuetify.display.xs" class="mt-2">
      <v-col cols="12">
        <v-expansion-panels variant="accordion">
          <v-expansion-panel title="Lot Code Details">
            <v-expansion-panel-text>
              <v-list density="compact">
                <v-list-item v-for="line in lines" :key="line.uuid">
                  <v-list-item-title>{{ line.item_name }}</v-list-item-title>
                  <v-list-item-subtitle>
                    Brewery: {{ getDetail(line.uuid)?.breweryLotCode || 'None' }}
                    <br>
                    Supplier: {{ getDetail(line.uuid)?.supplierLotCode || 'None' }}
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-col>
    </v-row>

    <!-- Notes -->
    <v-row>
      <v-col cols="12">
        <v-textarea
          auto-grow
          density="comfortable"
          hint="Optional notes about this receipt"
          label="Receipt notes"
          :model-value="notes"
          persistent-hint
          rows="2"
          @update:model-value="emit('update:notes', $event)"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
  import type { Ingredient, LineReceivingDetails, PurchaseOrderLine, StockLocation } from '@/types'

  const props = defineProps<{
    lines: PurchaseOrderLine[]
    lineDetails: LineReceivingDetails[]
    stockLocations: StockLocation[]
    ingredients: Ingredient[]
    supplierName: string
    newStatus: string
    notes: string
    error: string
  }>()

  const emit = defineEmits<{
    'update:notes': [value: string]
  }>()

  function getDetail (lineUuid: string): LineReceivingDetails | undefined {
    return props.lineDetails.find(d => d.lineUuid === lineUuid)
  }

  function getIngredientName (ingredientUuid: string | null): string {
    if (!ingredientUuid) return 'Unknown'
    return props.ingredients.find(i => i.uuid === ingredientUuid)?.name ?? 'Unknown'
  }

  function getLocationName (locationUuid: string): string {
    return props.stockLocations.find(l => l.uuid === locationUuid)?.name ?? 'Unknown'
  }
</script>
