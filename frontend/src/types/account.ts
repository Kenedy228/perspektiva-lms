import { type Role } from './auth'

export type AccountStatus = 'active' | 'blocked' | 'deleted'

export const ROLE_LABELS: Record<Role, string> = {
  admin: 'Администратор',
  creator: 'Создатель',
  student: 'Студент',
  organization: 'Организация',
}

export const STATUS_LABELS: Record<AccountStatus, string> = {
  active: 'Активен',
  blocked: 'Заблокирован',
  deleted: 'Удалён',
}

export const ACCOUNT_ROLES: Role[] = ['admin', 'creator', 'student', 'organization']

export type AccountView = {
  ID: string
  Login: string
  Role: Role
  Status: AccountStatus
  PersonID: string
}

export type AccountListFilter = {
  login?: string
  role?: Role | ''
  status?: AccountStatus | ''
  limit?: number
  offset?: number
}

export type CreateAccountRequest = {
  person_id: string
  login: string
  password: string
  role: Role
}
