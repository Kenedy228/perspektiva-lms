package commands

import (
	"context"
	"fmt"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/account/login"
)

type AuthenticateUseCase struct {
	r        accountports.Repository
	comparer accountdomain.PasswordComparer
}

func NewAuthenticateUseCase(r accountports.Repository, comparer accountdomain.PasswordComparer) *AuthenticateUseCase {
	if r == nil {
		panic("account authenticate usecase requires repository")
	}
	if comparer == nil {
		panic("account authenticate usecase requires password comparer")
	}

	return &AuthenticateUseCase{r: r, comparer: comparer}
}

type AuthenticateInput struct {
	Login    string
	Password string
}

type AuthenticateOutput struct {
	AccountID string
	PersonID  string
	Role      string
}

func (uc *AuthenticateUseCase) Execute(ctx context.Context, in AuthenticateInput) (*AuthenticateOutput, error) {
	l, err := login.New(in.Login)
	if err != nil {
		return nil, common.ErrInvalidCredentials
	}

	acc, err := uc.r.FindByLogin(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("find account by login: %w", err)
	}

	if !acc.IsActive() {
		return nil, common.ErrInvalidCredentials
	}

	if !uc.comparer.Compare(acc.PasswordHash(), in.Password) {
		return nil, common.ErrInvalidCredentials
	}

	return &AuthenticateOutput{
		AccountID: acc.ID().String(),
		PersonID:  acc.PersonID().String(),
		Role:      acc.Role().Kind().String(),
	}, nil
}
