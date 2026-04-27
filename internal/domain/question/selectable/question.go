package selectable

import (
	"math/rand/v2"
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

type Question struct {
	*base.Base
	options      []Option
	correctCount int
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validateOptions(params.Options); err != nil {
		return nil, err
	}

	cOptions := slices.Clone(params.Options)
	correctCount := countCorrect(cOptions)

	return &Question{
		Base:         base,
		options:      cOptions,
		correctCount: correctCount,
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Options() []Option {
	return slices.Clone(q.options)
}

func (q *Question) ShuffledOptions() []Option {
	cOptions := slices.Clone(q.options)

	rand.Shuffle(len(cOptions), func(i, j int) {
		cOptions[i], cOptions[j] = cOptions[j], cOptions[i]
	})

	return cOptions
}

func (q *Question) Type() question.Type {
	return question.TypeSelectable
}

func (q *Question) UpdateOptions(options []Option) error {
	if err := validateOptions(options); err != nil {
		return err
	}

	cOptions := slices.Clone(options)
	correctCount := countCorrect(cOptions)

	q.options = cOptions
	q.correctCount = correctCount
	q.Touch()
	return nil
}

func (q *Question) CheckAnswer(answer question.Answer) bool {
	cast, ok := answer.(Answer)
	if !ok {
		return false
	}

	optionIDs := cast.OptionIDs()

	if len(optionIDs) != q.correctCount {
		return false
	}

	for i := range optionIDs {
		for j := i + 1; j < len(optionIDs); j++ {
			if optionIDs[i] == optionIDs[j] {
				return false
			}
		}
	}

	for i := range optionIDs {
		isCorrect := slices.ContainsFunc(q.options, func(opt Option) bool {
			return optionIDs[i] == opt.ID() && opt.IsCorrect()
		})

		if !isCorrect {
			return false
		}
	}

	return true
}

func (q *Question) Clone() question.Question {
	return &Question{
		Base:    q.Base.Clone(),
		options: q.Options(),
	}
}
