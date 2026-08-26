<script setup lang="ts">
import { ref } from 'vue'

const citizenId = ref('')
</script>

<template>
  <main>
    <h1>Welfare Registration System</h1>

    <form @submit.prevent>
      <label>
        National ID
        <input v-model="citizenId" type="text" />
      </label>

      <p>เลขบัตรที่กรอก: {{ citizenId }}</p>

      <button type="submit">
        Submit
      </button>
    </form>
  </main>
</template>