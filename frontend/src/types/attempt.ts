export type AttemptStatus = 'in_progress' | 'finished' | 'expired' | 'cancelled' | 'interrupted'

export const ATTEMPT_STATUS_LABELS: Record<AttemptStatus, string> = {
  in_progress: 'В процессе',
  finished: 'Завершён',
  expired: 'Просрочен',
  cancelled: 'Отменён',
  interrupted: 'Прерван',
}

export type QuestionOption = {
  id: string
  text: string
}

export type MatchingPair = {
  prompt_id: string
  prompt: string
  match_id: string
  match: string
}

export type AttemptQuestion = {
  id: string
  type: 'selectable' | 'sequence' | 'matching' | 'short'
  title: string
  instruction: string
  options?: QuestionOption[]
  pairs?: MatchingPair[]
}

export type AnswerPayload = {
  type: string
  option_ids?: string[]
  matching_pairs?: Record<string, string>
  short_input?: string
}

export type AttemptView = {
  id: string
  enrollment_id: string
  quiz_id: string
  status: AttemptStatus
  started_at: string
  deadline_at: string
  finished_at: string
  questions_count: number
  answers_count: number
  answered_question_ids: string[]
  questions: AttemptQuestion[]
  /** Present only for non-in_progress attempts: question_id → score (0.0–1.0) */
  question_scores?: Record<string, number>
  /** Present only for non-in_progress attempts: sum(scores) / questions_count */
  total_score?: number
}

export type StartAttemptRequest = {
  account_id: string
  enrollment_id: string
  quiz_id: string
}

export type AttemptSummaryView = {
  id: string
  enrollment_id: string
  quiz_id: string
  status: AttemptStatus
  started_at: string
  deadline_at?: string
  finished_at?: string
  questions_count: number
  answers_count: number
}
