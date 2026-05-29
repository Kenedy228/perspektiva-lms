import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { BookMarked } from 'lucide-react'
import { listEnrollments } from '../../api/enrollments'
import { Badge } from '../../components/ui/Badge'
import { PageHeader } from '../../components/ui/PageHeader'
import { Pagination } from '../../components/ui/Pagination'
import type { EnrollmentStatus } from '../../types/enrollment'
import { ENROLLMENT_STATUS_LABELS } from '../../types/enrollment'
import styles from './MyEnrollmentsPage.module.css'

const LIMIT = 20

function statusBadgeVariant(status: string): 'success' | 'neutral' | 'warning' {
  if (status === 'active') return 'success'
  if (status === 'expired') return 'warning'
  return 'neutral'
}

function formatDate(dateStr: string) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('ru-RU')
}

function shortId(id: string) {
  return id.length > 8 ? `${id.slice(0, 8)}…` : id
}

export function MyEnrollmentsPage() {
  const [offset, setOffset] = useState(0)

  const enrollmentsQuery = useQuery({
    queryKey: ['my-enrollments', offset],
    queryFn: () => listEnrollments({ limit: LIMIT, offset }),
  })

  const enrollments = enrollmentsQuery.data ?? []
  const hasMore = enrollments.length >= LIMIT

  return (
    <>
      <PageHeader title="Мои зачисления" description="Список ваших зачислений на курсы." />

      {enrollmentsQuery.isLoading ? (
        <p className={styles.state}>Загрузка…</p>
      ) : enrollmentsQuery.isError ? (
        <p className={styles.stateError}>Не удалось загрузить зачисления.</p>
      ) : enrollments.length === 0 ? (
        <div className={styles.empty}>
          <BookMarked size={40} className={styles.emptyIcon} aria-hidden="true" />
          <p>У вас пока нет зачислений.</p>
        </div>
      ) : (
        <>
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID зачисления</th>
                  <th>ID курса</th>
                  <th>Активно с</th>
                  <th>Активно по</th>
                  <th>Статус</th>
                </tr>
              </thead>
              <tbody>
                {enrollments.map((e) => (
                  <tr key={e.id}>
                    <td className={styles.idCell}>{shortId(e.id)}</td>
                    <td className={styles.idCell}>{shortId(e.course_id)}</td>
                    <td className={styles.dateCell}>{formatDate(e.activated_at)}</td>
                    <td className={styles.dateCell}>{formatDate(e.deactivated_at)}</td>
                    <td>
                      <Badge variant={statusBadgeVariant(e.status)}>
                        {ENROLLMENT_STATUS_LABELS[e.status as EnrollmentStatus] ?? e.status}
                      </Badge>
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
