export type BankShortView = {
  ID: string
  Title: string
  QuestionsCount: number
}

export type SelectableOptionView = {
  ID: string
  Value: string
  IsCorrect: boolean
}

export type SequenceOptionView = {
  Value: string
}

export type MatchingPairView = {
  PromptID: string
  PromptText: string
  MatchID: string
  MatchText: string
}

export type ShortVariantView = {
  Value: string
}

export type BankQuestionView = {
  ID: string
  Type: 'selectable' | 'sequence' | 'matching' | 'short'
  Title: string
  SelectableOptions: SelectableOptionView[]
  SequenceOptions: SequenceOptionView[]
  MatchingPairs: MatchingPairView[]
  ShortVariants: ShortVariantView[]
}

export type BankDetailedView = {
  ID: string
  Title: string
  QuestionIDs: string[]
  Questions: BankQuestionView[]
}
