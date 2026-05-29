import type { AnswerPayload, AttemptSummaryView, AttemptView, StartAttemptRequest } from '../types/attempt'
import { request } from './client'

export function listAttemptsByEnrollment(enrollmentId: string, params?: { limit?: number; offset?: number }) {
  const qs = new URLSearchParams({ enrollment_id: enrollmentId })
  if (params?.limit != null) qs.set('limit', String(params.limit))
  if (params?.offset != null) qs.set('offset', String(params.offset))
  return request<AttemptSummaryView[]>(`/attempts?${qs}`)
}

export function getAttempt(id: string) {
  return request<AttemptView>(`/attempts/${id}`)
}

export function startAttempt(payload: StartAttemptRequest) {
  return request<{ id: string }>('/attempts', { method: 'POST', body: payload })
}

export function finishAttempt(id: string) {
  return request<{ id: string }>(`/attempts/${id}/finish`, { method: 'POST' })
}

export function cancelAttempt(id: string) {
  return request<{ id: string }>(`/attempts/${id}/cancel`, { method: 'POST' })
}

export function addAttemptAnswer(attemptId: string, questionId: string, payload: AnswerPayload) {
  return request<{ id: string }>(`/attempts/${attemptId}/answers/${questionId}`, {
    method: 'PUT',
    body: payload,
  })
}
