import { request } from './client'
import type { CreateQuizPayload, QuizView } from '../types/quizzes'

export function getQuiz(id: string) {
  return request<QuizView>(`/quizzes/${id}`)
}

export function createQuiz(payload: CreateQuizPayload) {
  return request<{ id: string }>('/quizzes', { method: 'POST', body: payload })
}

export function renameQuiz(id: string, title: string) {
  return request<{ id: string }>(`/quizzes/${id}`, { method: 'PATCH', body: { title } })
}

export function changeQuizLimits(id: string, max_attempts: number, time_limit_seconds: number) {
  return request<{ id: string }>(`/quizzes/${id}/limits`, {
    method: 'PATCH',
    body: { max_attempts, time_limit_seconds },
  })
}

export function changeQuizShuffle(id: string, shuffle_questions: boolean) {
  return request<{ id: string }>(`/quizzes/${id}/shuffle`, {
    method: 'PATCH',
    body: { shuffle_questions },
  })
}

export function replaceQuizSources(id: string, sources: CreateQuizPayload['sources']) {
  return request<{ id: string }>(`/quizzes/${id}/sources`, {
    method: 'PUT',
    body: { sources },
  })
}

export function deleteQuiz(id: string) {
  return request<void>(`/quizzes/${id}`, { method: 'DELETE' })
}
