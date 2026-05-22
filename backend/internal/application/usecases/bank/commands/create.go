package commands

import (
	"context"
	"fmt"

	bankports "gitflic.ru/lms/backend/internal/application/ports/bank"
	"gitflic.ru/lms/backend/internal/application/usecases/bank/common"
	bankdomain "gitflic.ru/lms/backend/internal/domain/bank"
)

type CreateUseCase struct {
	r     bankports.Repository
	audit bankports.AuditRecorder
}

func NewCreateUseCase(r bankports.Repository, audit bankports.AuditRecorder) *CreateUseCase {
	if r == nil {
		panic("bank create usecase requires repository")
	}
	return &CreateUseCase{r: r, audit: audit}
}

func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (*Output, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}

	t, err := buildTitle(in.Title)
	if err != nil {
		return nil, err
	}

	b, err := bankdomain.New(t)
	if err != nil {
		return nil, fmt.Errorf("create bank aggregate: %w", err)
	}

	if err := saveBank(ctx, uc.r, b); err != nil {
		return nil, err
	}

	if err := recordAudit(ctx, uc.audit, bankports.AuditActionCreate, b.ID().String(), in.ActorRole); err != nil {
		return nil, err
	}

	return &Output{ID: b.ID().String()}, nil
}
