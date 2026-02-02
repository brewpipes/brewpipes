<template>
  <v-card class="sub-card" variant="outlined">
    <v-card-text>
      <v-card class="sub-card mb-4" variant="tonal">
        <v-card-text>
          <div class="d-flex flex-wrap align-center justify-space-between ga-2">
            <div class="text-subtitle-2 font-weight-semibold">Quick Update</div>
            <v-btn
              append-icon="mdi-arrow-right"
              size="small"
              variant="text"
              @click="emit('open-extended')"
            >
              More
            </v-btn>
          </div>
          <v-divider class="my-3" />

          <v-row align="center" dense>
            <v-col class="d-flex align-center justify-center" cols="12" md="1">
              <v-menu
                v-model="observedAtMenu"
                :close-on-content-click="false"
                location="bottom"
              >
                <template #activator="{ props: menuProps }">
                  <v-tooltip location="top">
                    <template #activator="{ props: tooltipProps }">
                      <v-btn
                        v-bind="{ ...menuProps, ...tooltipProps }"
                        aria-label="Set observation time"
                        :color="reading.observed_at ? 'secondary' : 'primary'"
                        icon="mdi-clock-outline"
                        size="default"
                        variant="text"
                      />
                    </template>
                    <span>{{ observedAtLabel }}</span>
                  </v-tooltip>
                </template>
                <v-card>
                  <v-card-text>
                    <v-text-field
                      density="compact"
                      label="Observed at"
                      :model-value="reading.observed_at"
                      type="datetime-local"
                      @update:model-value="emit('update:reading', { ...reading, observed_at: $event })"
                    />
                  </v-card-text>
                  <v-card-actions class="justify-end">
                    <v-btn variant="text" @click="clearObservedAt">
                      Use now
                    </v-btn>
                    <v-btn variant="text" @click="observedAtMenu = false">Done</v-btn>
                  </v-card-actions>
                </v-card>
              </v-menu>
            </v-col>
            <v-col cols="12" md="2">
              <v-text-field
                density="compact"
                inputmode="decimal"
                label="Temp"
                :model-value="reading.temperature"
                :placeholder="temperaturePlaceholder"
                @update:model-value="emit('update:reading', { ...reading, temperature: $event })"
              />
            </v-col>
            <v-col cols="12" md="2">
              <v-text-field
                density="compact"
                inputmode="decimal"
                label="Gravity"
                :model-value="reading.gravity"
                :placeholder="gravityPlaceholder"
                @update:model-value="emit('update:reading', { ...reading, gravity: $event })"
              />
            </v-col>
            <v-col cols="12" md="5">
              <v-text-field
                density="compact"
                label="Notes"
                :model-value="reading.notes"
                placeholder="Aroma, flavor, observations"
                @update:model-value="emit('update:reading', { ...reading, notes: $event })"
              />
            </v-col>
            <v-col class="d-flex align-center justify-end" cols="12" md="1">
              <v-btn
                color="primary"
                :disabled="!readingReady"
                @click="emit('record')"
              >
                Add
              </v-btn>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>

      <v-timeline align="start" density="compact" side="end">
        <v-timeline-item
          v-for="event in events"
          :key="event.id"
          :dot-color="event.color"
          :icon="event.icon"
        >
          <template #opposite>
            <div class="text-caption text-medium-emphasis">
              {{ formatDateTime(event.at) }}
            </div>
          </template>
          <div class="text-subtitle-2 font-weight-semibold">
            {{ event.title }}
          </div>
          <div class="text-body-2 text-medium-emphasis">
            {{ event.subtitle }}
          </div>
        </v-timeline-item>

        <v-timeline-item v-if="events.length === 0" dot-color="grey">
          <div class="text-body-2 text-medium-emphasis">
            No timeline events yet.
          </div>
        </v-timeline-item>
      </v-timeline>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { TimelineEvent } from './types'
  import { computed, ref, watch } from 'vue'
  import { useFormatters } from '@/composables/useFormatters'

  export type TimelineReading = {
    observed_at: string
    temperature: string
    gravity: string
    notes: string
  }

  const props = defineProps<{
    events: TimelineEvent[]
    reading: TimelineReading
    readingReady: boolean
    temperatureUnit: 'f' | 'c'
    gravityUnit: 'sg' | 'plato'
  }>()

  const emit = defineEmits<{
    'update:reading': [reading: TimelineReading]
    'record': []
    'open-extended': []
  }>()

  const { formatDateTime } = useFormatters()

  const observedAtMenu = ref(false)

  const observedAtLabel = computed(() =>
    props.reading.observed_at ? formatDateTime(props.reading.observed_at) : 'Now',
  )

  const temperaturePlaceholder = computed(() =>
    props.temperatureUnit === 'f' ? '67F' : '19C',
  )

  const gravityPlaceholder = computed(() =>
    props.gravityUnit === 'sg' ? '1.056' : '13.8',
  )

  watch(observedAtMenu, isOpen => {
    if (isOpen && !props.reading.observed_at) {
      emit('update:reading', { ...props.reading, observed_at: nowInputValue() })
    }
  })

  function clearObservedAt () {
    emit('update:reading', { ...props.reading, observed_at: '' })
    observedAtMenu.value = false
  }

  function nowInputValue () {
    const now = new Date()
    const pad = (value: number) => String(value).padStart(2, '0')
    const year = now.getFullYear()
    const month = pad(now.getMonth() + 1)
    const day = pad(now.getDate())
    const hours = pad(now.getHours())
    const minutes = pad(now.getMinutes())
    return `${year}-${month}-${day}T${hours}:${minutes}`
  }
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}
</style>
