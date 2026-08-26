<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
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

function isValidThaiCitizenId(value: string): boolean {
  const digits = value.replace(/\D/g, '')
  if (digits.length !== 13) return false

  let sum = 0
  for (let index = 0; index < 12; index += 1) {
    sum += Number(digits[index]) * (13 - index)
  }

  const checkDigit = (11 - (sum % 11)) % 10
  return checkDigit === Number(digits[12])
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

      <label class="block space-y-2 text-sm font-medium">
        National ID
        <Input v-model="citizenId" type="text" />
      </label>

      <label class="block space-y-2 text-sm font-medium">
        Full name
        <Input v-model="fullName" type="text" />
      </label>

      <label class="block space-y-2 text-sm font-medium">
        Date of birth
        <Input v-model="dateOfBirth" type="date" />
      </label>

      <label class="block space-y-2 text-sm font-medium">
        Annual household income
        <Input v-model="annualIncome" type="number" min="0" />
      </label>

      <label class="block space-y-2 text-sm font-medium">
        Current address
        <Textarea v-model="currentAddress" :rows="3" />
      </label>

      <label class="flex items-center gap-2 text-sm">
        <Checkbox v-model:checked="consent" />
        I agree that the information provided is correct.
      </label>

      <ul v-if="validationErrors.length > 0" class="space-y-1 rounded-md bg-red-50 p-3 text-sm text-red-700">
        <li v-for="error in validationErrors" :key="error">{{ error }}</li>
      </ul>

      <p v-if="apiError" class="rounded-md bg-red-50 p-3 text-sm text-red-700">{{ apiError }}</p>

        <Button type="submit" :disabled="isSubmitting">
          {{ isSubmitting ? 'กำลังส่ง...' : 'Submit' }}
        </Button>
      </form>
    </Card>
  </main>
</template>