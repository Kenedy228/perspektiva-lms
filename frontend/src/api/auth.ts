import { type LoginRequest, type SessionResponse } from '../types/auth'
import { request } from './client'
import { clearSessionToken, setSessionToken } from './tokenStore'

export async function login(payload: LoginRequest) {
  const session = await request<SessionResponse>('/auth/login', {
    method: 'POST',
    body: payload,
    auth: false,
  })

  if (session.token) {
    setSessionToken(session.token)
  }

  return session
}

export async function getSession() {
  return request<SessionResponse>('/auth/session')
}

export async function logout() {
  await request<{ status: string }>('/auth/logout', { method: 'POST' })
  clearSessionToken()
}
