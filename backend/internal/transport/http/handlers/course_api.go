package handlers

import (
	"net/http"
	"time"

	coursecommands "gitflic.ru/lms/backend/internal/application/usecases/course/commands"
	coursequeries "gitflic.ru/lms/backend/internal/application/usecases/course/queries"
	"gitflic.ru/lms/backend/internal/domain/course/progress"
	"gitflic.ru/lms/backend/internal/transport/http/response"
	"github.com/google/uuid"
)

func (api *API) ListCourses(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	limit, offset := limitOffset(r)
	out, err := api.Courses.List.Execute(r.Context(), coursequeries.ListInput{
		ActorRole:     actor.role,
		AccountID:     actor.accountID,
		TitleContains: r.URL.Query().Get("title"),
		Status:        r.URL.Query().Get("status"),
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out.Views, response.Links{"self": {Href: r.URL.RequestURI(), Method: http.MethodGet}})
}

func (api *API) CreateCourse(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req CourseRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Courses.Create.Execute(r.Context(), coursecommands.CreateCourseInput{ActorRole: actor.role, Title: req.Title})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/courses", out.ID)
}

func (api *API) GetCourse(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	view, err := api.Courses.Query.GetDetailsByID(r.Context(), id)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, view, response.Links{
		"self":    {Href: r.URL.Path, Method: http.MethodGet},
		"ratings": {Href: r.URL.Path + "/ratings", Method: http.MethodGet},
	})
}

func (api *API) RenameCourse(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req CourseRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if err := api.Courses.Rename.Execute(r.Context(), coursecommands.RenameCourseInput{ActorRole: actor.role, CourseID: r.PathValue("id"), Title: req.Title}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) AddCourseBlock(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req CourseBlockRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Courses.AddBlock.Execute(r.Context(), coursecommands.AddBlockToCourseInput{
		ActorRole: actor.role,
		CourseID:  r.PathValue("courseID"),
		Title:     req.Title,
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/course-blocks", out.ID)
}

func (api *API) AddBlockElement(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req CourseElementRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Courses.AddElement.Execute(r.Context(), coursecommands.AddElementToBlockInput{
		ActorRole: actor.role,
		BlockID:   r.PathValue("blockID"),
		Title:     req.Title,
		Content: coursecommands.ElementContentInput{
			Type:           req.Type,
			FileName:       req.FileName,
			SizeBytes:      req.SizeBytes,
			QuizID:         req.QuizID,
			CompletionMode: req.CompletionMode,
		},
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/course-elements", out.ID)
}

func (api *API) MarkCourseProgress(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req MarkProgressRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	at := time.Now().UTC()
	if req.At != "" {
		parsed, err := parseDate(req.At)
		if err != nil {
			writeHandlerError(w, r, err)
			return
		}
		at = parsed
	}
	if err := api.Courses.Progress.Execute(r.Context(), coursecommands.MarkProgressInput{
		ActorRole:    actor.role,
		EnrollmentID: req.EnrollmentID,
		ElementID:    req.ElementID,
		MarkerType:   progress.MarkerType(req.MarkerType),
		At:           at,
	}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{
		"course_id":     r.PathValue("courseID"),
		"enrollment_id": req.EnrollmentID,
		"element_id":    req.ElementID,
	}, nil)
}

func (api *API) ListCourseRatings(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	limit, offset := limitOffset(r)
	out, err := api.Courses.Ratings.Execute(r.Context(), coursequeries.RatingsInput{ActorRole: actor.role, CourseID: r.PathValue("id"), Limit: limit, Offset: offset})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out.Views, response.Links{"self": {Href: r.URL.RequestURI(), Method: http.MethodGet}})
}
