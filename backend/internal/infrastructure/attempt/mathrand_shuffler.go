package attempt

import (
	"math/rand"
	"sync"
	"time"

	"gitflic.ru/lms/backend/internal/domain/question"
)

type MathRandQuestionShuffler struct {
	mu sync.Mutex
	r  *rand.Rand
}

func NewMathRandQuestionShuffler() *MathRandQuestionShuffler {
	return NewMathRandQuestionShufflerWithSource(rand.NewSource(time.Now().UnixNano()))
}

func NewMathRandQuestionShufflerWithSource(source rand.Source) *MathRandQuestionShuffler {
	if source == nil {
		panic("math rand question shuffler requires source")
	}

	return &MathRandQuestionShuffler{
		r: rand.New(source),
	}
}

func (s *MathRandQuestionShuffler) ShuffleQuestions(questions []question.Question) []question.Question {
	shuffled := make([]question.Question, len(questions))
	copy(shuffled, questions)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}
