import axios from 'axios'

export interface CreateApplicationPayload {
  citizenId: string
  fullName: string
  birthDate: string
  annualIncome: number
  currentAddress: string
}

export interface CreateApplicationResponse {
  referenceNumber: string
}

// เลขบัตรจำลองที่มีอยู่ในระบบแล้ว (เพื่อทดสอบ Duplicate Citizen ID)
export const ALREADY_REGISTERED_CITIZEN_ID = '1100400123450'

// เลขบัตรจำลองสำหรับทดสอบจำลองข้อผิดพลาดจากเซิร์ฟเวอร์
export const ERROR_TRIGGER_CITIZEN_ID = '1102500337895'

export function mockCreateApplication(payload: CreateApplicationPayload) {
  console.log('Mock create application payload:', payload)

  return new Promise<{ data: CreateApplicationResponse }>((resolve, reject) => {
    setTimeout(() => {
      // 1. ตรวจสอบกรณีเลขประจำตัวประชาชนซ้ำกับที่มีในระบบ
      if (payload.citizenId === ALREADY_REGISTERED_CITIZEN_ID) {
        reject(new Error('CITIZEN_ALREADY_REGISTERED'))
        return
      }

      // 2. ตรวจสอบกรณีจำลองเซิร์ฟเวอร์เกิดปัญหา
      if (payload.citizenId === ERROR_TRIGGER_CITIZEN_ID) {
        reject(new Error('INTERNAL_SERVER_ERROR'))
        return
      }

      // 3. ตรวจสอบกรณีรายได้เกินเกณฑ์ที่กำหนด (เช่น เกิน 100,000 บาท)
      if (payload.annualIncome > 100000) {
        reject(new Error('INCOME_EXCEEDS_LIMIT'))
        return
      }

      // 4. สำเร็จ: สุ่มสร้างเลขที่อ้างอิง
      const randomRefNum = Math.floor(100000 + Math.random() * 900000)
      resolve({
        data: {
          referenceNumber: `WRS-2026-${randomRefNum}`,
        },
      })
    }, 800)
  })
}

export function createApplication(payload: CreateApplicationPayload) {
  return axios.post<CreateApplicationResponse>('/api/v1/applications', payload)
}
