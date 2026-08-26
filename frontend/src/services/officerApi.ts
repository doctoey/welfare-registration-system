import axios from 'axios'

export interface OfficerApplication {
  id: number
  citizenId: string
  fullName: string
  annualIncome: number
  status: 'pending' | 'approved' | 'rejected'
  reason?: string
}

export function mockGetOfficerApplications() {
  return new Promise<{ data: OfficerApplication[] }>((resolve) => {
    setTimeout(() => {
      resolve({
        data: [
          { id: 1, citizenId: '1100400123458', fullName: 'สมชาย ใจดี', annualIncome: 54000, status: 'pending' },
          { id: 2, citizenId: '1101700205671', fullName: 'สมหญิง รักดี', annualIncome: 72000, status: 'approved' },
          { id: 3, citizenId: '1102500337896', fullName: 'วิชัย ตั้งใจ', annualIncome: 38000, status: 'rejected' },
        ],
      })
    }, 500)
  })
}

export function mockUpdateApplicationStatus(
  id: number,
  status: OfficerApplication['status'],
  reason: string,
) {
  return new Promise<{ data: OfficerApplication }>((resolve) => {
    setTimeout(() => {
      resolve({
        data: {
          id,
          citizenId: '',
          fullName: '',
          annualIncome: 0,
          status,
          reason,
        },
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
