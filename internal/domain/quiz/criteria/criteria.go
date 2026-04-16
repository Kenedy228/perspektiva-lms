package criteria

type Criteria interface {
	Type() CriteriaType
	QuestionCount() int
}
