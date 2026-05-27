package handlers

import (
	"net/http"
	"time"

	questioncommands "gitflic.ru/lms/backend/internal/application/usecases/question/commands"
	questiongrading "gitflic.ru/lms/backend/internal/application/usecases/question/grading"
	questdomain "gitflic.ru/lms/backend/internal/domain/question"
	matchinganswer "gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	selectableanswer "gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	sequenceanswer "gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	shortanswer "gitflic.ru/lms/backend/internal/domain/question/short/answer"
	typedanswer "gitflic.ru/lms/backend/internal/domain/question/typed/answer"
	"gitflic.ru/lms/backend/internal/transport/http/response"
	"github.com/google/uuid"
)

func (api *API) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req QuestionRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Questions.Create.Execute(r.Context(), questioncommands.CreateInput{
		ActorRole: actor.role, Type: req.Type, Title: req.Title, Attachment: toQuestionAttachment(req.Attachment),
		SelectableOptions: toSelectableInputs(req.SelectableOptions), SequenceOptions: toSequenceInputs(req.SequenceOptions),
		MatchingPairs: toMatchingInputs(req.MatchingPairs), TypedBlanks: toTypedInputs(req.TypedBlanks), ShortVariants: toShortInputs(req.ShortVariants),
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/questions", out.ID)
}

func (api *API) GetQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	q, err := api.Questions.Repository.FindByID(r.Context(), id)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]any{"id": q.ID().String(), "type": q.Type().String(), "title": q.Title().Value(), "instruction": q.Instruction()}, response.Links{"self": {Href: r.URL.Path, Method: http.MethodGet}, "grade": {Href: r.URL.Path + "/grade", Method: http.MethodPost}})
}

func (api *API) ChangeQuestionTitle(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req QuestionRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Questions.ChangeTitle.Execute(r.Context(), questioncommands.ChangeTitleInput{ActorRole: actor.role, QuestionID: r.PathValue("id"), Title: req.Title})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) ChangeQuestionContent(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req QuestionRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	var id string
	var err error
	switch questdomain.Type(req.Type) {
	case questdomain.TypeSelectable:
		out, e := api.Questions.Selectable.Execute(r.Context(), questioncommands.ChangeSelectableOptionsInput{ActorRole: actor.role, QuestionID: r.PathValue("id"), Options: toSelectableInputs(req.SelectableOptions)})
		if e == nil {
			id = out.ID
		}
		err = e
	case questdomain.TypeSequence:
		out, e := api.Questions.Sequence.Execute(r.Context(), questioncommands.ChangeSequenceOptionsInput{ActorRole: actor.role, QuestionID: r.PathValue("id"), Options: toSequenceInputs(req.SequenceOptions)})
		if e == nil {
			id = out.ID
		}
		err = e
	case questdomain.TypeMatching:
		out, e := api.Questions.Matching.Execute(r.Context(), questioncommands.ChangeMatchingPairsInput{ActorRole: actor.role, QuestionID: r.PathValue("id"), Pairs: toMatchingInputs(req.MatchingPairs)})
		if e == nil {
			id = out.ID
		}
		err = e
	case questdomain.TypeTyped:
		out, e := api.Questions.Typed.Execute(r.Context(), questioncommands.ChangeTypedContentInput{ActorRole: actor.role, QuestionID: r.PathValue("id"), Title: req.Title, Blanks: toTypedInputs(req.TypedBlanks)})
		if e == nil {
			id = out.ID
		}
		err = e
	case questdomain.TypeShort:
		out, e := api.Questions.Short.Execute(r.Context(), questioncommands.ChangeShortVariantsInput{ActorRole: actor.role, QuestionID: r.PathValue("id"), Variants: toShortInputs(req.ShortVariants)})
		if e == nil {
			id = out.ID
		}
		err = e
	default:
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_input", "unknown question type"))
		return
	}
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": id}, nil)
}

func (api *API) ChangeQuestionAttachment(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req AttachmentInput
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Questions.ChangeAttachment.Execute(r.Context(), questioncommands.ChangeAttachmentInput{ActorRole: actor.role, QuestionID: r.PathValue("id"), Attachment: questioncommands.AttachmentInput(req)})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) RemoveQuestionAttachment(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Questions.RemoveAttachment.Execute(r.Context(), questioncommands.RemoveAttachmentInput{ActorRole: actor.role, QuestionID: r.PathValue("id")})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) GradeQuestion(w http.ResponseWriter, r *http.Request) {
	var req AttemptAnswerRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	ans, err := buildAnswer(req)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	out, err := api.Questions.Grade.Execute(r.Context(), questiongrading.GradeInput{QuestionID: r.PathValue("id"), Answer: ans})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out, nil)
}

func buildAnswer(req AttemptAnswerRequest) (questdomain.Answer, error) {
	switch questdomain.Type(req.Type) {
	case questdomain.TypeSelectable:
		ids := make([]uuid.UUID, 0, len(req.OptionIDs))
		for _, raw := range req.OptionIDs {
			id, err := uuid.Parse(raw)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
		return selectableanswer.New(ids)
	case questdomain.TypeSequence:
		ids := make([]sequenceanswer.OptionID, 0, len(req.OptionIDs))
		for _, raw := range req.OptionIDs {
			id, err := uuid.Parse(raw)
			if err != nil {
				return nil, err
			}
			oid, err := sequenceanswer.NewOptionID(id)
			if err != nil {
				return nil, err
			}
			ids = append(ids, oid)
		}
		return sequenceanswer.New(ids)
	case questdomain.TypeMatching:
		pairs := make([]matchinganswer.Pair, 0, len(req.MatchingPairs))
		for p, m := range req.MatchingPairs {
			pid, err := uuid.Parse(p)
			if err != nil {
				return nil, err
			}
			mid, err := uuid.Parse(m)
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, matchinganswer.Pair{PromptID: pid, MatchID: mid})
		}
		return matchinganswer.New(pairs)
	case questdomain.TypeTyped:
		blanks := make([]typedanswer.AnswerBlank, 0, len(req.TypedBlanks))
		for p, v := range req.TypedBlanks {
			blanks = append(blanks, typedanswer.AnswerBlank{Placeholder: p, Variant: v})
		}
		return typedanswer.New(blanks)
	case questdomain.TypeShort:
		return shortanswer.New(req.ShortInput)
	default:
		return nil, questiongrading.ErrInvalidInput
	}
}

func toQuestionAttachment(in *AttachmentInput) *questioncommands.AttachmentInput {
	if in == nil {
		return nil
	}
	return &questioncommands.AttachmentInput{MediaType: in.MediaType, FileName: in.FileName, SizeBytes: in.SizeBytes}
}
func toSelectableInputs(in []SelectableOptionInput) []questioncommands.SelectableOptionInput {
	out := make([]questioncommands.SelectableOptionInput, len(in))
	for i := range in {
		out[i] = questioncommands.SelectableOptionInput{Text: in[i].Text, IsCorrect: in[i].IsCorrect}
	}
	return out
}
func toSequenceInputs(in []SequenceOptionInput) []questioncommands.SequenceOptionInput {
	out := make([]questioncommands.SequenceOptionInput, len(in))
	for i := range in {
		out[i] = questioncommands.SequenceOptionInput{Text: in[i].Text}
	}
	return out
}
func toMatchingInputs(in []MatchingPairInput) []questioncommands.MatchingPairInput {
	out := make([]questioncommands.MatchingPairInput, len(in))
	for i := range in {
		out[i] = questioncommands.MatchingPairInput{Prompt: in[i].Prompt, Match: in[i].Match}
	}
	return out
}
func toTypedInputs(in []TypedBlankInput) []questioncommands.TypedBlankInput {
	out := make([]questioncommands.TypedBlankInput, len(in))
	for i := range in {
		out[i] = questioncommands.TypedBlankInput{Placeholder: in[i].Placeholder, Variants: in[i].Variants}
	}
	return out
}
func toShortInputs(in []ShortVariantInput) []questioncommands.ShortVariantInput {
	out := make([]questioncommands.ShortVariantInput, len(in))
	for i := range in {
		out[i] = questioncommands.ShortVariantInput{Text: in[i].Text}
	}
	return out
}

var _ = time.Now
