<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { statusClass, statusLabel } from '@/lib/applicationStatus'
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
const actionError = ref('')

const filteredApplications = computed(() => {
  if (selectedStatus.value === 'all') return applications.value
  return applications.value.filter((application) => application.status === selectedStatus.value)
})

function openActionDialog(application: OfficerApplication) {
  selectedApplication.value = application
  selectedDecision.value = 'approved'
  reason.value = ''
  actionError.value = ''
  isDialogOpen.value = true
}

async function updateStatus() {
  if (!selectedApplication.value) return
  actionError.value = ''

  if (selectedDecision.value === 'rejected' && !reason.value.trim()) {
    actionError.value = 'กรุณาระบุเหตุผลเมื่อไม่อนุมัติ'
    return
  }

  isUpdating.value = true
  try {
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
    isDialogOpen.value = false
  } catch {
    actionError.value = 'ไม่สามารถบันทึกสถานะได้ กรุณาลองใหม่อีกครั้ง'
  } finally {
    isUpdating.value = false
  }
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
            <TableCell><Badge :class="statusClass(application.status)">{{ statusLabel(application.status) }}</Badge></TableCell>
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
          <dl v-if="selectedApplication" class="space-y-2 rounded-md bg-slate-50 p-3 text-sm">
            <div class="flex justify-between gap-4"><dt class="text-slate-500">รายได้ต่อปี</dt><dd>{{ selectedApplication.annualIncome.toLocaleString('th-TH') }} บาท</dd></div>
            <div class="flex justify-between gap-4"><dt class="text-slate-500">สถานะเดิม</dt><dd><Badge :class="statusClass(selectedApplication.status)">{{ statusLabel(selectedApplication.status) }}</Badge></dd></div>
          </dl>
          <p class="text-sm font-medium">เลือกการตัดสินใจ</p>
          <div class="grid grid-cols-2 gap-3">
            <Button
              type="button"
              variant="outline"
              :aria-pressed="selectedDecision === 'approved'"
              :class="selectedDecision === 'approved'
                ? 'border-green-600 bg-green-600 text-white hover:bg-green-700 hover:text-white'
                : 'border-green-200 text-green-700 hover:bg-green-50'"
              @click="selectedDecision = 'approved'"
            >
              Approve
            </Button>
            <Button
              type="button"
              variant="outline"
              :aria-pressed="selectedDecision === 'rejected'"
              :class="selectedDecision === 'rejected'
                ? 'border-red-600 bg-red-600 text-white hover:bg-red-700 hover:text-white'
                : 'border-red-200 text-red-700 hover:bg-red-50'"
              @click="selectedDecision = 'rejected'"
            >
              Reject
            </Button>
          </div>
          <p class="text-xs text-slate-500">เลือกแล้ว: {{ selectedDecision === 'approved' ? 'Approve' : 'Reject' }}</p>
          <Input v-model="reason" :placeholder="selectedDecision === 'rejected' ? 'เหตุผล (จำเป็น)' : 'เหตุผล (ถ้ามี)'" />
          <p v-if="actionError" class="rounded-md bg-red-50 p-3 text-sm text-red-700">{{ actionError }}</p>
        </div>

        <DialogFooter>
          <Button :disabled="isUpdating" @click="updateStatus">{{ isUpdating ? 'กำลังบันทึก...' : 'บันทึกสถานะ' }}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </main>
</template>