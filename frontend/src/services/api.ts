import type { GreetingResponse } from '../types/api'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

export async function fetchGreeting(name?: string): Promise<GreetingResponse> {
  const queryParams = new URLSearchParams()
  if (name) {
    queryParams.set('name', name)
  }

  const response = await fetch(`${API_BASE_URL}/api/v1/hello?${queryParams.toString()}`)

  if (!response.ok) {
    throw new Error(`Failed to fetch greeting: ${response.statusText}`)
  }

  return response.json()
}
