import axios from "axios";

export interface CreateApplicationPayload {
  citizenId: string;
  fullName: string;
  birthDate: string;
  annualIncome: number;
  currentAddress: string;
}

export interface CreateApplicationResponse {
  referenceNumber: string;
}

export interface ApplicationStatus {
  referenceNumber: string;
  citizenId: string;
  fullName: string;
  status: 'pending' | 'approved' | 'rejected';
  reason?: string;
}

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

export function mockGetApplicationStatus(citizenId: string) {
  return new Promise<{ data: ApplicationStatus }>((resolve, reject) => {
    setTimeout(() => {
      if (citizenId !== '1100400123459') {
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

export function mockCreateApplication(payload: CreateApplicationPayload) {
  console.log("Mock payload:", payload);

  return new Promise<{ data: CreateApplicationResponse }>((resolve) => {
    setTimeout(() => {
      resolve({
        data: {
          referenceNumber: "WRS-2026-000001",
        },
      });
    }, 800);
  });
}

export function createApplication(payload: CreateApplicationPayload) {
  return axios.post<CreateApplicationResponse>("/api/v1/applications", payload);
}
