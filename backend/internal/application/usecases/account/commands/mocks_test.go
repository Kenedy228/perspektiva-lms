package commands_test

import (
	"context"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*accountdomain.Account, error) {
	args := m.Called(ctx, id)

	var acc *accountdomain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*accountdomain.Account)
	}

	return acc, args.Error(1)
}

func (m *mockRepository) FindByLogin(ctx context.Context, l login.Login) (*accountdomain.Account, error) {
	args := m.Called(ctx, l)

	var acc *accountdomain.Account
	if args.Get(0) != nil {
		acc = args.Get(0).(*accountdomain.Account)
	}

	return acc, args.Error(1)
}

func (m *mockRepository) LoginExists(ctx context.Context, l login.Login, excludeAccountID uuid.UUID) (bool, error) {
	args := m.Called(ctx, l, excludeAccountID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepository) PersonHasAccount(ctx context.Context, personID uuid.UUID, excludeAccountID uuid.UUID) (bool, error) {
	args := m.Called(ctx, personID, excludeAccountID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepository) Save(ctx context.Context, acc *accountdomain.Account) error {
	args := m.Called(ctx, acc)
	return args.Error(0)
}

type mockPasswordHasher struct {
	mock.Mock
}

func (m *mockPasswordHasher) Hash(plain string) (passhash.Hash, error) {
	args := m.Called(plain)

	var h passhash.Hash
	if args.Get(0) != nil {
		h = args.Get(0).(passhash.Hash)
	}

	return h, args.Error(1)
}

type mockPasswordComparer struct {
	mock.Mock
}

func (m *mockPasswordComparer) Compare(hash passhash.Hash, plain string) bool {
	args := m.Called(hash, plain)
	return args.Bool(0)
}

type mockAuditRecorder struct {
	mock.Mock
}

func (m *mockAuditRecorder) RecordAccountAudit(ctx context.Context, event accountports.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
