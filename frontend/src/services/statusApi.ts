import axios from 'axios'

export const MOCK_CITIZEN_ID = '1100400123450'
export const MOCK_APPROVED_CITIZEN_ID = '1101700205673'
export const MOCK_REJECTED_CITIZEN_ID = '1102500337895'

export interface ApplicationStatus {
  referenceNumber: string
  citizenId: string
  fullName: string
  status: 'pending' | 'approved' | 'rejected'
  reason?: string
}

const mockStatusDatabase: Record<string, ApplicationStatus> = {
  [MOCK_CITIZEN_ID]: {
    referenceNumber: 'WRS-2026-000001',
    citizenId: MOCK_CITIZEN_ID,
    fullName: 'นายสมชาย ใจดี',
    status: 'pending',
  },
  [MOCK_APPROVED_CITIZEN_ID]: {
    referenceNumber: 'WRS-2026-000002',
    citizenId: MOCK_APPROVED_CITIZEN_ID,
    fullName: 'นางสาวสมหญิง รักดี',
    status: 'approved',
  },
  [MOCK_REJECTED_CITIZEN_ID]: {
    referenceNumber: 'WRS-2026-000003',
    citizenId: MOCK_REJECTED_CITIZEN_ID,
    fullName: 'นายวิชัย ตั้งใจ',
    status: 'rejected',
    reason: 'รายได้รวมของครัวเรือนเกินเกณฑ์ที่ระเบียบกำหนด หรือเอกสารไม่ครบถ้วน',
  },
}

export function mockGetApplicationStatus(citizenId: string) {
  return new Promise<{ data: ApplicationStatus }>((resolve, reject) => {
    setTimeout(() => {
      // ตรวจสอบข้อมูลใน mock database
      const found = mockStatusDatabase[citizenId]
      if (found) {
        resolve({ data: found })
        return
      }

      // หากไม่พบข้อมูล
      reject(new Error('APPLICATION_NOT_FOUND'))
    }, 600)
  })
}

export function getApplicationStatus(citizenId: string) {
  return axios.get<ApplicationStatus>(`/api/v1/applications/status/${citizenId}`)
}
