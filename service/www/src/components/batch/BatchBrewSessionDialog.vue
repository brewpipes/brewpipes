<template>
  <v-dialog max-width="680" :model-value="modelValue" persistent @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit brew session' : 'Add brew session' }}
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              density="comfortable"
              label="Brewed at"
              :model-value="form.brewed_at"
              :rules="[rules.required]"
              type="datetime-local"
              @update:model-value="updateForm('brewed_at', $event)"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-autocomplete
              clearable
              density="comfortable"
              item-title="name"
              item-value="uuid"
              :items="mashVesselOptions"
              label="Mash Vessel"
              :model-value="form.mash_vessel_uuid"
              @update:model-value="updateForm('mash_vessel_uuid', $event)"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-autocomplete
              clearable
              density="comfortable"
              item-title="name"
              item-value="uuid"
              :items="boilVesselOptions"
              label="Boil Vessel"
              :model-value="form.boil_vessel_uuid"
              @update:model-value="updateForm('boil_vessel_uuid', $event)"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-autocomplete
              clearable
              density="comfortable"
              hint="Select existing or create new"
              item-title="label"
              item-value="uuid"
              :items="wortVolumeOptions"
              label="Wort Volume"
              :model-value="form.wort_volume_uuid"
              persistent-hint
              @update:model-value="updateForm('wort_volume_uuid', $event)"
            >
              <template #no-data>
                <v-list-item>
                  <v-list-item-title>No volumes available</v-list-item-title>
                </v-list-item>
              </template>
              <template #append-item>
                <v-divider class="my-2" />
                <v-list-item @click="emit('create-volume')">
                  <template #prepend>
                    <v-icon icon="mdi-plus" />
                  </template>
                  <v-list-item-title>Create new wort volume</v-list-item-title>
                </v-list-item>
              </template>
            </v-autocomplete>
          </v-col>
          <v-col cols="12">
            <v-textarea
              auto-grow
              density="comfortable"
              label="Notes"
              :model-value="form.notes"
              placeholder="Mash temps, boil notes, etc."
              rows="2"
              @update:model-value="updateForm('notes', $event)"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="emit('update:modelValue', false)">
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          :disabled="!isValid"
          :loading="saving"
          @click="emit('submit')"
        >
          {{ isEditing ? 'Save changes' : 'Add session' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Volume as ProductionVolume, Vessel } from '@/types'
  import { computed } from 'vue'

  export type BrewSessionForm = {
    brewed_at: string
    mash_vessel_uuid: string | null
    boil_vessel_uuid: string | null
    wort_volume_uuid: string | null
    notes: string
  }

  const props = defineProps<{
    modelValue: boolean
    form: BrewSessionForm
    isEditing: boolean
    saving: boolean
    vessels: Vessel[]
    volumes: ProductionVolume[]
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'update:form': [form: BrewSessionForm]
    'submit': []
    'create-volume': []
  }>()

  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
  }

  const isValid = computed(() => {
    return props.form.brewed_at.trim().length > 0
  })

  const mashVesselOptions = computed(() =>
    props.vessels
      .filter(v => v.status === 'active' && v.type.toLowerCase().includes('mash'))
      .map(v => ({ uuid: v.uuid, name: v.name })),
  )

  const boilVesselOptions = computed(() =>
    props.vessels
      .filter(v => v.status === 'active' && (v.type.toLowerCase().includes('kettle') || v.type.toLowerCase().includes('boil')))
      .map(v => ({ uuid: v.uuid, name: v.name })),
  )

  const wortVolumeOptions = computed(() =>
    props.volumes.map(v => ({
      uuid: v.uuid,
      label: v.name ? `${v.name} (${v.amount} ${v.amount_unit})` : `Unnamed Volume (${v.amount} ${v.amount_unit})`,
    })),
  )

  function updateForm<K extends keyof BrewSessionForm> (key: K, value: BrewSessionForm[K]) {
    emit('update:form', { ...props.form, [key]: value })
  }
</script>
