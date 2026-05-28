package handlers

import (
	"net/http"
	"strconv"

	coursecommands "gitflic.ru/lms/backend/internal/application/usecases/course/commands"
	coursequeries "gitflic.ru/lms/backend/internal/application/usecases/course/queries"
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

func (api *API) RemoveCourseBlock(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := api.Courses.RemoveBlock.Execute(r.Context(), coursecommands.RemoveBlockFromCourseInput{
		ActorRole: actor.role,
		CourseID:  r.PathValue("courseID"),
		BlockID:   r.PathValue("blockID"),
	}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeNoContent(w)
}

func (api *API) MoveCourseBlock(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req MoveBlockRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if err := api.Courses.MoveBlock.Execute(r.Context(), coursecommands.MoveCourseBlockInput{
		ActorRole: actor.role,
		CourseID:  r.PathValue("courseID"),
		From:      req.From,
		To:        req.To,
	}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"course_id": r.PathValue("courseID")}, nil)
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

func (api *API) RemoveBlockElement(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := api.Courses.RemoveElement.Execute(r.Context(), coursecommands.RemoveElementFromBlockInput{
		ActorRole: actor.role,
		BlockID:   r.PathValue("blockID"),
		ElementID: r.PathValue("elementID"),
	}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeNoContent(w)
}

func (api *API) MoveBlockElement(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req MoveElementRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if err := api.Courses.MoveElement.Execute(r.Context(), coursecommands.MoveBlockElementInput{
		ActorRole: actor.role,
		BlockID:   r.PathValue("blockID"),
		From:      req.From,
		To:        req.To,
	}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"block_id": r.PathValue("blockID")}, nil)
}

func (api *API) ChangeElementCompletionMode(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req ChangeCompletionModeRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	if err := api.Courses.ChangeCompletionMode.Execute(r.Context(), coursecommands.ChangeElementCompletionModeInput{
		ActorRole:      actor.role,
		ElementID:      r.PathValue("elementID"),
		CompletionMode: req.CompletionMode,
	}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"element_id": r.PathValue("elementID")}, nil)
}

func (api *API) GetCourseProgress(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	enrollmentID := r.URL.Query().Get("enrollment_id")
	total, _ := strconv.Atoi(r.URL.Query().Get("total"))
	out, err := api.Courses.GetProgress.Execute(r.Context(), coursecommands.GetProgressInput{
		ActorRole:         actor.role,
		ActorPersonID:     actor.personID,
		EnrollmentID:      enrollmentID,
		TotalTrackedItems: total,
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	elementStrs := make([]string, len(out.CompletedElementIDs))
	for i, eid := range out.CompletedElementIDs {
		elementStrs[i] = eid.String()
	}
	writeOK(w, r, ProgressResponse{
		EnrollmentID:        enrollmentID,
		CompletedCount:      out.CompletedCount,
		Percent:             out.Percent,
		TotalTrackedItems:   out.TotalTrackedItems,
		CompletedElementIDs: elementStrs,
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

func (api *API) UploadElementContent(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_multipart", "failed to parse multipart form"))
		return
	}
	file, header, err := r.FormFile("content")
	if err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "missing_file", "content file is required"))
		return
	}
	defer file.Close()
	if err := api.Courses.UploadContent.Execute(r.Context(), coursecommands.UploadElementContentInput{
		ActorRole:   actor.role,
		ElementID:   r.PathValue("elementID"),
		ContentType: header.Header.Get("Content-Type"),
		Body:        file,
		Size:        header.Size,
	}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"element_id": r.PathValue("elementID")}, nil)
}

func (api *API) DownloadElementContent(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Courses.DownloadContent.Execute(r.Context(), coursecommands.DownloadElementContentInput{
		ActorRole: actor.role,
		ElementID: r.PathValue("elementID"),
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{
		"element_id":   r.PathValue("elementID"),
		"download_url": out.DownloadURL,
		"file_name":    out.FileName,
		"content_type": out.ContentType,
	}, nil)
}
