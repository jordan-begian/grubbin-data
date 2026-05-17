export interface PickupRequest {
  name: string
  lat: number
  lon: number
}

export interface DropoffRequest {
  lat: number
  lon: number
}

export interface EarningsRequest {
  tip: number
  base: number
  bonus?: number
}

export interface CreateDeliveryRequest {
  start: string
  end: string
  pickup: PickupRequest
  dropoff: DropoffRequest
  earnings: EarningsRequest
  note?: string
}

export interface PickupResponse {
  id: string
  name: string
  lat: number
  lon: number
}

export interface DropoffResponse {
  id: string
  lat: number
  lon: number
}

export interface EarningsResponse {
  id: string
  tip: number
  base: number
  bonus?: number
}

export interface DeliveryResponse {
  id: string
  user: string
  created: string
  updated?: string
  start: string
  end: string
  pickup: PickupResponse
  dropoff: DropoffResponse
  earnings: EarningsResponse
  note?: string
}

export interface DeliveryStats {
  total_time: number
  total_miles: number
  fuel_used?: number
  used_fuel_cost?: number
  average_delivery_time: number
  average_delivery_distance: number
  average_tip: number
  average_base_pay: number
}

export interface DeliveryListResponse {
  deliveries: DeliveryResponse[]
  stats: DeliveryStats
}

export interface UpdateDeliveryRequest {
  id: string
  start?: string
  end?: string
  pickup?: PickupRequest
  dropoff?: DropoffRequest
  earnings?: EarningsRequest
  note?: string | null
}
