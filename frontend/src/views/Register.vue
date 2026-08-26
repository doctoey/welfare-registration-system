<script setup lang="ts">
</script>

<template>
  <main>
    <h1>Welfare Registration System</h1>

    <form @submit.prevent>
      <label>
        National ID
        <input type="text" />
      </label>

      <button type="submit">
        Submit
      </button>
    </form>
  </main>
</template>