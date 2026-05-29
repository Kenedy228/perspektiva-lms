import { ArrowLeft, CheckCircle2, Circle, MinusCircle } from 'lucide-react'
import { Link, useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getAttempt } from '../../api/attempts'
import { Badge } from '../../components/ui/Badge'
import { PageHeader } from '../../components/ui/PageHeader'
import { ATTEMPT_STATUS_LABELS } from '../../types/attempt'
import type { AttemptStatus } from '../../types/attempt'
import styles from './AttemptAdminDetailPage.module.css'

function statusVariant(status: AttemptStatus): 'success' | 'warning' | 'danger' | 'neutral' {
  if (status === 'finished') return 'success'
  if (status === 'in_progress') return 'warning'
  if (status === 'expired' || status === 'cancelled' || status === 'interrupted') return 'danger'
  return 'neutral'
}

function formatDate(iso?: string) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString('ru-RU', { dateStyle: 'medium', timeStyle: 'short' })
}

function shortId(id: string) {
  return id.length > 8 ? `${id.slice(0, 8)}…` : id
}

const QUESTION_TYPE_LABELS: Record<string, string> = {
  selectable: 'Выбор ответа',
  sequence: 'Последовательность',
  matching: 'Сопоставление',
  short: 'Краткий ответ',
}

function ScoreIcon({ score }: { score: number | undefined }) {
  if (score === undefined) return <Circle size={18} className={styles.iconUnanswered} aria-label="Без ответа" />
  if (score >= 1) return <CheckCircle2 size={18} className={styles.iconCorrect} aria-label="Верно" />
  if (score > 0) return <MinusCircle size={18} className={styles.iconPartial} aria-label="Частично верно" />
  return <Circle size={18} className={styles.iconWrong} aria-label="Неверно" />
}

function scoreLabel(score: number | undefined, answered: boolean): string {
  if (!answered) return 'Без ответа'
  if (score === undefined) return 'Нет оценки'
  if (score >= 1) return 'Верно'
  if (score > 0) return `Частично (${Math.round(score * 100)}%)`
  return 'Неверно'
}

function TotalScoreBar({ totalScore, questionsCount }: { totalScore: number; questionsCount: number }) {
  const pct = Math.round(totalScore * 100)
  return (
    <div className={styles.scoreBar}>
      <div className={styles.scoreBarTrack}>
        <div className={styles.scoreBarFill} style={{ width: `${pct}%` }} />
      </div>
      <span className={styles.scoreBarText}>{pct}%</span>
      <span className={styles.scoreBarHint}>
        ({(totalScore * questionsCount).toFixed(1)} из {questionsCount} баллов)
      </span>
    </div>
  )
}

export function AttemptAdminDetailPage() {
  const { id = '' } = useParams<{ id: string }>()

  const query = useQuery({
    queryKey: ['attempt-admin', id],
    queryFn: () => getAttempt(id),
    enabled: Boolean(id),
  })

  const attempt = query.data

  if (query.isPending) {
    return <p className={styles.state}>Загрузка…</p>
  }

  if (query.isError || !attempt) {
    return <p className={styles.stateError}>Не удалось загрузить попытку.</p>
  }

  const answeredSet = new Set(attempt.answered_question_ids)
  const hasScores = attempt.question_scores != null

  return (
    <>
      <div className={styles.backRow}>
        <Link
          to={`/enrollments/${attempt.enrollment_id}/attempts`}
          className={styles.backLink}
        >
          <ArrowLeft size={15} aria-hidden="true" />
          Попытки зачисления
        </Link>
      </div>

      <PageHeader
        title="Просмотр попытки"
        description={`ID: ${shortId(id)}`}
      />

      <div className={styles.metaSection}>
        <dl className={styles.details}>
          <dt>Статус</dt>
          <dd>
            <Badge variant={statusVariant(attempt.status)}>
              {ATTEMPT_STATUS_LABELS[attempt.status] ?? attempt.status}
            </Badge>
          </dd>
          <dt>Начата</dt>
          <dd>{formatDate(attempt.started_at)}</dd>
          <dt>Дедлайн</dt>
          <dd>{formatDate(attempt.deadline_at)}</dd>
          <dt>Завершена</dt>
          <dd>{formatDate(attempt.finished_at)}</dd>
          <dt>Ответов</dt>
          <dd>{attempt.answers_count} из {attempt.questions_count}</dd>
          {hasScores && attempt.total_score != null && (
            <>
              <dt>Итоговый балл</dt>
              <dd>
                <TotalScoreBar
                  totalScore={attempt.total_score}
                  questionsCount={attempt.questions_count}
                />
              </dd>
            </>
          )}
          <dt>Зачисление</dt>
          <dd><code className={styles.mono}>{attempt.enrollment_id}</code></dd>
          <dt>Тест</dt>
          <dd><code className={styles.mono}>{attempt.quiz_id}</code></dd>
        </dl>
      </div>

      {attempt.questions.length > 0 && (
        <section className={styles.questionsSection}>
          <h2 className={styles.sectionTitle}>
            {hasScores ? 'Результаты по вопросам' : 'Вопросы'}
          </h2>
          <ul className={styles.questionList}>
            {attempt.questions.map((q, idx) => {
              const answered = answeredSet.has(q.id)
              const score = attempt.question_scores?.[q.id]
              return (
                <li
                  key={q.id}
                  className={[
                    styles.questionRow,
                    hasScores && answered && score !== undefined
                      ? score >= 1 ? styles.rowCorrect : score > 0 ? styles.rowPartial : styles.rowWrong
                      : '',
                  ].filter(Boolean).join(' ')}
                >
                  <div className={styles.questionStatus}>
                    <ScoreIcon score={answered ? score : undefined} />
                  </div>
                  <div className={styles.questionMeta}>
                    <span className={styles.questionNum}>{idx + 1}.</span>
                    <span className={styles.questionTitle}>{q.title}</span>
                    <span className={styles.questionType}>
                      {QUESTION_TYPE_LABELS[q.type] ?? q.type}
                    </span>
                  </div>
                  <div className={styles.questionResult}>
                    {scoreLabel(answered ? score : undefined, answered)}
                  </div>
                </li>
              )
            })}
          </ul>
        </section>
      )}
    </>
  )
}
