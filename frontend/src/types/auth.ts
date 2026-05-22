export type Role = 'admin' | 'creator' | 'student' | 'organization'

export type LoginRequest = {
  login: string
  password: string
}

export type SessionResponse = {
  token?: string
  token_type: 'Bearer'
  expires_at: number
  account: {
    id: string
  }
  person: {
    id: string
  }
  role: Role
}
