import { useQuery } from '@tanstack/react-query'
import { listAccounts } from '../../api/accounts'
import type { AccountListFilter } from '../../types/account'

export function useAccounts(filter: AccountListFilter) {
  return useQuery({
    queryKey: ['accounts', filter],
    queryFn: () => listAccounts(filter),
  })
}
