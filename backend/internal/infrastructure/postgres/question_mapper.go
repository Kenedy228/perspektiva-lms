package postgres

import (
	"encoding/json"
	"fmt"

	questdomain "gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	questiontitle "gitflic.ru/lms/backend/internal/domain/question/base/title"
	"gitflic.ru/lms/backend/internal/domain/question/matching"
	matchingpair "gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"gitflic.ru/lms/backend/internal/domain/question/selectable"
	selectableoption "gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"gitflic.ru/lms/backend/internal/domain/question/sequence"
	sequenceoption "gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	shortquestion "gitflic.ru/lms/backend/internal/domain/question/short"
	shortvariant "gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"github.com/google/uuid"
)

type questionPayload struct {
	Title      string          `json:"title"`
	Selectable []optionPayload `json:"selectable,omitempty"`
	Sequence   []string        `json:"sequence,omitempty"`
	Matching   []pairPayload   `json:"matching,omitempty"`
	Short      []string        `json:"short,omitempty"`
}

type optionPayload struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	IsCorrect bool   `json:"is_correct,omitempty"`
}

type pairPayload struct {
	PromptID   string `json:"prompt_id"`
	PromptText string `json:"prompt_text"`
	MatchID    string `json:"match_id"`
	MatchText  string `json:"match_text"`
}

func marshalQuestion(q questdomain.Question) ([]byte, error) {
	payload := questionPayload{Title: q.Title().Value()}

	switch typed := q.(type) {
	case *selectable.Question:
		for _, opt := range typed.Options() {
			payload.Selectable = append(payload.Selectable, optionPayload{
				ID:        opt.ID().String(),
				Value:     opt.Value(),
				IsCorrect: opt.IsCorrect(),
			})
		}
	case *sequence.Question:
		for _, opt := range typed.Options() {
			payload.Sequence = append(payload.Sequence, opt.Value())
		}
	case *matching.Question:
		for _, p := range typed.Pairs() {
			payload.Matching = append(payload.Matching, pairPayload{
				PromptID:   p.PromptID().String(),
				PromptText: p.Prompt().Value(),
				MatchID:    p.MatchID().String(),
				MatchText:  p.Match().Value(),
			})
		}
	case *shortquestion.Question:
		for _, v := range typed.Variants() {
			payload.Short = append(payload.Short, v.Value())
		}
	default:
		return nil, fmt.Errorf("%w: unsupported question type %T", ErrUnsupported, q)
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal question payload: %w", err)
	}
	return raw, nil
}

func unmarshalQuestion(id uuid.UUID, qType string, raw []byte) (questdomain.Question, error) {
	var payload questionPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal question payload: %w", err)
	}

	t, err := questiontitle.New(payload.Title)
	if err != nil {
		return nil, fmt.Errorf("restore question title: %w", err)
	}

	b, err := base.Restore(id, t)
	if err != nil {
		return nil, fmt.Errorf("restore question base: %w", err)
	}

	switch questdomain.Type(qType) {
	case questdomain.TypeSelectable:
		options, err := restoreSelectableOptions(payload.Selectable)
		if err != nil {
			return nil, err
		}
		return selectable.Restore(b, options)
	case questdomain.TypeSequence:
		options, err := restoreSequenceOptions(payload.Sequence)
		if err != nil {
			return nil, err
		}
		return sequence.Restore(b, options)
	case questdomain.TypeMatching:
		pairs, err := restorePairs(payload.Matching)
		if err != nil {
			return nil, err
		}
		return matching.Restore(b, pairs)
	case questdomain.TypeShort:
		variants, err := restoreShortVariants(payload.Short)
		if err != nil {
			return nil, err
		}
		return shortquestion.Restore(b, variants)
	default:
		return nil, fmt.Errorf("%w: unsupported question type %q", ErrUnsupported, qType)
	}
}

func restoreSelectableOptions(payloads []optionPayload) ([]selectableoption.Option, error) {
	options := make([]selectableoption.Option, 0, len(payloads))
	for _, payload := range payloads {
		id, err := uuid.Parse(payload.ID)
		if err != nil {
			return nil, err
		}
		opt, err := selectableoption.Restore(id, payload.Value, payload.IsCorrect)
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	return options, nil
}

func restoreSequenceOptions(values []string) ([]sequenceoption.Option, error) {
	options := make([]sequenceoption.Option, 0, len(values))
	for _, value := range values {
		opt, err := sequenceoption.New(value)
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	return options, nil
}

func restorePairs(payloads []pairPayload) ([]matchingpair.Pair, error) {
	pairs := make([]matchingpair.Pair, 0, len(payloads))
	for _, payload := range payloads {
		promptID, err := uuid.Parse(payload.PromptID)
		if err != nil {
			return nil, err
		}
		matchID, err := uuid.Parse(payload.MatchID)
		if err != nil {
			return nil, err
		}
		prompt, err := matchingpair.RestorePrompt(promptID, payload.PromptText)
		if err != nil {
			return nil, err
		}
		match, err := matchingpair.RestoreMatch(matchID, payload.MatchText)
		if err != nil {
			return nil, err
		}
		pairValue, err := matchingpair.New(prompt, match)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, pairValue)
	}
	return pairs, nil
}

func restoreShortVariants(values []string) ([]shortvariant.Variant, error) {
	variants := make([]shortvariant.Variant, 0, len(values))
	for _, value := range values {
		variant, err := shortvariant.New(value)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}
	return variants, nil
}
