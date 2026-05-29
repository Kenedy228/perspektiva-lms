import { ArrowLeft } from 'lucide-react'
import { Link, useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { listAttemptsByEnrollment } from '../../api/attempts'
import { Badge } from '../../components/ui/Badge'
import { PageHeader } from '../../components/ui/PageHeader'
import { Pagination } from '../../components/ui/Pagination'
import { ATTEMPT_STATUS_LABELS } from '../../types/attempt'
import type { AttemptStatus } from '../../types/attempt'
import { useState } from 'react'
import styles from './EnrollmentAttemptsPage.module.css'

const LIMIT = 20

function statusVariant(status: AttemptStatus): 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'finished') return 'success'
  if (status === 'in_progress') return 'warning'
  if (status === 'expired' || status === 'cancelled' || status === 'interrupted') return 'danger'
  return 'neutral'
}

function formatDate(iso?: string) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString('ru-RU', { dateStyle: 'short', timeStyle: 'short' })
}

function shortId(id: string) {
  return id.length > 8 ? `${id.slice(0, 8)}…` : id
}

export function EnrollmentAttemptsPage() {
  const { id: enrollmentId = '' } = useParams<{ id: string }>()
  const [offset, setOffset] = useState(0)

  const query = useQuery({
    queryKey: ['enrollment-attempts', enrollmentId, offset],
    queryFn: () => listAttemptsByEnrollment(enrollmentId, { limit: LIMIT, offset }),
    enabled: Boolean(enrollmentId),
  })

  const rows = query.data ?? []
  const hasMore = rows.length >= LIMIT

  return (
    <>
      <div className={styles.backRow}>
        <Link to="/enrollments" className={styles.backLink}>
          <ArrowLeft size={15} aria-hidden="true" />
          Зачисления
        </Link>
      </div>

      <PageHeader
        title="Попытки по зачислению"
        description={`Зачисление: ${shortId(enrollmentId)}`}
      />

      {query.isPending && <p className={styles.state}>Загрузка…</p>}
      {query.isError && <p className={styles.stateError}>Не удалось загрузить попытки.</p>}

      {!query.isPending && !query.isError && rows.length === 0 && (
        <p className={styles.empty}>По этому зачислению попыток ещё нет.</p>
      )}

      {rows.length > 0 && (
        <>
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ID попытки</th>
                  <th>Статус</th>
                  <th>Начата</th>
                  <th>Завершена</th>
                  <th>Ответов</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => (
                  <tr key={row.id}>
                    <td>
                      <code className={styles.uuid} title={row.id}>
                        {shortId(row.id)}
                      </code>
                    </td>
                    <td>
                      <Badge variant={statusVariant(row.status)}>
                        {ATTEMPT_STATUS_LABELS[row.status] ?? row.status}
                      </Badge>
                    </td>
                    <td>{formatDate(row.started_at)}</td>
                    <td>{formatDate(row.finished_at)}</td>
                    <td>{row.answers_count} / {row.questions_count}</td>
                    <td>
                      <Link to={`/attempts/${row.id}`} className={styles.viewLink}>
                        Просмотр
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {hasMore || offset > 0 ? (
            <Pagination offset={offset} limit={LIMIT} hasMore={hasMore} onOffsetChange={setOffset} />
          ) : null}
        </>
      )}
    </>
  )
}
