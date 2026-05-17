import type { LoginRequest, RegisterUserRequest, UserResponse } from '@/types/auth'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

export async function loginUser(credentials: LoginRequest): Promise<UserResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(credentials),
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.error || `Login failed: ${response.statusText}`)
  }

  return response.json()
}

export async function registerUser(data: RegisterUserRequest): Promise<UserResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(errorData.error || `Registration failed: ${response.statusText}`)
  }

  return response.json()
}
