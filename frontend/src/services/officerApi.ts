import axios from 'axios'

export interface OfficerApplication {
  id: number
  citizenId: string
  fullName: string
  annualIncome: number
  status: 'pending' | 'approved' | 'rejected'
  reason?: string
}

let mockApplicationsStore: OfficerApplication[] = [
  { id: 1, citizenId: '1100400123450', fullName: 'สมชาย ใจดี', annualIncome: 54000, status: 'pending' },
  { id: 2, citizenId: '1101700205673', fullName: 'สมหญิง รักดี', annualIncome: 72000, status: 'approved' },
  { id: 3, citizenId: '1102500337895', fullName: 'วิชัย ตั้งใจ', annualIncome: 38000, status: 'rejected', reason: 'เอกสารไม่ตรงตามความเป็นจริง' },
]

export function mockGetOfficerApplications() {
  return new Promise<{ data: OfficerApplication[] }>((resolve) => {
    setTimeout(() => {
      resolve({
        data: [...mockApplicationsStore],
      })
    }, 500)
  })
}

export function mockUpdateApplicationStatus(
  id: number,
  status: OfficerApplication['status'],
  reason: string,
) {
  return new Promise<{ data: OfficerApplication }>((resolve, reject) => {
    setTimeout(() => {
      // 1. ตรวจสอบเงื่อนไข Reject ต้องมี reason
      if (status === 'rejected' && (!reason || !reason.trim())) {
        reject(new Error('REASON_REQUIRED_FOR_REJECT'))
        return
      }

      // 2. ค้นหารายการคำร้อง
      const target = mockApplicationsStore.find((app) => app.id === id)
      if (!target) {
        reject(new Error('APPLICATION_NOT_FOUND'))
        return
      }

      // 3. อัปเดตข้อมูล
      target.status = status
      target.reason = status === 'rejected' ? reason.trim() : undefined

      resolve({
        data: { ...target },
      })
    }, 500)
  })
}

export function getOfficerApplications() {
  return axios.get<OfficerApplication[]>('/api/v1/officer/applications')
}

export function updateApplicationStatus(
  id: number,
  status: OfficerApplication['status'],
  reason: string,
) {
  return axios.patch<OfficerApplication>(`/api/v1/officer/applications/${id}/status`, {
    status,
    reason,
  })
}
