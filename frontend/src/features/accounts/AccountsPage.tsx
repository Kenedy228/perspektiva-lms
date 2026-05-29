import { useMutation, useQueryClient } from '@tanstack/react-query'
import { KeyRound, Lock, LockOpen, Pencil, Plus, Shield, Trash2 } from 'lucide-react'
import { FormEvent, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  activateAccount,
  blockAccount,
  changeAccountLogin,
  changeAccountPassword,
  changeAccountRole,
  createAccount,
  deleteAccount,
} from '../../api/accounts'
import { Badge } from '../../components/ui/Badge'
import { Button } from '../../components/ui/Button'
import { ConfirmDialog } from '../../components/ui/ConfirmDialog'
import { FormField } from '../../components/ui/FormField'
import { Modal } from '../../components/ui/Modal'
import { PageHeader } from '../../components/ui/PageHeader'
import { Pagination } from '../../components/ui/Pagination'
import { ApiError } from '../../types/api'
import type { AccountStatus, AccountView } from '../../types/account'
import {
  ACCOUNT_ROLES,
  ROLE_LABELS,
  STATUS_LABELS,
} from '../../types/account'
import type { Role } from '../../types/auth'
import { useDebounce } from '../../lib/hooks/useDebounce'
import { useAccounts } from './useAccounts'
import styles from './AccountsPage.module.css'

const LIMIT = 20

// ── Badge helpers ─────────────────────────────────────────────────────────────

function roleBadgeVariant(role: Role) {
  if (role === 'admin') return 'warning' as const
  if (role === 'creator') return 'neutral' as const
  if (role === 'student') return 'success' as const
  return 'neutral' as const
}

function statusBadgeVariant(status: AccountStatus) {
  if (status === 'active') return 'success' as const
  if (status === 'blocked') return 'danger' as const
  return 'neutral' as const
}

// ── Modal target types ────────────────────────────────────────────────────────

type LoginTarget = { id: string; login: string }
type PasswordTarget = { id: string; login: string }
type RoleTarget = { id: string; login: string; role: Role }
type ActionTarget = { id: string; login: string }

type CreateForm = { personId: string; login: string; password: string; role: Role }

function defaultCreate(): CreateForm {
  return { personId: '', login: '', password: '', role: 'student' }
}

// ── Page ──────────────────────────────────────────────────────────────────────

export function AccountsPage() {
  const queryClient = useQueryClient()

  // Filters
  const [loginSearch, setLoginSearch] = useState('')
  const [roleFilter, setRoleFilter] = useState<Role | ''>('')
  const [statusFilter, setStatusFilter] = useState<AccountStatus | ''>('')
  const [offset, setOffset] = useState(0)
  const debouncedLogin = useDebounce(loginSearch, 350)

  const filter = {
    login: debouncedLogin || undefined,
    role: roleFilter || undefined,
    status: statusFilter || undefined,
    limit: LIMIT,
    offset,
  }

  const { data: accounts = [], isLoading, isFetching, error } = useAccounts(filter)

  // Modal state
  const [createOpen, setCreateOpen] = useState(false)
  const [createForm, setCreateForm] = useState<CreateForm>(defaultCreate)
  const [createError, setCreateError] = useState<string | null>(null)

  const [loginTarget, setLoginTarget] = useState<LoginTarget | null>(null)
  const [newLogin, setNewLogin] = useState('')
  const [loginError, setLoginError] = useState<string | null>(null)

  const [passwordTarget, setPasswordTarget] = useState<PasswordTarget | null>(null)
  const [newPassword, setNewPassword] = useState('')
  const [passwordError, setPasswordError] = useState<string | null>(null)

  const [roleTarget, setRoleTarget] = useState<RoleTarget | null>(null)
  const [roleError, setRoleError] = useState<string | null>(null)

  const [blockTarget, setBlockTarget] = useState<ActionTarget | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<ActionTarget | null>(null)

  function invalidate() {
    void queryClient.invalidateQueries({ queryKey: ['accounts'] })
  }

  function extractMessage(caught: unknown, fallback: string) {
    return caught instanceof ApiError ? caught.message : fallback
  }

  // ── Mutations ──────────────────────────────────────────────────────────────

  const createMutation = useMutation({
    mutationFn: createAccount,
    onSuccess: () => { invalidate(); setCreateOpen(false); setCreateForm(defaultCreate()); setCreateError(null) },
  })

  const loginMutation = useMutation({
    mutationFn: ({ id, login }: { id: string; login: string }) => changeAccountLogin(id, login),
    onSuccess: () => { invalidate(); setLoginTarget(null); setNewLogin(''); setLoginError(null) },
  })

  const passwordMutation = useMutation({
    mutationFn: ({ id, password }: { id: string; password: string }) => changeAccountPassword(id, password),
    onSuccess: () => { invalidate(); setPasswordTarget(null); setNewPassword(''); setPasswordError(null) },
  })

  const roleMutation = useMutation({
    mutationFn: ({ id, role }: { id: string; role: Role }) => changeAccountRole(id, role),
    onSuccess: () => { invalidate(); setRoleTarget(null); setRoleError(null) },
  })

  const blockMutation = useMutation({
    mutationFn: (id: string) => blockAccount(id),
    onSuccess: () => { invalidate(); setBlockTarget(null) },
  })

  const activateMutation = useMutation({
    mutationFn: (id: string) => activateAccount(id),
    onSuccess: () => invalidate(),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteAccount(id),
    onSuccess: () => { invalidate(); setDeleteTarget(null) },
  })

  // ── Handlers ──────────────────────────────────────────────────────────────

  async function handleCreate(event: FormEvent) {
    event.preventDefault()
    setCreateError(null)
    try {
      await createMutation.mutateAsync({
        person_id: createForm.personId.trim(),
        login: createForm.login.trim(),
        password: createForm.password,
        role: createForm.role,
      })
    } catch (caught) {
      setCreateError(extractMessage(caught, 'Не удалось создать учётную запись'))
    }
  }

  async function handleChangeLogin(event: FormEvent) {
    event.preventDefault()
    if (!loginTarget) return
    setLoginError(null)
    try {
      await loginMutation.mutateAsync({ id: loginTarget.id, login: newLogin.trim() })
    } catch (caught) {
      setLoginError(extractMessage(caught, 'Не удалось изменить логин'))
    }
  }

  async function handleChangePassword(event: FormEvent) {
    event.preventDefault()
    if (!passwordTarget) return
    setPasswordError(null)
    try {
      await passwordMutation.mutateAsync({ id: passwordTarget.id, password: newPassword })
    } catch (caught) {
      setPasswordError(extractMessage(caught, 'Не удалось изменить пароль'))
    }
  }

  async function handleChangeRole(event: FormEvent) {
    event.preventDefault()
    if (!roleTarget) return
    setRoleError(null)
    try {
      await roleMutation.mutateAsync({ id: roleTarget.id, role: roleTarget.role })
    } catch (caught) {
      setRoleError(extractMessage(caught, 'Не удалось изменить роль'))
    }
  }

  const hasMore = accounts.length >= LIMIT

  return (
    <>
      <PageHeader
        title="Учётные записи"
        description="Управляйте учётными записями, ролями и статусами пользователей."
        actions={
          <Button onClick={() => { setCreateOpen(true); setCreateError(null) }}>
            <Plus size={16} aria-hidden="true" />
            Создать
          </Button>
        }
      />

      {/* ── Toolbar ── */}
      <div className={styles.toolbar}>
        <input
          className={styles.searchInput}
          value={loginSearch}
          onChange={(e) => { setLoginSearch(e.target.value); setOffset(0) }}
          placeholder="Поиск по логину…"
          aria-label="Поиск по логину"
        />
        <select
          className={styles.filterSelect}
          value={roleFilter}
          onChange={(e) => { setRoleFilter(e.target.value as Role | ''); setOffset(0) }}
          aria-label="Фильтр по роли"
        >
          <option value="">Все роли</option>
          {ACCOUNT_ROLES.map((r) => (
            <option key={r} value={r}>{ROLE_LABELS[r]}</option>
          ))}
        </select>
        <select
          className={styles.filterSelect}
          value={statusFilter}
          onChange={(e) => { setStatusFilter(e.target.value as AccountStatus | ''); setOffset(0) }}
          aria-label="Фильтр по статусу"
        >
          <option value="">Все статусы</option>
          <option value="active">{STATUS_LABELS.active}</option>
          <option value="blocked">{STATUS_LABELS.blocked}</option>
          <option value="deleted">{STATUS_LABELS.deleted}</option>
        </select>
      </div>

      {/* ── Table ── */}
      {isLoading || isFetching ? (
        <p className={styles.state}>Загрузка…</p>
      ) : error ? (
        <p className={styles.stateError}>
          {error instanceof ApiError ? error.message : 'Не удалось загрузить учётные записи'}
        </p>
      ) : accounts.length === 0 ? (
        <p className={styles.state}>Учётные записи не найдены</p>
      ) : (
        <>
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>Логин</th>
                  <th>Роль</th>
                  <th>Статус</th>
                  <th>Сотрудник</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                {accounts.map((account) => (
                  <AccountRow
                    key={account.ID}
                    account={account}
                    activating={activateMutation.isPending && activateMutation.variables === account.ID}
                    onChangeLogin={() => { setLoginTarget({ id: account.ID, login: account.Login }); setNewLogin(account.Login); setLoginError(null) }}
                    onChangePassword={() => { setPasswordTarget({ id: account.ID, login: account.Login }); setNewPassword(''); setPasswordError(null) }}
                    onChangeRole={() => { setRoleTarget({ id: account.ID, login: account.Login, role: account.Role }); setRoleError(null) }}
                    onBlock={() => setBlockTarget({ id: account.ID, login: account.Login })}
                    onActivate={() => void activateMutation.mutateAsync(account.ID)}
                    onDelete={() => setDeleteTarget({ id: account.ID, login: account.Login })}
                  />
                ))}
              </tbody>
            </table>
          </div>
          <Pagination offset={offset} limit={LIMIT} hasMore={hasMore} onOffsetChange={setOffset} />
        </>
      )}

      {/* ── Модал создания ── */}
      <Modal open={createOpen} onClose={() => setCreateOpen(false)} title="Создать учётную запись">
        <form className={styles.form} onSubmit={handleCreate} noValidate>
          <FormField label="ID сотрудника" htmlFor="ca-person" required>
            <input
              id="ca-person"
              value={createForm.personId}
              onChange={(e) => setCreateForm((f) => ({ ...f, personId: e.target.value }))}
              placeholder="UUID из раздела Сотрудники"
              required
              autoFocus
            />
          </FormField>
          <FormField label="Логин" htmlFor="ca-login" required>
            <input
              id="ca-login"
              value={createForm.login}
              onChange={(e) => setCreateForm((f) => ({ ...f, login: e.target.value }))}
              placeholder="Введите логин"
              autoComplete="off"
              required
            />
          </FormField>
          <FormField label="Пароль" htmlFor="ca-password" required>
            <input
              id="ca-password"
              type="password"
              value={createForm.password}
              onChange={(e) => setCreateForm((f) => ({ ...f, password: e.target.value }))}
              placeholder="Временный пароль"
              autoComplete="new-password"
              required
            />
          </FormField>
          <FormField label="Роль" htmlFor="ca-role" required>
            <select
              id="ca-role"
              value={createForm.role}
              onChange={(e) => setCreateForm((f) => ({ ...f, role: e.target.value as Role }))}
            >
              {ACCOUNT_ROLES.map((r) => (
                <option key={r} value={r}>{ROLE_LABELS[r]}</option>
              ))}
            </select>
          </FormField>
          {createError ? <p className={styles.formError}>{createError}</p> : null}
          <div className={styles.formActions}>
            <Button variant="secondary" type="button" onClick={() => setCreateOpen(false)} disabled={createMutation.isPending}>Отмена</Button>
            <Button type="submit" disabled={createMutation.isPending}>{createMutation.isPending ? 'Создание…' : 'Создать'}</Button>
          </div>
        </form>
      </Modal>

      {/* ── Изменить логин ── */}
      <Modal open={loginTarget !== null} onClose={() => setLoginTarget(null)} title="Изменить логин" size="sm">
        <form className={styles.form} onSubmit={handleChangeLogin} noValidate>
          <p className={styles.modalHint}>Аккаунт: <strong>{loginTarget?.login}</strong></p>
          <FormField label="Новый логин" htmlFor="nl-login" required>
            <input id="nl-login" value={newLogin} onChange={(e) => setNewLogin(e.target.value)} required autoFocus autoComplete="off" />
          </FormField>
          {loginError ? <p className={styles.formError}>{loginError}</p> : null}
          <div className={styles.formActions}>
            <Button variant="secondary" type="button" onClick={() => setLoginTarget(null)} disabled={loginMutation.isPending}>Отмена</Button>
            <Button type="submit" disabled={loginMutation.isPending}>{loginMutation.isPending ? 'Сохранение…' : 'Сохранить'}</Button>
          </div>
        </form>
      </Modal>

      {/* ── Изменить пароль ── */}
      <Modal open={passwordTarget !== null} onClose={() => setPasswordTarget(null)} title="Изменить пароль" size="sm">
        <form className={styles.form} onSubmit={handleChangePassword} noValidate>
          <p className={styles.modalHint}>Аккаунт: <strong>{passwordTarget?.login}</strong></p>
          <FormField label="Новый пароль" htmlFor="np-pass" required>
            <input id="np-pass" type="password" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} required autoFocus autoComplete="new-password" placeholder="Введите новый пароль" />
          </FormField>
          {passwordError ? <p className={styles.formError}>{passwordError}</p> : null}
          <div className={styles.formActions}>
            <Button variant="secondary" type="button" onClick={() => setPasswordTarget(null)} disabled={passwordMutation.isPending}>Отмена</Button>
            <Button type="submit" disabled={passwordMutation.isPending}>{passwordMutation.isPending ? 'Сохранение…' : 'Сохранить'}</Button>
          </div>
        </form>
      </Modal>

      {/* ── Изменить роль ── */}
      <Modal open={roleTarget !== null} onClose={() => setRoleTarget(null)} title="Изменить роль" size="sm">
        <form className={styles.form} onSubmit={handleChangeRole} noValidate>
          <p className={styles.modalHint}>Аккаунт: <strong>{roleTarget?.login}</strong></p>
          <FormField label="Новая роль" htmlFor="nr-role" required>
            <select
              id="nr-role"
              value={roleTarget?.role ?? 'student'}
              onChange={(e) => setRoleTarget((t) => t ? { ...t, role: e.target.value as Role } : null)}
            >
              {ACCOUNT_ROLES.map((r) => (
                <option key={r} value={r}>{ROLE_LABELS[r]}</option>
              ))}
            </select>
          </FormField>
          {roleError ? <p className={styles.formError}>{roleError}</p> : null}
          <div className={styles.formActions}>
            <Button variant="secondary" type="button" onClick={() => setRoleTarget(null)} disabled={roleMutation.isPending}>Отмена</Button>
            <Button type="submit" disabled={roleMutation.isPending}>{roleMutation.isPending ? 'Сохранение…' : 'Сохранить'}</Button>
          </div>
        </form>
      </Modal>

      {/* ── Блокировка ── */}
      <ConfirmDialog
        open={blockTarget !== null}
        onClose={() => setBlockTarget(null)}
        onConfirm={() => { if (blockTarget) void blockMutation.mutateAsync(blockTarget.id) }}
        title="Заблокировать аккаунт"
        message={`Аккаунт «${blockTarget?.login ?? ''}» будет заблокирован. Пользователь не сможет войти в систему.`}
        confirmLabel="Заблокировать"
        danger
        isPending={blockMutation.isPending}
      />

      {/* ── Удаление ── */}
      <ConfirmDialog
        open={deleteTarget !== null}
        onClose={() => setDeleteTarget(null)}
        onConfirm={() => { if (deleteTarget) void deleteMutation.mutateAsync(deleteTarget.id) }}
        title="Удалить учётную запись"
        message={`Вы уверены, что хотите удалить аккаунт «${deleteTarget?.login ?? ''}»? Это действие нельзя отменить.`}
        confirmLabel="Удалить"
        danger
        isPending={deleteMutation.isPending}
      />
    </>
  )
}

// ── AccountRow ────────────────────────────────────────────────────────────────

function AccountRow({
  account,
  activating,
  onChangeLogin,
  onChangePassword,
  onChangeRole,
  onBlock,
  onActivate,
  onDelete,
}: {
  account: AccountView
  activating: boolean
  onChangeLogin: () => void
  onChangePassword: () => void
  onChangeRole: () => void
  onBlock: () => void
  onActivate: () => void
  onDelete: () => void
}) {
  return (
    <tr>
      <td className={styles.loginCell}>{account.Login}</td>
      <td>
        <Badge variant={roleBadgeVariant(account.Role)}>
          {ROLE_LABELS[account.Role] ?? account.Role}
        </Badge>
      </td>
      <td>
        <Badge variant={statusBadgeVariant(account.Status)}>
          {STATUS_LABELS[account.Status] ?? account.Status}
        </Badge>
      </td>
      <td className={styles.personCell}>
        {account.PersonID ? (
          <Link to={`/people/${account.PersonID}`} className={styles.personLink}>
            {account.PersonID.slice(0, 8)}…
          </Link>
        ) : (
          <span className={styles.noValue}>—</span>
        )}
      </td>
      <td>
        <div className={styles.rowActions}>
          <Button variant="secondary" onClick={onChangeLogin} title="Изменить логин">
            <Pencil size={14} aria-hidden="true" />
          </Button>
          <Button variant="secondary" onClick={onChangePassword} title="Изменить пароль">
            <KeyRound size={14} aria-hidden="true" />
          </Button>
          <Button variant="secondary" onClick={onChangeRole} title="Изменить роль">
            <Shield size={14} aria-hidden="true" />
          </Button>
          {account.Status === 'active' ? (
            <Button variant="secondary" onClick={onBlock} title="Заблокировать">
              <Lock size={14} aria-hidden="true" />
            </Button>
          ) : account.Status === 'blocked' ? (
            <Button variant="secondary" onClick={onActivate} disabled={activating} title="Активировать">
              <LockOpen size={14} aria-hidden="true" />
            </Button>
          ) : null}
          <Button variant="secondary" onClick={onDelete} title="Удалить">
            <Trash2 size={14} aria-hidden="true" />
          </Button>
        </div>
      </td>
    </tr>
  )
}
