package commands

import (
	"context"

	"gitflic.ru/lms/internal/application/ports/person"
	"github.com/google/uuid"
)

type DetachProfileUseCase struct {
	r person.Repository
}

func NewDetachProfileUseCase(r person.Repository) *DetachProfileUseCase {
	return &DetachProfileUseCase{
		r: r,
	}
}

type DetachProfileInput struct {
	PersonID string
}

type DetachProfileOutput struct {
	PersonID string
}

func (uc *DetachProfileUseCase) Execute(ctx context.Context, in DetachProfileInput) (*DetachProfileOutput, error) {
	pID, err := uuid.Parse(in.PersonID)
	if err != nil {
		return nil, err
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return nil, err
	}

	if p.HasProfile() {
		p.DetachProfile()
		if err := uc.r.Save(ctx, p); err != nil {
			return nil, err
		}
	}

	return &DetachProfileOutput{
		PersonID: p.ID().String(),
	}, nil
}
