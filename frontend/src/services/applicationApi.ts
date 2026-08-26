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

export function mockGetApplicationStatus(citizenId: string) {
  return new Promise<{ data: ApplicationStatus }>((resolve) => {
    setTimeout(() => {
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
