import { ApiError, type ApiEnvelope, type ApiErrorPayload } from '../types/api'
import { getSessionToken } from './tokenStore'

const baseURL = import.meta.env.VITE_API_BASE_URL ?? '/api'

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  body?: unknown
  auth?: boolean
}

export async function upload<T>(path: string, file: File): Promise<T> {
  const formData = new FormData()
  formData.append('content', file)

  const headers = new Headers()
  headers.set('Accept', 'application/json')
  const token = getSessionToken()
  if (token) headers.set('Authorization', `Bearer ${token}`)

  const response = await fetch(`${baseURL}${path}`, {
    method: 'PUT',
    headers,
    body: formData,
  })

  const payload = await readJSON(response)
  if (!response.ok) throw toApiError(response.status, payload)

  const envelope = payload as ApiEnvelope<T>
  return envelope.data as T
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const headers = new Headers()
  headers.set('Accept', 'application/json')

  if (options.body !== undefined) {
    headers.set('Content-Type', 'application/json')
  }

  if (options.auth !== false) {
    const token = getSessionToken()
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
  }

  const response = await fetch(`${baseURL}${path}`, {
    method: options.method ?? 'GET',
    headers,
    body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
  })

  const payload = await readJSON(response)

  if (!response.ok) {
    throw toApiError(response.status, payload)
  }

  const envelope = payload as ApiEnvelope<T>
  return envelope.data as T
}

async function readJSON(response: Response): Promise<unknown> {
  if (response.status === 204) {
    return undefined
  }

  const text = await response.text()
  if (!text) {
    return undefined
  }

  return JSON.parse(text) as unknown
}

function toApiError(fallbackStatus: number, payload: unknown): ApiError {
  if (isApiErrorPayload(payload)) {
    return new ApiError(payload.error.status, payload.error.code, payload.error.message, payload.error.details)
  }

  return new ApiError(fallbackStatus, 'http_error', 'HTTP-запрос завершился ошибкой')
}

function isApiErrorPayload(value: unknown): value is ApiErrorPayload {
  if (typeof value !== 'object' || value === null || !('error' in value)) {
    return false
  }

  const error = (value as { error: unknown }).error
  return typeof error === 'object' && error !== null && 'status' in error && 'code' in error && 'message' in error
}
