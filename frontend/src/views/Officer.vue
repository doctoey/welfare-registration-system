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
import {
  ShieldCheck,
  Users,
  Clock,
  CheckCircle2,
  XCircle,
  Search,
  AlertCircle,
  Eye,
  RefreshCw,
} from 'lucide-vue-next'
import {
  getOfficerApplications,
  updateApplicationStatus,
  type OfficerApplication,
} from '@/services/officerApi'

const applications = ref<OfficerApplication[]>([])
const selectedStatus = ref('all')
const searchQuery = ref('')
const isLoading = ref(true)

const isActionDialogOpen = ref(false)
const selectedApplication = ref<OfficerApplication | null>(null)
const selectedDecision = ref<'approved' | 'rejected'>('approved')
const reason = ref('')
const isUpdating = ref(false)
const actionError = ref('')

const isDetailModalOpen = ref(false)
const detailApplication = ref<OfficerApplication | null>(null)

function formatThaiCitizenId(raw: string): string {
  const digits = raw.replace(/\D/g, '').slice(0, 13)
  const parts = [
    digits.slice(0, 1),
    digits.slice(1, 5),
    digits.slice(5, 10),
    digits.slice(10, 12),
    digits.slice(12, 13),
  ]
  return parts.filter(Boolean).join('-')
}

const stats = computed(() => {
  const total = applications.value.length
  const pending = applications.value.filter((a) => a.status === 'pending').length
  const approved = applications.value.filter((a) => a.status === 'approved').length
  const rejected = applications.value.filter((a) => a.status === 'rejected').length
  return { total, pending, approved, rejected }
})

const filteredApplications = computed(() => {
  return applications.value.filter((application) => {
    const matchesStatus = selectedStatus.value === 'all' || application.status === selectedStatus.value
    const query = searchQuery.value.trim().toLowerCase()
    const matchesQuery =
      !query ||
      application.fullName.toLowerCase().includes(query) ||
      application.citizenId.includes(query) ||
      application.referenceNumber.toLowerCase().includes(query)
    return matchesStatus && matchesQuery
  })
})

function openDetailModal(application: OfficerApplication) {
  detailApplication.value = application
  isDetailModalOpen.value = true
}

function openActionDialog(application: OfficerApplication) {
  selectedApplication.value = application
  selectedDecision.value = 'approved'
  reason.value = application.reason || ''
  actionError.value = ''
  isActionDialogOpen.value = true
}

async function loadData() {
  isLoading.value = true
  try {
    const response = await getOfficerApplications()
    applications.value = response.data
  } catch (error) {
    console.error('Failed to load officer applications:', error)
  } finally {
    isLoading.value = false
  }
}

async function updateStatus() {
  if (!selectedApplication.value) return
  actionError.value = ''

  if (selectedDecision.value === 'rejected' && !reason.value.trim()) {
    actionError.value = 'กรุณาระบุเหตุผลเมื่อเลือกไม่อนุมัติ'
    return
  }

  isUpdating.value = true
  try {
    const response = await updateApplicationStatus(
      selectedApplication.value.id,
      selectedDecision.value,
      reason.value.trim(),
    )
    const application = applications.value.find((item) => item.id === selectedApplication.value?.id)
    if (application) {
      application.status = response.data.status
      application.reason = response.data.reason
    }
    isActionDialogOpen.value = false
  } catch (err: any) {
    const errorMsg = err.response?.data || err.message || ''
    if (typeof errorMsg === 'string' && errorMsg.includes('REASON_REQUIRED_FOR_REJECT')) {
      actionError.value = 'กรุณาระบุเหตุผลเมื่อเลือกไม่อนุมัติ'
    } else {
      actionError.value = 'ไม่สามารถบันทึกสถานะได้ กรุณาลองใหม่อีกครั้ง'
    }
  } finally {
    isUpdating.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <main class="py-6 px-4 sm:px-6 space-y-6">
    <div class="mx-auto max-w-5xl grid grid-cols-2 sm:grid-cols-4 gap-4">
      <Card class="p-4 bg-white shadow-sm border-slate-200">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-50 text-blue-600">
            <Users class="h-5 w-5" />
          </div>
          <div>
            <p class="text-xs text-slate-500 font-medium">คำร้องทั้งหมด</p>
            <p class="text-xl font-bold text-slate-800">{{ stats.total }}</p>
          </div>
        </div>
      </Card>

      <Card class="p-4 bg-white shadow-sm border-slate-200">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-50 text-amber-600">
            <Clock class="h-5 w-5" />
          </div>
          <div>
            <p class="text-xs text-slate-500 font-medium">รอตรวจสอบ</p>
            <p class="text-xl font-bold text-amber-600">{{ stats.pending }}</p>
          </div>
        </div>
      </Card>

      <Card class="p-4 bg-white shadow-sm border-slate-200">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600">
            <CheckCircle2 class="h-5 w-5" />
          </div>
          <div>
            <p class="text-xs text-slate-500 font-medium">อนุมัติแล้ว</p>
            <p class="text-xl font-bold text-emerald-600">{{ stats.approved }}</p>
          </div>
        </div>
      </Card>

      <Card class="p-4 bg-white shadow-sm border-slate-200">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-rose-50 text-rose-600">
            <XCircle class="h-5 w-5" />
          </div>
          <div>
            <p class="text-xs text-slate-500 font-medium">ไม่อนุมัติ</p>
            <p class="text-xl font-bold text-rose-600">{{ stats.rejected }}</p>
          </div>
        </div>
      </Card>
    </div>

    <Card class="mx-auto max-w-5xl overflow-hidden bg-white shadow-sm border-slate-200">
      <div class="border-b border-slate-100 bg-slate-50/70 p-6 flex flex-wrap items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-600 text-white shadow-sm">
            <ShieldCheck class="h-6 w-6" />
          </div>
          <div>
            <h1 class="text-xl font-bold text-slate-800">ระบบพิจารณาคำร้องสำหรับเจ้าหน้าที่</h1>
            <p class="text-xs sm:text-sm text-slate-500">Officer Welfare Review Portal</p>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <Button variant="outline" size="sm" class="h-9 px-2.5 text-xs text-slate-600" @click="loadData">
            <RefreshCw class="h-3.5 w-3.5 mr-1" />
            รีเฟรช
          </Button>

          <div class="relative w-44 sm:w-56">
            <Input
              v-model="searchQuery"
              placeholder="ค้นหาชื่อ หรือเลขบัตร..."
              class="h-9 text-xs pl-8"
            />
            <Search class="absolute left-2.5 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-slate-400" />
          </div>

          <Select v-model="selectedStatus">
            <SelectTrigger class="w-36 h-9 text-xs"><SelectValue placeholder="กรองสถานะ" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">ทุกสถานะ</SelectItem>
              <SelectItem value="pending">รอตรวจสอบ</SelectItem>
              <SelectItem value="approved">อนุมัติแล้ว</SelectItem>
              <SelectItem value="rejected">ไม่อนุมัติ</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div class="p-6">
        <div v-if="isLoading" class="py-12 text-center text-sm text-slate-500 flex flex-col items-center gap-2">
          <svg class="animate-spin h-6 w-6 text-blue-600" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          กำลังโหลดข้อมูลคำร้อง...
        </div>

        <div v-else class="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow class="bg-slate-50/50">
                <TableHead class="font-semibold text-slate-700">เลขอ้างอิง</TableHead>
                <TableHead class="font-semibold text-slate-700">ชื่อ - นามสกุล</TableHead>
                <TableHead class="font-semibold text-slate-700">เลขประจำตัวประชาชน</TableHead>
                <TableHead class="font-semibold text-slate-700">รายได้/ปี</TableHead>
                <TableHead class="font-semibold text-slate-700">สถานะ</TableHead>
                <TableHead class="text-right font-semibold text-slate-700">การดำเนินการ</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="application in filteredApplications" :key="application.id" class="hover:bg-slate-50/60">
                <TableCell class="font-mono text-xs font-semibold text-blue-700">{{ application.referenceNumber }}</TableCell>
                <TableCell class="font-medium text-slate-800">{{ application.fullName }}</TableCell>
                <TableCell class="font-mono text-slate-600 text-xs">{{ formatThaiCitizenId(application.citizenId) }}</TableCell>
                <TableCell class="font-mono text-slate-600 text-xs">{{ application.annualIncome.toLocaleString('th-TH') }} ฿</TableCell>
                <TableCell><Badge :class="statusClass(application.status)">{{ statusLabel(application.status) }}</Badge></TableCell>
                <TableCell class="text-right">
                  <div class="flex items-center justify-end gap-1.5">
                    <Button
                      size="sm"
                      variant="ghost"
                      class="h-8 px-2 text-xs text-slate-600 hover:text-blue-700 hover:bg-blue-50"
                      title="ดูรายละเอียด"
                      @click="openDetailModal(application)"
                    >
                      <Eye class="h-3.5 w-3.5" />
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      class="h-8 text-xs hover:bg-blue-50 hover:text-blue-700 hover:border-blue-200"
                      @click="openActionDialog(application)"
                    >
                      พิจารณา
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
              <TableRow v-if="filteredApplications.length === 0">
                <TableCell colspan="6" class="py-10 text-center text-slate-400">
                  ไม่พบรายการคำร้องที่ตรงกับเงื่อนไขการค้นหา
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </div>
    </Card>

    <Dialog v-model:open="isDetailModalOpen">
      <DialogContent class="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle class="text-lg font-bold text-slate-800">รายละเอียดคำร้องสวัสดิการ</DialogTitle>
          <DialogDescription v-if="detailApplication" class="font-mono text-xs text-blue-700 font-semibold">
            {{ detailApplication.referenceNumber }}
          </DialogDescription>
        </DialogHeader>

        <div v-if="detailApplication" class="space-y-4 py-2 text-sm">
          <div class="rounded-lg bg-slate-50 p-4 border border-slate-100 divide-y divide-slate-200/60">
            <div class="pb-2.5 flex justify-between">
              <span class="text-slate-500 font-medium">ชื่อ-นามสกุล</span>
              <span class="font-semibold text-slate-800">{{ detailApplication.fullName }}</span>
            </div>
            <div class="py-2.5 flex justify-between">
              <span class="text-slate-500 font-medium">เลขประจำตัวประชาชน</span>
              <span class="font-mono text-slate-800">{{ formatThaiCitizenId(detailApplication.citizenId) }}</span>
            </div>
            <div class="py-2.5 flex justify-between">
              <span class="text-slate-500 font-medium">วันเกิด</span>
              <span class="text-slate-800">{{ detailApplication.birthDate || '-' }}</span>
            </div>
            <div class="py-2.5 flex justify-between">
              <span class="text-slate-500 font-medium">รายได้รวมครัวเรือนต่อปี</span>
              <span class="font-mono font-semibold text-slate-800">{{ detailApplication.annualIncome.toLocaleString('th-TH') }} บาท</span>
            </div>
            <div class="py-2.5 flex justify-between items-center">
              <span class="text-slate-500 font-medium">สถานะ</span>
              <Badge :class="statusClass(detailApplication.status)">{{ statusLabel(detailApplication.status) }}</Badge>
            </div>
            <div v-if="detailApplication.reason" class="py-2.5 flex justify-between text-red-600">
              <span class="font-medium">เหตุผลที่ไม่อนุมัติ</span>
              <span class="text-right">{{ detailApplication.reason }}</span>
            </div>
          </div>

          <div class="rounded-lg bg-slate-50 p-4 border border-slate-100 space-y-1">
            <span class="text-xs font-semibold text-slate-500 block">ที่อยู่ปัจจุบัน</span>
            <p class="text-slate-700 text-xs sm:text-sm leading-relaxed">{{ detailApplication.currentAddress }}</p>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="isDetailModalOpen = false">ปิดหน้าต่าง</Button>
          <Button
            v-if="detailApplication"
            class="bg-blue-600 text-white"
            @click="isDetailModalOpen = false; openActionDialog(detailApplication)"
          >
            พิจารณาคำร้องนี้
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog v-model:open="isActionDialogOpen">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle class="text-lg font-bold text-slate-800">พิจารณาคำร้องสวัสดิการ</DialogTitle>
          <DialogDescription v-if="selectedApplication" class="text-xs">
            ผู้ยื่นคำร้อง: <span class="font-semibold text-slate-700">{{ selectedApplication.fullName }}</span>
            (เลขบัตร: {{ formatThaiCitizenId(selectedApplication.citizenId) }})
          </DialogDescription>
        </DialogHeader>

        <div class="space-y-4 py-2">
          <dl v-if="selectedApplication" class="space-y-2 rounded-lg bg-slate-50 p-3.5 text-xs sm:text-sm border border-slate-100">
            <div class="flex justify-between gap-4"><dt class="text-slate-500">รายได้รวมครัวเรือน/ปี</dt><dd class="font-semibold text-slate-800">{{ selectedApplication.annualIncome.toLocaleString('th-TH') }} บาท</dd></div>
            <div class="flex justify-between gap-4"><dt class="text-slate-500">สถานะปัจจุบัน</dt><dd><Badge :class="statusClass(selectedApplication.status)">{{ statusLabel(selectedApplication.status) }}</Badge></dd></div>
          </dl>

          <div class="space-y-2">
            <Label class="text-xs font-semibold text-slate-700 uppercase tracking-wider">ผลการพิจารณา <span class="text-red-600">*</span></Label>
            <div class="grid grid-cols-2 gap-3">
              <Button
                type="button"
                variant="outline"
                :class="selectedDecision === 'approved'
                  ? 'border-emerald-600 bg-emerald-600 text-white hover:bg-emerald-700 hover:text-white'
                  : 'border-emerald-200 text-emerald-700 hover:bg-emerald-50'"
                @click="selectedDecision = 'approved'"
              >
                <CheckCircle2 class="mr-1.5 h-4 w-4" />
                อนุมัติ (Approve)
              </Button>
              <Button
                type="button"
                variant="outline"
                :class="selectedDecision === 'rejected'
                  ? 'border-rose-600 bg-rose-600 text-white hover:bg-rose-700 hover:text-white'
                  : 'border-rose-200 text-rose-700 hover:bg-rose-50'"
                @click="selectedDecision = 'rejected'"
              >
                <XCircle class="mr-1.5 h-4 w-4" />
                ไม่อนุมัติ (Reject)
              </Button>
            </div>
          </div>

          <div class="space-y-1.5">
            <Label for="review-reason" class="text-xs font-semibold text-slate-700">
              {{ selectedDecision === 'rejected' ? 'ระบุเหตุผลที่ไม่อนุมัติ (จำเป็น)' : 'บันทึกเพิ่มเติม (ไม่บังคับ)' }}
              <span v-if="selectedDecision === 'rejected'" class="text-red-600">*</span>
            </Label>
            <Input
              id="review-reason"
              v-model="reason"
              :placeholder="selectedDecision === 'rejected' ? 'เช่น รายได้เกินเกณฑ์ที่กำหนด' : 'ข้อความบันทึกช่วยจำ'"
            />
          </div>

          <div v-if="actionError" class="rounded-md bg-red-50 p-2.5 text-xs text-red-700 flex items-center gap-1.5 border border-red-200">
            <AlertCircle class="h-4 w-4 shrink-0 text-red-600" />
            <span>{{ actionError }}</span>
          </div>
        </div>

        <DialogFooter class="gap-2 sm:gap-0">
          <Button variant="outline" @click="isActionDialogOpen = false">ยกเลิก</Button>
          <Button :disabled="isUpdating" class="bg-blue-600 hover:bg-blue-700 text-white" @click="updateStatus">
            {{ isUpdating ? 'กำลังบันทึก...' : 'บันทึกผลการพิจารณา' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </main>
</template>