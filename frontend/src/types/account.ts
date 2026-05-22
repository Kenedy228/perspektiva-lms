import { type Role } from './auth'

export type AccountStatus = 'active' | 'blocked' | 'deleted'

export type AccountView = {
  ID: string
  Login: string
  Role: Role
  Status: AccountStatus
  PersonID: string
}

export type CreateAccountRequest = {
  person_id: string
  login: string
  password: string
  role: Role
}
