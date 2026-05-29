import { useQuery } from '@tanstack/react-query'
import { listEnrollments } from '../../api/enrollments'
import type { EnrollmentListFilter } from '../../types/enrollment'

export function useEnrollments(filter: EnrollmentListFilter) {
  return useQuery({
    queryKey: ['enrollments', filter],
    queryFn: () => listEnrollments(filter),
  })
}
