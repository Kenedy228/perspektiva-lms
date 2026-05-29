import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Building2, Pencil, Plus, Trash2 } from 'lucide-react'
import { FormEvent, useState } from 'react'
import {
  changeOrganizationINN,
  createOrganization,
  deleteOrganization,
  renameOrganization,
} from '../../api/organizations'
import { Badge } from '../../components/ui/Badge'
import { Button } from '../../components/ui/Button'
import { ConfirmDialog } from '../../components/ui/ConfirmDialog'
import { FormField } from '../../components/ui/FormField'
import { Modal } from '../../components/ui/Modal'
import { PageHeader } from '../../components/ui/PageHeader'
import { Pagination } from '../../components/ui/Pagination'
import { useDebounce } from '../../lib/hooks/useDebounce'
import type { InnType, OrganizationShortView } from '../../types/organizations'
import { INN_TYPE_LABELS, INN_TYPES } from '../../types/organizations'
import { ApiError } from '../../types/api'
import { type SearchMode, useOrganizations } from './useOrganizations'
import styles from './OrganizationsPage.module.css'

const LIMIT = 20

type CreateForm = { name: string; inn: string; inn_type: InnType }
type RenameForm = { id: string; name: string }
type INNForm = { id: string; inn: string; inn_type: InnType }
type DeleteTarget = { id: string; name: string }

function defaultCreateForm(): CreateForm {
  return { name: '', inn: '', inn_type: 'legal entity' }
}

export function OrganizationsPage() {
  const queryClient = useQueryClient()

  const [search, setSearch] = useState('')
  const [searchMode, setSearchMode] = useState<SearchMode>('name')
  const [offset, setOffset] = useState(0)

  const debouncedSearch = useDebounce(search, 350)

  const { data: orgs = [], isLoading, error, isFetching } = useOrganizations(
    debouncedSearch,
    searchMode,
    LIMIT,
    offset,
  )

  const [createOpen, setCreateOpen] = useState(false)
  const [createForm, setCreateForm] = useState<CreateForm>(defaultCreateForm)
  const [createError, setCreateError] = useState<string | null>(null)

  const [renameTarget, setRenameTarget] = useState<RenameForm | null>(null)
  const [renameError, setRenameError] = useState<string | null>(null)

  const [innTarget, setINNTarget] = useState<INNForm | null>(null)
  const [innError, setINNError] = useState<string | null>(null)

  const [deleteTarget, setDeleteTarget] = useState<DeleteTarget | null>(null)

  function invalidate() {
    void queryClient.invalidateQueries({ queryKey: ['organizations'] })
  }

  const createMutation = useMutation({
    mutationFn: createOrganization,
    onSuccess: () => {
      invalidate()
      setCreateOpen(false)
      setCreateForm(defaultCreateForm())
      setCreateError(null)
    },
  })

  const renameMutation = useMutation({
    mutationFn: ({ id, name }: { id: string; name: string }) => renameOrganization(id, name),
    onSuccess: () => {
      invalidate()
      setRenameTarget(null)
      setRenameError(null)
    },
  })

  const innMutation = useMutation({
    mutationFn: ({ id, inn, inn_type }: { id: string; inn: string; inn_type: InnType }) =>
      changeOrganizationINN(id, inn, inn_type),
    onSuccess: () => {
      invalidate()
      setINNTarget(null)
      setINNError(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteOrganization(id),
    onSuccess: () => {
      invalidate()
      setDeleteTarget(null)
    },
  })

  function extractMessage(caught: unknown, fallback: string) {
    return caught instanceof ApiError ? caught.message : fallback
  }

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setCreateError(null)
    const payload = createForm.inn.trim()
      ? { name: createForm.name.trim(), inn: createForm.inn.trim(), inn_type: createForm.inn_type }
      : { name: createForm.name.trim() }
    try {
      await createMutation.mutateAsync(payload)
    } catch (caught) {
      setCreateError(extractMessage(caught, 'Не удалось создать организацию'))
    }
  }

  async function handleRename(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!renameTarget) return
    setRenameError(null)
    try {
      await renameMutation.mutateAsync({ id: renameTarget.id, name: renameTarget.name.trim() })
    } catch (caught) {
      setRenameError(extractMessage(caught, 'Не удалось переименовать организацию'))
    }
  }

  async function handleChangeINN(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!innTarget) return
    setINNError(null)
    try {
      await innMutation.mutateAsync({
        id: innTarget.id,
        inn: innTarget.inn.trim(),
        inn_type: innTarget.inn_type,
      })
    } catch (caught) {
      setINNError(extractMessage(caught, 'Не удалось изменить ИНН'))
    }
  }

  function openRename(org: OrganizationShortView) {
    setRenameTarget({ id: org.ID, name: org.OrganizationName })
    setRenameError(null)
  }

  function openINN(org: OrganizationShortView) {
    setINNTarget({ id: org.ID, inn: '', inn_type: 'legal entity' })
    setINNError(null)
  }

  const showEmpty = debouncedSearch.trim().length === 0
  const hasMore = orgs.length >= LIMIT

  return (
    <>
      <PageHeader
        title="Организации"
        description="Управляйте записями организаций, их названиями и ИНН."
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
          placeholder={searchMode === 'inn' ? 'Поиск по ИНН…' : 'Поиск по названию…'}
          aria-label="Поиск организации"
        />
        <div className={styles.modeToggle}>
          <button
            type="button"
            className={`${styles.modeBtn} ${searchMode === 'name' ? styles.modeBtnActive : ''}`}
            onClick={() => { setSearchMode('name'); setOffset(0) }}
          >
            По названию
          </button>
          <button
            type="button"
            className={`${styles.modeBtn} ${searchMode === 'inn' ? styles.modeBtnActive : ''}`}
            onClick={() => { setSearchMode('inn'); setOffset(0) }}
          >
            По ИНН
          </button>
        </div>
      </div>

      {isLoading || isFetching ? (
        <p className={styles.state}>Загрузка…</p>
      ) : error ? (
        <p className={styles.stateError}>
          {error instanceof ApiError ? error.message : 'Не удалось загрузить организации'}
        </p>
      ) : showEmpty ? (
        <div className={styles.emptyHint}>
          <Building2 size={40} className={styles.emptyIcon} aria-hidden="true" />
          <p>Введите название или ИНН для поиска организаций</p>
        </div>
      ) : orgs.length === 0 ? (
        <p className={styles.state}>Организации не найдены</p>
      ) : (
        <>
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Название</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                {orgs.map((org) => (
                  <tr key={org.ID}>
                    <td className={styles.nameCell}>{org.OrganizationName}</td>
                    <td>
                      <div className={styles.rowActions}>
                        <Button
                          variant="secondary"
                          onClick={() => openRename(org)}
                          title="Переименовать"
                        >
                          <Pencil size={15} aria-hidden="true" />
                          Переименовать
                        </Button>
                        <Button
                          variant="secondary"
                          onClick={() => openINN(org)}
                          title="Изменить ИНН"
                        >
                          <Badge variant="neutral">ИНН</Badge>
                          Изменить ИНН
                        </Button>
                        <Button
                          variant="secondary"
                          onClick={() => setDeleteTarget({ id: org.ID, name: org.OrganizationName })}
                          title="Удалить"
                        >
                          <Trash2 size={15} aria-hidden="true" />
                        </Button>
                      </div>
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

      {/* ── Модал создания ── */}
      <Modal
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        title="Создать организацию"
      >
        <form className={styles.form} onSubmit={handleCreate} noValidate>
          <FormField label="Название" htmlFor="create-name" required>
            <input
              id="create-name"
              value={createForm.name}
              onChange={(e) => setCreateForm((f) => ({ ...f, name: e.target.value }))}
              placeholder="ООО «Пример»"
              required
              autoFocus
            />
          </FormField>

          <FormField label="ИНН" htmlFor="create-inn">
            <input
              id="create-inn"
              value={createForm.inn}
              onChange={(e) => setCreateForm((f) => ({ ...f, inn: e.target.value }))}
              placeholder="Необязательно"
            />
          </FormField>

          {createForm.inn.trim() ? (
            <FormField label="Тип ИНН" htmlFor="create-inn-type" required>
              <select
                id="create-inn-type"
                value={createForm.inn_type}
                onChange={(e) =>
                  setCreateForm((f) => ({ ...f, inn_type: e.target.value as InnType }))
                }
              >
                {INN_TYPES.map((t) => (
                  <option key={t} value={t}>
                    {INN_TYPE_LABELS[t]}
                  </option>
                ))}
              </select>
            </FormField>
          ) : null}

          {createError ? <p className={styles.formError}>{createError}</p> : null}

          <div className={styles.formActions}>
            <Button
              variant="secondary"
              type="button"
              onClick={() => setCreateOpen(false)}
              disabled={createMutation.isPending}
            >
              Отмена
            </Button>
            <Button type="submit" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Создание…' : 'Создать'}
            </Button>
          </div>
        </form>
      </Modal>

      {/* ── Модал переименования ── */}
      <Modal
        open={renameTarget !== null}
        onClose={() => setRenameTarget(null)}
        title="Переименовать организацию"
        size="sm"
      >
        <form className={styles.form} onSubmit={handleRename} noValidate>
          <FormField label="Новое название" htmlFor="rename-name" required>
            <input
              id="rename-name"
              value={renameTarget?.name ?? ''}
              onChange={(e) =>
                setRenameTarget((t) => (t ? { ...t, name: e.target.value } : null))
              }
              placeholder="Введите название"
              required
              autoFocus
            />
          </FormField>

          {renameError ? <p className={styles.formError}>{renameError}</p> : null}

          <div className={styles.formActions}>
            <Button
              variant="secondary"
              type="button"
              onClick={() => setRenameTarget(null)}
              disabled={renameMutation.isPending}
            >
              Отмена
            </Button>
            <Button type="submit" disabled={renameMutation.isPending}>
              {renameMutation.isPending ? 'Сохранение…' : 'Сохранить'}
            </Button>
          </div>
        </form>
      </Modal>

      {/* ── Модал изменения ИНН ── */}
      <Modal
        open={innTarget !== null}
        onClose={() => setINNTarget(null)}
        title="Изменить ИНН"
        size="sm"
      >
        <form className={styles.form} onSubmit={handleChangeINN} noValidate>
          <FormField label="ИНН" htmlFor="inn-value" required>
            <input
              id="inn-value"
              value={innTarget?.inn ?? ''}
              onChange={(e) =>
                setINNTarget((t) => (t ? { ...t, inn: e.target.value } : null))
              }
              placeholder="Введите ИНН"
              required
              autoFocus
            />
          </FormField>

          <FormField label="Тип ИНН" htmlFor="inn-type" required>
            <select
              id="inn-type"
              value={innTarget?.inn_type ?? 'legal entity'}
              onChange={(e) =>
                setINNTarget((t) => (t ? { ...t, inn_type: e.target.value as InnType } : null))
              }
            >
              {INN_TYPES.map((t) => (
                <option key={t} value={t}>
                  {INN_TYPE_LABELS[t]}
                </option>
              ))}
            </select>
          </FormField>

          {innError ? <p className={styles.formError}>{innError}</p> : null}

          <div className={styles.formActions}>
            <Button
              variant="secondary"
              type="button"
              onClick={() => setINNTarget(null)}
              disabled={innMutation.isPending}
            >
              Отмена
            </Button>
            <Button type="submit" disabled={innMutation.isPending}>
              {innMutation.isPending ? 'Сохранение…' : 'Сохранить'}
            </Button>
          </div>
        </form>
      </Modal>

      {/* ── Подтверждение удаления ── */}
      <ConfirmDialog
        open={deleteTarget !== null}
        onClose={() => setDeleteTarget(null)}
        onConfirm={() => {
          if (deleteTarget) void deleteMutation.mutateAsync(deleteTarget.id)
        }}
        title="Удалить организацию"
        message={`Вы уверены, что хотите удалить «${deleteTarget?.name ?? ''}»? Это действие нельзя отменить.`}
        confirmLabel="Удалить"
        danger
        isPending={deleteMutation.isPending}
      />
    </>
  )
}
