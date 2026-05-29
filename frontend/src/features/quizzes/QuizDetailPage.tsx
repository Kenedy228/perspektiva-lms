import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Check, Copy, Pencil, Plus, Trash2 } from 'lucide-react'
import {
  changeQuizLimits,
  changeQuizShuffle,
  deleteQuiz,
  getQuiz,
  renameQuiz,
  replaceQuizSources,
} from '../../api/quizzes'
import { listBanks } from '../../api/banks'
import { Button } from '../../components/ui/Button'
import { ConfirmDialog } from '../../components/ui/ConfirmDialog'
import { FormField } from '../../components/ui/FormField'
import { Modal } from '../../components/ui/Modal'
import { ApiError } from '../../types/api'
import type { QuizSourcePayload, QuizView } from '../../types/quizzes'
import styles from './QuizDetailPage.module.css'

type SourceDraft = {
  _k: string
  bank_id: string
  criteria_type: 'random' | 'manual'
  question_count: number
  question_ids: string
}

function uid() {
  return Math.random().toString(36).slice(2)
}

function viewToSourceDraft(s: QuizView['sources'][number]): SourceDraft {
  return {
    _k: uid(),
    bank_id: s.bank_id,
    criteria_type: s.criteria_type,
    question_count: s.question_count ?? 5,
    question_ids: s.question_ids?.join(', ') ?? '',
  }
}

function sourcesToPayload(sources: SourceDraft[]): QuizSourcePayload[] {
  return sources.map((s) => ({
    bank_id: s.bank_id,
    criteria_type: s.criteria_type,
    question_count: s.criteria_type === 'random' ? s.question_count : undefined,
    question_ids:
      s.criteria_type === 'manual'
        ? s.question_ids
            .split(',')
            .map((x) => x.trim())
            .filter(Boolean)
        : undefined,
  }))
}

export function QuizDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const quizQuery = useQuery({
    queryKey: ['quiz', id],
    queryFn: () => getQuiz(id!),
    enabled: Boolean(id),
  })
  const quiz = quizQuery.data

  const banksQuery = useQuery({
    queryKey: ['banks'],
    queryFn: () => listBanks({ limit: 200 }),
  })
  const banks = banksQuery.data ?? []
  const bankMap = Object.fromEntries(banks.map((b) => [b.ID, b.Title]))

  // Copy ID state
  const [copied, setCopied] = useState(false)

  // Modal state
  const [renameOpen, setRenameOpen] = useState(false)
  const [limitsOpen, setLimitsOpen] = useState(false)
  const [sourcesOpen, setSourcesOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [formError, setFormError] = useState<string | null>(null)

  // Form state
  const [newTitle, setNewTitle] = useState('')
  const [newMaxAttempts, setNewMaxAttempts] = useState(1)
  const [newTimeLimit, setNewTimeLimit] = useState(0)
  const [sourceDrafts, setSourceDrafts] = useState<SourceDraft[]>([])

  function inv() {
    return queryClient.invalidateQueries({ queryKey: ['quiz', id] })
  }

  const renameMut = useMutation({
    mutationFn: (title: string) => renameQuiz(id!, title),
    onSuccess: async () => { await inv(); setRenameOpen(false) },
    onError: (e) => setFormError(e instanceof ApiError ? e.message : 'Не удалось переименовать'),
  })

  const limitsMut = useMutation({
    mutationFn: () => changeQuizLimits(id!, newMaxAttempts, newTimeLimit),
    onSuccess: async () => { await inv(); setLimitsOpen(false) },
    onError: (e) => setFormError(e instanceof ApiError ? e.message : 'Не удалось сохранить'),
  })

  const shuffleMut = useMutation({
    mutationFn: (v: boolean) => changeQuizShuffle(id!, v),
    onSuccess: () => inv(),
  })

  const sourcesMut = useMutation({
    mutationFn: (s: SourceDraft[]) => replaceQuizSources(id!, sourcesToPayload(s)),
    onSuccess: async () => { await inv(); setSourcesOpen(false) },
    onError: (e) => setFormError(e instanceof ApiError ? e.message : 'Не удалось сохранить'),
  })

  const deleteMut = useMutation({
    mutationFn: () => deleteQuiz(id!),
    onSuccess: async () => {
      navigate('/quizzes')
    },
  })

  function copyId() {
    if (!id) return
    void navigator.clipboard.writeText(id)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  function openRename() {
    setNewTitle(quiz?.title ?? '')
    setFormError(null)
    setRenameOpen(true)
  }

  function openLimits() {
    setNewMaxAttempts(quiz?.max_attempts ?? 1)
    setNewTimeLimit(quiz?.time_limit_seconds ?? 0)
    setFormError(null)
    setLimitsOpen(true)
  }

  function openSources() {
    setSourceDrafts(quiz?.sources.map(viewToSourceDraft) ?? [])
    setFormError(null)
    setSourcesOpen(true)
  }

  function updateSource(k: string, patch: Partial<SourceDraft>) {
    setSourceDrafts((prev) => prev.map((s) => (s._k === k ? { ...s, ...patch } : s)))
  }

  function removeSource(k: string) {
    setSourceDrafts((prev) => prev.filter((s) => s._k !== k))
  }

  if (quizQuery.isPending) return <p className={styles.state}>Загрузка…</p>
  if (quizQuery.isError || !quiz) return <p className={styles.stateError}>Тест не найден.</p>

  return (
    <>
      <Link to="/quizzes" className={styles.backLink}>
        <ArrowLeft size={15} aria-hidden="true" />
        Тесты
      </Link>

      {/* Header */}
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <h1 className={styles.title}>{quiz.title}</h1>
          <div className={styles.idRow}>
            <span className={styles.idLabel}>ID:</span>
            <code className={styles.idValue}>{id}</code>
            <button type="button" className={styles.copyBtn} onClick={copyId} title="Скопировать ID">
              {copied ? <Check size={14} /> : <Copy size={14} />}
            </button>
            {copied && <span className={styles.copiedHint}>Скопировано!</span>}
          </div>
        </div>
        <div className={styles.headerActions}>
          <Button variant="secondary" onClick={openRename}>
            <Pencil size={15} /> Переименовать
          </Button>
          <Button variant="secondary" onClick={() => setDeleteOpen(true)}>
            <Trash2 size={15} /> Удалить
          </Button>
        </div>
      </div>

      {/* Settings */}
      <div className={styles.cards}>
        {/* Limits */}
        <section className={styles.card}>
          <div className={styles.cardHead}>
            <h2 className={styles.cardTitle}>Параметры</h2>
            <button type="button" className={styles.editBtn} onClick={openLimits}>
              <Pencil size={14} />
            </button>
          </div>
          <dl className={styles.dl}>
            <div>
              <dt>Попыток</dt>
              <dd>{quiz.max_attempts === 0 ? 'Без ограничений' : quiz.max_attempts}</dd>
            </div>
            <div>
              <dt>Лимит времени</dt>
              <dd>{quiz.time_limit_seconds === 0 ? 'Без лимита' : `${quiz.time_limit_seconds} сек`}</dd>
            </div>
            <div>
              <dt>Перемешивать вопросы</dt>
              <dd>
                <label className={styles.toggleLabel}>
                  <input
                    type="checkbox"
                    checked={quiz.shuffle_questions}
                    onChange={(e) => shuffleMut.mutate(e.target.checked)}
                    disabled={shuffleMut.isPending}
                  />
                  {quiz.shuffle_questions ? 'Да' : 'Нет'}
                </label>
              </dd>
            </div>
          </dl>
        </section>

        {/* Sources */}
        <section className={styles.card}>
          <div className={styles.cardHead}>
            <h2 className={styles.cardTitle}>Источники вопросов</h2>
            <button type="button" className={styles.editBtn} onClick={openSources}>
              <Pencil size={14} />
            </button>
          </div>
          {quiz.sources.length === 0 ? (
            <p className={styles.hint}>Источники не настроены.</p>
          ) : (
            <ul className={styles.sourceList}>
              {quiz.sources.map((s, i) => (
                <li key={i} className={styles.sourceItem}>
                  <span className={styles.sourceBank}>
                    {bankMap[s.bank_id] ?? <code>{s.bank_id}</code>}
                  </span>
                  <span className={styles.sourceMeta}>
                    {s.criteria_type === 'random'
                      ? `Случайные ${s.question_count} вопр.`
                      : `Конкретные: ${s.question_ids?.length ?? 0} вопр.`}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

      {/* ── Rename modal ── */}
      <Modal open={renameOpen} onClose={() => setRenameOpen(false)} title="Переименовать тест" size="sm">
        <div className={styles.form}>
          <FormField label="Название" htmlFor="quiz-rename" required>
            <input
              id="quiz-rename"
              className={styles.input}
              value={newTitle}
              onChange={(e) => setNewTitle(e.target.value)}
            />
          </FormField>
          {formError && <p className={styles.formError}>{formError}</p>}
          <div className={styles.formActions}>
            <Button variant="secondary" onClick={() => setRenameOpen(false)}>Отмена</Button>
            <Button
              disabled={!newTitle.trim() || renameMut.isPending}
              onClick={() => renameMut.mutate(newTitle.trim())}
            >
              {renameMut.isPending ? 'Сохранение…' : 'Сохранить'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* ── Limits modal ── */}
      <Modal open={limitsOpen} onClose={() => setLimitsOpen(false)} title="Параметры теста" size="sm">
        <div className={styles.form}>
          <FormField label="Макс. попыток (0 = без ограничений)" htmlFor="quiz-attempts">
            <input
              id="quiz-attempts"
              type="number"
              min={0}
              className={styles.input}
              value={newMaxAttempts}
              onChange={(e) => setNewMaxAttempts(Number(e.target.value))}
            />
          </FormField>
          <FormField label="Лимит времени, сек (0 = без лимита)" htmlFor="quiz-time">
            <input
              id="quiz-time"
              type="number"
              min={0}
              className={styles.input}
              value={newTimeLimit}
              onChange={(e) => setNewTimeLimit(Number(e.target.value))}
            />
          </FormField>
          {formError && <p className={styles.formError}>{formError}</p>}
          <div className={styles.formActions}>
            <Button variant="secondary" onClick={() => setLimitsOpen(false)}>Отмена</Button>
            <Button disabled={limitsMut.isPending} onClick={() => limitsMut.mutate()}>
              {limitsMut.isPending ? 'Сохранение…' : 'Сохранить'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* ── Sources modal ── */}
      <Modal open={sourcesOpen} onClose={() => setSourcesOpen(false)} title="Источники вопросов">
        <div className={styles.form}>
          <div className={styles.sourcesHeader}>
            <span className={styles.sourcesLabel}>Источники</span>
            <Button
              variant="secondary"
              onClick={() =>
                setSourceDrafts((prev) => [
                  ...prev,
                  { _k: uid(), bank_id: '', criteria_type: 'random', question_count: 5, question_ids: '' },
                ])
              }
            >
              <Plus size={13} /> Добавить
            </Button>
          </div>
          {sourceDrafts.map((s, i) => (
            <div key={s._k} className={styles.sourceRow}>
              <span className={styles.sourceNum}>{i + 1}.</span>
              <div className={styles.sourceFields}>
                <select
                  className={styles.select}
                  value={s.bank_id}
                  onChange={(e) => updateSource(s._k, { bank_id: e.target.value })}
                >
                  <option value="">— Банк —</option>
                  {banks.map((b) => (
                    <option key={b.ID} value={b.ID}>{b.Title}</option>
                  ))}
                </select>
                <select
                  className={styles.select}
                  value={s.criteria_type}
                  onChange={(e) =>
                    updateSource(s._k, { criteria_type: e.target.value as 'random' | 'manual' })
                  }
                >
                  <option value="random">Случайные N</option>
                  <option value="manual">Конкретные ID</option>
                </select>
                {s.criteria_type === 'random' ? (
                  <input
                    type="number"
                    min={1}
                    className={styles.input}
                    value={s.question_count}
                    onChange={(e) => updateSource(s._k, { question_count: Number(e.target.value) })}
                  />
                ) : (
                  <input
                    className={styles.input}
                    value={s.question_ids}
                    onChange={(e) => updateSource(s._k, { question_ids: e.target.value })}
                    placeholder="ID через запятую"
                  />
                )}
              </div>
              <button
                type="button"
                className={styles.removeBtn}
                onClick={() => removeSource(s._k)}
                aria-label="Удалить источник"
              >
                <Trash2 size={14} />
              </button>
            </div>
          ))}
          {formError && <p className={styles.formError}>{formError}</p>}
          <div className={styles.formActions}>
            <Button variant="secondary" onClick={() => setSourcesOpen(false)}>Отмена</Button>
            <Button
              disabled={sourceDrafts.length === 0 || sourcesMut.isPending}
              onClick={() => {
                setFormError(null)
                const invalid = sourceDrafts.find((s) => !s.bank_id)
                if (invalid) { setFormError('Выберите банк для каждого источника'); return }
                sourcesMut.mutate(sourceDrafts)
              }}
            >
              {sourcesMut.isPending ? 'Сохранение…' : 'Сохранить'}
            </Button>
          </div>
        </div>
      </Modal>

      {/* ── Delete confirm ── */}
      <ConfirmDialog
        open={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        onConfirm={() => deleteMut.mutate()}
        title="Удалить тест"
        message={`Тест «${quiz.title}» будет удалён безвозвратно.`}
        confirmLabel="Удалить"
        danger
        isPending={deleteMut.isPending}
      />
    </>
  )
}
