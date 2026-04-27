package criteria

type Criteria interface {
	Type() Type
	QuestionCount() int
}
