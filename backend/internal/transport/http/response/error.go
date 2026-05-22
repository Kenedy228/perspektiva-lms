package response

import (
	"net/http"
)

// Error is the standard API error response body.
type Error struct {
	Status  int            `json:"status"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
	Links   Links          `json:"_links,omitempty"`
}

// ErrorEnvelope wraps an error using the same top-level contract as successful responses.
type ErrorEnvelope struct {
	Error Error `json:"error"`
}

// ErrorOption customizes an Error.
type ErrorOption func(*Error)

// NewError creates a response error.
func NewError(status int, code, message string, opts ...ErrorOption) Error {
	err := Error{
		Status:  status,
		Code:    code,
		Message: localizedMessage(code, message),
	}
	for _, opt := range opts {
		opt(&err)
	}
	return err
}

func localizedMessage(code, fallback string) string {
	switch code {
	case "unauthorized":
		return "Требуется авторизация."
	case "invalid_token":
		return "Сессионный токен недействителен."
	case "expired_token":
		return "Срок действия сессии истек."
	case "invalid_credentials":
		return "Неверный логин или пароль."
	case "forbidden":
		return "Операция запрещена."
	case "not_found":
		return "Ресурс не найден."
	case "invalid_json":
		return "Тело запроса содержит некорректный JSON."
	case "invalid_input":
		return "Некорректные данные запроса."
	case "invalid_date":
		return "Дата должна быть указана в формате ГГГГ-ММ-ДД."
	case "conflict":
		return "Запрос конфликтует с текущим состоянием ресурса."
	case "internal_error":
		return "Внутренняя ошибка сервера."
	case "http_error":
		return "HTTP-запрос завершился ошибкой."
	default:
		if fallback != "" {
			return fallback
		}
		return "Запрос не может быть обработан."
	}
}

// WithDetail adds a structured error detail.
func WithDetail(key string, value any) ErrorOption {
	return func(err *Error) {
		if err.Details == nil {
			err.Details = make(map[string]any)
		}
		err.Details[key] = value
	}
}

// WithLinks adds HATEOAS links to an error response.
func WithLinks(links Links) ErrorOption {
	return func(err *Error) {
		err.Links = links
	}
}

// WriteError writes a standard error response.
func WriteError(w http.ResponseWriter, r *http.Request, err Error) {
	if err.Links == nil {
		err.Links = Links{"self": {Href: r.URL.Path, Method: r.Method}}
	}
	WriteJSON(w, err.Status, ErrorEnvelope{Error: err})
}
