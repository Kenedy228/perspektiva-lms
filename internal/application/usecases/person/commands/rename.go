package commands

import (
	"context"

	personrepo "gitflic.ru/lms/internal/application/ports/person"
	personname "gitflic.ru/lms/internal/domain/person/name"
	"github.com/google/uuid"
)

type RenameUseCase struct {
	r personrepo.Repository
}

func NewRenameUseCase(r personrepo.Repository) *RenameUseCase {
	return &RenameUseCase{
		r: r,
	}
}

type RenameInput struct {
	PersonID   string
	FirstName  string
	LastName   string
	MiddleName string
}

type RenameOutput struct {
	ID string
}

func (uc *RenameUseCase) Execute(ctx context.Context, in RenameInput) (*RenameOutput, error) {
	pName, err := personname.New(in.FirstName, in.LastName, in.MiddleName)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(in.PersonID)
	if err != nil {
		return nil, err
	}

	p, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Rename(pName)

	err = uc.r.Save(ctx, p)
	if err != nil {
		return nil, err
	}

	return &RenameOutput{
		ID: p.ID().String(),
	}, nil
}
