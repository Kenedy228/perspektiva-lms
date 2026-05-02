package organization

import (
	"gitflic.ru/lms/internal/domain/organization/inn"
	"gitflic.ru/lms/internal/domain/organization/orgname"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Organization struct {
	id   uuid.UUID
	inn  inn.INN
	name orgname.Name
}

func New(inn inn.INN, name orgname.Name) (*Organization, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Organization{
		id:   id,
		inn:  inn,
		name: name,
	}, nil
}

func (o *Organization) Name() orgname.Name {
	return o.name
}

func (o *Organization) INN() inn.INN {
	return o.inn
}

func (o *Organization) Rename(name orgname.Name) {
	o.name = name
}

func (o *Organization) ChangeINN(inn inn.INN) {
	o.inn = inn
}
