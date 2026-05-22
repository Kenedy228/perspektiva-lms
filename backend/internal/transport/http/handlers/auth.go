package handlers

import (
	"errors"
	"net/http"

	accountcommands "gitflic.ru/lms/backend/internal/application/usecases/account/commands"
	accountcommon "gitflic.ru/lms/backend/internal/application/usecases/account/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/transport/http/middleware"
	"gitflic.ru/lms/backend/internal/transport/http/response"
	"gitflic.ru/lms/backend/internal/transport/http/session"
	"github.com/google/uuid"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	authenticate *accountcommands.AuthenticateUseCase
	sessions     *session.Manager
}

// NewAuthHandler creates an authentication handler.
func NewAuthHandler(authenticate *accountcommands.AuthenticateUseCase, sessions *session.Manager) *AuthHandler {
	if authenticate == nil {
		panic("auth handler requires authenticate usecase")
	}
	if sessions == nil {
		panic("auth handler requires session manager")
	}
	return &AuthHandler{authenticate: authenticate, sessions: sessions}
}

// Login authenticates credentials and returns a signed bearer session token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}

	out, err := h.authenticate.Execute(r.Context(), accountcommands.AuthenticateInput{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, accountcommon.ErrInvalidCredentials) {
			response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "invalid_credentials", "login or password is invalid"))
			return
		}
		response.WriteError(w, r, response.NewError(http.StatusInternalServerError, "internal_error", "authentication failed"))
		return
	}

	accountID, err := uuid.Parse(out.AccountID)
	if err != nil {
		response.WriteError(w, r, response.NewError(http.StatusInternalServerError, "internal_error", "account id is invalid"))
		return
	}
	personID, err := uuid.Parse(out.PersonID)
	if err != nil {
		response.WriteError(w, r, response.NewError(http.StatusInternalServerError, "internal_error", "person id is invalid"))
		return
	}
	roleType, err := role.ParseType(out.Role)
	if err != nil {
		response.WriteError(w, r, response.NewError(http.StatusInternalServerError, "internal_error", "role is invalid"))
		return
	}
	accountRole, err := role.New(roleType)
	if err != nil {
		response.WriteError(w, r, response.NewError(http.StatusInternalServerError, "internal_error", "role is invalid"))
		return
	}

	token, claims, err := h.sessions.Issue(accountID, personID, accountRole)
	if err != nil {
		response.WriteError(w, r, response.NewError(http.StatusInternalServerError, "internal_error", "session creation failed"))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{
		Data: SessionResponse{
			Token:     token,
			TokenType: "Bearer",
			ExpiresAt: claims.ExpiresAt,
			Account:   AccountRef{ID: out.AccountID},
			Person:    PersonRef{ID: out.PersonID},
			Role:      out.Role,
		},
		Links: response.Links{
			"self":    response.Link{Href: "/auth/session", Method: http.MethodGet},
			"logout":  response.Link{Href: "/auth/logout", Method: http.MethodPost},
			"account": response.Link{Href: "/accounts/" + out.AccountID, Method: http.MethodGet},
			"person":  response.Link{Href: "/persons/" + out.PersonID, Method: http.MethodGet},
			"courses": response.Link{Href: "/courses", Method: http.MethodGet},
		},
	})
}

// Logout acknowledges client-side bearer token disposal.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, response.Envelope{
		Data: map[string]string{"status": "logged_out"},
		Links: response.Links{
			"login": response.Link{Href: "/auth/login", Method: http.MethodPost},
		},
	})
}

// Session returns the authenticated session claims propagated by middleware.
func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	response.WriteJSON(w, http.StatusOK, response.Envelope{
		Data: SessionResponse{
			TokenType: "Bearer",
			ExpiresAt: claims.ExpiresAt,
			Account:   AccountRef{ID: claims.AccountID.String()},
			Person:    PersonRef{ID: claims.PersonID.String()},
			Role:      claims.Role.Kind().String(),
		},
		Links: response.Links{
			"self":    response.Link{Href: "/auth/session", Method: http.MethodGet},
			"logout":  response.Link{Href: "/auth/logout", Method: http.MethodPost},
			"account": response.Link{Href: "/accounts/" + claims.AccountID.String(), Method: http.MethodGet},
			"person":  response.Link{Href: "/persons/" + claims.PersonID.String(), Method: http.MethodGet},
		},
	})
}
