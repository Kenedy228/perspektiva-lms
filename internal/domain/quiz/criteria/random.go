package criteria

type Random struct {
	questionCount int
}

func NewRandom(questionCount int) (Criteria, error) {
	if err := validateQuestionCount(questionCount); err != nil {
		return nil, err
	}

	return Random{
		questionCount: questionCount,
	}, nil
}

func (c Random) QuestionCount() int {
	return c.questionCount
}

func (c Random) Type() Type {
	return TypeRandom
}
