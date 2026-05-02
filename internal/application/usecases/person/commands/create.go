package commands

import (
	"context"

	personrepo "gitflic.ru/lms/internal/application/ports/person"
	persondomain "gitflic.ru/lms/internal/domain/person"
	personname "gitflic.ru/lms/internal/domain/person/name"
)

type CreateUseCase struct {
	r personrepo.Repository
}

func NewCreateUseCase(r personrepo.Repository) *CreateUseCase {
	return &CreateUseCase{
		r: r,
	}
}

type CreateInput struct {
	FirstName  string
	LastName   string
	MiddleName string
}

type CreateOutput struct {
	ID string
}

func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (*CreateOutput, error) {
	pName, err := personname.New(in.FirstName, in.LastName, in.MiddleName)
	if err != nil {
		return nil, err
	}

	p, err := persondomain.New(pName)
	if err != nil {
		return nil, err
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, err
	}

	return &CreateOutput{
		ID: p.ID().String(),
	}, nil
}
