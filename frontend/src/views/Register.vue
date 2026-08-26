<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
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
    <section v-if="submitted">
      <h1>ส่งคำร้องสำเร็จ</h1>
      <p>หมายเลขอ้างอิง: {{ referenceNumber }}</p>
    </section>

    <form v-else @submit.prevent="handleSubmit">
      <label>
        National ID
        <Input v-model="citizenId" type="text" />
      </label>

      <label>
        Full name
        <Input v-model="fullName" type="text" />
      </label>

      <label>
        Date of birth
        <Input v-model="dateOfBirth" type="date" />
      </label>

      <label>
        Annual household income
        <Input v-model="annualIncome" type="number" min="0" />
      </label>

      <label>
        Current address
        <Textarea v-model="currentAddress" :rows="3" />
      </label>

      <label>
        <input v-model="consent" type="checkbox" />
        I agree that the information provided is correct.
      </label>

      <p>เลขบัตรที่กรอก: {{ citizenId }}</p>
      <p>ชื่อ: {{ fullName }}</p>
      <p>วันเกิด: {{ dateOfBirth }}</p>
      <p>รายได้ต่อปี: {{ annualIncome }}</p>
      <p>ที่อยู่: {{ currentAddress }}</p>
      <p>ยินยอม: {{ consent }}</p>

      <ul v-if="validationErrors.length > 0">
        <li v-for="error in validationErrors" :key="error">{{ error }}</li>
      </ul>

      <p v-if="apiError">{{ apiError }}</p>

      <Button type="submit" :disabled="isSubmitting">
        {{ isSubmitting ? 'กำลังส่ง...' : 'Submit' }}
      </Button>
    </form>
  </main>
</template>