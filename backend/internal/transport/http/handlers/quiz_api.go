package handlers

import (
	"net/http"

	quizcommands "gitflic.ru/lms/backend/internal/application/usecases/quiz/commands"
	quizcommon "gitflic.ru/lms/backend/internal/application/usecases/quiz/common"
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"gitflic.ru/lms/backend/internal/transport/http/response"
	"github.com/google/uuid"
)

func (api *API) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req QuizRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Quizzes.Create.Execute(r.Context(), quizcommands.CreateInput{
		ActorRole:        actor.role,
		Title:            req.Title,
		MaxAttempts:      req.MaxAttempts,
		TimeLimitSeconds: req.TimeLimitSeconds,
		ShuffleQuestions: req.ShuffleQuestions,
		Sources:          toQuizSourceInputs(req.Sources),
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/quizzes", out.ID)
}

func (api *API) GetQuiz(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	q, err := api.Quizzes.Repository.FindByID(r.Context(), id)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]any{
		"id":                 q.ID().String(),
		"title":              q.Title().Value(),
		"max_attempts":       q.Attempts().Count(),
		"time_limit_seconds": q.Time().Seconds(),
		"shuffle_questions":  q.ShuffleQuestions(),
		"sources":            quizSourcesView(q.Sources()),
	}, response.Links{
		"self":    {Href: r.URL.Path, Method: http.MethodGet},
		"attempt": {Href: "/attempts", Method: http.MethodPost},
	})
}

func (api *API) RenameQuiz(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req QuizRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if err := api.Quizzes.Rename.Execute(r.Context(), quizcommands.RenameInput{ActorRole: actor.role, QuizID: r.PathValue("id"), Title: req.Title}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) ChangeQuizLimits(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req QuizRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	err := api.Quizzes.ChangeLimits.Execute(r.Context(), quizcommands.ChangeLimitsInput{
		ActorRole:        actor.role,
		QuizID:           r.PathValue("id"),
		MaxAttempts:      req.MaxAttempts,
		TimeLimitSeconds: req.TimeLimitSeconds,
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) ChangeQuizShuffle(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req QuizRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if err := api.Quizzes.ChangeShuffle.Execute(r.Context(), quizcommands.ChangeShufflePolicyInput{ActorRole: actor.role, QuizID: r.PathValue("id"), ShuffleQuestions: req.ShuffleQuestions}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) ReplaceQuizSources(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req QuizRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if err := api.Quizzes.Replace.Execute(r.Context(), quizcommands.ReplaceSourcesInput{ActorRole: actor.role, QuizID: r.PathValue("id"), Sources: toQuizSourceInputs(req.Sources)}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := quizcommon.RequireManager(actor.role); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	if err := api.Quizzes.Repository.DeleteByID(r.Context(), id); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeNoContent(w)
}

func toQuizSourceInputs(in []QuizSource) []quizcommands.SourceInput {
	out := make([]quizcommands.SourceInput, len(in))
	for i := range in {
		out[i] = quizcommands.SourceInput{
			BankID:        in[i].BankID,
			CriteriaType:  in[i].CriteriaType,
			QuestionCount: in[i].QuestionCount,
			QuestionIDs:   in[i].QuestionIDs,
		}
	}
	return out
}

func quizSourcesView(in []source.Source) []QuizSource {
	out := make([]QuizSource, 0, len(in))
	for i := range in {
		view := QuizSource{
			BankID:        in[i].BankID().String(),
			CriteriaType:  in[i].Criteria().Type().String(),
			QuestionCount: in[i].Criteria().QuestionCount(),
		}
		if manual, ok := in[i].Criteria().(criteria.Manual); ok {
			ids := manual.QuestionIDs()
			view.QuestionIDs = make([]string, 0, len(ids))
			for _, id := range ids {
				view.QuestionIDs = append(view.QuestionIDs, id.String())
			}
		}
		out = append(out, view)
	}
	return out
}
