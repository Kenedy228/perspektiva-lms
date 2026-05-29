import type { StudentRatingView } from '../types/courses'
import { request } from './client'

type StudentStatisticsFilter = {
  account_id?: string
  organization_id?: string
  limit?: number
  offset?: number
}

export function listStudentStatistics(filter: StudentStatisticsFilter) {
  const params = new URLSearchParams()
  if (filter.account_id) params.set('account_id', filter.account_id)
  if (filter.organization_id) params.set('organization_id', filter.organization_id)
  if (filter.limit != null) params.set('limit', String(filter.limit))
  if (filter.offset != null) params.set('offset', String(filter.offset))
  const qs = params.toString()
  return request<StudentRatingView[]>(`/statistics/students${qs ? `?${qs}` : ''}`)
}
