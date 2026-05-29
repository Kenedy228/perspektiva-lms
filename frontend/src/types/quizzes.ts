// GET /quizzes/{id} returns snake_case (map[string]any in Go handler)
export type QuizView = {
  id: string
  title: string
  max_attempts: number
  time_limit_seconds: number
  shuffle_questions: boolean
  sources: QuizSourceView[]
}

export type QuizSourceView = {
  bank_id: string
  criteria_type: 'random' | 'manual'
  question_count: number
  question_ids: string[]
}

export type CreateQuizPayload = {
  title: string
  max_attempts: number
  time_limit_seconds: number
  shuffle_questions: boolean
  sources: QuizSourcePayload[]
}

export type QuizSourcePayload = {
  bank_id: string
  criteria_type: 'random' | 'manual'
  question_count?: number
  question_ids?: string[]
}
