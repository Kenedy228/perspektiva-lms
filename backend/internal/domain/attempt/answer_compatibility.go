package attempt

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question"
	matchinganswer "gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	selectableanswer "gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	sequenceanswer "gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	shortanswer "gitflic.ru/lms/backend/internal/domain/question/short/answer"
	typedanswer "gitflic.ru/lms/backend/internal/domain/question/typed/answer"
)

func validateAnswerForQuestion(q question.Question, ans question.Answer) error {
	if q == nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	switch q.Type() {
	case question.TypeSelectable:
		if _, ok := ans.(selectableanswer.Answer); ok {
			return nil
		}
	case question.TypeMatching:
		if _, ok := ans.(matchinganswer.Answer); ok {
			return nil
		}
	case question.TypeSequence:
		if _, ok := ans.(sequenceanswer.Answer); ok {
			return nil
		}
	case question.TypeTyped:
		if _, ok := ans.(typedanswer.Answer); ok {
			return nil
		}
	case question.TypeShort:
		if _, ok := ans.(shortanswer.Answer); ok {
			return nil
		}
	}

	return fmt.Errorf("%w: invalid value (%q)", ErrInvalid, q.Type())
}
