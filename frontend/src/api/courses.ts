import type {
  CompletionMode,
  CourseDetailedView,
  CourseShortView,
  CreateCourseRequest,
  ElementDownloadInfo,
  ElementType,
  StudentRatingView,
} from '../types/courses'
import { request, upload } from './client'

export function getCourses(params?: { title?: string; limit?: number; offset?: number }) {
  const qs = params
    ? new URLSearchParams(
        Object.entries(params)
          .filter(([, v]) => v != null && v !== '')
          .map(([k, v]) => [k, String(v)]),
      ).toString()
    : ''
  return request<CourseShortView[]>(`/courses${qs ? `?${qs}` : ''}`)
}

export function getCourse(id: string) {
  return request<CourseDetailedView>(`/courses/${id}`)
}

export function createCourse(payload: CreateCourseRequest) {
  return request<{ id: string }>('/courses', { method: 'POST', body: payload })
}

export function renameCourse(courseId: string, title: string) {
  return request<{ id: string }>(`/courses/${courseId}`, { method: 'PATCH', body: { title } })
}

export function addCourseBlock(courseId: string, title: string) {
  return request<{ id: string }>(`/courses/${courseId}/blocks`, { method: 'POST', body: { title } })
}

export function removeCourseBlock(courseId: string, blockId: string) {
  return request<void>(`/courses/${courseId}/blocks/${blockId}`, { method: 'DELETE' })
}

export function moveCourseBlock(courseId: string, from: number, to: number) {
  return request<{ course_id: string }>(`/courses/${courseId}/blocks/move`, {
    method: 'PATCH',
    body: { from, to },
  })
}

export type CreateBlockElementPayload = {
  title: string
  type: ElementType
  file_name?: string
  size_bytes?: number
  quiz_id?: string
  completion_mode?: CompletionMode
}

export function addBlockElement(blockId: string, payload: CreateBlockElementPayload) {
  return request<{ id: string }>(`/blocks/${blockId}/elements`, { method: 'POST', body: payload })
}

export function removeBlockElement(blockId: string, elementId: string) {
  return request<void>(`/blocks/${blockId}/elements/${elementId}`, { method: 'DELETE' })
}

export function moveBlockElement(blockId: string, from: number, to: number) {
  return request<{ block_id: string }>(`/blocks/${blockId}/elements/move`, {
    method: 'PATCH',
    body: { from, to },
  })
}

export function changeElementCompletionMode(elementId: string, completionMode: CompletionMode) {
  return request<{ element_id: string }>(`/elements/${elementId}/completion-mode`, {
    method: 'PATCH',
    body: { completion_mode: completionMode },
  })
}

export function uploadElementContent(elementId: string, file: File) {
  return upload<{ element_id: string }>(`/elements/${elementId}/content`, file)
}

export function getElementDownloadURL(elementId: string) {
  return request<ElementDownloadInfo>(`/elements/${elementId}/download`)
}

export function listCourseRatings(courseId: string, limit: number, offset: number) {
  return request<StudentRatingView[]>(
    `/courses/${courseId}/ratings?limit=${limit}&offset=${offset}`,
  )
}
