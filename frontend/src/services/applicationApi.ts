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

export function mockCreateApplication(payload: CreateApplicationPayload) {
  console.log('Mock payload:', payload)

  return new Promise<{ data: CreateApplicationResponse }>((resolve) => {
    setTimeout(() => {
      resolve({
        data: {
          referenceNumber: 'WRS-2026-000001',
        },
      })
    }, 800)
  })
}

export function createApplication(payload: CreateApplicationPayload) {
  return axios.post<CreateApplicationResponse>('/api/v1/applications', payload)
}
