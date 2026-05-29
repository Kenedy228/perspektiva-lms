import type {
  CreateEnrollmentRequest,
  EnrollmentListFilter,
  EnrollmentView,
} from '../types/enrollment'
import { request } from './client'

export function listEnrollments(filter: EnrollmentListFilter = {}) {
  const params = new URLSearchParams()
  if (filter.account_id) params.set('account_id', filter.account_id)
  if (filter.course_id) params.set('course_id', filter.course_id)
  if (filter.limit != null) params.set('limit', String(filter.limit))
  if (filter.offset != null) params.set('offset', String(filter.offset))
  const qs = params.toString()
  return request<EnrollmentView[]>(`/enrollments${qs ? `?${qs}` : ''}`)
}

export function getEnrollment(id: string) {
  return request<EnrollmentView>(`/enrollments/${id}`)
}

export function createEnrollment(payload: CreateEnrollmentRequest) {
  return request<{ id: string }>('/enrollments', { method: 'POST', body: payload })
}
