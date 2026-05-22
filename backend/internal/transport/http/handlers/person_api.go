package handlers

import (
	"net/http"
	"time"

	personcommands "gitflic.ru/lms/backend/internal/application/usecases/person/commands"
	personqueries "gitflic.ru/lms/backend/internal/application/usecases/person/queries"
	"gitflic.ru/lms/backend/internal/transport/http/response"
)

func (api *API) ListPersons(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	limit, offset := limitOffset(r)
	var data any
	var err error
	switch {
	case r.URL.Query().Get("organization_id") != "":
		out, e := api.Persons.ListOrg.Execute(r.Context(), personqueries.ListByOrganizationIDInput{ActorRole: actor.role, OrganizationID: r.URL.Query().Get("organization_id"), Limit: limit, Offset: offset})
		if e == nil {
			data = out.Views
		}
		err = e
	case r.URL.Query().Get("snils") != "":
		out, e := api.Persons.ListSNILS.Execute(r.Context(), personqueries.ListBySnilsInput{ActorRole: actor.role, Snils: r.URL.Query().Get("snils"), Limit: limit, Offset: offset})
		if e == nil {
			data = out.Views
		}
		err = e
	default:
		out, e := api.Persons.ListLastName.Execute(r.Context(), personqueries.ListByLastnameInput{ActorRole: actor.role, LastName: r.URL.Query().Get("last_name"), Limit: limit, Offset: offset})
		if e == nil {
			data = out.Views
		}
		err = e
	}
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, data, response.Links{"self": {Href: r.URL.RequestURI(), Method: http.MethodGet}, "create": {Href: "/persons", Method: http.MethodPost}})
}

func (api *API) CreatePerson(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req struct {
		PersonRequest
		Profile *ProfileRequest `json:"profile,omitempty"`
	}
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if req.Profile != nil {
		dob, err := parseDate(req.Profile.DateOfBirth)
		if err != nil {
			writeHandlerError(w, r, err)
			return
		}
		out, err := api.Persons.CreateWithProfile.Execute(r.Context(), personcommands.CreateWithProfileInput{ActorRole: actor.role, FirstName: req.FirstName, LastName: req.LastName, MiddleName: req.MiddleName, Snils: req.Profile.SNILS, DateOfBirth: dob, JobTitle: req.Profile.JobTitle, Education: req.Profile.Education, OrganizationID: req.Profile.OrganizationID})
		if err != nil {
			writeHandlerError(w, r, err)
			return
		}
		writeCreated(w, "/persons", out.ID)
		return
	}
	out, err := api.Persons.Create.Execute(r.Context(), personcommands.CreateInput{ActorRole: actor.role, FirstName: req.FirstName, LastName: req.LastName, MiddleName: req.MiddleName})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/persons", out.ID)
}

func (api *API) GetPerson(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Persons.Get.Execute(r.Context(), personqueries.GetDetailsByIDInput{ActorRole: actor.role, ID: r.PathValue("id")})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out.View, response.Links{"self": {Href: r.URL.Path, Method: http.MethodGet}, "profile": {Href: r.URL.Path + "/profile", Method: http.MethodPut}})
}

func (api *API) RenamePerson(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req PersonRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Persons.Rename.Execute(r.Context(), personcommands.RenameInput{ActorRole: actor.role, PersonID: r.PathValue("id"), FirstName: req.FirstName, LastName: req.LastName, MiddleName: req.MiddleName})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) ReplacePersonProfile(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	req, dob, ok := decodeProfile(w, r)
	if !ok {
		return
	}
	out, err := api.Persons.ReplaceProfile.Execute(r.Context(), personcommands.ReplaceProfileInput{ActorRole: actor.role, PersonID: r.PathValue("id"), Snils: req.SNILS, DateOfBirth: dob, JobTitle: req.JobTitle, Education: req.Education, OrganizationID: req.OrganizationID})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) ChangePersonProfile(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	req, dob, ok := decodeProfile(w, r)
	if !ok {
		return
	}
	if req.SNILS != "" {
		if _, err := api.Persons.ChangeSNILS.Execute(r.Context(), personcommands.ChangeSNILSInput{ActorRole: actor.role, PersonID: r.PathValue("id"), Snils: req.SNILS}); err != nil {
			writeHandlerError(w, r, err)
			return
		}
	}
	if !dob.IsZero() {
		if _, err := api.Persons.ChangeDOB.Execute(r.Context(), personcommands.ChangeDateOfBirthInput{ActorRole: actor.role, PersonID: r.PathValue("id"), DateOfBirth: dob}); err != nil {
			writeHandlerError(w, r, err)
			return
		}
	}
	if req.JobTitle != "" {
		if _, err := api.Persons.ChangeJobTitle.Execute(r.Context(), personcommands.ChangeJobTitleInput{ActorRole: actor.role, PersonID: r.PathValue("id"), JobTitle: req.JobTitle}); err != nil {
			writeHandlerError(w, r, err)
			return
		}
	}
	if req.Education != "" {
		if _, err := api.Persons.ChangeEducation.Execute(r.Context(), personcommands.ChangeEducationInput{ActorRole: actor.role, PersonID: r.PathValue("id"), Education: req.Education}); err != nil {
			writeHandlerError(w, r, err)
			return
		}
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) DetachPersonProfile(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Persons.DetachProfile.Execute(r.Context(), personcommands.DetachProfileInput{ActorRole: actor.role, PersonID: r.PathValue("id")})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.PersonID}, nil)
}

func (api *API) AssignPersonOrganization(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req ProfileRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Persons.AssignOrg.Execute(r.Context(), personcommands.AssignOrganizationInput{ActorRole: actor.role, PersonID: r.PathValue("id"), OrganizationID: req.OrganizationID})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) RemovePersonOrganization(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Persons.RemoveOrg.Execute(r.Context(), personcommands.RemoveOrganizationInput{ActorRole: actor.role, PersonID: r.PathValue("id")})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) DeletePerson(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := api.Persons.Delete.Execute(r.Context(), personcommands.DeleteByIDInput{ActorRole: actor.role, PersonID: r.PathValue("id")}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeNoContent(w)
}

func decodeProfile(w http.ResponseWriter, r *http.Request) (ProfileRequest, time.Time, bool) {
	var req ProfileRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return req, time.Time{}, false
	}
	dob, err := parseDate(req.DateOfBirth)
	if err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_date", "date_of_birth must use YYYY-MM-DD"))
		return req, time.Time{}, false
	}
	return req, dob, true
}
