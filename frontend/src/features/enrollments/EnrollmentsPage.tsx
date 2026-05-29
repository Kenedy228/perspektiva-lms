import { useMutation, useQueryClient } from '@tanstack/react-query'
import { BookMarked, Plus } from 'lucide-react'
import { FormEvent, useState } from 'react'
import { Link } from 'react-router-dom'
import { createEnrollment } from '../../api/enrollments'
import { Badge } from '../../components/ui/Badge'
import { Button } from '../../components/ui/Button'
import { FormField } from '../../components/ui/FormField'
import { Modal } from '../../components/ui/Modal'
import { PageHeader } from '../../components/ui/PageHeader'
import { Pagination } from '../../components/ui/Pagination'
import { ApiError } from '../../types/api'
import type { EnrollmentStatus, EnrollmentView } from '../../types/enrollment'
import { ENROLLMENT_STATUS_LABELS } from '../../types/enrollment'
import { useEnrollments } from './useEnrollments'
import styles from './EnrollmentsPage.module.css'

const LIMIT = 20

type CreateForm = {
  accountId: string
  courseId: string
  activatedAt: string
  deactivatedAt: string
}

function defaultCreate(): CreateForm {
  return { accountId: '', courseId: '', activatedAt: '', deactivatedAt: '' }
}

function statusBadgeVariant(status: string): 'success' | 'neutral' | 'warning' {
  if (status === 'active') return 'success'
  if (status === 'expired') return 'warning'
  return 'neutral'
}

function shortId(id: string) {
  return id.length > 8 ? `${id.slice(0, 8)}…` : id
}

export function EnrollmentsPage() {
  const queryClient = useQueryClient()

  const [accountIdFilter, setAccountIdFilter] = useState('')
  const [courseIdFilter, setCourseIdFilter] = useState('')
  const [offset, setOffset] = useState(0)

  const filter = {
    account_id: accountIdFilter.trim() || undefined,
    course_id: courseIdFilter.trim() || undefined,
    limit: LIMIT,
    offset,
  }

  const { data: enrollments = [], isLoading, isFetching, error } = useEnrollments(filter)

  const [createOpen, setCreateOpen] = useState(false)
  const [createForm, setCreateForm] = useState<CreateForm>(defaultCreate)
  const [createError, setCreateError] = useState<string | null>(null)

  function invalidate() {
    void queryClient.invalidateQueries({ queryKey: ['enrollments'] })
  }

  const createMutation = useMutation({
    mutationFn: createEnrollment,
    onSuccess: () => {
      invalidate()
      setCreateOpen(false)
      setCreateForm(defaultCreate())
      setCreateError(null)
    },
  })

  async function handleCreate(event: FormEvent) {
    event.preventDefault()
    setCreateError(null)
    const payload = {
      account_id: createForm.accountId.trim(),
      course_id: createForm.courseId.trim(),
      activated_at: createForm.activatedAt || undefined,
      deactivated_at: createForm.deactivatedAt || undefined,
    }
    try {
      await createMutation.mutateAsync(payload)
    } catch (caught) {
      setCreateError(caught instanceof ApiError ? caught.message : 'Не удалось создать зачисление')
    }
  }

  function applyFilter() {
    setOffset(0)
  }

  const hasMore = enrollments.length >= LIMIT

  return (
    <>
      <PageHeader
        title="Зачисления"
        description="Управляйте зачислениями студентов на курсы."
        actions={
          <Button onClick={() => { setCreateOpen(true); setCreateError(null) }}>
            <Plus size={16} aria-hidden="true" />
            Создать
          </Button>
        }
      />

      {/* ── Фильтры ── */}
      <div className={styles.toolbar}>
        <div className={styles.filterGroup}>
          <label htmlFor="f-account" className={styles.filterLabel}>ID аккаунта</label>
          <input
            id="f-account"
            className={styles.filterInput}
            value={accountIdFilter}
            onChange={(e) => setAccountIdFilter(e.target.value)}
            onBlur={applyFilter}
            onKeyDown={(e) => { if (e.key === 'Enter') applyFilter() }}
            placeholder="UUID аккаунта…"
          />
        </div>
        <div className={styles.filterGroup}>
          <label htmlFor="f-course" className={styles.filterLabel}>ID курса</label>
          <input
            id="f-course"
            className={styles.filterInput}
            value={courseIdFilter}
            onChange={(e) => setCourseIdFilter(e.target.value)}
            onBlur={applyFilter}
            onKeyDown={(e) => { if (e.key === 'Enter') applyFilter() }}
            placeholder="UUID курса…"
          />
        </div>
        {(accountIdFilter || courseIdFilter) ? (
          <button
            type="button"
            className={styles.clearBtn}
            onClick={() => { setAccountIdFilter(''); setCourseIdFilter(''); setOffset(0) }}
          >
            Сбросить
          </button>
        ) : null}
      </div>

      {/* ── Контент ── */}
      {isLoading || isFetching ? (
        <p className={styles.state}>Загрузка…</p>
      ) : error ? (
        <p className={styles.stateError}>
          {error instanceof ApiError ? error.message : 'Не удалось загрузить зачисления'}
        </p>
      ) : enrollments.length === 0 ? (
        <div className={styles.emptyHint}>
          <BookMarked size={40} className={styles.emptyIcon} aria-hidden="true" />
          <p>Зачисления не найдены</p>
        </div>
      ) : (
        <>
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Статус</th>
                  <th>Аккаунт</th>
                  <th>Курс</th>
                  <th>Дата начала</th>
                  <th>Дата окончания</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {enrollments.map((enrollment) => (
                  <EnrollmentRow key={enrollment.id} enrollment={enrollment} />
                ))}
              </tbody>
            </table>
          </div>
          <Pagination offset={offset} limit={LIMIT} hasMore={hasMore} onOffsetChange={setOffset} />
        </>
      )}

      {/* ── Модал создания ── */}
      <Modal open={createOpen} onClose={() => setCreateOpen(false)} title="Создать зачисление">
        <form className={styles.form} onSubmit={handleCreate} noValidate>
          <FormField label="ID аккаунта" htmlFor="ce-account" required>
            <input
              id="ce-account"
              value={createForm.accountId}
              onChange={(e) => setCreateForm((f) => ({ ...f, accountId: e.target.value }))}
              placeholder="UUID из раздела Учётные записи"
              required
              autoFocus
            />
          </FormField>
          <FormField label="ID курса" htmlFor="ce-course" required>
            <input
              id="ce-course"
              value={createForm.courseId}
              onChange={(e) => setCreateForm((f) => ({ ...f, courseId: e.target.value }))}
              placeholder="UUID из раздела Курсы"
              required
            />
          </FormField>
          <FormField label="Дата начала" htmlFor="ce-start">
            <input
              id="ce-start"
              type="date"
              value={createForm.activatedAt}
              onChange={(e) => setCreateForm((f) => ({ ...f, activatedAt: e.target.value }))}
            />
          </FormField>
          <FormField label="Дата окончания" htmlFor="ce-end">
            <input
              id="ce-end"
              type="date"
              value={createForm.deactivatedAt}
              onChange={(e) => setCreateForm((f) => ({ ...f, deactivatedAt: e.target.value }))}
            />
          </FormField>

          {createError ? <p className={styles.formError}>{createError}</p> : null}

          <div className={styles.formActions}>
            <Button variant="secondary" type="button" onClick={() => setCreateOpen(false)} disabled={createMutation.isPending}>
              Отмена
            </Button>
            <Button type="submit" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Создание…' : 'Создать'}
            </Button>
          </div>
        </form>
      </Modal>
    </>
  )
}

function EnrollmentRow({ enrollment }: { enrollment: EnrollmentView }) {
  const statusLabel =
    ENROLLMENT_STATUS_LABELS[enrollment.status as EnrollmentStatus] ?? enrollment.status_title

  return (
    <tr>
      <td>
        <Badge variant={statusBadgeVariant(enrollment.status)}>{statusLabel}</Badge>
      </td>
      <td className={styles.idCell}>
        <Link to={`/accounts`} className={styles.idLink} title={enrollment.account_id}>
          {shortId(enrollment.account_id)}
        </Link>
      </td>
      <td className={styles.idCell}>
        <Link to={`/courses`} className={styles.idLink} title={enrollment.course_id}>
          {shortId(enrollment.course_id)}
        </Link>
      </td>
      <td className={styles.dateCell}>{enrollment.activated_at || '—'}</td>
      <td className={styles.dateCell}>{enrollment.deactivated_at || '—'}</td>
      <td>
        <Link to={`/enrollments/${enrollment.id}/attempts`} className={styles.idLink}>
          Попытки
        </Link>
      </td>
    </tr>
  )
}
