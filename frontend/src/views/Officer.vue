<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { mockGetOfficerApplications, type OfficerApplication } from '../services/applicationApi'

const applications = ref<OfficerApplication[]>([])
const selectedStatus = ref('all')
const isLoading = ref(true)

const filteredApplications = computed(() => {
  if (selectedStatus.value === 'all') return applications.value
  return applications.value.filter((application) => application.status === selectedStatus.value)
})

onMounted(async () => {
  const response = await mockGetOfficerApplications()
  applications.value = response.data
  isLoading.value = false
})
</script>

<template>
  <main class="space-y-5">
    <Card class="mx-auto max-w-5xl p-6">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl font-semibold">Officer Applications</h1>
          <p class="text-sm text-slate-600">รายการคำร้องสวัสดิการ</p>
        </div>
        <Select v-model="selectedStatus">
          <SelectTrigger class="w-44"><SelectValue placeholder="กรองสถานะ" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="all">ทุกสถานะ</SelectItem>
            <SelectItem value="pending">รอตรวจสอบ</SelectItem>
            <SelectItem value="approved">อนุมัติแล้ว</SelectItem>
            <SelectItem value="rejected">ไม่อนุมัติ</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <p v-if="isLoading" class="py-8 text-center text-sm text-slate-500">กำลังโหลด...</p>
      <Table v-else class="mt-6">
        <TableHeader>
          <TableRow><TableHead>ชื่อ</TableHead><TableHead>เลขบัตรประชาชน</TableHead><TableHead>รายได้ต่อปี</TableHead><TableHead>สถานะ</TableHead><TableHead class="text-right">Action</TableHead></TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="application in filteredApplications" :key="application.id">
            <TableCell>{{ application.fullName }}</TableCell>
            <TableCell>{{ application.citizenId }}</TableCell>
            <TableCell>{{ application.annualIncome.toLocaleString('th-TH') }}</TableCell>
            <TableCell class="capitalize">{{ application.status }}</TableCell>
            <TableCell class="text-right"><Button size="sm" variant="outline">ดูรายละเอียด</Button></TableCell>
          </TableRow>
          <TableRow v-if="filteredApplications.length === 0"><TableCell colspan="5" class="py-8 text-center">ไม่พบรายการ</TableCell></TableRow>
        </TableBody>
      </Table>
    </Card>
  </main>
</template>