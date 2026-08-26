import axios from 'axios'

export const MOCK_CITIZEN_ID = '1100400123450'

export interface ApplicationStatus {
  referenceNumber: string
  citizenId: string
  fullName: string
  status: 'pending' | 'approved' | 'rejected'
  reason?: string
}

export function mockGetApplicationStatus(citizenId: string) {
  return new Promise<{ data: ApplicationStatus }>((resolve, reject) => {
    setTimeout(() => {
      if (citizenId !== MOCK_CITIZEN_ID) {
        reject(new Error('APPLICATION_NOT_FOUND'))
        return
      }

      resolve({
        data: {
          referenceNumber: 'WRS-2026-000001',
          citizenId,
          fullName: 'สมชาย ใจดี',
          status: 'pending',
        },
      })
    }, 600)
  })
}

export function getApplicationStatus(citizenId: string) {
  return axios.get<ApplicationStatus>(`/api/v1/applications/status/${citizenId}`)
}
