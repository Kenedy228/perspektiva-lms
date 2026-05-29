import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { listStudentStatistics } from '../../api/statistics'
import { Badge } from '../../components/ui/Badge'
import { Button } from '../../components/ui/Button'
import { PageHeader } from '../../components/ui/PageHeader'
import { Pagination } from '../../components/ui/Pagination'
import styles from './StatisticsPage.module.css'

const LIMIT = 20

function shortUuid(id: string) {
  return id.length > 8 ? `${id.slice(0, 8)}…` : id
}

function ProgressBar({ percent }: { percent: number }) {
  return (
    <div className={styles.progressWrap}>
      <div className={styles.progressTrack}>
        <div className={styles.progressFill} style={{ width: `${Math.min(100, percent)}%` }} />
      </div>
      <span className={styles.progressText}>{percent}%</span>
    </div>
  )
}

export function StatisticsPage() {
  // Draft filter state — applied only on submit/blur
  const [draftAccountId, setDraftAccountId] = useState('')
  const [draftOrgId, setDraftOrgId] = useState('')

  // Applied filter state — triggers the query
  const [accountId, setAccountId] = useState('')
  const [orgId, setOrgId] = useState('')
  const [offset, setOffset] = useState(0)

  function applyFilters() {
    setAccountId(draftAccountId.trim())
    setOrgId(draftOrgId.trim())
    setOffset(0)
  }

  function resetFilters() {
    setDraftAccountId('')
    setDraftOrgId('')
    setAccountId('')
    setOrgId('')
    setOffset(0)
  }

  const statsQuery = useQuery({
    queryKey: ['statistics', accountId, orgId, offset],
    queryFn: () =>
      listStudentStatistics({
        account_id: accountId || undefined,
        organization_id: orgId || undefined,
        limit: LIMIT,
        offset,
      }),
  })

  const rows = statsQuery.data ?? []
  const hasMore = rows.length >= LIMIT

  return (
    <>
      <PageHeader
        title="Статистика студентов"
        description="Прогресс прохождения курсов по всем зачислениям."
      />

      {/* Toolbar */}
      <div className={styles.toolbar}>
        <input
          className={styles.filterInput}
          value={draftAccountId}
          onChange={(e) => setDraftAccountId(e.target.value)}
          onBlur={applyFilters}
          onKeyDown={(e) => { if (e.key === 'Enter') applyFilters() }}
          placeholder="ID аккаунта (Enter для применения)"
          aria-label="Фильтр по ID аккаунта"
        />
        <input
          className={styles.filterInput}
          value={draftOrgId}
          onChange={(e) => setDraftOrgId(e.target.value)}
          onBlur={applyFilters}
          onKeyDown={(e) => { if (e.key === 'Enter') applyFilters() }}
          placeholder="ID организации (Enter)"
          aria-label="Фильтр по ID организации"
        />
        <Button onClick={applyFilters}>Применить</Button>
        {(accountId || orgId) && (
          <Button variant="secondary" onClick={resetFilters}>
            Сбросить
          </Button>
        )}
      </div>

      {/* Active filters */}
      {(accountId || orgId) && (
        <div className={styles.activeFilters}>
          {accountId && (
            <Badge variant="neutral">
              Аккаунт: {shortUuid(accountId)}
            </Badge>
          )}
          {orgId && (
            <Badge variant="neutral">
              Организация: {shortUuid(orgId)}
            </Badge>
          )}
        </div>
      )}

      {/* States */}
      {statsQuery.isPending && <p className={styles.state}>Загрузка…</p>}
      {statsQuery.isError && (
        <p className={styles.stateError}>Не удалось загрузить статистику.</p>
      )}

      {/* Empty state */}
      {!statsQuery.isPending && !statsQuery.isError && rows.length === 0 && (
        <p className={styles.empty}>Данные не найдены.</p>
      )}

      {/* Table */}
      {rows.length > 0 && (
        <>
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Аккаунт</th>
                  <th>Курс</th>
                  <th>Зачисление</th>
                  <th>Прогресс</th>
                  <th>Выполнено</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => (
                  <tr key={row.EnrollmentID}>
                    <td>
                      <code className={styles.uuid} title={row.AccountID}>
                        {shortUuid(row.AccountID)}
                      </code>
                    </td>
                    <td>
                      <code className={styles.uuid} title={row.CourseID}>
                        {shortUuid(row.CourseID)}
                      </code>
                    </td>
                    <td>
                      <code className={styles.uuid} title={row.EnrollmentID}>
                        {shortUuid(row.EnrollmentID)}
                      </code>
                    </td>
                    <td>
                      <ProgressBar percent={row.CompletionPercent} />
                    </td>
                    <td className={styles.completedCell}>
                      {row.CompletedItems} / {row.TotalItems}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <Pagination
            offset={offset}
            limit={LIMIT}
            hasMore={hasMore}
            onOffsetChange={setOffset}
          />
        </>
      )}
    </>
  )
}
