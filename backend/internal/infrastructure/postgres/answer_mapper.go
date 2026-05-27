package postgres

import (
	"encoding/json"
	"fmt"

	questdomain "gitflic.ru/lms/backend/internal/domain/question"
	matchinganswer "gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	selectableanswer "gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	sequenceanswer "gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	shortanswer "gitflic.ru/lms/backend/internal/domain/question/short/answer"
	typedanswer "gitflic.ru/lms/backend/internal/domain/question/typed/answer"
	"github.com/google/uuid"
)

type answerPayload struct {
	SelectableOptionIDs []string                  `json:"selectable_option_ids,omitempty"`
	SequenceOptionIDs   []string                  `json:"sequence_option_ids,omitempty"`
	MatchingPairs       []answerPairPayload       `json:"matching_pairs,omitempty"`
	TypedBlanks         []typedanswer.AnswerBlank `json:"typed_blanks,omitempty"`
	ShortInput          string                    `json:"short_input,omitempty"`
}

type answerPairPayload struct {
	PromptID string `json:"prompt_id"`
	MatchID  string `json:"match_id"`
}

func marshalAnswer(ans questdomain.Answer) ([]byte, error) {
	var payload answerPayload
	switch typed := ans.(type) {
	case selectableanswer.Answer:
		for _, id := range typed.OptionIDs() {
			payload.SelectableOptionIDs = append(payload.SelectableOptionIDs, id.String())
		}
	case sequenceanswer.Answer:
		for _, id := range typed.OptionIDs() {
			payload.SequenceOptionIDs = append(payload.SequenceOptionIDs, id.ID().String())
		}
	case matchinganswer.Answer:
		for _, p := range typed.Pairs() {
			payload.MatchingPairs = append(payload.MatchingPairs, answerPairPayload{
				PromptID: p.PromptID.String(),
				MatchID:  p.MatchID.String(),
			})
		}
	case typedanswer.Answer:
		payload.TypedBlanks = typed.Blanks()
	case shortanswer.Answer:
		payload.ShortInput = typed.Value()
	default:
		return nil, fmt.Errorf("%w: unsupported answer type %T", ErrUnsupported, ans)
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal answer payload: %w", err)
	}
	return raw, nil
}

func unmarshalAnswer(qType questdomain.Type, raw []byte) (questdomain.Answer, error) {
	var payload answerPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal answer payload: %w", err)
	}

	switch qType {
	case questdomain.TypeSelectable:
		ids := make([]uuid.UUID, 0, len(payload.SelectableOptionIDs))
		for _, rawID := range payload.SelectableOptionIDs {
			id, err := uuid.Parse(rawID)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
		return selectableanswer.New(ids)
	case questdomain.TypeSequence:
		ids := make([]sequenceanswer.OptionID, 0, len(payload.SequenceOptionIDs))
		for _, rawID := range payload.SequenceOptionIDs {
			id, err := uuid.Parse(rawID)
			if err != nil {
				return nil, err
			}
			optionID, err := sequenceanswer.NewOptionID(id)
			if err != nil {
				return nil, err
			}
			ids = append(ids, optionID)
		}
		return sequenceanswer.New(ids)
	case questdomain.TypeMatching:
		pairs := make([]matchinganswer.Pair, 0, len(payload.MatchingPairs))
		for _, rawPair := range payload.MatchingPairs {
			promptID, err := uuid.Parse(rawPair.PromptID)
			if err != nil {
				return nil, err
			}
			matchID, err := uuid.Parse(rawPair.MatchID)
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, matchinganswer.Pair{PromptID: promptID, MatchID: matchID})
		}
		return matchinganswer.New(pairs)
	case questdomain.TypeTyped:
		return typedanswer.New(payload.TypedBlanks)
	case questdomain.TypeShort:
		return shortanswer.New(payload.ShortInput)
	default:
		return nil, fmt.Errorf("%w: unsupported answer question type %q", ErrUnsupported, qType)
	}
}
