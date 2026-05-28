package answer

type AnswerBlank struct {
	Placeholder string
	Variant     string
}

type Answer struct {
	blanks []AnswerBlank
}

func New(blanks []AnswerBlank) (Answer, error) {
	return Answer{blanks: append([]AnswerBlank(nil), blanks...)}, nil
}
