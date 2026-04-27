package question

type Answer interface {
	IsEmpty() bool
	Clone() Answer
}
