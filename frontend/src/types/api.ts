export type Link = {
  href: string
  method?: string
}

export type Links = Record<string, Link>

export type ApiEnvelope<T> = {
  data?: T
  _links?: Links
}

export type ApiErrorPayload = {
  error: {
    status: number
    code: string
    message: string
    details?: Record<string, unknown>
    _links?: Links
  }
}

export class ApiError extends Error {
  readonly status: number
  readonly code: string
  readonly details?: Record<string, unknown>

  constructor(status: number, code: string, message: string, details?: Record<string, unknown>) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.details = details
  }
}
