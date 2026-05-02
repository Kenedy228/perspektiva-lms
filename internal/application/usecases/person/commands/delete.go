package commands

import (
	"context"

	personrepo "gitflic.ru/lms/internal/application/ports/person"
	"github.com/google/uuid"
)

type DeleteByIDUseCase struct {
	r personrepo.Repository
}

func NewDeleteByIDUseCase(r personrepo.Repository) *DeleteByIDUseCase {
	return &DeleteByIDUseCase{
		r: r,
	}
}

type DeleteByIDInput struct {
	PersonID string
}

func (uc *DeleteByIDUseCase) Execute(ctx context.Context, in DeleteByIDInput) error {
	pID, err := uuid.Parse(in.PersonID)
	if err != nil {
		return err
	}

	if err := uc.r.DeleteByID(ctx, pID); err != nil {
		return err
	}

	return nil
}
