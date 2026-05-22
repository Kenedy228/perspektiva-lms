package handlers

import (
	"net/http"

	orgcommands "gitflic.ru/lms/backend/internal/application/usecases/organization/commands"
	orgqueries "gitflic.ru/lms/backend/internal/application/usecases/organization/queries"
	"gitflic.ru/lms/backend/internal/transport/http/response"
)

func (api *API) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	limit, offset := limitOffset(r)
	var data any
	var err error
	if inn := r.URL.Query().Get("inn"); inn != "" {
		out, e := api.Organizations.ListINN.Execute(r.Context(), orgqueries.ListByINNInput{ActorRole: actor.role, INN: inn, Limit: limit, Offset: offset})
		if e == nil {
			data = out.Views
		}
		err = e
	} else {
		out, e := api.Organizations.ListName.Execute(r.Context(), orgqueries.ListByNameInput{ActorRole: actor.role, Name: r.URL.Query().Get("name"), Limit: limit, Offset: offset})
		if e == nil {
			data = out.Views
		}
		err = e
	}
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, data, response.Links{"self": {Href: r.URL.RequestURI(), Method: http.MethodGet}, "create": {Href: "/organizations", Method: http.MethodPost}})
}

func (api *API) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req OrganizationRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Organizations.Create.Execute(r.Context(), orgcommands.CreateInput{ActorRole: actor.role, Name: req.Name, INN: req.INN, INNType: req.INNType})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/organizations", out.ID)
}

func (api *API) GetOrganization(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Organizations.Get.Execute(r.Context(), orgqueries.GetDetailsByIDInput{ActorRole: actor.role, ID: r.PathValue("id")})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out.View, response.Links{"self": {Href: r.URL.Path, Method: http.MethodGet}, "rename": {Href: r.URL.Path, Method: http.MethodPatch}, "change_inn": {Href: r.URL.Path + "/inn", Method: http.MethodPatch}})
}

func (api *API) RenameOrganization(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req OrganizationRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Organizations.Rename.Execute(r.Context(), orgcommands.RenameInput{ActorRole: actor.role, OrganizationID: r.PathValue("id"), Name: req.Name})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) ChangeOrganizationINN(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req OrganizationRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Organizations.ChangeINN.Execute(r.Context(), orgcommands.ChangeINNInput{ActorRole: actor.role, OrganizationID: r.PathValue("id"), INN: req.INN, INNType: req.INNType})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := api.Organizations.Delete.Execute(r.Context(), orgcommands.DeleteByIDInput{ActorRole: actor.role, OrganizationID: r.PathValue("id")}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeNoContent(w)
}
