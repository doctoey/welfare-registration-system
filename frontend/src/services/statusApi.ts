import axios from "axios";
import { sharedApplicationsStore } from "./mockStore";

export interface ApplicationStatus {
  referenceNumber: string;
  citizenId: string;
  fullName: string;
  status: "pending" | "approved" | "rejected";
  reason?: string;
  submittedAt?: string;
}

export function mockGetApplicationStatus(citizenId: string) {
  return new Promise<{ data: ApplicationStatus }>((resolve, reject) => {
    setTimeout(() => {
      const found = sharedApplicationsStore.find(
        (app) => app.citizenId === citizenId,
      );
      if (found) {
        resolve({
          data: {
            referenceNumber: found.referenceNumber,
            citizenId: found.citizenId,
            fullName: found.fullName,
            status: found.status,
            reason: found.reason,
            submittedAt: found.submittedAt,
          },
        });
        return;
      }

      reject(new Error("APPLICATION_NOT_FOUND"));
    }, 500);
  });
}

export function getApplicationStatus(
  citizenId: string,
  params?: { birthDate?: string; ref?: string },
) {
  return axios.get<ApplicationStatus>(
    `/api/v1/applications/status/${citizenId}`,
    {
      params,
    },
  );
}
