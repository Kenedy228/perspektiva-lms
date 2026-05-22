package organization

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/organization/inn"
	"gitflic.ru/lms/backend/internal/domain/organization/name"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		inn inn.INN
		n   name.Name
	}
	tests := []struct {
		name    string
		args    args
		want    *Organization
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "некорректное имя",
			args: args{
				inn: innFixture(t),
				n:   name.Name{},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "некорректный ИНН",
			args: args{
				inn: inn.INN{},
				n:   nameFixture(t),
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "некорректное имя и ИНН",
			args: args{
				inn: inn.INN{},
				n:   name.Name{},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "корректные значения",
			args: args{
				inn: innFixture(t),
				n:   nameFixture(t),
			},
			want: &Organization{
				inn: innFixture(t),
				n:   nameFixture(t),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.inn, tt.args.n)
			if !tt.wantErr(t, err, fmt.Sprintf("New(%v, %v)", tt.args.inn, tt.args.n)) {
				return
			}

			if tt.want != nil {
				assert.NotEqual(t, uuid.Nil, got.id)
				assert.Equal(t, tt.want.n, got.n)
				assert.Equal(t, tt.want.inn, got.inn)
			}
		})
	}
}

func TestOrganization_ChangeINN(t *testing.T) {
	type fields struct {
		id  uuid.UUID
		inn inn.INN
		n   name.Name
	}
	type args struct {
		inn inn.INN
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "корректное значение",
			fields: fields{
				id:  uuid.UUID{},
				inn: inn.INN{},
				n:   name.Name{},
			},
			args: args{
				inn: innFixture(t),
			},
			wantErr: assert.NoError,
		},
		{
			name: "некорректное значение",
			fields: fields{
				id:  uuid.UUID{},
				inn: inn.INN{},
				n:   name.Name{},
			},
			args: args{
				inn: inn.INN{},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Organization{
				id:  tt.fields.id,
				inn: tt.fields.inn,
				n:   tt.fields.n,
			}
			tt.wantErr(t, o.ChangeINN(tt.args.inn), fmt.Sprintf("ChangeINN(%v)", tt.args.inn))
		})
	}
}

func TestOrganization_ChangeName(t *testing.T) {
	type fields struct {
		id  uuid.UUID
		inn inn.INN
		n   name.Name
	}
	type args struct {
		n name.Name
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "корректное значение",
			fields: fields{
				id:  uuid.UUID{},
				inn: inn.INN{},
				n:   name.Name{},
			},
			args: args{
				n: nameFixture(t),
			},
			wantErr: assert.NoError,
		},
		{
			name: "некорректное значение",
			fields: fields{
				id:  uuid.UUID{},
				inn: inn.INN{},
				n:   name.Name{},
			},
			args: args{
				n: name.Name{},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Organization{
				id:  tt.fields.id,
				inn: tt.fields.inn,
				n:   tt.fields.n,
			}
			tt.wantErr(t, o.ChangeName(tt.args.n), fmt.Sprintf("ChangeName(%v)", tt.args.n))
		})
	}
}

func TestOrganization_ID(t *testing.T) {
	type fields struct {
		id  uuid.UUID
		inn inn.INN
		n   name.Name
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "возвращает id как есть",
			fields: fields{
				id:  idFixture,
				inn: inn.INN{},
				n:   name.Name{},
			},
			want: idFixture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Organization{
				id:  tt.fields.id,
				inn: tt.fields.inn,
				n:   tt.fields.n,
			}
			assert.Equalf(t, tt.want, o.ID(), "ID()")
		})
	}
}

func TestOrganization_INN(t *testing.T) {
	type fields struct {
		id  uuid.UUID
		inn inn.INN
		n   name.Name
	}
	tests := []struct {
		name   string
		fields fields
		want   inn.INN
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				id:  uuid.UUID{},
				inn: innFixture(t),
				n:   name.Name{},
			},
			want: innFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Organization{
				id:  tt.fields.id,
				inn: tt.fields.inn,
				n:   tt.fields.n,
			}
			assert.Equalf(t, tt.want, o.INN(), "INN()")
		})
	}
}

func TestOrganization_Name(t *testing.T) {
	type fields struct {
		id  uuid.UUID
		inn inn.INN
		n   name.Name
	}
	tests := []struct {
		name   string
		fields fields
		want   name.Name
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				id:  uuid.UUID{},
				inn: inn.INN{},
				n:   nameFixture(t),
			},
			want: nameFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Organization{
				id:  tt.fields.id,
				inn: tt.fields.inn,
				n:   tt.fields.n,
			}
			assert.Equalf(t, tt.want, o.Name(), "Name()")
		})
	}
}

func TestRestore(t *testing.T) {
	type args struct {
		id  uuid.UUID
		inn inn.INN
		n   name.Name
	}
	tests := []struct {
		name    string
		args    args
		want    *Organization
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "некорректный id",
			args: args{
				id:  uuid.UUID{},
				inn: innFixture(t),
				n:   nameFixture(t),
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "некорректное имя",
			args: args{
				id:  idFixture,
				inn: innFixture(t),
				n:   name.Name{},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "некорректный ИНН",
			args: args{
				id:  idFixture,
				inn: inn.INN{},
				n:   nameFixture(t),
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "некорректное имя и ИНН",
			args: args{
				id:  idFixture,
				inn: inn.INN{},
				n:   name.Name{},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "корректные значения",
			args: args{
				id:  idFixture,
				inn: innFixture(t),
				n:   nameFixture(t),
			},
			want: &Organization{
				id:  idFixture,
				inn: innFixture(t),
				n:   nameFixture(t),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Restore(tt.args.id, tt.args.inn, tt.args.n)
			if !tt.wantErr(t, err, fmt.Sprintf("Restore(%v, %v, %v)", tt.args.id, tt.args.inn, tt.args.n)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Restore(%v, %v, %v)", tt.args.id, tt.args.inn, tt.args.n)
		})
	}
}
