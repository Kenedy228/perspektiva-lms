import { useMutation, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Pencil, Trash2, UserRound } from 'lucide-react'
import { FormEvent, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import {
  assignPersonOrganization,
  detachPersonProfile,
  deletePerson,
  removePersonOrganization,
  renamePerson,
  replacePersonProfile,
} from '../../api/people'
import { Badge } from '../../components/ui/Badge'
import { Button } from '../../components/ui/Button'
import { ConfirmDialog } from '../../components/ui/ConfirmDialog'
import { FormField } from '../../components/ui/FormField'
import { Modal } from '../../components/ui/Modal'
import { ApiError } from '../../types/api'
import type { ProfileRequest } from '../../types/person'
import { usePerson } from './usePerson'
import styles from './PersonDetailPage.module.css'

type RenameForm = { firstName: string; lastName: string; middleName: string }
type ProfileForm = ProfileRequest

function defaultProfile(snils = '', dob = '', jobTitle = '', education = '', orgId = ''): ProfileForm {
  return { snils, date_of_birth: dob, job_title: jobTitle, education, organization_id: orgId }
}

export function PersonDetailPage() {
  const { id = '' } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data: person, isLoading, error } = usePerson(id)

  const [renameOpen, setRenameOpen] = useState(false)
  const [renameForm, setRenameForm] = useState<RenameForm>({ firstName: '', lastName: '', middleName: '' })
  const [renameError, setRenameError] = useState<string | null>(null)

  const [profileOpen, setProfileOpen] = useState(false)
  const [profileForm, setProfileForm] = useState<ProfileForm>(defaultProfile())
  const [profileError, setProfileError] = useState<string | null>(null)

  const [assignOrgOpen, setAssignOrgOpen] = useState(false)
  const [orgId, setOrgId] = useState('')
  const [assignOrgError, setAssignOrgError] = useState<string | null>(null)

  const [detachProfileConfirm, setDetachProfileConfirm] = useState(false)
  const [removeOrgConfirm, setRemoveOrgConfirm] = useState(false)
  const [deleteConfirm, setDeleteConfirm] = useState(false)

  function invalidate() {
    void queryClient.invalidateQueries({ queryKey: ['person', id] })
    void queryClient.invalidateQueries({ queryKey: ['people'] })
  }

  const renameMutation = useMutation({
    mutationFn: (form: RenameForm) =>
      renamePerson(id, {
        first_name: form.firstName.trim(),
        last_name: form.lastName.trim(),
        middle_name: form.middleName.trim() || undefined,
      }),
    onSuccess: () => { invalidate(); setRenameOpen(false); setRenameError(null) },
  })

  const profileMutation = useMutation({
    mutationFn: (form: ProfileForm) => replacePersonProfile(id, form),
    onSuccess: () => { invalidate(); setProfileOpen(false); setProfileError(null) },
  })

  const detachProfileMutation = useMutation({
    mutationFn: () => detachPersonProfile(id),
    onSuccess: () => { invalidate(); setDetachProfileConfirm(false) },
  })

  const assignOrgMutation = useMutation({
    mutationFn: (orgID: string) => assignPersonOrganization(id, orgID),
    onSuccess: () => { invalidate(); setAssignOrgOpen(false); setOrgId(''); setAssignOrgError(null) },
  })

  const removeOrgMutation = useMutation({
    mutationFn: () => removePersonOrganization(id),
    onSuccess: () => { invalidate(); setRemoveOrgConfirm(false) },
  })

  const deleteMutation = useMutation({
    mutationFn: () => deletePerson(id),
    onSuccess: () => { void queryClient.invalidateQueries({ queryKey: ['people'] }); navigate('/people', { replace: true }) },
  })

  function extractMessage(caught: unknown, fallback: string) {
    return caught instanceof ApiError ? caught.message : fallback
  }

  async function handleRename(event: FormEvent) {
    event.preventDefault()
    setRenameError(null)
    try {
      await renameMutation.mutateAsync(renameForm)
    } catch (caught) {
      setRenameError(extractMessage(caught, 'Не удалось изменить ФИО'))
    }
  }

  async function handleProfile(event: FormEvent) {
    event.preventDefault()
    setProfileError(null)
    try {
      await profileMutation.mutateAsync(profileForm)
    } catch (caught) {
      setProfileError(extractMessage(caught, 'Не удалось сохранить профиль'))
    }
  }

  async function handleAssignOrg(event: FormEvent) {
    event.preventDefault()
    setAssignOrgError(null)
    try {
      await assignOrgMutation.mutateAsync(orgId.trim())
    } catch (caught) {
      setAssignOrgError(extractMessage(caught, 'Не удалось назначить организацию'))
    }
  }

  const hasProfile = Boolean(person?.Snils)

  if (isLoading) {
    return <p className={styles.state}>Загрузка…</p>
  }

  if (error || !person) {
    return (
      <div className={styles.state}>
        <p className={styles.stateError}>
          {error instanceof ApiError ? error.message : 'Сотрудник не найден'}
        </p>
        <Button as={Link} to="/people" variant="secondary">
          <ArrowLeft size={16} aria-hidden="true" />
          К списку
        </Button>
      </div>
    )
  }

  const fullName = [person.LastName, person.FirstName, person.MiddleName].filter(Boolean).join(' ')

  return (
    <>
      <div className={styles.backRow}>
        <Button as={Link} to="/people" variant="secondary">
          <ArrowLeft size={16} aria-hidden="true" />
          Сотрудники
        </Button>
      </div>

      <div className={styles.header}>
        <div className={styles.avatar}>
          <UserRound size={32} aria-hidden="true" />
        </div>
        <div>
          <h1 className={styles.name}>{fullName}</h1>
          {person.OrganizationName ? (
            <p className={styles.orgName}>{person.OrganizationName}</p>
          ) : null}
        </div>
        <div className={styles.headerActions}>
          <Button
            variant="secondary"
            onClick={() => { setRenameForm({ firstName: person.FirstName, lastName: person.LastName, middleName: person.MiddleName }); setRenameOpen(true) }}
          >
            <Pencil size={15} aria-hidden="true" />
            Изменить ФИО
          </Button>
          <Button
            variant="secondary"
            onClick={() => setDeleteConfirm(true)}
          >
            <Trash2 size={15} aria-hidden="true" />
            Удалить
          </Button>
        </div>
      </div>

      {/* ── Профиль ── */}
      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2 className={styles.sectionTitle}>Профиль</h2>
          <div className={styles.sectionActions}>
            <Button
              variant="secondary"
              onClick={() => {
                setProfileForm(defaultProfile(person.Snils, person.DateOfBirth, person.JobTitle, person.Education))
                setProfileOpen(true)
              }}
            >
              <Pencil size={15} aria-hidden="true" />
              {hasProfile ? 'Изменить' : 'Добавить профиль'}
            </Button>
            {hasProfile ? (
              <Button variant="secondary" onClick={() => setDetachProfileConfirm(true)}>
                <Trash2 size={15} aria-hidden="true" />
                Открепить
              </Button>
            ) : null}
          </div>
        </div>

        {hasProfile ? (
          <dl className={styles.grid}>
            <ProfileRow label="СНИЛС" value={person.Snils} />
            <ProfileRow label="Дата рождения" value={person.DateOfBirth} />
            <ProfileRow label="Должность" value={person.JobTitle} />
            <ProfileRow label="Образование" value={person.Education} />
          </dl>
        ) : (
          <p className={styles.noProfile}>Профиль не добавлен</p>
        )}
      </section>

      {/* ── Организация ── */}
      <section className={styles.section}>
        <div className={styles.sectionHead}>
          <h2 className={styles.sectionTitle}>Организация</h2>
          <div className={styles.sectionActions}>
            <Button variant="secondary" onClick={() => { setOrgId(''); setAssignOrgOpen(true) }}>
              <Pencil size={15} aria-hidden="true" />
              {person.OrganizationName ? 'Изменить' : 'Назначить'}
            </Button>
            {person.OrganizationName ? (
              <Button variant="secondary" onClick={() => setRemoveOrgConfirm(true)}>
                <Trash2 size={15} aria-hidden="true" />
                Открепить
              </Button>
            ) : null}
          </div>
        </div>

        {person.OrganizationName ? (
          <p className={styles.orgValue}>{person.OrganizationName}</p>
        ) : (
          <p className={styles.noProfile}>Организация не назначена</p>
        )}
      </section>

      {/* ── Статус профиля ── */}
      <div className={styles.statusRow}>
        <Badge variant={hasProfile ? 'success' : 'neutral'}>
          {hasProfile ? 'Профиль заполнен' : 'Без профиля'}
        </Badge>
      </div>

      {/* ── Модалы ── */}
      <Modal open={renameOpen} onClose={() => setRenameOpen(false)} title="Изменить ФИО" size="sm">
        <form className={styles.form} onSubmit={handleRename} noValidate>
          <FormField label="Фамилия" htmlFor="r-last" required>
            <input id="r-last" value={renameForm.lastName} onChange={(e) => setRenameForm((f) => ({ ...f, lastName: e.target.value }))} required autoFocus />
          </FormField>
          <FormField label="Имя" htmlFor="r-first" required>
            <input id="r-first" value={renameForm.firstName} onChange={(e) => setRenameForm((f) => ({ ...f, firstName: e.target.value }))} required />
          </FormField>
          <FormField label="Отчество" htmlFor="r-middle">
            <input id="r-middle" value={renameForm.middleName} onChange={(e) => setRenameForm((f) => ({ ...f, middleName: e.target.value }))} />
          </FormField>
          {renameError ? <p className={styles.formError}>{renameError}</p> : null}
          <div className={styles.formActions}>
            <Button variant="secondary" type="button" onClick={() => setRenameOpen(false)} disabled={renameMutation.isPending}>Отмена</Button>
            <Button type="submit" disabled={renameMutation.isPending}>{renameMutation.isPending ? 'Сохранение…' : 'Сохранить'}</Button>
          </div>
        </form>
      </Modal>

      <Modal open={profileOpen} onClose={() => setProfileOpen(false)} title={hasProfile ? 'Изменить профиль' : 'Добавить профиль'}>
        <form className={styles.form} onSubmit={handleProfile} noValidate>
          <FormField label="СНИЛС" htmlFor="p-snils" required>
            <input id="p-snils" value={profileForm.snils} onChange={(e) => setProfileForm((f) => ({ ...f, snils: e.target.value }))} placeholder="000-000-000 00" required autoFocus />
          </FormField>
          <FormField label="Дата рождения" htmlFor="p-dob" required>
            <input id="p-dob" type="date" value={profileForm.date_of_birth} onChange={(e) => setProfileForm((f) => ({ ...f, date_of_birth: e.target.value }))} required />
          </FormField>
          <FormField label="Должность" htmlFor="p-job">
            <input id="p-job" value={profileForm.job_title ?? ''} onChange={(e) => setProfileForm((f) => ({ ...f, job_title: e.target.value }))} placeholder="Необязательно" />
          </FormField>
          <FormField label="Образование" htmlFor="p-edu">
            <input id="p-edu" value={profileForm.education ?? ''} onChange={(e) => setProfileForm((f) => ({ ...f, education: e.target.value }))} placeholder="Необязательно" />
          </FormField>
          <FormField label="ID организации" htmlFor="p-org">
            <input id="p-org" value={profileForm.organization_id ?? ''} onChange={(e) => setProfileForm((f) => ({ ...f, organization_id: e.target.value }))} placeholder="UUID организации (необязательно)" />
          </FormField>
          {profileError ? <p className={styles.formError}>{profileError}</p> : null}
          <div className={styles.formActions}>
            <Button variant="secondary" type="button" onClick={() => setProfileOpen(false)} disabled={profileMutation.isPending}>Отмена</Button>
            <Button type="submit" disabled={profileMutation.isPending}>{profileMutation.isPending ? 'Сохранение…' : 'Сохранить'}</Button>
          </div>
        </form>
      </Modal>

      <Modal open={assignOrgOpen} onClose={() => setAssignOrgOpen(false)} title="Назначить организацию" size="sm">
        <form className={styles.form} onSubmit={handleAssignOrg} noValidate>
          <FormField label="ID организации" htmlFor="org-id" required>
            <input id="org-id" value={orgId} onChange={(e) => setOrgId(e.target.value)} placeholder="UUID из раздела Организации" required autoFocus />
          </FormField>
          {assignOrgError ? <p className={styles.formError}>{assignOrgError}</p> : null}
          <div className={styles.formActions}>
            <Button variant="secondary" type="button" onClick={() => setAssignOrgOpen(false)} disabled={assignOrgMutation.isPending}>Отмена</Button>
            <Button type="submit" disabled={assignOrgMutation.isPending}>{assignOrgMutation.isPending ? 'Назначение…' : 'Назначить'}</Button>
          </div>
        </form>
      </Modal>

      <ConfirmDialog
        open={detachProfileConfirm}
        onClose={() => setDetachProfileConfirm(false)}
        onConfirm={() => void detachProfileMutation.mutateAsync()}
        title="Открепить профиль"
        message="Данные профиля (СНИЛС, дата рождения, должность) будут удалены. Продолжить?"
        confirmLabel="Открепить"
        danger
        isPending={detachProfileMutation.isPending}
      />

      <ConfirmDialog
        open={removeOrgConfirm}
        onClose={() => setRemoveOrgConfirm(false)}
        onConfirm={() => void removeOrgMutation.mutateAsync()}
        title="Открепить организацию"
        message={`Сотрудник будет откреплён от организации «${person.OrganizationName}». Продолжить?`}
        confirmLabel="Открепить"
        danger
        isPending={removeOrgMutation.isPending}
      />

      <ConfirmDialog
        open={deleteConfirm}
        onClose={() => setDeleteConfirm(false)}
        onConfirm={() => void deleteMutation.mutateAsync()}
        title="Удалить сотрудника"
        message={`Вы уверены, что хотите удалить «${fullName}»? Это действие нельзя отменить.`}
        confirmLabel="Удалить"
        danger
        isPending={deleteMutation.isPending}
      />
    </>
  )
}

function ProfileRow({ label, value }: { label: string; value: string }) {
  return (
    <>
      <dt className={styles.dt}>{label}</dt>
      <dd className={styles.dd}>{value || <span className={styles.empty}>—</span>}</dd>
    </>
  )
}
