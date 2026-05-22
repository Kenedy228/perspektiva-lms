package commands

import (
	"context"
	"fmt"

	questports "gitflic.ru/lms/backend/internal/application/ports/question"
	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	qmatching "gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	selectableoption "gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	qsequence "gitflic.ru/lms/backend/internal/domain/question/sequence"
	sequenceoption "gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	qshort "gitflic.ru/lms/backend/internal/domain/question/short"
	shortvariant "gitflic.ru/lms/backend/internal/domain/question/short/variant"
	qtyped "gitflic.ru/lms/backend/internal/domain/question/typed"
	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
	"gitflic.ru/lms/backend/internal/domain/shared/media"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
)

func parseRequiredUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", field, err)
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: %s is required", common.ErrInvalidInput, field)
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
		return nil, fmt.Errorf("find question: %w", err)
	}

	return q, nil
}

func buildAttachment(in AttachmentInput) (attachment.Attachment, error) {
	f, err := file.New(in.FileName, in.SizeBytes)
	if err != nil {
		return attachment.Attachment{}, fmt.Errorf("create attachment file: %w", err)
	}

	m, err := media.New(media.Type(in.MediaType), f)
	if err != nil {
		return attachment.Attachment{}, fmt.Errorf("create attachment media: %w", err)
	}

	att, err := attachment.New(m)
	if err != nil {
		return attachment.Attachment{}, fmt.Errorf("create question attachment: %w", err)
	}

	return att, nil
}

func buildSelectableOptions(items []SelectableOptionInput) ([]selectableoption.Option, error) {
	options := make([]selectableoption.Option, 0, len(items))
	for i := range items {
		t, err := text.New(items[i].Text)
		if err != nil {
			return nil, fmt.Errorf("create selectable option text: %w", err)
		}
		opt, err := selectableoption.New(t, items[i].IsCorrect)
		if err != nil {
			return nil, fmt.Errorf("create selectable option: %w", err)
		}
		options = append(options, opt)
	}
	return options, nil
}

func buildSequenceOptions(items []SequenceOptionInput) ([]sequenceoption.Option, error) {
	options := make([]sequenceoption.Option, 0, len(items))
	for i := range items {
		t, err := text.New(items[i].Text)
		if err != nil {
			return nil, fmt.Errorf("create sequence option text: %w", err)
		}
		opt, err := sequenceoption.New(t)
		if err != nil {
			return nil, fmt.Errorf("create sequence option: %w", err)
		}
		options = append(options, opt)
	}
	return options, nil
}

func buildMatchingPairs(items []MatchingPairInput) ([]pair.Pair, error) {
	pairs := make([]pair.Pair, 0, len(items))
	for i := range items {
		promptText, err := text.New(items[i].Prompt)
		if err != nil {
			return nil, fmt.Errorf("create matching prompt text: %w", err)
		}
		matchText, err := text.New(items[i].Match)
		if err != nil {
			return nil, fmt.Errorf("create matching match text: %w", err)
		}
		prompt, err := pair.NewPrompt(promptText)
		if err != nil {
			return nil, fmt.Errorf("create matching prompt: %w", err)
		}
		match, err := pair.NewMatch(matchText)
		if err != nil {
			return nil, fmt.Errorf("create matching match: %w", err)
		}
		p, err := pair.New(prompt, match)
		if err != nil {
			return nil, fmt.Errorf("create matching pair: %w", err)
		}
		pairs = append(pairs, p)
	}
	return pairs, nil
}

func buildTypedBlanks(items []TypedBlankInput) ([]blank.Blank, error) {
	blanks := make([]blank.Blank, 0, len(items))
	for i := range items {
		variants := make([]text.Text, 0, len(items[i].Variants))
		for j := range items[i].Variants {
			t, err := text.New(items[i].Variants[j])
			if err != nil {
				return nil, fmt.Errorf("create typed blank variant text: %w", err)
			}
			variants = append(variants, t)
		}
		b, err := blank.New(items[i].Placeholder, variants)
		if err != nil {
			return nil, fmt.Errorf("create typed blank: %w", err)
		}
		blanks = append(blanks, b)
	}
	return blanks, nil
}

func buildShortVariants(items []ShortVariantInput) ([]shortvariant.Variant, error) {
	variants := make([]shortvariant.Variant, 0, len(items))
	for i := range items {
		t, err := text.New(items[i].Text)
		if err != nil {
			return nil, fmt.Errorf("create short variant text: %w", err)
		}
		v, err := shortvariant.New(t)
		if err != nil {
			return nil, fmt.Errorf("create short variant: %w", err)
		}
		variants = append(variants, v)
	}
	return variants, nil
}

func createQuestion(qType question.Type, t title.Title, in CreateInput) (question.Question, error) {
	switch qType {
	case question.TypeSelectable:
		options, err := buildSelectableOptions(in.SelectableOptions)
		if err != nil {
			return nil, err
		}
		return qselectable.New(t, options)
	case question.TypeSequence:
		options, err := buildSequenceOptions(in.SequenceOptions)
		if err != nil {
			return nil, err
		}
		return qsequence.New(t, options)
	case question.TypeMatching:
		pairs, err := buildMatchingPairs(in.MatchingPairs)
		if err != nil {
			return nil, err
		}
		return qmatching.New(t, pairs)
	case question.TypeTyped:
		blanks, err := buildTypedBlanks(in.TypedBlanks)
		if err != nil {
			return nil, err
		}
		return qtyped.New(t, blanks)
	case question.TypeShort:
		variants, err := buildShortVariants(in.ShortVariants)
		if err != nil {
			return nil, err
		}
		return qshort.New(t, variants)
	default:
		return nil, fmt.Errorf("%w: unsupported question type", common.ErrInvalidInput)
	}
}
