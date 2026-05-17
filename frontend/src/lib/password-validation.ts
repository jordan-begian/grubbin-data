/**
 * Password validation that mirrors backend/internal/core/password.go
 * Pure function - same input always produces same output.
 */

export interface PasswordValidationResult {
  isValid: boolean
  requirements: PasswordRequirement[]
}

export interface PasswordRequirement {
  label: string
  met: boolean
}

const commonPasswords = [
  "password", "123456", "qwerty", "admin", "letmein",
  "welcome", "monkey", "dragon", "master", "sunshine",
]

const commonPasswordPatterns = [
  /^password/i,
  /^123456/i,
  /^qwerty/i,
  /^admin/i,
  /^letmein/i,
]

const sequentialPattern = /(012|123|234|345|456|567|678|789|890|abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)/i

function isCommonPassword(password: string): boolean {
  for (const common of commonPasswords) {
    const regex = new RegExp(common, "i")
    if (regex.test(password)) {
      return true
    }
  }

  for (const pattern of commonPasswordPatterns) {
    if (pattern.test(password)) {
      return true
    }
  }

  if (sequentialPattern.test(password)) {
    return true
  }

  return false
}

export function validatePassword(password: string): PasswordValidationResult {
  const requirements: PasswordRequirement[] = []

  // Length check
  requirements.push({
    label: "At least 12 characters",
    met: password.length >= 12,
  })

  // Character type checks
  let hasUpper = false
  let hasLower = false
  let hasDigit = false
  let hasSpecial = false

  for (const char of password) {
    if (/[A-Z]/.test(char)) hasUpper = true
    else if (/[a-z]/.test(char)) hasLower = true
    else if (/\d/.test(char)) hasDigit = true
    else if (/[^\w\s]/.test(char)) hasSpecial = true
  }

  requirements.push({
    label: "At least one uppercase letter",
    met: hasUpper,
  })

  requirements.push({
    label: "At least one lowercase letter",
    met: hasLower,
  })

  requirements.push({
    label: "At least one digit",
    met: hasDigit,
  })

  requirements.push({
    label: "At least one special character",
    met: hasSpecial,
  })

  // Common password check
  requirements.push({
    label: "Not a common or weak password",
    met: password.length > 0 && !isCommonPassword(password),
  })

  const isValid = requirements.every((r) => r.met)

  return { isValid, requirements }
}
