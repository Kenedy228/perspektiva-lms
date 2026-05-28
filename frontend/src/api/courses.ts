import { type CourseDetailedView, type CourseShortView, type CreateCourseRequest } from '../types/courses'
import { request } from './client'

export function getCourses() {
  return request<CourseShortView[]>('/courses')
}

export function getCourse(id: string) {
  return request<CourseDetailedView>(`/courses/${id}`)
}

export function createCourse(payload: CreateCourseRequest) {
  return request<{ id: string }>('/courses', {
    method: 'POST',
    body: payload,
  })
}

export function renameCourse(courseId: string, payload: CreateCourseRequest) {
  return request<{ id: string }>(`/courses/${courseId}`, {
    method: 'PATCH',
    body: payload,
  })
}

export function addCourseBlock(courseId: string, title: string) {
  return request<{ id: string }>(`/courses/${courseId}/blocks`, {
    method: 'POST',
    body: { title },
  })
}

export function removeCourseBlock(courseId: string, blockId: string) {
  return request<void>(`/courses/${courseId}/blocks/${blockId}`, {
    method: 'DELETE',
  })
}

export type CreateBlockElementPayload = {
  title: string
  type: 'test' | 'lecture_material' | 'download_file'
  file_name?: string
  size_bytes?: number
  quiz_id?: string
  completion_mode?: 'none' | 'manual'
}

export function addBlockElement(blockId: string, payload: CreateBlockElementPayload) {
  return request<{ id: string }>(`/blocks/${blockId}/elements`, {
    method: 'POST',
    body: payload,
  })
}

export function removeBlockElement(blockId: string, elementId: string) {
  return request<void>(`/blocks/${blockId}/elements/${elementId}`, {
    method: 'DELETE',
  })
}

export function changeElementCompletionMode(elementId: string, completionMode: 'none' | 'manual') {
  return request<{ element_id: string }>(`/elements/${elementId}/completion-mode`, {
    method: 'PATCH',
    body: { completion_mode: completionMode },
  })
}
