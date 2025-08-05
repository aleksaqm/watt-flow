export interface Household {
  id: number
  cadastral_number: string
  city?: string
  street?: string
  number?: string
  floor: string
  suite: string
  status?: string
}

export interface HouseholdFull {
  id: number
  cadastral_number: string
  city?: string
  street?: string
  number?: string
  floor: string
  status: string
  suite: string
  device_address: string
  owner_name: string
  sq_footage: number
  latitude: number
  longitude: number
}
