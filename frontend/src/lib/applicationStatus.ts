export type ApplicationStatusValue = 'pending' | 'approved' | 'rejected'

const labels: Record<ApplicationStatusValue, string> = {
  pending: 'รอตรวจสอบ',
  approved: 'อนุมัติแล้ว',
  rejected: 'ไม่อนุมัติ',
}

const classes: Record<ApplicationStatusValue, string> = {
  pending: 'border-amber-200 bg-amber-50 text-amber-700',
  approved: 'border-green-200 bg-green-50 text-green-700',
  rejected: 'border-red-200 bg-red-50 text-red-700',
}

export function statusLabel(value: ApplicationStatusValue) {
  return labels[value]
}

export function statusClass(value: ApplicationStatusValue) {
  return classes[value]
}