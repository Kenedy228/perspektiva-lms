package commands

import (
	"context"
	"fmt"

	questports "gitflic.ru/lms/backend/internal/application/ports/question"
	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	questiontitle "gitflic.ru/lms/backend/internal/domain/question/base/title"
	qmatching "gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	selectableoption "gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	qsequence "gitflic.ru/lms/backend/internal/domain/question/sequence"
	sequenceoption "gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	qshort "gitflic.ru/lms/backend/internal/domain/question/short"
	shortvariant "gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"github.com/google/uuid"
)

func parseRequiredUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("разбор %s: %w", field, err)
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: поле %s обязательно", common.ErrInvalidInput, field)
	}
	return id, nil
}

func loadQuestion(ctx context.Context, r questports.Repository, id string) (question.Question, error) {
	qID, err := parseRequiredUUID(id, "question id")
	if err != nil {
		return nil, err
	}

	q, err := r.FindByID(ctx, qID)
	if err != nil {
		return nil, fmt.Errorf("загрузка вопроса: %w", err)
	}

	return q, nil
}

func buildSelectableOptions(items []SelectableOptionInput) ([]selectableoption.Option, error) {
	options := make([]selectableoption.Option, 0, len(items))
	for i := range items {
		opt, err := selectableoption.New(items[i].Text, items[i].IsCorrect)
		if err != nil {
			return nil, fmt.Errorf("создание варианта ответа с выбором: %w", err)
		}
		options = append(options, opt)
	}
	return options, nil
}

func buildSequenceOptions(items []SequenceOptionInput) ([]sequenceoption.Option, error) {
	options := make([]sequenceoption.Option, 0, len(items))
	for i := range items {
		opt, err := sequenceoption.New(items[i].Text)
		if err != nil {
			return nil, fmt.Errorf("создание варианта последовательности: %w", err)
		}
		options = append(options, opt)
	}
	return options, nil
}

func buildMatchingPairs(items []MatchingPairInput) ([]pair.Pair, error) {
	pairs := make([]pair.Pair, 0, len(items))
	for i := range items {
		prompt, err := pair.NewPrompt(items[i].Prompt)
		if err != nil {
			return nil, fmt.Errorf("создание текста в паре сопоставления: %w", err)
		}
		match, err := pair.NewMatch(items[i].Match)
		if err != nil {
			return nil, fmt.Errorf("создание соответствия в паре: %w", err)
		}
		p, err := pair.New(prompt, match)
		if err != nil {
			return nil, fmt.Errorf("создание пары сопоставления: %w", err)
		}
		pairs = append(pairs, p)
	}
	return pairs, nil
}

func buildShortVariants(items []ShortVariantInput) ([]shortvariant.Variant, error) {
	variants := make([]shortvariant.Variant, 0, len(items))
	for i := range items {
		v, err := shortvariant.New(items[i].Text)
		if err != nil {
			return nil, fmt.Errorf("создание варианта короткого ответа: %w", err)
		}
		variants = append(variants, v)
	}
	return variants, nil
}

func createQuestion(qType question.Type, t questiontitle.Title, in CreateInput) (question.Question, error) {
	b, err := base.New(t)
	if err != nil {
		return nil, fmt.Errorf("создание базовой части вопроса: %w", err)
	}

	switch qType {
	case question.TypeSelectable:
		options, err := buildSelectableOptions(in.SelectableOptions)
		if err != nil {
			return nil, err
		}
		return qselectable.New(b, options)
	case question.TypeSequence:
		options, err := buildSequenceOptions(in.SequenceOptions)
		if err != nil {
			return nil, err
		}
		return qsequence.New(b, options)
	case question.TypeMatching:
		pairs, err := buildMatchingPairs(in.MatchingPairs)
		if err != nil {
			return nil, err
		}
		return qmatching.New(b, pairs)
	case question.TypeShort:
		variants, err := buildShortVariants(in.ShortVariants)
		if err != nil {
			return nil, err
		}
		return qshort.New(b, variants)
	default:
		return nil, fmt.Errorf("%w: неподдерживаемый тип вопроса", common.ErrInvalidInput)
	}
}
