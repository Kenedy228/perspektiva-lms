import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, ArrowDown, ArrowUp, CheckCircle, Clock, XCircle } from 'lucide-react'
import { addAttemptAnswer, cancelAttempt, finishAttempt, getAttempt } from '../../api/attempts'
import { Badge } from '../../components/ui/Badge'
import { Button } from '../../components/ui/Button'
import { ConfirmDialog } from '../../components/ui/ConfirmDialog'
import { ApiError } from '../../types/api'
import type { AnswerPayload, AttemptQuestion, AttemptStatus, MatchingPair, QuestionOption } from '../../types/attempt'
import { ATTEMPT_STATUS_LABELS } from '../../types/attempt'
import styles from './AttemptPage.module.css'

function statusBadgeVariant(status: AttemptStatus): 'success' | 'warning' | 'neutral' | 'danger' {
  if (status === 'finished') return 'success'
  if (status === 'in_progress') return 'warning'
  if (status === 'expired') return 'danger'
  return 'neutral'
}

function scoreLabel(score: number | undefined, answered: boolean): string {
  if (!answered) return 'Без ответа'
  if (score === undefined) return 'Отвечен'
  if (score >= 1) return 'Верно'
  if (score > 0) return `Частично (${Math.round(score * 100)}%)`
  return 'Неверно'
}

function scoreBadgeVariant(score: number | undefined, answered: boolean): 'success' | 'warning' | 'danger' | 'neutral' {
  if (!answered) return 'neutral'
  if (score === undefined) return 'success'
  if (score >= 1) return 'success'
  if (score > 0) return 'warning'
  return 'danger'
}

function isValidDate(dateStr: string) {
  if (!dateStr) return false
  return new Date(dateStr).getFullYear() > 2000
}

function formatDateTime(dateStr: string) {
  return new Date(dateStr).toLocaleString('ru-RU')
}

function formatCountdown(ms: number) {
  if (ms <= 0) return '0:00'
  const totalSec = Math.floor(ms / 1000)
  const h = Math.floor(totalSec / 3600)
  const m = Math.floor((totalSec % 3600) / 60)
  const s = totalSec % 60
  if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  return `${m}:${String(s).padStart(2, '0')}`
}

function shortId(id: string) {
  return id.length > 8 ? `${id.slice(0, 8)}…` : id
}

// ── Answer draft types ──────────────────────────────────────────────────────

type SelectableDraft = { type: 'selectable'; option_ids: Set<string> }
type SequenceDraft = { type: 'sequence'; ordered: QuestionOption[] }
type MatchingDraft = { type: 'matching'; pairs: Record<string, string> }
type ShortDraft = { type: 'short'; text: string }
type AnswerDraft = SelectableDraft | SequenceDraft | MatchingDraft | ShortDraft

function draftToPayload(draft: AnswerDraft): AnswerPayload {
  switch (draft.type) {
    case 'selectable':
      return { type: 'selectable', option_ids: Array.from(draft.option_ids) }
    case 'sequence':
      return { type: 'sequence', option_ids: draft.ordered.map((o) => o.id) }
    case 'matching':
      return { type: 'matching', matching_pairs: draft.pairs }
    case 'short':
      return { type: 'short', short_input: draft.text }
  }
}

function initDraft(q: AttemptQuestion): AnswerDraft {
  if (q.type === 'selectable') return { type: 'selectable', option_ids: new Set() }
  if (q.type === 'sequence') return { type: 'sequence', ordered: q.options ?? [] }
  if (q.type === 'matching') return { type: 'matching', pairs: {} }
  return { type: 'short', text: '' }
}

// ── Question answer components ──────────────────────────────────────────────

function SelectableAnswerForm({
  options,
  draft,
  onChange,
}: {
  options: QuestionOption[]
  draft: SelectableDraft
  onChange: (d: SelectableDraft) => void
}) {
  function toggle(id: string) {
    const next = new Set(draft.option_ids)
    if (next.has(id)) next.delete(id)
    else next.add(id)
    onChange({ ...draft, option_ids: next })
  }
  return (
    <ul className={styles.optionList}>
      {options.map((opt) => (
        <li key={opt.id}>
          <label className={styles.optionLabel}>
            <input
              type="checkbox"
              checked={draft.option_ids.has(opt.id)}
              onChange={() => toggle(opt.id)}
            />
            <span>{opt.text}</span>
          </label>
        </li>
      ))}
    </ul>
  )
}

function SequenceAnswerForm({
  draft,
  onChange,
}: {
  draft: SequenceDraft
  onChange: (d: SequenceDraft) => void
}) {
  function move(from: number, to: number) {
    const next = [...draft.ordered]
    const [item] = next.splice(from, 1)
    next.splice(to, 0, item)
    onChange({ ...draft, ordered: next })
  }
  return (
    <ul className={styles.sequenceList}>
      {draft.ordered.map((opt, i) => (
        <li key={opt.id} className={styles.sequenceItem}>
          <span className={styles.sequencePos}>{i + 1}</span>
          <span className={styles.sequenceText}>{opt.text}</span>
          <div className={styles.sequenceButtons}>
            <button
              type="button"
              className={styles.seqBtn}
              onClick={() => move(i, i - 1)}
              disabled={i === 0}
              aria-label="Переместить выше"
            >
              <ArrowUp size={14} />
            </button>
            <button
              type="button"
              className={styles.seqBtn}
              onClick={() => move(i, i + 1)}
              disabled={i === draft.ordered.length - 1}
              aria-label="Переместить ниже"
            >
              <ArrowDown size={14} />
            </button>
          </div>
        </li>
      ))}
    </ul>
  )
}

function MatchingAnswerForm({
  pairs,
  draft,
  onChange,
}: {
  pairs: MatchingPair[]
  draft: MatchingDraft
  onChange: (d: MatchingDraft) => void
}) {
  const matches = pairs.map((p) => ({ id: p.match_id, text: p.match }))
  return (
    <table className={styles.matchTable}>
      <thead>
        <tr>
          <th>Понятие</th>
          <th>Соответствие</th>
        </tr>
      </thead>
      <tbody>
        {pairs.map((p) => (
          <tr key={p.prompt_id}>
            <td className={styles.matchPrompt}>{p.prompt}</td>
            <td>
              <select
                className={styles.matchSelect}
                value={draft.pairs[p.prompt_id] ?? ''}
                onChange={(e) =>
                  onChange({ ...draft, pairs: { ...draft.pairs, [p.prompt_id]: e.target.value } })
                }
              >
                <option value="">— выберите —</option>
                {matches.map((m) => (
                  <option key={m.id} value={m.id}>{m.text}</option>
                ))}
              </select>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

function ShortAnswerForm({
  draft,
  onChange,
}: {
  draft: ShortDraft
  onChange: (d: ShortDraft) => void
}) {
  return (
    <input
      className={styles.shortInput}
      type="text"
      value={draft.text}
      onChange={(e) => onChange({ ...draft, text: e.target.value })}
      placeholder="Введите ответ"
    />
  )
}

// ── Question card ────────────────────────────────────────────────────────────

function QuestionCard({
  question,
  index,
  attemptId,
  answered,
  onAnswered,
  active,
  isCompleted,
  completedScore,
}: {
  question: AttemptQuestion
  index: number
  attemptId: string
  answered: boolean
  onAnswered: () => void
  active: boolean
  isCompleted: boolean
  completedScore?: number
}) {
  const [draft, setDraft] = useState<AnswerDraft>(() => initDraft(question))
  const [saveError, setSaveError] = useState<string | null>(null)

  const saveMutation = useMutation({
    mutationFn: () => addAttemptAnswer(attemptId, question.id, draftToPayload(draft)),
    onSuccess: () => { setSaveError(null); onAnswered() },
    onError: (err) => {
      setSaveError(err instanceof ApiError ? err.message : 'Не удалось сохранить ответ')
    },
  })

  const resultVariant = scoreBadgeVariant(completedScore, answered)

  return (
    <div className={`${styles.questionCard} ${answered ? styles.questionAnswered : ''}`}>
      <div className={styles.questionHeader}>
        <span className={styles.questionNum}>{index + 1}</span>
        <div className={styles.questionMeta}>
          <p className={styles.questionTitle}>{question.title}</p>
          <p className={styles.questionInstruction}>{question.instruction}</p>
        </div>
        {isCompleted ? (
          <Badge variant={resultVariant}>
            {scoreLabel(completedScore, answered)}
          </Badge>
        ) : answered ? (
          <Badge variant="success">Отвечен</Badge>
        ) : null}
      </div>

      {active && !isCompleted ? (
        <div className={styles.questionBody}>
          {question.type === 'selectable' && question.options ? (
            <SelectableAnswerForm
              options={question.options}
              draft={draft as SelectableDraft}
              onChange={setDraft}
            />
          ) : question.type === 'sequence' && question.options ? (
            <SequenceAnswerForm
              draft={draft as SequenceDraft}
              onChange={setDraft}
            />
          ) : question.type === 'matching' && question.pairs ? (
            <MatchingAnswerForm
              pairs={question.pairs}
              draft={draft as MatchingDraft}
              onChange={setDraft}
            />
          ) : question.type === 'short' ? (
            <ShortAnswerForm
              draft={draft as ShortDraft}
              onChange={setDraft}
            />
          ) : null}

          {saveError ? <p className={styles.saveError}>{saveError}</p> : null}

          <div className={styles.questionActions}>
            <Button
              onClick={() => void saveMutation.mutateAsync()}
              disabled={saveMutation.isPending}
            >
              {saveMutation.isPending ? 'Сохранение…' : 'Сохранить ответ'}
            </Button>
          </div>
        </div>
      ) : active && isCompleted ? (
        <div className={styles.questionBody}>
          <p className={styles.completedHint}>
            {answered
              ? 'Ответ был записан до завершения теста.'
              : 'Ответ не был дан до завершения теста.'}
          </p>
        </div>
      ) : null}
    </div>
  )
}

// ── Main page ────────────────────────────────────────────────────────────────

export function AttemptPage() {
  const { id = '' } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const [finishConfirm, setFinishConfirm] = useState(false)
  const [cancelConfirm, setCancelConfirm] = useState(false)
  const [actionError, setActionError] = useState<string | null>(null)
  const [now, setNow] = useState(() => Date.now())
  const [openQuestion, setOpenQuestion] = useState<string | null>(null)

  const { data: attempt, isLoading, error } = useQuery({
    queryKey: ['attempt', id],
    queryFn: () => getAttempt(id),
    refetchInterval: (query) => query.state.data?.status === 'in_progress' ? 30_000 : false,
  })

  useEffect(() => {
    if (attempt?.status !== 'in_progress') return
    if (!isValidDate(attempt.deadline_at)) return
    const timer = setInterval(() => setNow(Date.now()), 1000)
    return () => clearInterval(timer)
  }, [attempt?.status, attempt?.deadline_at])

  function invalidate() {
    void queryClient.invalidateQueries({ queryKey: ['attempt', id] })
  }

  const finishMutation = useMutation({
    mutationFn: () => finishAttempt(id),
    onSuccess: () => { invalidate(); setFinishConfirm(false); setActionError(null) },
    onError: (err) => {
      setActionError(err instanceof ApiError ? err.message : 'Не удалось завершить попытку')
      setFinishConfirm(false)
    },
  })

  const cancelMutation = useMutation({
    mutationFn: () => cancelAttempt(id),
    onSuccess: () => { invalidate(); setCancelConfirm(false); setActionError(null) },
    onError: (err) => {
      setActionError(err instanceof ApiError ? err.message : 'Не удалось отменить попытку')
      setCancelConfirm(false)
    },
  })

  if (isLoading) return <p className={styles.state}>Загрузка…</p>

  if (error || !attempt) {
    return (
      <div className={styles.state}>
        <p className={styles.stateError}>
          {error instanceof ApiError ? error.message : 'Попытка не найдена'}
        </p>
        <Button variant="secondary" onClick={() => navigate(-1)}>
          <ArrowLeft size={16} aria-hidden="true" />
          Назад
        </Button>
      </div>
    )
  }

  const isActive = attempt.status === 'in_progress'
  const isCompleted = !isActive
  const hasDeadline = isValidDate(attempt.deadline_at)
  const deadlineMs = hasDeadline ? new Date(attempt.deadline_at).getTime() - now : 0
  const isOverdue = hasDeadline && deadlineMs <= 0
  const answeredSet = new Set(attempt.answered_question_ids)
  const hasScores = attempt.question_scores != null
  const totalScore = attempt.total_score ?? 0

  return (
    <>
      <div className={styles.backRow}>
        <Button variant="secondary" onClick={() => navigate(-1)}>
          <ArrowLeft size={16} aria-hidden="true" />
          Назад
        </Button>
      </div>

      <div className={styles.header}>
        <div>
          <h1 className={styles.title}>Попытка</h1>
          <p className={styles.subtitle}>ID: {shortId(attempt.id)}</p>
        </div>
        <Badge variant={statusBadgeVariant(attempt.status)}>
          {ATTEMPT_STATUS_LABELS[attempt.status] ?? attempt.status}
        </Badge>
      </div>

      {actionError ? <p className={styles.actionError}>{actionError}</p> : null}

      {/* ── Прогресс ── */}
      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Прогресс</h2>
        <div className={styles.progressRow}>
          <div className={styles.progressTrack}>
            <div
              className={styles.progressFill}
              style={{
                width: attempt.questions_count > 0
                  ? `${Math.round((attempt.answers_count / attempt.questions_count) * 100)}%`
                  : '0%',
              }}
            />
          </div>
          <span className={styles.progressText}>
            {attempt.answers_count} / {attempt.questions_count}
          </span>
        </div>
      </section>

      {/* ── Таймер ── */}
      {isActive && hasDeadline ? (
        <section className={styles.section}>
          <div className={isOverdue ? styles.timerExpired : styles.timer}>
            <Clock size={18} aria-hidden="true" />
            {isOverdue ? 'Время истекло' : `Осталось: ${formatCountdown(deadlineMs)}`}
          </div>
        </section>
      ) : null}

      {/* ── Детали (collapsed) ── */}
      <section className={styles.section}>
        <dl className={styles.details}>
          <dt>Начато</dt>
          <dd>{formatDateTime(attempt.started_at)}</dd>
          {hasDeadline ? (<><dt>Дедлайн</dt><dd>{formatDateTime(attempt.deadline_at)}</dd></>) : null}
          {isValidDate(attempt.finished_at) ? (<><dt>Завершено</dt><dd>{formatDateTime(attempt.finished_at)}</dd></>) : null}
        </dl>
      </section>

      {/* ── Итоговый результат ── */}
      {isCompleted && hasScores ? (
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>Итоговый результат</h2>
          <div className={styles.resultBlock}>
            <div className={styles.progressRow}>
              <div className={styles.progressTrack}>
                <div
                  className={styles.progressFill}
                  style={{ width: `${Math.round(totalScore * 100)}%` }}
                />
              </div>
              <span className={styles.progressText}>{Math.round(totalScore * 100)}%</span>
            </div>
            <p className={styles.progressHint}>
              {(totalScore * attempt.questions_count).toFixed(1)} из {attempt.questions_count} баллов
            </p>
          </div>
        </section>
      ) : null}

      {/* ── Вопросы ── */}
      {attempt.questions.length > 0 ? (
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>
            {isCompleted ? 'Результаты по вопросам' : 'Вопросы'}
          </h2>
          <div className={styles.questions}>
            {attempt.questions.map((q, i) => {
              const answered = answeredSet.has(q.id)
              const qScore = attempt.question_scores?.[q.id]
              return (
                <div key={q.id}>
                  <button
                    type="button"
                    className={styles.questionToggle}
                    onClick={() => setOpenQuestion(openQuestion === q.id ? null : q.id)}
                  >
                    <span className={styles.questionToggleNum}>{i + 1}</span>
                    <span className={styles.questionToggleTitle}>{q.title}</span>
                    {isCompleted ? (
                      <Badge variant={scoreBadgeVariant(qScore, answered)}>
                        {scoreLabel(qScore, answered)}
                      </Badge>
                    ) : answered ? (
                      <Badge variant="success">✓</Badge>
                    ) : null}
                    <span className={styles.questionToggleArrow}>
                      {openQuestion === q.id ? '▲' : '▼'}
                    </span>
                  </button>
                  {openQuestion === q.id ? (
                    <QuestionCard
                      question={q}
                      index={i}
                      attemptId={id}
                      answered={answered}
                      onAnswered={invalidate}
                      active
                      isCompleted={isCompleted}
                      completedScore={qScore}
                    />
                  ) : null}
                </div>
              )
            })}
          </div>
        </section>
      ) : null}

      {/* ── Действия ── */}
      {isActive ? (
        <div className={styles.actions}>
          <Button onClick={() => setFinishConfirm(true)}>
            <CheckCircle size={16} aria-hidden="true" />
            Завершить тест
          </Button>
          <Button variant="secondary" onClick={() => setCancelConfirm(true)}>
            <XCircle size={16} aria-hidden="true" />
            Отменить попытку
          </Button>
        </div>
      ) : null}

      <ConfirmDialog
        open={finishConfirm}
        onClose={() => setFinishConfirm(false)}
        onConfirm={() => void finishMutation.mutateAsync()}
        title="Завершить тест"
        message={`Отвечено ${attempt.answers_count} из ${attempt.questions_count} вопросов. Завершить попытку?`}
        confirmLabel="Завершить"
        isPending={finishMutation.isPending}
      />

      <ConfirmDialog
        open={cancelConfirm}
        onClose={() => setCancelConfirm(false)}
        onConfirm={() => void cancelMutation.mutateAsync()}
        title="Отменить попытку"
        message="Попытка будет отменена. Это действие нельзя отменить."
        confirmLabel="Отменить попытку"
        danger
        isPending={cancelMutation.isPending}
      />
    </>
  )
}
