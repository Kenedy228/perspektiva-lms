package http

import (
	"log/slog"
	"net/http"

	"gitflic.ru/lms/backend/internal/transport/http/handlers"
	"gitflic.ru/lms/backend/internal/transport/http/middleware"
	"gitflic.ru/lms/backend/internal/transport/http/response"
	"gitflic.ru/lms/backend/internal/transport/http/session"
)

// ServerConfig contains dependencies needed by the HTTP composition root.
type ServerConfig struct {
	Logger   *slog.Logger
	Session  *session.Manager
	Handlers Handlers
}

// Handlers groups endpoint handlers registered by the HTTP server.
type Handlers struct {
	Auth *handlers.AuthHandler
	API  *handlers.API
}

// Server owns route registration and presentation middleware composition.
type Server struct {
	logger   *slog.Logger
	session  *session.Manager
	handlers Handlers
}

// NewServer creates the HTTP presentation server.
func NewServer(cfg ServerConfig) *Server {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	if cfg.Session == nil {
		panic("http server requires session manager")
	}
	if cfg.Handlers.Auth == nil {
		panic("http server requires auth handler")
	}
	if cfg.Handlers.API == nil {
		panic("http server requires api handlers")
	}
	return &Server{
		logger:   logger,
		session:  cfg.Session,
		handlers: cfg.Handlers,
	}
}

// Handler returns the composed HTTP handler with all routes and middleware.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /", s.index)

	mux.HandleFunc("POST /auth/login", s.handlers.Auth.Login)
	mux.HandleFunc("POST /auth/logout", s.requireAuth(s.handlers.Auth.Logout))
	mux.HandleFunc("GET /auth/session", s.requireAuth(s.handlers.Auth.Session))

	api := s.handlers.API
	mux.HandleFunc("GET /accounts", s.requireAuth(api.ListAccounts))
	mux.HandleFunc("POST /accounts", s.requireAuth(api.CreateAccount))
	mux.HandleFunc("GET /accounts/{id}", s.requireAuth(api.GetAccount))
	mux.HandleFunc("PATCH /accounts/{id}/login", s.requireAuth(api.ChangeAccountLogin))
	mux.HandleFunc("PATCH /accounts/{id}/password", s.requireAuth(api.ChangeAccountPassword))
	mux.HandleFunc("PATCH /accounts/{id}/role", s.requireAuth(api.ChangeAccountRole))
	mux.HandleFunc("POST /accounts/{id}/block", s.requireAuth(api.BlockAccount))
	mux.HandleFunc("POST /accounts/{id}/activate", s.requireAuth(api.ActivateAccount))
	mux.HandleFunc("DELETE /accounts/{id}", s.requireAuth(api.DeleteAccount))

	mux.HandleFunc("GET /organizations", s.requireAuth(api.ListOrganizations))
	mux.HandleFunc("POST /organizations", s.requireAuth(api.CreateOrganization))
	mux.HandleFunc("GET /organizations/{id}", s.requireAuth(api.GetOrganization))
	mux.HandleFunc("PATCH /organizations/{id}", s.requireAuth(api.RenameOrganization))
	mux.HandleFunc("PATCH /organizations/{id}/inn", s.requireAuth(api.ChangeOrganizationINN))
	mux.HandleFunc("DELETE /organizations/{id}", s.requireAuth(api.DeleteOrganization))

	mux.HandleFunc("GET /persons", s.requireAuth(api.ListPersons))
	mux.HandleFunc("POST /persons", s.requireAuth(api.CreatePerson))
	mux.HandleFunc("GET /persons/{id}", s.requireAuth(api.GetPerson))
	mux.HandleFunc("PATCH /persons/{id}", s.requireAuth(api.RenamePerson))
	mux.HandleFunc("PUT /persons/{id}/profile", s.requireAuth(api.ReplacePersonProfile))
	mux.HandleFunc("PATCH /persons/{id}/profile", s.requireAuth(api.ChangePersonProfile))
	mux.HandleFunc("DELETE /persons/{id}/profile", s.requireAuth(api.DetachPersonProfile))
	mux.HandleFunc("PUT /persons/{id}/organization", s.requireAuth(api.AssignPersonOrganization))
	mux.HandleFunc("DELETE /persons/{id}/organization", s.requireAuth(api.RemovePersonOrganization))
	mux.HandleFunc("DELETE /persons/{id}", s.requireAuth(api.DeletePerson))

	mux.HandleFunc("GET /banks", s.requireAuth(api.ListBanks))
	mux.HandleFunc("POST /banks", s.requireAuth(api.CreateBank))
	mux.HandleFunc("GET /banks/{id}", s.requireAuth(api.GetBank))
	mux.HandleFunc("PATCH /banks/{id}", s.requireAuth(api.RenameBank))
	mux.HandleFunc("DELETE /banks/{id}", s.requireAuth(api.DeleteBank))
	mux.HandleFunc("POST /banks/{id}/questions", s.requireAuth(api.AddBankQuestions))
	mux.HandleFunc("DELETE /banks/{id}/questions", s.requireAuth(api.RemoveBankQuestions))

	mux.HandleFunc("POST /questions", s.requireAuth(api.CreateQuestion))
	mux.HandleFunc("GET /questions/{id}", s.requireAuth(api.GetQuestion))
	mux.HandleFunc("PATCH /questions/{id}", s.requireAuth(api.ChangeQuestionTitle))
	mux.HandleFunc("PUT /questions/{id}/content", s.requireAuth(api.ChangeQuestionContent))
	mux.HandleFunc("DELETE /questions/{id}", s.requireAuth(api.DeleteQuestion))
	mux.HandleFunc("POST /questions/{id}/grade", s.requireAuth(api.GradeQuestion))

	mux.HandleFunc("POST /quizzes", s.requireAuth(api.CreateQuiz))
	mux.HandleFunc("GET /quizzes/{id}", s.requireAuth(api.GetQuiz))
	mux.HandleFunc("PATCH /quizzes/{id}", s.requireAuth(api.RenameQuiz))
	mux.HandleFunc("PATCH /quizzes/{id}/limits", s.requireAuth(api.ChangeQuizLimits))
	mux.HandleFunc("PATCH /quizzes/{id}/shuffle", s.requireAuth(api.ChangeQuizShuffle))
	mux.HandleFunc("PUT /quizzes/{id}/sources", s.requireAuth(api.ReplaceQuizSources))
	mux.HandleFunc("DELETE /quizzes/{id}", s.requireAuth(api.DeleteQuiz))

	mux.HandleFunc("GET /attempts", s.requireAuth(api.ListAttempts))
	mux.HandleFunc("POST /attempts", s.requireAuth(api.StartAttempt))
	mux.HandleFunc("GET /attempts/{id}", s.requireAuth(api.GetAttempt))
	mux.HandleFunc("PUT /attempts/{id}/answers/{questionID}", s.requireAuth(api.AddAttemptAnswer))
	mux.HandleFunc("POST /attempts/{id}/finish", s.requireAuth(api.FinishAttempt))
	mux.HandleFunc("POST /attempts/{id}/cancel", s.requireAuth(api.CancelAttempt))

	mux.HandleFunc("GET /courses", s.requireAuth(api.ListCourses))
	mux.HandleFunc("POST /courses", s.requireAuth(api.CreateCourse))
	mux.HandleFunc("GET /courses/{id}", s.requireAuth(api.GetCourse))
	mux.HandleFunc("PATCH /courses/{id}", s.requireAuth(api.RenameCourse))
	mux.HandleFunc("POST /courses/{courseID}/blocks", s.requireAuth(api.AddCourseBlock))
	mux.HandleFunc("DELETE /courses/{courseID}/blocks/{blockID}", s.requireAuth(api.RemoveCourseBlock))
	mux.HandleFunc("PATCH /courses/{courseID}/blocks/move", s.requireAuth(api.MoveCourseBlock))
	mux.HandleFunc("POST /blocks/{blockID}/elements", s.requireAuth(api.AddBlockElement))
	mux.HandleFunc("DELETE /blocks/{blockID}/elements/{elementID}", s.requireAuth(api.RemoveBlockElement))
	mux.HandleFunc("PATCH /blocks/{blockID}/elements/move", s.requireAuth(api.MoveBlockElement))
	mux.HandleFunc("PATCH /elements/{elementID}/completion-mode", s.requireAuth(api.ChangeElementCompletionMode))
	mux.HandleFunc("PUT /elements/{elementID}/content", s.requireAuth(api.UploadElementContent))
	mux.HandleFunc("GET /elements/{elementID}/download", s.requireAuth(api.DownloadElementContent))
	mux.HandleFunc("GET /courses/{courseID}/progress", s.requireAuth(api.GetCourseProgress))
	mux.HandleFunc("GET /courses/{id}/ratings", s.requireAuth(api.ListCourseRatings))

	mux.HandleFunc("GET /enrollments", s.requireAuth(api.ListEnrollments))
	mux.HandleFunc("POST /enrollments", s.requireAuth(api.CreateEnrollment))
	mux.HandleFunc("GET /enrollments/{id}", s.requireAuth(api.GetEnrollment))
	mux.HandleFunc("GET /statistics/students", s.requireAuth(api.ListStudentStatistics))

	return middleware.Recover(s.logger)(
		middleware.RequestID(
			middleware.SecurityHeaders(mux),
		),
	)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, response.Envelope{
		Data: map[string]string{"status": "ok"},
		Links: response.Links{
			"self": response.Link{Href: "/healthz"},
			"root": response.Link{Href: "/"},
		},
	})
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, response.Envelope{
		Data: map[string]string{"name": "lms-api"},
		Links: response.Links{
			"self":          response.Link{Href: "/"},
			"login":         response.Link{Href: "/auth/login", Method: http.MethodPost},
			"accounts":      response.Link{Href: "/accounts", Method: http.MethodGet},
			"organizations": response.Link{Href: "/organizations", Method: http.MethodGet},
			"persons":       response.Link{Href: "/persons", Method: http.MethodGet},
			"courses":       response.Link{Href: "/courses", Method: http.MethodGet},
		},
	})
}

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return middleware.Auth(s.session)(next)
}
