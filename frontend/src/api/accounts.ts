import { type AccountView, type CreateAccountRequest } from '../types/account'
import { request } from './client'

export function getAccounts() {
  return request<AccountView[]>('/accounts')
}

export function createAccount(payload: CreateAccountRequest) {
  return request<{ id: string }>('/accounts', {
    method: 'POST',
    body: payload,
  })
}
