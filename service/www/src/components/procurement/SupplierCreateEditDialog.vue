<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 720"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit supplier' : 'Create supplier' }}
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.name" label="Name" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.contact_name" label="Contact name" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.email" label="Email" type="email" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.phone" label="Phone" type="tel" />
          </v-col>
          <v-col cols="12">
            <v-text-field v-model="form.address_line1" label="Address line 1" />
          </v-col>
          <v-col cols="12">
            <v-text-field v-model="form.address_line2" label="Address line 2" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.city" label="City" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.region" label="Region" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.postal_code" label="Postal code" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.country" label="Country" />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!form.name.trim()"
          :loading="saving"
          @click="handleSubmit"
        >
          {{ isEditing ? 'Save changes' : 'Add supplier' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Supplier } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import { normalizeText } from '@/utils/normalize'

  export interface SupplierFormData {
    name: string
    contact_name: string | null
    email: string | null
    phone: string | null
    address_line1: string | null
    address_line2: string | null
    city: string | null
    region: string | null
    postal_code: string | null
    country: string | null
  }

  const props = defineProps<{
    modelValue: boolean
    editSupplier?: Supplier | null
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: SupplierFormData]
  }>()

  const form = reactive({
    name: '',
    contact_name: '',
    email: '',
    phone: '',
    address_line1: '',
    address_line2: '',
    city: '',
    region: '',
    postal_code: '',
    country: '',
  })

  const isEditing = computed(() => !!props.editSupplier)

  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        resetForm()
      }
    },
  )

  function resetForm () {
    if (props.editSupplier) {
      form.name = props.editSupplier.name
      form.contact_name = props.editSupplier.contact_name ?? ''
      form.email = props.editSupplier.email ?? ''
      form.phone = props.editSupplier.phone ?? ''
      form.address_line1 = props.editSupplier.address_line1 ?? ''
      form.address_line2 = props.editSupplier.address_line2 ?? ''
      form.city = props.editSupplier.city ?? ''
      form.region = props.editSupplier.region ?? ''
      form.postal_code = props.editSupplier.postal_code ?? ''
      form.country = props.editSupplier.country ?? ''
    } else {
      form.name = ''
      form.contact_name = ''
      form.email = ''
      form.phone = ''
      form.address_line1 = ''
      form.address_line2 = ''
      form.city = ''
      form.region = ''
      form.postal_code = ''
      form.country = ''
    }
  }

  function handleSubmit () {
    if (!form.name.trim()) return

    const payload: SupplierFormData = {
      name: form.name.trim(),
      contact_name: normalizeText(form.contact_name),
      email: normalizeText(form.email),
      phone: normalizeText(form.phone),
      address_line1: normalizeText(form.address_line1),
      address_line2: normalizeText(form.address_line2),
      city: normalizeText(form.city),
      region: normalizeText(form.region),
      postal_code: normalizeText(form.postal_code),
      country: normalizeText(form.country),
    }
    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
