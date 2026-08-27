import axios from "axios";
import {
  SERVER_ERROR_TRIGGER_ID,
  sharedApplicationsStore,
  type ApplicationRecord,
} from "./mockStore";

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

export function mockCreateApplication(payload: CreateApplicationPayload) {
  return new Promise<{ data: CreateApplicationResponse }>((resolve, reject) => {
    setTimeout(() => {
      if (payload.citizenId === SERVER_ERROR_TRIGGER_ID) {
        reject(new Error("INTERNAL_SERVER_ERROR"));
        return;
      }

      const isAlreadyRegistered = sharedApplicationsStore.some(
        (app) => app.citizenId === payload.citizenId,
      );
      if (isAlreadyRegistered) {
        reject(new Error("CITIZEN_ALREADY_REGISTERED"));
        return;
      }

      const refNum = `WRS-2026-${Math.floor(100000 + Math.random() * 900000)}`;
      const newRecord: ApplicationRecord = {
        id: sharedApplicationsStore.length + 1,
        referenceNumber: refNum,
        citizenId: payload.citizenId,
        fullName: payload.fullName,
        birthDate: payload.birthDate,
        annualIncome: payload.annualIncome,
        currentAddress: payload.currentAddress,
        status: "pending",
        submittedAt: new Date().toISOString(),
      };

      sharedApplicationsStore.unshift(newRecord);

      resolve({
        data: {
          referenceNumber: refNum,
        },
      });
    }, 700);
  });
}

export function createApplication(payload: CreateApplicationPayload) {
  return axios.post<CreateApplicationResponse>("/api/v1/applications", payload);
}
