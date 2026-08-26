export function isValidThaiCitizenId(value: string): boolean {
  const digits = value.replace(/\D/g, '')
  if (digits.length !== 13) return false

  let sum = 0
  for (let index = 0; index < 12; index += 1) {
    sum += Number(digits[index]) * (13 - index)
  }

  return (11 - (sum % 11)) % 10 === Number(digits[12])
}