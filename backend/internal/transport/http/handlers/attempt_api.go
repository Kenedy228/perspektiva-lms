package handlers

import (
	"net/http"
	"time"

	attemptcommands "gitflic.ru/lms/backend/internal/application/usecases/attempt/commands"
	"gitflic.ru/lms/backend/internal/transport/http/response"
	"github.com/google/uuid"
)

func (api *API) StartAttempt(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req AttemptStartRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Attempts.Start.Execute(r.Context(), attemptcommands.StartInput{
		ActorRole:    actor.role,
		AccountID:    req.AccountID,
		EnrollmentID: req.EnrollmentID,
		QuizID:       req.QuizID,
		Now:          time.Now().UTC(),
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/attempts", out.ID)
}

func (api *API) GetAttempt(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	a, err := api.Attempts.Repository.FindByID(r.Context(), id)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]any{
		"id":              a.ID().String(),
		"enrollment_id":   a.EnrollmentID().String(),
		"quiz_id":         a.QuizID().String(),
		"status":          a.Status().String(),
		"started_at":      a.StartedAt(),
		"deadline_at":     a.DeadlineAt(),
		"finished_at":     a.FinishedAt(),
		"questions_count": a.CountItems(),
		"answers_count":   a.CountAnswers(),
	}, response.Links{
		"self":   {Href: r.URL.Path, Method: http.MethodGet},
		"finish": {Href: r.URL.Path + "/finish", Method: http.MethodPost},
		"cancel": {Href: r.URL.Path + "/cancel", Method: http.MethodPost},
	})
}

func (api *API) AddAttemptAnswer(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
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
	err = api.Attempts.Answer.Execute(r.Context(), attemptcommands.AddAnswerInput{
		ActorRole:  actor.role,
		AttemptID:  r.PathValue("id"),
		QuestionID: r.PathValue("questionID"),
		Answer:     ans,
		AnsweredAt: time.Now().UTC(),
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) FinishAttempt(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := api.Attempts.Finish.Execute(r.Context(), attemptcommands.FinishInput{ActorRole: actor.role, AttemptID: r.PathValue("id"), FinishedAt: time.Now().UTC()}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) CancelAttempt(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := api.Attempts.Cancel.Execute(r.Context(), attemptcommands.CancelInput{ActorRole: actor.role, AttemptID: r.PathValue("id"), CancelledAt: time.Now().UTC()}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}
