import { useQuery } from '@tanstack/react-query'
import { listOrganizationsByINN, listOrganizationsByName } from '../../api/organizations'

export type SearchMode = 'name' | 'inn'

export function useOrganizations(
  search: string,
  mode: SearchMode,
  limit: number,
  offset: number,
) {
  const trimmed = search.trim()

  return useQuery({
    queryKey: ['organizations', mode, trimmed, limit, offset],
    queryFn: () =>
      mode === 'inn'
        ? listOrganizationsByINN(trimmed, limit, offset)
        : listOrganizationsByName(trimmed, limit, offset),
    enabled: trimmed.length > 0,
  })
}
