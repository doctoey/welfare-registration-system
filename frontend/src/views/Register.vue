<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { isValidThaiCitizenId } from '@/lib/thaiCitizenId'
import { Textarea } from '@/components/ui/textarea'
import { mockCreateApplication } from '../services/applicationApi'

const citizenId = ref('')
const fullName = ref('')
const dateOfBirth = ref('')
const annualIncome = ref('')
const currentAddress = ref('')
const consent = ref(false)
const validationErrors = ref<string[]>([])
const isSubmitting = ref(false)
const apiError = ref('')
const submitted = ref(false)
const referenceNumber = ref('')

function hasValidationError(message: string) {
  return validationErrors.value.includes(message)
}

async function handleSubmit() {
  validationErrors.value = []
  apiError.value = ''

  const citizenIdDigits = citizenId.value.replace(/\D/g, '')
  if (!isValidThaiCitizenId(citizenIdDigits)) {
    validationErrors.value.push('National ID is invalid')
  }
  if (!fullName.value.trim()) validationErrors.value.push('Full name is required')
  if (!dateOfBirth.value) validationErrors.value.push('Date of birth is required')
  if (!annualIncome.value || Number(annualIncome.value) < 0) {
    validationErrors.value.push('Annual income is required')
  }
  if (!currentAddress.value.trim()) validationErrors.value.push('Current address is required')
  if (!consent.value) validationErrors.value.push('Consent is required')

  if (validationErrors.value.length > 0) return

  const payload = {
    citizenId: citizenIdDigits,
    fullName: fullName.value.trim(),
    birthDate: dateOfBirth.value,
    annualIncome: Number(annualIncome.value),
    currentAddress: currentAddress.value.trim(),
  }

  isSubmitting.value = true

  try {
    const response = await mockCreateApplication(payload)
    referenceNumber.value = response.data.referenceNumber
    submitted.value = true
  } catch (error) {
    console.error(error)
    apiError.value = 'ไม่สามารถส่งคำร้องได้ กรุณาลองใหม่อีกครั้ง'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main>
    <Card v-if="submitted" class="mx-auto max-w-2xl p-8 text-center">
      <h1 class="text-2xl font-semibold">ส่งคำร้องสำเร็จ</h1>
      <p class="mt-3 text-slate-600">หมายเลขอ้างอิง: {{ referenceNumber }}</p>
    </Card>

    <Card v-else class="mx-auto max-w-2xl">
      <form class="space-y-5 p-6" @submit.prevent="handleSubmit">
        <h1 class="text-2xl font-semibold">Welfare Registration System</h1>

        <div class="space-y-2">
          <Label for="citizen-id">National ID <span class="text-red-600">*</span></Label>
          <Input id="citizen-id" v-model="citizenId" type="text" :aria-invalid="hasValidationError('National ID is invalid')" />
          <p v-if="hasValidationError('National ID is invalid')" class="text-sm text-red-600">National ID must be valid 13-digit Thai ID</p>
        </div>

        <div class="space-y-2">
          <Label for="full-name">Full name <span class="text-red-600">*</span></Label>
          <Input id="full-name" v-model="fullName" :aria-invalid="hasValidationError('Full name is required')" />
          <p v-if="hasValidationError('Full name is required')" class="text-sm text-red-600">Full name is required</p>
        </div>

        <div class="space-y-2">
          <Label for="date-of-birth">Date of birth <span class="text-red-600">*</span></Label>
          <Input id="date-of-birth" v-model="dateOfBirth" type="date" :aria-invalid="hasValidationError('Date of birth is required')" />
          <p v-if="hasValidationError('Date of birth is required')" class="text-sm text-red-600">Date of birth is required</p>
        </div>

        <div class="space-y-2">
          <Label for="annual-income">Annual household income <span class="text-red-600">*</span></Label>
          <Input id="annual-income" v-model="annualIncome" type="number" min="0" :aria-invalid="hasValidationError('Annual income is required')" />
          <p v-if="hasValidationError('Annual income is required')" class="text-sm text-red-600">Annual income is required</p>
        </div>

        <div class="space-y-2">
          <Label for="current-address">Current address <span class="text-red-600">*</span></Label>
          <Textarea id="current-address" v-model="currentAddress" :rows="3" :aria-invalid="hasValidationError('Current address is required')" />
          <p v-if="hasValidationError('Current address is required')" class="text-sm text-red-600">Current address is required</p>
        </div>

        <div class="flex items-center gap-2">
          <Checkbox id="consent" v-model:checked="consent" />
          <Label for="consent">I agree that the information provided is correct. <span class="text-red-600">*</span></Label>
        </div>

      <p v-if="apiError" class="rounded-md bg-red-50 p-3 text-sm text-red-700">{{ apiError }}</p>

        <Button type="submit" :disabled="isSubmitting">
          {{ isSubmitting ? 'กำลังส่ง...' : 'Submit' }}
        </Button>
      </form>
    </Card>
  </main>
</template>