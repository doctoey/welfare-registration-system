import axios from 'axios'
import { sharedApplicationsStore } from './mockStore'

export interface OfficerApplication {
  id: number
  referenceNumber: string
  citizenId: string
  fullName: string
  birthDate: string
  annualIncome: number
  currentAddress: string
  status: 'pending' | 'approved' | 'rejected'
  reason?: string
  submittedAt: string
}

export function mockGetOfficerApplications() {
  return new Promise<{ data: OfficerApplication[] }>((resolve) => {
    setTimeout(() => {
      resolve({
        data: [...sharedApplicationsStore],
      })
    }, 400)
  })
}

export function mockUpdateApplicationStatus(
  id: number,
  status: OfficerApplication['status'],
  reason: string,
) {
  return new Promise<{ data: OfficerApplication }>((resolve, reject) => {
    setTimeout(() => {
      if (status === 'rejected' && (!reason || !reason.trim())) {
        reject(new Error('REASON_REQUIRED_FOR_REJECT'))
        return
      }

      const target = sharedApplicationsStore.find((app) => app.id === id)
      if (!target) {
        reject(new Error('APPLICATION_NOT_FOUND'))
        return
      }

      target.status = status
      target.reason = status === 'rejected' ? reason.trim() : undefined

      resolve({
        data: { ...target },
      })
    }, 400)
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
  return axios.patch<OfficerApplication>(
    `/api/v1/officer/applications/${id}/status`,
    {
      status,
      reason,
    },
  )
}
