package handlers

import (
	"context"
	"net/http"
	"strconv"

	bankcommands "gitflic.ru/lms/backend/internal/application/usecases/bank/commands"
	bankqueries "gitflic.ru/lms/backend/internal/application/usecases/bank/queries"
	"gitflic.ru/lms/backend/internal/transport/http/response"
)

func (api *API) ListBanks(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	limit, offset := limitOffset(r)
	minQ, _ := strconv.Atoi(r.URL.Query().Get("min_questions"))
	maxQ, _ := strconv.Atoi(r.URL.Query().Get("max_questions"))
	out, err := api.Banks.List.Execute(r.Context(), bankqueries.ListInput{
		ActorRole: actor.role, TitleContains: r.URL.Query().Get("title"), QuestionID: r.URL.Query().Get("question_id"),
		MinQuestions: minQ, MaxQuestions: maxQ, Limit: limit, Offset: offset,
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out.Views, response.Links{"self": {Href: r.URL.RequestURI(), Method: http.MethodGet}, "create": {Href: "/banks", Method: http.MethodPost}})
}

func (api *API) CreateBank(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req BankRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Banks.Create.Execute(r.Context(), bankcommands.CreateInput{ActorRole: actor.role, Title: req.Title})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/banks", out.ID)
}

func (api *API) GetBank(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Banks.Get.Execute(r.Context(), bankqueries.GetDetailsByIDInput{ActorRole: actor.role, BankID: r.PathValue("id")})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out, response.Links{"self": {Href: r.URL.Path, Method: http.MethodGet}, "questions": {Href: r.URL.Path + "/questions", Method: http.MethodPost}})
}

func (api *API) RenameBank(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req BankRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if err := api.Banks.Rename.Execute(r.Context(), bankcommands.RenameInput{ActorRole: actor.role, BankID: r.PathValue("id"), Title: req.Title}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) AddBankQuestions(w http.ResponseWriter, r *http.Request) {
	api.changeBankQuestions(w, r, api.Banks.Add.Execute)
}

func (api *API) RemoveBankQuestions(w http.ResponseWriter, r *http.Request) {
	api.changeBankQuestions(w, r, api.Banks.Remove.Execute)
}

func (api *API) changeBankQuestions(w http.ResponseWriter, r *http.Request, fn func(context.Context, bankcommands.QuestionIDsInput) error) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req QuestionIDsRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if err := fn(r.Context(), bankcommands.QuestionIDsInput{ActorRole: actor.role, BankID: r.PathValue("id"), QuestionIDs: req.QuestionIDs}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}
