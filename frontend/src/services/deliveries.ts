import type {
  CreateDeliveryRequest,
  DeliveryResponse,
  DeliveryListResponse,
  UpdateDeliveryRequest,
} from '@/types/delivery'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

export async function createDelivery(
  userId: string,
  data: CreateDeliveryRequest
): Promise<DeliveryResponse> {
  await delay(600)

  const response = await fetch(
    `${API_BASE_URL}/api/v1/users/${userId}/deliveries`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    }
  )

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(
      errorData.error || `Failed to create delivery: ${response.statusText}`
    )
  }

  return response.json()
}

export async function getDeliveries(
  userId: string
): Promise<DeliveryListResponse> {
  await delay(400)

  const response = await fetch(
    `${API_BASE_URL}/api/v1/users/${userId}/deliveries`
  )

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(
      errorData.error || `Failed to fetch deliveries: ${response.statusText}`
    )
  }

  return response.json()
}

export async function updateDeliveries(
  userId: string,
  updates: UpdateDeliveryRequest[]
): Promise<void> {
  await delay(600)

  const response = await fetch(
    `${API_BASE_URL}/api/v1/users/${userId}/deliveries`,
    {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(updates),
    }
  )

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(
      errorData.error || `Failed to update deliveries: ${response.statusText}`
    )
  }
}

export async function deleteDeliveries(
  userId: string,
  deliveryIds: string[]
): Promise<void> {
  await delay(600)

  const params = new URLSearchParams()
  deliveryIds.forEach((id) => params.append('id', id))

  const response = await fetch(
    `${API_BASE_URL}/api/v1/users/${userId}/deliveries?${params.toString()}`,
    {
      method: 'DELETE',
    }
  )

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new Error(
      errorData.error || `Failed to delete deliveries: ${response.statusText}`
    )
  }
}
