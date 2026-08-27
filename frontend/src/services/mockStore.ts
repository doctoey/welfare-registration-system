export interface ApplicationRecord {
  id: number;
  referenceNumber: string;
  citizenId: string;
  fullName: string;
  birthDate: string;
  annualIncome: number;
  currentAddress: string;
  status: "pending" | "approved" | "rejected";
  reason?: string;
  submittedAt: string;
}

export const sharedApplicationsStore: ApplicationRecord[] = [
  {
    id: 1,
    referenceNumber: "WRS-2026-000001",
    citizenId: "1100400123450",
    fullName: "สมชาย ใจดี",
    birthDate: "1990-05-12",
    annualIncome: 54000,
    currentAddress: "99/1 แขวงดินแดง เขตดินแดง กทม. 10400",
    status: "pending",
    submittedAt: "2026-08-25T09:30:00Z",
  },
  {
    id: 2,
    referenceNumber: "WRS-2026-000002",
    citizenId: "1101700205673",
    fullName: "สมหญิง รักดี",
    birthDate: "1988-11-20",
    annualIncome: 72000,
    currentAddress: "12/4 ต.สุเทพ อ.เมือง จ.เชียงใหม่ 50200",
    status: "approved",
    submittedAt: "2026-08-20T14:15:00Z",
  },
  {
    id: 3,
    referenceNumber: "WRS-2026-000003",
    citizenId: "1102500337895",
    fullName: "วิชัย ตั้งใจ",
    birthDate: "1985-02-14",
    annualIncome: 142000,
    currentAddress: "45 ถ.มิตรภาพ ต.ในเมือง อ.เมือง จ.ขอนแก่น 40000",
    status: "rejected",
    reason: "รายได้รวมของครัวเรือนเกินเกณฑ์ที่กำหนด (ไม่เกิน 100,000 บาท/ปี)",
    submittedAt: "2026-08-18T11:00:00Z",
  },
];

export const SAMPLE_NEW_VALID_IDS = [
  {
    id: "3100600929594",
    name: "มาลี สุขใจ",
    income: 58000,
    desc: "เลขใหม่ (รายได้ผ่านเกณฑ์)",
    type: "new",
  },
  {
    id: "1100400123450",
    name: "สมชาย ใจดี",
    income: 54000,
    desc: "เลขซ้ำในระบบ (จะแจ้งเตือน)",
    type: "duplicate",
  },
  {
    id: "1103700286802",
    name: "ประเสริฐ คงคำ",
    income: 142000,
    desc: "เลขใหม่ (รายได้เกินเกณฑ์)",
    type: "high_income",
  },
];

export const SERVER_ERROR_TRIGGER_ID = "9999999999999";
