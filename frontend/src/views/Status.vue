<script setup lang="ts">
import { ref, computed } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { statusClass, statusLabel } from '@/lib/applicationStatus'
import { isValidThaiCitizenId } from '@/lib/thaiCitizenId'
import {
  Search,
  AlertCircle,
  CheckCircle2,
  Clock,
  XCircle,
  ShieldCheck,
  Sparkles,
} from 'lucide-vue-next'
import { mockGetApplicationStatus, type ApplicationStatus } from '@/services/statusApi'

const citizenId = ref('')
const status = ref<ApplicationStatus | null>(null)
const errorMessage = ref('')
const isSearching = ref(false)

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

const formattedCitizenId = computed({
  get: () => formatThaiCitizenId(citizenId.value),
  set: (val: string) => {
    citizenId.value = val.replace(/\D/g, '').slice(0, 13)
  },
})

function fillMockId(id: string) {
  citizenId.value = id
  errorMessage.value = ''
  status.value = null
}

async function searchStatus() {
  errorMessage.value = ''
  status.value = null

  const citizenIdDigits = citizenId.value.replace(/\D/g, '')
  if (!citizenIdDigits) {
    errorMessage.value = 'กรุณาระบุเลขประจำตัวประชาชน'
    return
  }

  if (!isValidThaiCitizenId(citizenIdDigits)) {
    errorMessage.value = 'เลขประจำตัวประชาชนไม่ถูกต้องตามรูปแบบ (Checksum ไม่ถูกต้อง)'
    return
  }

  isSearching.value = true
  try {
    const response = await mockGetApplicationStatus(citizenIdDigits)
    status.value = response.data
  } catch (error) {
    errorMessage.value =
      error instanceof Error && error.message === 'APPLICATION_NOT_FOUND'
        ? 'ไม่พบข้อมูลคำร้องของเลขบัตรประชาชนนี้ในระบบ กรุณาตรวจสอบเลขบัตรหรือยื่นคำร้องใหม่'
        : 'ไม่สามารถค้นหาสถานะได้ กรุณาลองใหม่อีกครั้ง'
  } finally {
    isSearching.value = false
  }
}
</script>

<template>
  <main class="py-6 px-4 sm:px-6 space-y-6">
    <!-- Search Form Card -->
    <Card class="mx-auto max-w-3xl overflow-hidden bg-white shadow-sm border-slate-200">
      <!-- Header Banner -->
      <div class="border-b border-slate-100 bg-slate-50/70 p-6">
        <div class="flex items-center gap-3">
          <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-600 text-white shadow-sm">
            <Search class="h-6 w-6" />
          </div>
          <div>
            <h1 class="text-xl font-bold text-slate-800">ติดตามสถานะคำร้องสวัสดิการ</h1>
            <p class="text-xs sm:text-sm text-slate-500">Track Welfare Application Status</p>
          </div>
        </div>
      </div>

      <div class="p-6 sm:p-8 space-y-6">
        <form class="space-y-4" @submit.prevent="searchStatus">
          <div class="space-y-2">
            <Label for="status-citizen-id" class="text-sm font-medium text-slate-700">
              เลขประจำตัวประชาชน (National ID) <span class="text-red-600">*</span>
            </Label>
            <div class="relative">
              <Input
                id="status-citizen-id"
                v-model="formattedCitizenId"
                type="text"
                inputmode="numeric"
                maxlength="17"
                placeholder="X-XXXX-XXXXX-XX-X"
                class="font-mono text-base tracking-wider pr-10"
              />
              <button
                v-if="citizenId"
                type="button"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition cursor-pointer p-0.5"
                title="ล้างข้อมูล"
                @click="citizenId = ''; status = null; errorMessage = ''"
              >
                <XCircle class="h-4 w-4" />
              </button>
              <Search
                v-else
                class="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400 pointer-events-none"
              />
            </div>

            <!-- Quick Selection Mock Chips -->
            <div class="rounded-lg bg-slate-50 p-3 border border-slate-100 text-xs text-slate-600 space-y-1.5">
              <span class="font-semibold text-slate-700 flex items-center gap-1">
                <Sparkles class="h-3.5 w-3.5 text-blue-600" />
                คลิกเลือกเลขตัวอย่างในระบบเพื่อทดสอบค้นหา:
              </span>
              <div class="flex flex-wrap gap-2 pt-1">
                <button
                  type="button"
                  class="rounded bg-amber-50 px-2 py-1 text-amber-700 border border-amber-200 hover:bg-amber-100 transition shadow-2xs font-mono"
                  @click="fillMockId('1100400123450')"
                >
                  1-1004-00123-45-0 (รอตรวจสอบ)
                </button>
                <button
                  type="button"
                  class="rounded bg-emerald-50 px-2 py-1 text-emerald-700 border border-emerald-200 hover:bg-emerald-100 transition shadow-2xs font-mono"
                  @click="fillMockId('1101700205673')"
                >
                  1-1017-00205-67-3 (อนุมัติแล้ว)
                </button>
                <button
                  type="button"
                  class="rounded bg-rose-50 px-2 py-1 text-rose-700 border border-rose-200 hover:bg-rose-100 transition shadow-2xs font-mono"
                  @click="fillMockId('1102500337895')"
                >
                  1-1025-00337-89-5 (ไม่อนุมัติ)
                </button>
              </div>
            </div>
          </div>

          <!-- Error Alert -->
          <div
            v-if="errorMessage"
            class="rounded-lg bg-red-50 border border-red-200 p-4 text-sm text-red-700 flex items-center gap-2"
          >
            <AlertCircle class="h-5 w-5 shrink-0 text-red-600" />
            <span>{{ errorMessage }}</span>
          </div>

          <Button
            type="submit"
            :disabled="isSearching"
            class="w-full bg-blue-600 hover:bg-blue-700 text-white py-2.5 font-medium shadow-sm transition"
          >
            <span v-if="isSearching" class="flex items-center justify-center gap-2">
              <svg class="animate-spin h-4 w-4 text-white" viewBox="0 0 24 24" fill="none">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              กำลังค้นหา...
            </span>
            <span v-else class="flex items-center justify-center gap-2">
              <Search class="h-4 w-4" />
              ตรวจสอบสถานะคำร้อง
            </span>
          </Button>
        </form>
      </div>
    </Card>

    <!-- Search Result Card with Timeline Step Progress -->
    <Card v-if="status" class="mx-auto max-w-3xl overflow-hidden bg-white shadow-sm border-slate-200">
      <div class="border-b border-slate-100 bg-slate-50/70 px-6 py-4 flex items-center justify-between">
        <h2 class="text-base font-bold text-slate-800 flex items-center gap-2">
          <ShieldCheck class="h-5 w-5 text-blue-600" />
          ผลการตรวจสอบสถานะคำร้อง
        </h2>
        <Badge :class="statusClass(status.status)" class="text-xs px-2.5 py-0.5">
          {{ statusLabel(status.status) }}
        </Badge>
      </div>

      <div class="p-6 sm:p-8 space-y-6">
        <!-- Progress Stepper (Timeline) -->
        <div class="relative py-2">
          <div class="flex items-center justify-between">
            <!-- Step 1 -->
            <div class="flex flex-col items-center flex-1">
              <div
                class="flex h-8 w-8 items-center justify-center rounded-full bg-emerald-600 text-white text-xs font-bold shadow-2xs"
              >
                <CheckCircle2 class="h-5 w-5" />
              </div>
              <span class="text-xs font-medium text-slate-800 mt-1.5">1. ยื่นคำร้อง</span>
              <span class="text-[10px] text-slate-400">สำเร็จ</span>
            </div>

            <div class="h-0.5 flex-1 bg-emerald-500 mx-2"></div>

            <!-- Step 2 -->
            <div class="flex flex-col items-center flex-1">
              <div
                class="flex h-8 w-8 items-center justify-center rounded-full text-xs font-bold shadow-2xs"
                :class="
                  status.status === 'pending'
                    ? 'bg-amber-500 text-white animate-pulse'
                    : 'bg-emerald-600 text-white'
                "
              >
                <Clock v-if="status.status === 'pending'" class="h-5 w-5" />
                <CheckCircle2 v-else class="h-5 w-5" />
              </div>
              <span class="text-xs font-medium text-slate-800 mt-1.5">2. เจ้าหน้าที่ตรวจสอบ</span>
              <span
                class="text-[10px]"
                :class="status.status === 'pending' ? 'text-amber-600 font-semibold' : 'text-slate-400'"
              >
                {{ status.status === 'pending' ? 'กำลังดำเนินการ' : 'ตรวจสอบแล้ว' }}
              </span>
            </div>

            <div
              class="h-0.5 flex-1 mx-2"
              :class="
                status.status === 'pending'
                  ? 'bg-slate-200'
                  : status.status === 'approved'
                    ? 'bg-emerald-500'
                    : 'bg-rose-500'
              "
            ></div>

            <!-- Step 3 -->
            <div class="flex flex-col items-center flex-1">
              <div
                class="flex h-8 w-8 items-center justify-center rounded-full text-xs font-bold shadow-2xs"
                :class="
                  status.status === 'approved'
                    ? 'bg-emerald-600 text-white'
                    : status.status === 'rejected'
                      ? 'bg-rose-600 text-white'
                      : 'bg-slate-200 text-slate-500'
                "
              >
                <CheckCircle2 v-if="status.status === 'approved'" class="h-5 w-5" />
                <XCircle v-else-if="status.status === 'rejected'" class="h-5 w-5" />
                <span v-else>3</span>
              </div>
              <span class="text-xs font-medium text-slate-800 mt-1.5">3. ผลการพิจารณา</span>
              <span
                class="text-[10px] font-semibold"
                :class="
                  status.status === 'approved'
                    ? 'text-emerald-600'
                    : status.status === 'rejected'
                      ? 'text-rose-600'
                      : 'text-slate-400'
                "
              >
                {{ statusLabel(status.status) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Status Specific Banner -->
        <div
          v-if="status.status === 'approved'"
          class="rounded-xl border border-emerald-200 bg-emerald-50/70 p-4 flex items-start gap-3"
        >
          <CheckCircle2 class="h-6 w-6 text-emerald-600 shrink-0 mt-0.5" />
          <div>
            <h3 class="font-semibold text-emerald-900">คำร้องได้รับการอนุมัติเรียบร้อยแล้ว</h3>
            <p class="text-xs text-emerald-700 mt-0.5">ท่านผ่านเกณฑ์การพิจารณาคุณสมบัติเพื่อรับสิทธิสวัสดิการแห่งรัฐ</p>
          </div>
        </div>

        <div
          v-else-if="status.status === 'rejected'"
          class="rounded-xl border border-rose-200 bg-rose-50/70 p-4 flex items-start gap-3"
        >
          <XCircle class="h-6 w-6 text-rose-600 shrink-0 mt-0.5" />
          <div>
            <h3 class="font-semibold text-rose-900">คำร้องไม่ผ่านการอนุมัติ</h3>
            <p class="text-xs text-rose-700 mt-0.5">{{ status.reason || 'ไม่ผ่านเกณฑ์การพิจารณาคุณสมบัติ' }}</p>
          </div>
        </div>

        <div v-else class="rounded-xl border border-amber-200 bg-amber-50/70 p-4 flex items-start gap-3">
          <Clock class="h-6 w-6 text-amber-600 shrink-0 mt-0.5" />
          <div>
            <h3 class="font-semibold text-amber-900">อยู่ระหว่างการตรวจสอบข้อมูล</h3>
            <p class="text-xs text-amber-700 mt-0.5">
              เจ้าหน้าที่กำลังดำเนินการตรวจสอบคุณสมบัติและเอกสารหลักฐานของท่าน
            </p>
          </div>
        </div>

        <!-- Details List -->
        <div class="rounded-xl border border-slate-200 bg-slate-50/50 p-4 divide-y divide-slate-100">
          <div class="py-2.5 flex justify-between items-center text-sm">
            <span class="text-slate-500 font-medium">ชื่อ-นามสกุล</span>
            <span class="font-semibold text-slate-800">{{ status.fullName }}</span>
          </div>
          <div class="py-2.5 flex justify-between items-center text-sm">
            <span class="text-slate-500 font-medium">เลขประจำตัวประชาชน</span>
            <span class="font-mono text-slate-800">{{ formatThaiCitizenId(status.citizenId) }}</span>
          </div>
          <div class="py-2.5 flex justify-between items-center text-sm">
            <span class="text-slate-500 font-medium">หมายเลขอ้างอิง</span>
            <span class="font-mono font-bold text-blue-700">{{ status.referenceNumber }}</span>
          </div>
          <div class="py-2.5 flex justify-between items-center text-sm">
            <span class="text-slate-500 font-medium">สถานะคำร้อง</span>
            <Badge :class="statusClass(status.status)">{{ statusLabel(status.status) }}</Badge>
          </div>
        </div>
      </div>
    </Card>
  </main>
</template>