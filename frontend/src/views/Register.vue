<script setup lang="ts">
import { ref } from 'vue'

const citizenId = ref('')
const fullName = ref('')
const dateOfBirth = ref('')
const annualIncome = ref('')
const currentAddress = ref('')
const consent = ref(false)
</script>

<template>
  <main>
    <h1>Welfare Registration System</h1>

    <form @submit.prevent>
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

      <button type="submit">
        Submit
      </button>
    </form>
  </main>
</template>