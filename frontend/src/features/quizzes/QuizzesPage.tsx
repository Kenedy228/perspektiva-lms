import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation, useQuery } from '@tanstack/react-query'
import { ClipboardList, Plus, Trash2 } from 'lucide-react'
import { createQuiz } from '../../api/quizzes'
import { listBanks } from '../../api/banks'
import { Button } from '../../components/ui/Button'
import { FormField } from '../../components/ui/FormField'
import { Modal } from '../../components/ui/Modal'
import { PageHeader } from '../../components/ui/PageHeader'
import { ApiError } from '../../types/api'
import type { CreateQuizPayload, QuizSourcePayload } from '../../types/quizzes'
import styles from './QuizzesPage.module.css'

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

function emptySource(): SourceDraft {
  return { _k: uid(), bank_id: '', criteria_type: 'random', question_count: 5, question_ids: '' }
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

export function QuizzesPage() {
  const navigate = useNavigate()
  const [open, setOpen] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const [title, setTitle] = useState('')
  const [maxAttempts, setMaxAttempts] = useState(1)
  const [timeLimit, setTimeLimit] = useState(0)
  const [shuffle, setShuffle] = useState(false)
  const [sources, setSources] = useState<SourceDraft[]>([emptySource()])

  const banksQuery = useQuery({
    queryKey: ['banks'],
    queryFn: () => listBanks({ limit: 200 }),
  })
  const banks = banksQuery.data ?? []

  const createMut = useMutation({
    mutationFn: (payload: CreateQuizPayload) => createQuiz(payload),
    onSuccess: (data) => {
      setOpen(false)
      navigate(`/quizzes/${data.id}`)
    },
    onError: (e) => setError(e instanceof ApiError ? e.message : 'Не удалось создать тест'),
  })

  function handleCreate() {
    setError(null)
    const invalidSrc = sources.find((s) => !s.bank_id)
    if (invalidSrc) {
      setError('Выберите банк для каждого источника')
      return
    }
    createMut.mutate({
      title: title.trim(),
      max_attempts: maxAttempts,
      time_limit_seconds: timeLimit,
      shuffle_questions: shuffle,
      sources: sourcesToPayload(sources),
    })
  }

  function openCreate() {
    setTitle('')
    setMaxAttempts(1)
    setTimeLimit(0)
    setShuffle(false)
    setSources([emptySource()])
    setError(null)
    setOpen(true)
  }

  function updateSource(k: string, patch: Partial<SourceDraft>) {
    setSources((prev) => prev.map((s) => (s._k === k ? { ...s, ...patch } : s)))
  }

  function removeSource(k: string) {
    setSources((prev) => prev.filter((s) => s._k !== k))
  }

  return (
    <>
      <PageHeader
        title="Тесты"
        description="Управление тестами. Скопируйте ID теста для привязки к элементу курса."
      />

      <div className={styles.cta}>
        <div className={styles.ctaIcon}>
          <ClipboardList size={40} />
        </div>
        <p className={styles.ctaText}>
          Создайте тест, настройте его параметры и источники вопросов из банков. После создания
          скопируйте ID теста и укажите его при добавлении элемента типа «тест» в курс.
        </p>
        <Button onClick={openCreate}>
          <Plus size={15} /> Создать тест
        </Button>
      </div>

      {/* ── Create quiz modal ── */}
      <Modal open={open} onClose={() => setOpen(false)} title="Создать тест">
        <div className={styles.form}>
          <FormField label="Название теста" htmlFor="quiz-title" required>
            <input
              id="quiz-title"
              className={styles.input}
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Введите название"
            />
          </FormField>

          <div className={styles.row}>
            <FormField label="Попыток (0 = без ограничений)" htmlFor="quiz-attempts">
              <input
                id="quiz-attempts"
                type="number"
                min={0}
                className={styles.input}
                value={maxAttempts}
                onChange={(e) => setMaxAttempts(Number(e.target.value))}
              />
            </FormField>
            <FormField label="Лимит времени, сек (0 = без лимита)" htmlFor="quiz-time">
              <input
                id="quiz-time"
                type="number"
                min={0}
                className={styles.input}
                value={timeLimit}
                onChange={(e) => setTimeLimit(Number(e.target.value))}
              />
            </FormField>
          </div>

          <label className={styles.checkLabel}>
            <input type="checkbox" checked={shuffle} onChange={(e) => setShuffle(e.target.checked)} />
            Перемешивать вопросы
          </label>

          <div className={styles.sourcesSection}>
            <div className={styles.sourcesHeader}>
              <span className={styles.sourcesLabel}>Источники вопросов</span>
              <Button
                variant="secondary"
                onClick={() => setSources((prev) => [...prev, emptySource()])}
              >
                <Plus size={13} /> Добавить источник
              </Button>
            </div>

            {sources.map((s, i) => (
              <div key={s._k} className={styles.sourceRow}>
                <span className={styles.sourceNum}>{i + 1}.</span>
                <div className={styles.sourceFields}>
                  <select
                    className={styles.select}
                    value={s.bank_id}
                    onChange={(e) => updateSource(s._k, { bank_id: e.target.value })}
                  >
                    <option value="">— Выберите банк —</option>
                    {banks.map((b) => (
                      <option key={b.ID} value={b.ID}>
                        {b.Title} ({b.QuestionsCount} вопр.)
                      </option>
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
                      className={`${styles.input} ${styles.countInput}`}
                      value={s.question_count}
                      onChange={(e) => updateSource(s._k, { question_count: Number(e.target.value) })}
                      placeholder="Количество"
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

            {sources.length === 0 && (
              <p className={styles.hint}>Нет источников. Тест без источников нельзя создать.</p>
            )}
          </div>

          {error && <p className={styles.error}>{error}</p>}

          <div className={styles.actions}>
            <Button variant="secondary" onClick={() => setOpen(false)}>
              Отмена
            </Button>
            <Button
              disabled={!title.trim() || sources.length === 0 || createMut.isPending}
              onClick={handleCreate}
            >
              {createMut.isPending ? 'Создание…' : 'Создать тест'}
            </Button>
          </div>
        </div>
      </Modal>
    </>
  )
}
