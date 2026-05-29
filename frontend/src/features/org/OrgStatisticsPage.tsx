import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { listStudentStatistics } from '../../api/statistics'
import { PageHeader } from '../../components/ui/PageHeader'
import { Pagination } from '../../components/ui/Pagination'
import styles from './OrgStatisticsPage.module.css'

const LIMIT = 20

function shortId(id: string) {
  return id.length > 8 ? `${id.slice(0, 8)}…` : id
}

function ProgressBar({ percent }: { percent: number }) {
  const clamped = Math.min(100, Math.max(0, percent))
  return (
    <div className={styles.progressWrap}>
      <div className={styles.progressTrack}>
        <div className={styles.progressFill} style={{ width: `${clamped}%` }} />
      </div>
      <span className={styles.progressText}>{percent}%</span>
    </div>
  )
}

export function OrgStatisticsPage() {
  const [offset, setOffset] = useState(0)

  const query = useQuery({
    queryKey: ['org-statistics', offset],
    queryFn: () => listStudentStatistics({ limit: LIMIT, offset }),
  })

  const rows = query.data ?? []
  const hasMore = rows.length >= LIMIT

  return (
    <>
      <PageHeader
        title="Статистика студентов"
        description="Прогресс прохождения курсов студентами вашей организации."
      />

      {query.isPending && <p className={styles.state}>Загрузка…</p>}
      {query.isError && <p className={styles.stateError}>Не удалось загрузить статистику.</p>}

      {!query.isPending && !query.isError && rows.length === 0 && (
        <p className={styles.empty}>Данные не найдены. Убедитесь, что ваш профиль привязан к организации.</p>
      )}

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
                        {shortId(row.AccountID)}
                      </code>
                    </td>
                    <td>
                      <code className={styles.uuid} title={row.CourseID}>
                        {shortId(row.CourseID)}
                      </code>
                    </td>
                    <td>
                      <code className={styles.uuid} title={row.EnrollmentID}>
                        {shortId(row.EnrollmentID)}
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

          {hasMore || offset > 0 ? (
            <Pagination
              offset={offset}
              limit={LIMIT}
              hasMore={hasMore}
              onOffsetChange={setOffset}
            />
          ) : null}
        </>
      )}
    </>
  )
}
