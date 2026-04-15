package quiz

import (
	"gitflic.ru/lms/internal/domain/question"
	"math/rand/v2"
)

type ManualStrategy struct {
	questions []question.Question
}

func New(questions []question.Question) (ManualStrategy, error) {
	if len(questions) == 0 {
		return ManualStrategy{}, nil
	}

	for i := range questions {
		if questions[i] == nil {
			return ManualStrategy{}, nil
		}
	}

	qCopy := make([]question.Question, 0, len(questions))
	for i := range questions {
		qCopy = append(qCopy, questions[i])
	}
}

func (s ManualStrategy) CountQuestions() int {
	return len(s.questions)
}

func (s ManualStrategy) Questions() []question.Question {
	return s.questions
}

func (s ManualStrategy) Prepare() []question.Question {
	rand.Shuffle(len(s.questions), func(i, j int) {
		s.questions[i], s.questions[j] = s.questions[j], s.questions[i]
	})

	return s.questions
}
