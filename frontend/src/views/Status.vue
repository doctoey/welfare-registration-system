<script setup lang="ts">
import { ref } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { statusClass, statusLabel } from '@/lib/applicationStatus'
import { isValidThaiCitizenId } from '@/lib/thaiCitizenId'
import { MOCK_CITIZEN_ID, mockGetApplicationStatus, type ApplicationStatus } from '../services/applicationApi'

const citizenId = ref('')
const status = ref<ApplicationStatus | null>(null)
const errorMessage = ref('')
const isSearching = ref(false)

async function searchStatus() {
  errorMessage.value = ''
  status.value = null

  const citizenIdDigits = citizenId.value.replace(/\D/g, '')
  if (!isValidThaiCitizenId(citizenIdDigits)) {
    errorMessage.value = 'National ID is invalid'
    return
  }

  isSearching.value = true
  try {
    const response = await mockGetApplicationStatus(citizenIdDigits)
    status.value = response.data
  } catch (error) {
    errorMessage.value = error instanceof Error && error.message === 'APPLICATION_NOT_FOUND'
      ? 'ไม่พบข้อมูลคำร้องของเลขบัตรประชาชนนี้'
      : 'ไม่สามารถค้นหาสถานะได้ กรุณาลองใหม่อีกครั้ง'
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
          <p class="text-xs text-slate-500">Mock test ID: {{ MOCK_CITIZEN_ID }}</p>
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
        <div class="flex justify-between gap-4"><dt class="text-slate-500">สถานะ</dt><dd><Badge :class="statusClass(status.status)">{{ statusLabel(status.status) }}</Badge></dd></div>
      </dl>
    </Card>
  </main>
</template>