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
