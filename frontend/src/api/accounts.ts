import type { AccountListFilter, AccountView, CreateAccountRequest } from '../types/account'
import type { Role } from '../types/auth'
import { request } from './client'

export function listAccounts(filter: AccountListFilter = {}) {
  const params = new URLSearchParams()
  if (filter.login) params.set('login', filter.login)
  if (filter.role) params.set('role', filter.role)
  if (filter.status) params.set('status', filter.status)
  if (filter.limit != null) params.set('limit', String(filter.limit))
  if (filter.offset != null) params.set('offset', String(filter.offset))
  const qs = params.toString()
  return request<AccountView[]>(`/accounts${qs ? `?${qs}` : ''}`)
}

export function getAccount(id: string) {
  return request<AccountView>(`/accounts/${id}`)
}

export function createAccount(payload: CreateAccountRequest) {
  return request<{ id: string }>('/accounts', { method: 'POST', body: payload })
}

export function changeAccountLogin(id: string, login: string) {
  return request<{ id: string }>(`/accounts/${id}/login`, { method: 'PATCH', body: { login } })
}

export function changeAccountPassword(id: string, password: string) {
  return request<{ id: string }>(`/accounts/${id}/password`, { method: 'PATCH', body: { password } })
}

export function changeAccountRole(id: string, role: Role) {
  return request<{ id: string }>(`/accounts/${id}/role`, { method: 'PATCH', body: { role } })
}

export function blockAccount(id: string) {
  return request<{ id: string }>(`/accounts/${id}/block`, { method: 'POST' })
}

export function activateAccount(id: string) {
  return request<{ id: string }>(`/accounts/${id}/activate`, { method: 'POST' })
}

export function deleteAccount(id: string) {
  return request<void>(`/accounts/${id}`, { method: 'DELETE' })
}
