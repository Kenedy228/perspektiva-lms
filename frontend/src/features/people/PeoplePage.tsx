import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus, Trash2, UserRound } from 'lucide-react'
import { FormEvent, useState } from 'react'
import { Link } from 'react-router-dom'
import { createPerson, deletePerson } from '../../api/people'
import { Button } from '../../components/ui/Button'
import { ConfirmDialog } from '../../components/ui/ConfirmDialog'
import { FormField } from '../../components/ui/FormField'
import { Modal } from '../../components/ui/Modal'
import { PageHeader } from '../../components/ui/PageHeader'
import { Pagination } from '../../components/ui/Pagination'
import { ApiError } from '../../types/api'
import type { PersonShortView } from '../../types/person'
import { useDebounce } from '../../lib/hooks/useDebounce'
import { type PeopleSearchMode, usePeople } from './usePeople'
import styles from './PeoplePage.module.css'

const LIMIT = 20

type CreateForm = { firstName: string; lastName: string; middleName: string }
type DeleteTarget = { id: string; name: string }

function defaultCreate(): CreateForm {
  return { firstName: '', lastName: '', middleName: '' }
}

export function PeoplePage() {
  const queryClient = useQueryClient()

  const [search, setSearch] = useState('')
  const [searchMode, setSearchMode] = useState<PeopleSearchMode>('last_name')
  const [offset, setOffset] = useState(0)
  const debouncedSearch = useDebounce(search, 350)

  const { data: people = [], isLoading, isFetching, error } = usePeople(
    debouncedSearch,
    searchMode,
    LIMIT,
    offset,
  )

  const [createOpen, setCreateOpen] = useState(false)
  const [createForm, setCreateForm] = useState<CreateForm>(defaultCreate)
  const [createError, setCreateError] = useState<string | null>(null)

  const [deleteTarget, setDeleteTarget] = useState<DeleteTarget | null>(null)

  function invalidate() {
    void queryClient.invalidateQueries({ queryKey: ['people'] })
  }

  const createMutation = useMutation({
    mutationFn: createPerson,
    onSuccess: () => {
      invalidate()
      setCreateOpen(false)
      setCreateForm(defaultCreate())
      setCreateError(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deletePerson(id),
    onSuccess: () => {
      invalidate()
      setDeleteTarget(null)
    },
  })

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setCreateError(null)
    try {
      await createMutation.mutateAsync({
        first_name: createForm.firstName.trim(),
        last_name: createForm.lastName.trim(),
        middle_name: createForm.middleName.trim() || undefined,
      })
    } catch (caught) {
      setCreateError(caught instanceof ApiError ? caught.message : 'Не удалось создать сотрудника')
    }
  }

  function switchMode(mode: PeopleSearchMode) {
    setSearchMode(mode)
    setSearch('')
    setOffset(0)
  }

  const isSnilsMode = searchMode === 'snils'
  const showSnilsHint = isSnilsMode && debouncedSearch.trim().length === 0
  const hasMore = people.length >= LIMIT

  return (
    <>
      <PageHeader
        title="Сотрудники"
        description="Управляйте персональными данными сотрудников организации."
        actions={
          <Button onClick={() => { setCreateOpen(true); setCreateError(null) }}>
            <Plus size={16} aria-hidden="true" />
            Создать
          </Button>
        }
      />

      <div className={styles.toolbar}>
        <input
          className={styles.searchInput}
          value={search}
          onChange={(e) => { setSearch(e.target.value); setOffset(0) }}
          placeholder={isSnilsMode ? 'Поиск по СНИЛС…' : 'Поиск по фамилии…'}
          aria-label="Поиск сотрудника"
        />
        <div className={styles.modeToggle}>
          <button
            type="button"
            className={`${styles.modeBtn} ${searchMode === 'last_name' ? styles.modeBtnActive : ''}`}
            onClick={() => switchMode('last_name')}
          >
            По фамилии
          </button>
          <button
            type="button"
            className={`${styles.modeBtn} ${searchMode === 'snils' ? styles.modeBtnActive : ''}`}
            onClick={() => switchMode('snils')}
          >
            По СНИЛС
          </button>
        </div>
      </div>

      {isLoading || isFetching ? (
        <p className={styles.state}>Загрузка…</p>
      ) : error ? (
        <p className={styles.stateError}>
          {error instanceof ApiError ? error.message : 'Не удалось загрузить сотрудников'}
        </p>
      ) : showSnilsHint ? (
        <div className={styles.emptyHint}>
          <UserRound size={40} className={styles.emptyIcon} aria-hidden="true" />
          <p>Введите СНИЛС для поиска</p>
        </div>
      ) : people.length === 0 ? (
        <p className={styles.state}>Сотрудники не найдены</p>
      ) : (
        <>
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>ФИО</th>
                  <th>Организация</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                {people.map((person) => (
                  <PersonRow
                    key={person.ID}
                    person={person}
                    onDelete={() => setDeleteTarget({ id: person.ID, name: person.FullName })}
                  />
                ))}
              </tbody>
            </table>
          </div>
          <Pagination offset={offset} limit={LIMIT} hasMore={hasMore} onOffsetChange={setOffset} />
        </>
      )}

      {/* ── Модал создания ── */}
      <Modal open={createOpen} onClose={() => setCreateOpen(false)} title="Создать сотрудника">
        <form className={styles.form} onSubmit={handleCreate} noValidate>
          <FormField label="Фамилия" htmlFor="c-last" required>
            <input
              id="c-last"
              value={createForm.lastName}
              onChange={(e) => setCreateForm((f) => ({ ...f, lastName: e.target.value }))}
              placeholder="Иванов"
              required
              autoFocus
            />
          </FormField>
          <FormField label="Имя" htmlFor="c-first" required>
            <input
              id="c-first"
              value={createForm.firstName}
              onChange={(e) => setCreateForm((f) => ({ ...f, firstName: e.target.value }))}
              placeholder="Иван"
              required
            />
          </FormField>
          <FormField label="Отчество" htmlFor="c-middle">
            <input
              id="c-middle"
              value={createForm.middleName}
              onChange={(e) => setCreateForm((f) => ({ ...f, middleName: e.target.value }))}
              placeholder="Иванович (необязательно)"
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

      {/* ── Подтверждение удаления ── */}
      <ConfirmDialog
        open={deleteTarget !== null}
        onClose={() => setDeleteTarget(null)}
        onConfirm={() => { if (deleteTarget) void deleteMutation.mutateAsync(deleteTarget.id) }}
        title="Удалить сотрудника"
        message={`Вы уверены, что хотите удалить «${deleteTarget?.name ?? ''}»? Это действие нельзя отменить.`}
        confirmLabel="Удалить"
        danger
        isPending={deleteMutation.isPending}
      />
    </>
  )
}

function PersonRow({
  person,
  onDelete,
}: {
  person: PersonShortView
  onDelete: () => void
}) {
  return (
    <tr>
      <td>
        <Link to={`/people/${person.ID}`} className={styles.nameLink}>
          {person.FullName}
        </Link>
      </td>
      <td className={styles.orgCell}>
        {person.OrganizationName || <span className={styles.noOrg}>—</span>}
      </td>
      <td>
        <div className={styles.rowActions}>
          <Button as={Link} to={`/people/${person.ID}`} variant="secondary">
            Открыть
          </Button>
          <Button variant="secondary" onClick={onDelete} title="Удалить">
            <Trash2 size={15} aria-hidden="true" />
          </Button>
        </div>
      </td>
    </tr>
  )
}
