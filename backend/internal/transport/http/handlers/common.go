package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	accountcommon "gitflic.ru/lms/backend/internal/application/usecases/account/common"
	attemptcommon "gitflic.ru/lms/backend/internal/application/usecases/attempt/common"
	bankcommon "gitflic.ru/lms/backend/internal/application/usecases/bank/common"
	coursecommon "gitflic.ru/lms/backend/internal/application/usecases/course/common"
	enrollmentcommon "gitflic.ru/lms/backend/internal/application/usecases/enrollment/common"
	orgcommon "gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	personcommon "gitflic.ru/lms/backend/internal/application/usecases/person/common"
	questioncommon "gitflic.ru/lms/backend/internal/application/usecases/question/common"
	questiongrading "gitflic.ru/lms/backend/internal/application/usecases/question/grading"
	quizcommon "gitflic.ru/lms/backend/internal/application/usecases/quiz/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/transport/http/middleware"
	"gitflic.ru/lms/backend/internal/transport/http/response"
)

func actorRole(r *http.Request) (roleValue roleCarrier, ok bool) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		return roleCarrier{}, false
	}
	return roleCarrier{
		accountID: claims.AccountID.String(),
		personID:  claims.PersonID.String(),
		role:      claims.Role,
	}, true
}

type roleCarrier struct {
	accountID string
	personID  string
	role      role.Role
}

func limitOffset(r *http.Request) (int, int) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	return limit, offset
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", value)
}

func writeCreated(w http.ResponseWriter, path, id string) {
	response.WriteJSON(w, http.StatusCreated, response.Envelope{
		Data: map[string]string{"id": id},
		Links: response.Links{
			"self": response.Link{Href: path + "/" + id, Method: http.MethodGet},
		},
	})
}

func writeOK(w http.ResponseWriter, r *http.Request, data any, links response.Links) {
	if links == nil {
		links = response.Links{"self": response.Link{Href: r.URL.Path, Method: r.Method}}
	}
	response.WriteJSON(w, http.StatusOK, response.Envelope{Data: data, Links: links})
}

func writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeHandlerError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	status := http.StatusInternalServerError
	code := "internal_error"
	message := "Запрос не может быть обработан."

	switch {
	case errors.Is(err, sql.ErrNoRows):
		status, code, message = http.StatusNotFound, "not_found", "Ресурс не найден."
	case errors.Is(err, accountcommon.ErrForbidden),
		errors.Is(err, orgcommon.ErrForbidden),
		errors.Is(err, personcommon.ErrForbidden),
		errors.Is(err, bankcommon.ErrForbidden),
		errors.Is(err, questioncommon.ErrForbidden),
		errors.Is(err, quizcommon.ErrForbidden),
		errors.Is(err, attemptcommon.ErrForbidden),
		errors.Is(err, coursecommon.ErrForbidden),
		errors.Is(err, enrollmentcommon.ErrForbidden):
		status, code, message = http.StatusForbidden, "forbidden", "Операция запрещена."
	case errors.Is(err, accountcommon.ErrInvalidInput),
		errors.Is(err, orgcommon.ErrInvalidInput),
		errors.Is(err, personcommon.ErrInvalidInput),
		errors.Is(err, bankcommon.ErrInvalidInput),
		errors.Is(err, questioncommon.ErrInvalidInput),
		errors.Is(err, quizcommon.ErrInvalidInput),
		errors.Is(err, attemptcommon.ErrInvalidInput),
		errors.Is(err, coursecommon.ErrInvalidInput),
		errors.Is(err, enrollmentcommon.ErrInvalidInput),
		errors.Is(err, questiongrading.ErrInvalidInput):
		status, code, message = http.StatusBadRequest, "invalid_input", "Некорректные данные запроса."
	case errors.Is(err, attemptcommon.ErrLimitReached), errors.Is(err, enrollmentcommon.ErrConflict):
		status, code, message = http.StatusConflict, "conflict", "Запрос конфликтует с текущим состоянием ресурса."
	}

	response.WriteError(w, r, response.NewError(status, code, message))
}
