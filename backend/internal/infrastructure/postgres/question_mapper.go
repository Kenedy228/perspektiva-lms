package postgres

import (
	"encoding/json"
	"fmt"

	questdomain "gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/question/matching"
	matchingpair "gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"gitflic.ru/lms/backend/internal/domain/question/selectable"
	selectableoption "gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"gitflic.ru/lms/backend/internal/domain/question/sequence"
	sequenceoption "gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	shortquestion "gitflic.ru/lms/backend/internal/domain/question/short"
	shortvariant "gitflic.ru/lms/backend/internal/domain/question/short/variant"
	typedquestion "gitflic.ru/lms/backend/internal/domain/question/typed"
	typedblank "gitflic.ru/lms/backend/internal/domain/question/typed/blank"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
	"gitflic.ru/lms/backend/internal/domain/shared/media"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
)

type questionPayload struct {
	Title      string             `json:"title"`
	Attachment *attachmentPayload `json:"attachment,omitempty"`
	Selectable []optionPayload    `json:"selectable,omitempty"`
	Sequence   []optionPayload    `json:"sequence,omitempty"`
	Matching   []pairPayload      `json:"matching,omitempty"`
	Short      []string           `json:"short,omitempty"`
	Typed      []blankPayload     `json:"typed,omitempty"`
}

type attachmentPayload struct {
	MediaType string `json:"media_type"`
	FileName  string `json:"file_name"`
	SizeBytes int64  `json:"size_bytes"`
}

type optionPayload struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct,omitempty"`
}

type pairPayload struct {
	PromptID   string `json:"prompt_id"`
	PromptText string `json:"prompt_text"`
	MatchID    string `json:"match_id"`
	MatchText  string `json:"match_text"`
}

type blankPayload struct {
	Placeholder string   `json:"placeholder"`
	Variants    []string `json:"variants"`
}

func marshalQuestion(q questdomain.Question) ([]byte, error) {
	payload := questionPayload{Title: q.Title().Value()}
	if att, ok := q.Attachment(); ok {
		payload.Attachment = marshalAttachment(att)
	}

	switch typed := q.(type) {
	case *selectable.Question:
		for _, opt := range typed.Options() {
			payload.Selectable = append(payload.Selectable, optionPayload{
				ID:        opt.ID().String(),
				Text:      opt.Text().Value(),
				IsCorrect: opt.IsCorrect(),
			})
		}
	case *sequence.Question:
		for _, opt := range typed.Options() {
			payload.Sequence = append(payload.Sequence, optionPayload{
				ID:   opt.ID().String(),
				Text: opt.Text().Value(),
			})
		}
	case *matching.Question:
		for _, p := range typed.Pairs() {
			payload.Matching = append(payload.Matching, pairPayload{
				PromptID:   p.PromptID().String(),
				PromptText: p.Prompt().Text().Value(),
				MatchID:    p.MatchID().String(),
				MatchText:  p.Match().Text().Value(),
			})
		}
	case *shortquestion.Question:
		for _, v := range typed.Variants() {
			payload.Short = append(payload.Short, v.Text().Value())
		}
	case *typedquestion.Question:
		for _, b := range typed.Blanks() {
			payload.Typed = append(payload.Typed, blankPayload{
				Placeholder: b.Placeholder(),
				Variants:    b.VariantsValues(),
			})
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

	t, err := title.New(payload.Title)
	if err != nil {
		return nil, fmt.Errorf("restore question title: %w", err)
	}
	att, err := unmarshalAttachment(payload.Attachment)
	if err != nil {
		return nil, err
	}

	switch questdomain.Type(qType) {
	case questdomain.TypeSelectable:
		options, err := restoreSelectableOptions(payload.Selectable)
		if err != nil {
			return nil, err
		}
		return selectable.Restore(id, t, att, options)
	case questdomain.TypeSequence:
		options, err := restoreSequenceOptions(payload.Sequence)
		if err != nil {
			return nil, err
		}
		return sequence.Restore(id, t, att, options)
	case questdomain.TypeMatching:
		pairs, err := restorePairs(payload.Matching)
		if err != nil {
			return nil, err
		}
		return matching.Restore(id, t, att, pairs)
	case questdomain.TypeShort:
		variants, err := restoreShortVariants(payload.Short)
		if err != nil {
			return nil, err
		}
		return shortquestion.Restore(id, t, att, variants)
	case questdomain.TypeTyped:
		blanks, err := restoreTypedBlanks(payload.Typed)
		if err != nil {
			return nil, err
		}
		return typedquestion.Restore(id, t, att, blanks)
	default:
		return nil, fmt.Errorf("%w: unsupported question type %q", ErrUnsupported, qType)
	}
}

func marshalAttachment(att attachment.Attachment) *attachmentPayload {
	m := att.Media()
	f := m.File()
	return &attachmentPayload{
		MediaType: m.Type().String(),
		FileName:  f.Name(),
		SizeBytes: f.SizeBytes(),
	}
}

func unmarshalAttachment(payload *attachmentPayload) (*attachment.Attachment, error) {
	if payload == nil {
		return nil, nil
	}
	f, err := file.New(payload.FileName, payload.SizeBytes)
	if err != nil {
		return nil, fmt.Errorf("restore attachment file: %w", err)
	}
	m, err := media.New(media.Type(payload.MediaType), f)
	if err != nil {
		return nil, fmt.Errorf("restore attachment media: %w", err)
	}
	att, err := attachment.New(m)
	if err != nil {
		return nil, fmt.Errorf("restore attachment: %w", err)
	}
	return &att, nil
}

func restoreSelectableOptions(payloads []optionPayload) ([]selectableoption.Option, error) {
	options := make([]selectableoption.Option, 0, len(payloads))
	for _, payload := range payloads {
		id, err := uuid.Parse(payload.ID)
		if err != nil {
			return nil, err
		}
		txt, err := text.New(payload.Text)
		if err != nil {
			return nil, err
		}
		opt, err := selectableoption.Restore(id, txt, payload.IsCorrect)
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	return options, nil
}

func restoreSequenceOptions(payloads []optionPayload) ([]sequenceoption.Option, error) {
	options := make([]sequenceoption.Option, 0, len(payloads))
	for _, payload := range payloads {
		id, err := uuid.Parse(payload.ID)
		if err != nil {
			return nil, err
		}
		txt, err := text.New(payload.Text)
		if err != nil {
			return nil, err
		}
		opt, err := sequenceoption.Restore(id, txt)
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
		promptText, err := text.New(payload.PromptText)
		if err != nil {
			return nil, err
		}
		matchText, err := text.New(payload.MatchText)
		if err != nil {
			return nil, err
		}
		prompt, err := matchingpair.RestorePrompt(promptID, promptText)
		if err != nil {
			return nil, err
		}
		match, err := matchingpair.RestoreMatch(matchID, matchText)
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
		txt, err := text.New(value)
		if err != nil {
			return nil, err
		}
		variant, err := shortvariant.New(txt)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}
	return variants, nil
}

func restoreTypedBlanks(payloads []blankPayload) ([]typedblank.Blank, error) {
	blanks := make([]typedblank.Blank, 0, len(payloads))
	for _, payload := range payloads {
		variants := make([]text.Text, 0, len(payload.Variants))
		for _, value := range payload.Variants {
			txt, err := text.New(value)
			if err != nil {
				return nil, err
			}
			variants = append(variants, txt)
		}
		blank, err := typedblank.New(payload.Placeholder, variants)
		if err != nil {
			return nil, err
		}
		blanks = append(blanks, blank)
	}
	return blanks, nil
}
