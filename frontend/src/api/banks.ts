import { request } from './client'
import type { BankDetailedView, BankShortView } from '../types/banks'

type QuestionRequest = {
  type: 'selectable' | 'sequence' | 'matching' | 'short'
  title: string
  selectable_options?: Array<{ text: string; is_correct: boolean }>
  sequence_options?: Array<{ text: string }>
  matching_pairs?: Array<{ prompt: string; match: string }>
  short_variants?: Array<{ text: string }>
}

export function listBanks(params?: { title?: string; limit?: number; offset?: number }) {
  const query = new URLSearchParams()
  if (params?.title) query.set('title', params.title)
  if (params?.limit) query.set('limit', String(params.limit))
  if (params?.offset) query.set('offset', String(params.offset))
  const suffix = query.toString() ? `?${query.toString()}` : ''
  return request<BankShortView[]>(`/banks${suffix}`)
}

export function getBank(id: string) {
  return request<BankDetailedView>(`/banks/${id}`)
}

export function createBank(title: string) {
  return request<{ id: string }>('/banks', { method: 'POST', body: { title } })
}

export function renameBank(id: string, title: string) {
  return request<{ id: string }>(`/banks/${id}`, { method: 'PATCH', body: { title } })
}

export function deleteBank(id: string) {
  return request<void>(`/banks/${id}`, { method: 'DELETE' })
}

export function createQuestion(payload: QuestionRequest) {
  return request<{ id: string }>('/questions', { method: 'POST', body: payload })
}

export function updateQuestionTitle(id: string, title: string) {
  return request<{ id: string }>(`/questions/${id}`, { method: 'PATCH', body: { title } })
}

export function updateQuestionContent(id: string, payload: QuestionRequest) {
  return request<{ id: string }>(`/questions/${id}/content`, { method: 'PUT', body: payload })
}

export function deleteQuestion(id: string) {
  return request<void>(`/questions/${id}`, { method: 'DELETE' })
}

export function addQuestionsToBank(bankId: string, questionIDs: string[]) {
  return request<{ id: string }>(`/banks/${bankId}/questions`, { method: 'POST', body: { question_ids: questionIDs } })
}

export function removeQuestionsFromBank(bankId: string, questionIDs: string[]) {
  return request<{ id: string }>(`/banks/${bankId}/questions`, {
    method: 'DELETE',
    body: { question_ids: questionIDs },
  })
}
