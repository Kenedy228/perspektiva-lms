export type EnrollmentStatus = 'active' | 'inactive' | 'expired'

export const ENROLLMENT_STATUS_LABELS: Record<EnrollmentStatus, string> = {
  active: 'Активно',
  inactive: 'Неактивно',
  expired: 'Истекло',
}

// Fields are snake_case because EnrollmentView in Go has explicit json tags.
export type EnrollmentView = {
  id: string
  account_id: string
  course_id: string
  activated_at: string
  deactivated_at: string
  status: string
  status_title: string
}

export type EnrollmentListFilter = {
  account_id?: string
  course_id?: string
  limit?: number
  offset?: number
}

export type CreateEnrollmentRequest = {
  account_id: string
  course_id: string
  activated_at?: string
  deactivated_at?: string
}
