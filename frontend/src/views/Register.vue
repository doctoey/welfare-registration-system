<script setup lang="ts">
import { ref } from 'vue'

const citizenId = ref('')
const fullName = ref('')
const dateOfBirth = ref('')
const annualIncome = ref('')
const currentAddress = ref('')
const consent = ref(false)
const validationErrors = ref<string[]>([])

function handleSubmit() {
  validationErrors.value = []

  if (citizenId.value.trim().length !== 13) {
    validationErrors.value.push('National ID must contain 13 digits')
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
    citizenId: citizenId.value,
    fullName: fullName.value.trim(),
    birthDate: dateOfBirth.value,
    annualIncome: Number(annualIncome.value),
    currentAddress: currentAddress.value.trim(),
  }

  console.log('Application payload:', payload)
}
</script>

<template>
  <main>
    <h1>Welfare Registration System</h1>

    <form @submit.prevent="handleSubmit">
      <label>
        National ID
        <input v-model="citizenId" type="text" />
      </label>

      <label>
        Full name
        <input v-model="fullName" type="text" />
      </label>

      <label>
        Date of birth
        <input v-model="dateOfBirth" type="date" />
      </label>

      <label>
        Annual household income
        <input v-model="annualIncome" type="number" min="0" />
      </label>

      <label>
        Current address
        <textarea v-model="currentAddress" rows="3"></textarea>
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

      <button type="submit">
        Submit
      </button>
    </form>
  </main>
</template>