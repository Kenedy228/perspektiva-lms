package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

var (
	// ErrInvalidToken is returned when a token cannot be verified.
	ErrInvalidToken = errors.New("invalid session token")
	// ErrExpiredToken is returned when a token is valid but expired.
	ErrExpiredToken = errors.New("expired session token")
)

// Claims describes authenticated account identity propagated to handlers.
type Claims struct {
	AccountID uuid.UUID `json:"account_id"`
	PersonID  uuid.UUID `json:"person_id"`
	Role      role.Role `json:"-"`
	RoleValue string    `json:"role"`
	IssuedAt  int64     `json:"iat"`
	ExpiresAt int64     `json:"exp"`
}

// Manager creates and verifies stateless HMAC session tokens.
type Manager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

// NewManager creates a session manager.
func NewManager(secret []byte, ttl time.Duration) *Manager {
	if len(secret) == 0 {
		panic("session manager requires secret")
	}
	if ttl <= 0 {
		panic("session manager requires positive ttl")
	}
	return &Manager{
		secret: secret,
		ttl:    ttl,
		now:    time.Now,
	}
}

// Issue creates a signed token for an authenticated account.
func (m *Manager) Issue(accountID, personID uuid.UUID, accountRole role.Role) (string, Claims, error) {
	now := m.now().UTC()
	claims := Claims{
		AccountID: accountID,
		PersonID:  personID,
		Role:      accountRole,
		RoleValue: accountRole.Kind().String(),
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(m.ttl).Unix(),
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", Claims{}, fmt.Errorf("marshal session claims: %w", err)
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := m.sign(encodedPayload)
	return encodedPayload + "." + signature, claims, nil
}

// Verify checks token signature and expiration.
func (m *Manager) Verify(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return Claims{}, ErrInvalidToken
	}
	expected := m.sign(parts[0])
	if !hmac.Equal([]byte(expected), []byte(parts[1])) {
		return Claims{}, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}
	if m.now().UTC().Unix() >= claims.ExpiresAt {
		return Claims{}, ErrExpiredToken
	}
	roleType, err := role.ParseType(claims.RoleValue)
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	claims.Role, err = role.New(roleType)
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	return claims, nil
}

func (m *Manager) sign(payload string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
