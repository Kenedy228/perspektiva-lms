package handlers

import (
	"net/http"
	"time"

	coursequeries "gitflic.ru/lms/backend/internal/application/usecases/course/queries"
	enrollmentcommands "gitflic.ru/lms/backend/internal/application/usecases/enrollment/commands"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/transport/http/response"
)

func (api *API) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req EnrollmentRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	activatedAt, err := parseOptionalDate(req.ActivatedAt)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	deactivatedAt, err := parseOptionalDate(req.DeactivatedAt)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	out, err := api.Enrollments.Create.Execute(r.Context(), enrollmentcommands.CreateInput{
		ActorRole:     actor.role,
		CourseID:      req.CourseID,
		VersionID:     req.VersionID,
		AccountID:     req.AccountID,
		ActivatedAt:   activatedAt,
		DeactivatedAt: deactivatedAt,
		Now:           time.Now().UTC(),
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/enrollments", out.ID)
}

func (api *API) ListStudentStatistics(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	limit, offset := limitOffset(r)
	accountID := actor.accountID
	organizationID := r.URL.Query().Get("organization_id")
	if actor.role.Kind() == role.TypeAdmin {
		accountID = r.URL.Query().Get("account_id")
	}
	out, err := api.Courses.Statistics.Execute(r.Context(), coursequeries.StudentStatisticsInput{
		ActorRole:      actor.role,
		AccountID:      accountID,
		OrganizationID: organizationID,
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out.Views, response.Links{"self": {Href: r.URL.RequestURI(), Method: http.MethodGet}})
}

func parseOptionalDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return parseDate(value)
}
