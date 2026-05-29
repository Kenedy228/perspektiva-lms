import { useQuery } from '@tanstack/react-query'
import { getSession } from '../../api/auth'
import { ApiError } from '../../types/api'

export function useSession() {
  const query = useQuery({
    queryKey: ['session'],
    queryFn: getSession,
    staleTime: 5 * 60 * 1000,
    retry: (_failureCount, error) => {
      if (error instanceof ApiError && error.status === 401) return false
      return false
    },
  })

  const isUnauthenticated =
    query.isError && query.error instanceof ApiError && query.error.status === 401

  const role = query.data?.role ?? null

  return {
    session: query.data ?? null,
    isLoading: query.isPending,
    isAuthenticated: query.isSuccess,
    isUnauthenticated,
    role,
    accountId: query.data?.account.id ?? null,
    isAdmin: role === 'admin',
    isCreator: role === 'creator',
    isStudent: role === 'student',
    isOrganization: role === 'organization',
    isManager: role === 'admin' || role === 'creator',
  }
}
