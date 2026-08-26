<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { mockGetOfficerApplications, mockUpdateApplicationStatus, type OfficerApplication } from '../services/applicationApi'

const applications = ref<OfficerApplication[]>([])
const selectedStatus = ref('all')
const isLoading = ref(true)
const isDialogOpen = ref(false)
const selectedApplication = ref<OfficerApplication | null>(null)
const selectedDecision = ref<'approved' | 'rejected'>('approved')
const reason = ref('')
const isUpdating = ref(false)

const filteredApplications = computed(() => {
  if (selectedStatus.value === 'all') return applications.value
  return applications.value.filter((application) => application.status === selectedStatus.value)
})

function openActionDialog(application: OfficerApplication) {
  selectedApplication.value = application
  selectedDecision.value = 'approved'
  reason.value = ''
  isDialogOpen.value = true
}

async function updateStatus() {
  if (!selectedApplication.value) return

  isUpdating.value = true
  const response = await mockUpdateApplicationStatus(
    selectedApplication.value.id,
    selectedDecision.value,
    reason.value.trim(),
  )
  const application = applications.value.find((item) => item.id === selectedApplication.value?.id)
  if (application) {
    application.status = response.data.status
    application.reason = response.data.reason
  }
  isUpdating.value = false
  isDialogOpen.value = false
}

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
            <TableCell class="text-right"><Button size="sm" variant="outline" @click="openActionDialog(application)">ดำเนินการ</Button></TableCell>
          </TableRow>
          <TableRow v-if="filteredApplications.length === 0"><TableCell colspan="5" class="py-8 text-center">ไม่พบรายการ</TableCell></TableRow>
        </TableBody>
      </Table>
    </Card>

    <Dialog v-model:open="isDialogOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>ดำเนินการคำร้อง</DialogTitle>
          <DialogDescription v-if="selectedApplication">{{ selectedApplication.fullName }} ({{ selectedApplication.citizenId }})</DialogDescription>
        </DialogHeader>

        <div class="space-y-3 py-4">
          <div class="flex gap-2">
            <Button :variant="selectedDecision === 'approved' ? 'default' : 'outline'" @click="selectedDecision = 'approved'">Approve</Button>
            <Button :variant="selectedDecision === 'rejected' ? 'default' : 'outline'" @click="selectedDecision = 'rejected'">Reject</Button>
          </div>
          <Input v-model="reason" placeholder="เหตุผล (ถ้ามี)" />
        </div>

        <DialogFooter>
          <Button :disabled="isUpdating" @click="updateStatus">{{ isUpdating ? 'กำลังบันทึก...' : 'บันทึกสถานะ' }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </main>
</template>