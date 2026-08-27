<script setup lang="ts">
import { ref, computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { isValidThaiCitizenId } from '@/lib/thaiCitizenId'
import { Textarea } from '@/components/ui/textarea'
import {
  AlertCircle,
  CheckCircle2,
  XCircle,
  FileText,
  User,
  Home,
  ShieldCheck,
  Copy,
  Check,
  Sparkles,
} from 'lucide-vue-next'
import { mockCreateApplication } from '../services/applicationApi'
import { SAMPLE_NEW_VALID_IDS } from '../services/mockStore'

const citizenId = ref('')
const fullName = ref('')
const dateOfBirth = ref('')
const annualIncome = ref('')
const currentAddress = ref('')
const consent = ref(false)
const validationErrors = ref<string[]>([])
const isSubmitting = ref(false)
const apiError = ref('')
const submitted = ref(false)
const referenceNumber = ref('')
const isCopied = ref(false)

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


const idStatus = computed(() => {
  const digits = citizenId.value.replace(/\D/g, '')
  if (digits.length === 0) return 'empty'
  if (digits.length < 13) return 'incomplete'
  return isValidThaiCitizenId(digits) ? 'valid' : 'invalid'
})


const formattedIncome = computed({
  get: () => {
    if (!annualIncome.value) return ''
    const num = parseInt(annualIncome.value.replace(/\D/g, ''), 10)
    return isNaN(num) ? '' : num.toLocaleString('th-TH')
  },
  set: (val: string) => {
    const raw = val.replace(/\D/g, '')
    annualIncome.value = raw
  },
})

function hasValidationError(message: string) {
  return validationErrors.value.includes(message)
}

function fillSampleData(sample: (typeof SAMPLE_NEW_VALID_IDS)[0]) {
  citizenId.value = sample.id
  fullName.value = sample.name
  dateOfBirth.value = '1992-06-15'
  annualIncome.value = String(sample.income)
  currentAddress.value = '123/45 หมู่ 6 ต.บ้านใหม่ อ.เมือง จ.ปทุมธานี 12000'
  consent.value = true
  validationErrors.value = []
  apiError.value = ''
}

function resetForm() {
  citizenId.value = ''
  fullName.value = ''
  dateOfBirth.value = ''
  annualIncome.value = ''
  currentAddress.value = ''
  consent.value = false
  validationErrors.value = []
  apiError.value = ''
  referenceNumber.value = ''
  submitted.value = false
  isCopied.value = false
}

async function copyReference() {
  if (!referenceNumber.value) return
  try {
    await navigator.clipboard.writeText(referenceNumber.value)
    isCopied.value = true
    setTimeout(() => {
      isCopied.value = false
    }, 2000)
  } catch (e) {
    console.error('Failed to copy', e)
  }
}

async function handleSubmit() {
  if (isSubmitting.value) return

  validationErrors.value = []
  apiError.value = ''

  const citizenIdDigits = citizenId.value.replace(/\D/g, '')
  if (!isValidThaiCitizenId(citizenIdDigits)) {
    validationErrors.value.push('National ID is invalid')
  }
  if (!fullName.value.trim()) validationErrors.value.push('Full name is required')
  if (!dateOfBirth.value) validationErrors.value.push('Date of birth is required')
  if (!annualIncome.value || Number(annualIncome.value) < 0) {
    validationErrors.value.push('Annual income is required')
  }
  if (!currentAddress.value.trim()) validationErrors.value.push('Current address is required')
  if (!consent.value) validationErrors.value.push('Consent is required')

  if (validationErrors.value.length > 0) return

  const payload = {
    citizenId: citizenIdDigits,
    fullName: fullName.value.trim(),
    birthDate: dateOfBirth.value,
    annualIncome: Number(annualIncome.value),
    currentAddress: currentAddress.value.trim(),
  }

  isSubmitting.value = true

  try {
    const response = await mockCreateApplication(payload)
    referenceNumber.value = response.data.referenceNumber
    submitted.value = true
  } catch (error) {
    console.error(error)
    if (error instanceof Error) {
      if (error.message === 'CITIZEN_ALREADY_REGISTERED') {
        apiError.value = 'เลขประจำตัวประชาชนนี้ได้ลงทะเบียนในระบบแล้ว ไม่สามารถลงทะเบียนซ้ำได้'
      } else if (error.message === 'INCOME_EXCEEDS_LIMIT') {
        apiError.value = 'รายได้เกินเกณฑ์ที่กำหนดสำหรับการลงทะเบียน (ต้องไม่เกิน 100,000 บาท/ปี)'
      } else {
        apiError.value = 'ไม่สามารถส่งคำร้องได้ กรุณาลองใหม่อีกครั้ง'
      }
    } else {
      apiError.value = 'เกิดข้อผิดพลาดที่ไม่ทราบสาเหตุ กรุณาลองใหม่อีกครั้ง'
    }
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main class="py-6 px-4 sm:px-6">
    <!-- Success State -->
    <Card v-if="submitted" class="mx-auto max-w-2xl p-8 text-center border-emerald-200 bg-white shadow-sm">
      <div
        class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-emerald-50 border-2 border-emerald-200 text-emerald-600">
        <CheckCircle2 class="h-8 w-8" />
      </div>
      <h1 class="text-2xl font-bold text-slate-800">ส่งคำร้องสำเร็จ</h1>
      <p class="text-sm text-slate-500 mt-1">Application Submitted Successfully</p>
      <p class="mt-4 text-sm text-slate-600 max-w-md mx-auto">
        ระบบได้รับข้อมูลคำร้องสวัสดิการของท่านเรียบร้อยแล้ว ข้อมูลถูกบันทึกและพร้อมสำหรับเจ้าหน้าที่พิจารณา
      </p>

      <!-- Reference Number Box with Copy Action -->
      <div
        class="mt-6 inline-flex flex-col items-center justify-center rounded-xl bg-slate-50 p-4 border border-slate-200 shadow-2xs">
        <span class="text-xs font-semibold text-slate-500 uppercase tracking-wider block">หมายเลขอ้างอิง (Reference
          Number)</span>
        <div class="mt-1 flex items-center gap-3">
          <span class="font-mono text-2xl font-bold text-blue-700 tracking-wide">{{ referenceNumber }}</span>
          <button type="button"
            class="flex items-center gap-1 rounded-md bg-white px-2.5 py-1 text-xs font-medium text-slate-700 border border-slate-200 hover:bg-slate-100 shadow-2xs transition"
            @click="copyReference">
            <component :is="isCopied ? Check : Copy" class="h-3.5 w-3.5"
              :class="isCopied ? 'text-emerald-600' : 'text-slate-500'" />
            <span>{{ isCopied ? 'คัดลอกแล้ว' : 'คัดลอก' }}</span>
          </button>
        </div>
      </div>

      <div class="mt-8 flex flex-wrap justify-center gap-3">
        <Button variant="outline" @click="resetForm">
          ยื่นคำร้องใหม่อีกครั้ง
        </Button>
        <RouterLink to="/status">
          <Button class="bg-blue-600 hover:bg-blue-700 text-white">
            ไปหน้าตรวจสอบสถานะ
          </Button>
        </RouterLink>
      </div>
    </Card>

    <!-- Form State -->
    <Card v-else class="mx-auto max-w-3xl overflow-hidden bg-white shadow-sm border-slate-200">
      <!-- Header Banner -->
      <div class="border-b border-slate-100 bg-slate-50/70 p-6 flex flex-wrap items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-blue-600 text-white shadow-sm">
            <FileText class="h-6 w-6" />
          </div>
          <div>
            <h1 class="text-xl font-bold text-slate-800">ระบบลงทะเบียนสวัสดิการแห่งรัฐ</h1>
            <p class="text-xs sm:text-sm text-slate-500">Welfare Registration System</p>
          </div>
        </div>

        <!-- Quick Sample Fill Buttons for various test cases -->
        <div class="flex flex-wrap items-center gap-2">
          <span class="text-xs text-slate-400">ทดสอบเคส:</span>
          <button type="button"
            class="flex items-center gap-1 rounded-md bg-emerald-50 px-2.5 py-1 text-xs font-medium text-emerald-700 border border-emerald-200 hover:bg-emerald-100 transition shadow-2xs"
            @click="fillSampleData(SAMPLE_NEW_VALID_IDS[0])">
            <Sparkles class="h-3 w-3 text-emerald-600" />
            <span>มาลี (เลขใหม่)</span>
          </button>
          <button type="button"
            class="flex items-center gap-1 rounded-md bg-amber-50 px-2.5 py-1 text-xs font-medium text-amber-700 border border-amber-200 hover:bg-amber-100 transition shadow-2xs"
            @click="fillSampleData(SAMPLE_NEW_VALID_IDS[1])">
            <AlertCircle class="h-3 w-3 text-amber-600" />
            <span>สมชาย (เลขซ้ำในระบบ)</span>
          </button>
        </div>
      </div>

      <!-- Notice Bar -->
      <div class="border-b border-amber-200 bg-amber-50 px-6 py-3 flex items-start gap-2.5">
        <AlertCircle class="h-4 w-4 text-amber-600 shrink-0 mt-0.5" />
        <p class="text-xs text-amber-800 leading-relaxed">
          <span class="font-semibold">ข้อควรทราบ:</span>
          กรุณากรอกข้อมูลให้ตรงตามความเป็นจริงและตรงกับบัตรประจำตัวประชาชน เพื่อประโยชน์ในการพิจารณาสิทธิ์
        </p>
      </div>

      <div class="p-6 sm:p-8 space-y-8">
        <!-- Section 1: Personal Information -->
        <section class="space-y-4">
          <div class="flex items-center gap-2 border-b border-slate-100 pb-2">
            <div
              class="flex h-6 w-6 items-center justify-center rounded-full bg-blue-100 text-xs font-semibold text-blue-700">
              1
            </div>
            <h2 class="text-base font-semibold text-slate-800 flex items-center gap-1.5">
              <User class="h-4 w-4 text-slate-400" />
              ข้อมูลส่วนตัว <span class="text-xs font-normal text-slate-400">(Personal Information)</span>
            </h2>
          </div>

          <div class="space-y-4">
            <!-- National ID with Auto-formatting and Live Checksum Indicator -->
            <div class="space-y-1.5">
              <Label for="citizen-id" class="text-sm font-medium text-slate-700">
                เลขประจำตัวประชาชน (National ID) <span class="text-red-600">*</span>
              </Label>
              <div class="relative">
                <Input id="citizen-id" v-model="formattedCitizenId" type="text" inputmode="numeric" maxlength="17"
                  placeholder="X-XXXX-XXXXX-XX-X" :aria-invalid="hasValidationError('National ID is invalid')"
                  class="font-mono text-base tracking-wider pr-10" />
                <!-- Live indicator / Clear button -->
                <div class="absolute right-3 top-1/2 -translate-y-1/2 flex items-center">
                  <CheckCircle2 v-if="idStatus === 'valid'" class="h-5 w-5 text-emerald-500 pointer-events-none" />
                  <button v-else-if="idStatus === 'invalid'" type="button"
                    class="text-rose-500 hover:text-rose-700 transition cursor-pointer p-0.5 rounded"
                    title="ล้างข้อมูลเลขบัตร" @click="citizenId = ''">
                    <XCircle class="h-5 w-5" />
                  </button>
                </div>
              </div>

              <div class="flex items-center justify-between text-xs text-slate-400">
                <span>กรอกตัวเลข 13 หลักตามบัตรประชาชน</span>
                <span>{{ citizenId.length }}/13</span>
              </div>

              <p v-if="hasValidationError('National ID is invalid')"
                class="text-xs text-red-600 flex items-center gap-1">
                <AlertCircle class="h-3.5 w-3.5 shrink-0" />
                กรุณาระบุเลขประจำตัวประชาชน 13 หลักให้ถูกต้องตามระบบตรวจสอบ (Checksum)
              </p>
            </div>

            <!-- Full Name & Date of Birth (Grid) -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <Label for="full-name" class="text-sm font-medium text-slate-700">
                  ชื่อ - นามสกุล (Full Name) <span class="text-red-600">*</span>
                </Label>
                <Input id="full-name" v-model="fullName" placeholder="เช่น นายสมชาย ใจดี"
                  :aria-invalid="hasValidationError('Full name is required')" />
                <p v-if="hasValidationError('Full name is required')"
                  class="text-xs text-red-600 flex items-center gap-1">
                  <AlertCircle class="h-3.5 w-3.5 shrink-0" />
                  กรุณาระบุชื่อ-นามสกุล
                </p>
              </div>

              <div class="space-y-1.5">
                <Label for="date-of-birth" class="text-sm font-medium text-slate-700">
                  วัน/เดือน/ปี เกิด (Date of Birth) <span class="text-red-600">*</span>
                </Label>
                <Input id="date-of-birth" v-model="dateOfBirth" type="date"
                  :aria-invalid="hasValidationError('Date of birth is required')" />
                <p v-if="hasValidationError('Date of birth is required')"
                  class="text-xs text-red-600 flex items-center gap-1">
                  <AlertCircle class="h-3.5 w-3.5 shrink-0" />
                  กรุณาเลือกวันเดือนปีเกิด
                </p>
              </div>
            </div>
          </div>
        </section>

        <!-- Section 2: Financial & Address -->
        <section class="space-y-4">
          <div class="flex items-center gap-2 border-b border-slate-100 pb-2">
            <div
              class="flex h-6 w-6 items-center justify-center rounded-full bg-blue-100 text-xs font-semibold text-blue-700">
              2
            </div>
            <h2 class="text-base font-semibold text-slate-800 flex items-center gap-1.5">
              <Home class="h-4 w-4 text-slate-400" />
              ข้อมูลรายได้และที่อยู่ <span class="text-xs font-normal text-slate-400">(Financial & Address)</span>
            </h2>
          </div>

          <div class="space-y-4">
            <!-- Annual Income with Commas -->
            <div class="space-y-1.5">
              <Label for="annual-income" class="text-sm font-medium text-slate-700">
                รายได้รวมครัวเรือนต่อปี (Annual Household Income) <span class="text-red-600">*</span>
              </Label>
              <div class="relative flex">
                <Input id="annual-income" v-model="formattedIncome" type="text" inputmode="numeric"
                  placeholder="เช่น 54,000" :aria-invalid="hasValidationError('Annual income is required')"
                  class="rounded-r-none border-r-0 font-mono text-right" />
                <span
                  class="inline-flex items-center rounded-r-md border border-l-0 border-slate-300 bg-slate-100 px-3 text-xs font-medium text-slate-600 select-none">
                  บาท / THB
                </span>
              </div>
              <p class="text-xs text-slate-400">รายได้รวมของทุกคนในครอบครัวต่อปี (เกณฑ์สิทธิ์: ไม่เกิน 100,000 บาท)</p>
              <p v-if="hasValidationError('Annual income is required')"
                class="text-xs text-red-600 flex items-center gap-1">
                <AlertCircle class="h-3.5 w-3.5 shrink-0" />
                กรุณาระบุรายได้รวมต่อปีเป็นตัวเลขที่ถูกต้อง (0 บาทขึ้นไป)
              </p>
            </div>

            <!-- Current Address -->
            <div class="space-y-1.5">
              <Label for="current-address" class="text-sm font-medium text-slate-700">
                ที่อยู่ปัจจุบัน (Current Address) <span class="text-red-600">*</span>
              </Label>
              <Textarea id="current-address" v-model="currentAddress" :rows="3"
                placeholder="บ้านเลขที่, หมู่, ถนน, ตำบล/แขวง, อำเภอ/เขต, จังหวัด, รหัสไปรษณีย์"
                :aria-invalid="hasValidationError('Current address is required')" />
              <p v-if="hasValidationError('Current address is required')"
                class="text-xs text-red-600 flex items-center gap-1">
                <AlertCircle class="h-3.5 w-3.5 shrink-0" />
                กรุณาระบุที่อยู่ปัจจุบันให้ครบถ้วน
              </p>
            </div>
          </div>
        </section>

        <!-- Section 3: Consent & Declaration -->
        <section class="space-y-4">
          <div class="flex items-center gap-2 border-b border-slate-100 pb-2">
            <div
              class="flex h-6 w-6 items-center justify-center rounded-full bg-blue-100 text-xs font-semibold text-blue-700">
              3
            </div>
            <h2 class="text-base font-semibold text-slate-800 flex items-center gap-1.5">
              <ShieldCheck class="h-4 w-4 text-slate-400" />
              การยินยอมและรับรองข้อมูล <span class="text-xs font-normal text-slate-400">(Declaration & Consent)</span>
            </h2>
          </div>

          <div :class="[
            'rounded-lg border p-4 transition-colors flex items-start gap-3 cursor-pointer',
            consent ? 'border-blue-200 bg-blue-50/50' : 'border-slate-200 bg-slate-50/50'
          ]" @click="consent = !consent">
            <input id="consent" v-model="consent" type="checkbox"
              class="mt-1 h-4 w-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500 cursor-pointer"
              @click.stop />
            <Label for="consent" class="text-xs sm:text-sm text-slate-700 leading-relaxed cursor-pointer font-normal">
              ข้าพเจ้าขอรับรองว่าข้อความและข้อมูลที่ได้ระบุไว้ในคำร้องนี้เป็นความจริงทุกประการ
              และยินยอมให้หน่วยงานตรวจสอบความถูกต้องของข้อมูล <span class="text-red-600 font-semibold">*</span>
            </Label>
          </div>
          <p v-if="hasValidationError('Consent is required')" class="text-xs text-red-600 flex items-center gap-1">
            <AlertCircle class="h-3.5 w-3.5 shrink-0" />
            กรุณากดยินยอมและรับรองความถูกต้องของข้อมูลก่อนส่งคำร้อง
          </p>
        </section>

        <!-- Error Message -->
        <div v-if="apiError"
          class="rounded-lg bg-red-50 border border-red-200 p-4 text-sm text-red-700 flex items-center gap-2">
          <AlertCircle class="h-5 w-5 shrink-0 text-red-600" />
          <span>{{ apiError }}</span>
        </div>

        <!-- Submit Button -->
        <div class="pt-2">
          <Button type="button" :disabled="isSubmitting"
            class="w-full bg-blue-600 hover:bg-blue-700 text-white py-2.5 font-medium shadow-sm transition-all"
            @click="handleSubmit">
            <span v-if="isSubmitting" class="flex items-center justify-center gap-2">
              <svg class="animate-spin h-4 w-4 text-white" viewBox="0 0 24 24" fill="none">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>
              กำลังส่งข้อมูลคำร้อง...
            </span>
            <span v-else>ส่งคำร้องลงทะเบียน (Submit Application)</span>
          </Button>
        </div>
      </div>
    </Card>
  </main>
</template>