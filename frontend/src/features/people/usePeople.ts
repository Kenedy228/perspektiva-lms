import { useQuery } from '@tanstack/react-query'
import { listPersonsByLastName, listPersonsBySnils } from '../../api/people'

export type PeopleSearchMode = 'last_name' | 'snils'

export function usePeople(
  search: string,
  mode: PeopleSearchMode,
  limit: number,
  offset: number,
) {
  const trimmed = search.trim()

  return useQuery({
    queryKey: ['people', mode, trimmed, limit, offset],
    queryFn: () =>
      mode === 'snils'
        ? listPersonsBySnils(trimmed, limit, offset)
        : listPersonsByLastName(trimmed, limit, offset),
    // last_name mode: always enabled (empty string → all persons)
    // snils mode: only when something is typed
    enabled: mode === 'last_name' || trimmed.length > 0,
  })
}
