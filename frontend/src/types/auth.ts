export interface Vehicle {
  id: string
  name: string
  average_mpg: number
}

export interface Profile {
  id: string
  first_name: string
  last_name: string
  vehicle?: Vehicle
}

export interface UserResponse {
  id: string
  username: string
  created: string
  updated?: string
  profile?: Profile
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterUserRequest {
  username: string
  password: string
  first_name: string
  last_name: string
  vehicle_name?: string
  vehicle_mpg?: number
}
