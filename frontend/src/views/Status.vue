<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { mockGetApplicationStatus, type ApplicationStatus } from '../services/applicationApi'

const citizenId = ref('')
const status = ref<ApplicationStatus | null>(null)
const errorMessage = ref('')
const isSearching = ref(false)

async function searchStatus() {
  errorMessage.value = ''
  status.value = null

  if (citizenId.value.trim().length !== 13) {
    errorMessage.value = 'National ID must contain 13 digits'
    return
  }

  isSearching.value = true
  try {
    const response = await mockGetApplicationStatus(citizenId.value.trim())
    status.value = response.data
  } catch {
    errorMessage.value = 'ไม่สามารถค้นหาสถานะได้ กรุณาลองใหม่อีกครั้ง'
  } finally {
    isSearching.value = false
  }
}
</script>

<template>
  <main class="space-y-5">
    <Card class="mx-auto max-w-2xl p-6">
      <h1 class="text-2xl font-semibold">ติดตามสถานะคำร้อง</h1>
      <p class="mt-1 text-sm text-slate-600">กรอกเลขบัตรประชาชน 13 หลัก</p>

      <form class="mt-6 space-y-4" @submit.prevent="searchStatus">
        <div class="space-y-2">
          <Label for="status-citizen-id">National ID</Label>
          <Input id="status-citizen-id" v-model="citizenId" inputmode="numeric" maxlength="13" />
        </div>

        <p v-if="errorMessage" class="rounded-md bg-red-50 p-3 text-sm text-red-700">{{ errorMessage }}</p>

        <Button type="submit" :disabled="isSearching">
          {{ isSearching ? 'กำลังค้นหา...' : 'ค้นหาสถานะ' }}
        </Button>
      </form>
    </Card>

    <Card v-if="status" class="mx-auto max-w-2xl p-6">
      <h2 class="text-xl font-semibold">ผลการค้นหา</h2>
      <dl class="mt-4 space-y-2 text-sm">
        <div class="flex justify-between gap-4"><dt class="text-slate-500">ชื่อ</dt><dd>{{ status.fullName }}</dd></div>
        <div class="flex justify-between gap-4"><dt class="text-slate-500">เลขอ้างอิง</dt><dd>{{ status.referenceNumber }}</dd></div>
        <div class="flex justify-between gap-4"><dt class="text-slate-500">สถานะ</dt><dd class="font-semibold capitalize">{{ status.status }}</dd></div>
      </dl>
    </Card>
  </main>
</template>