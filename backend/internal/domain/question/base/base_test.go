package base

import (
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBase_ChangeTitle(t *testing.T) {
	type fields struct {
		id    uuid.UUID
		title title.Title
	}
	type args struct {
		t title.Title
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "пустой заголовок",
			fields: fields{
				id:    uuid.UUID{},
				title: title.Title{},
			},
			args: args{
				t: title.Title{},
			},
			wantErr: true,
		},
		{
			name: "непустой заголовок",
			fields: fields{
				id:    uuid.UUID{},
				title: title.Title{},
			},
			args: args{
				t: titleFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:    tt.fields.id,
				title: tt.fields.title,
			}
			if err := b.ChangeTitle(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("ChangeTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBase_Clone(t *testing.T) {
	type fields struct {
		id    uuid.UUID
		title title.Title
	}
	tests := []struct {
		name   string
		fields fields
		want   *Base
	}{
		{
			name: "клонирует значения как есть",
			fields: fields{
				id:    idFixture,
				title: titleFixture(t),
			},
			want: &Base{
				id:    idFixture,
				title: titleFixture(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:    tt.fields.id,
				title: tt.fields.title,
			}
			if got := b.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase_ID(t *testing.T) {
	type fields struct {
		id    uuid.UUID
		title title.Title
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				id:    idFixture,
				title: titleFixture(t),
			},
			want: idFixture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:    tt.fields.id,
				title: tt.fields.title,
			}
			if got := b.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase_Title(t *testing.T) {
	type fields struct {
		id    uuid.UUID
		title title.Title
	}
	tests := []struct {
		name   string
		fields fields
		want   title.Title
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				id:    idFixture,
				title: titleFixture(t),
			},
			want: titleFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Base{
				id:    tt.fields.id,
				title: tt.fields.title,
			}
			if got := b.Title(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Title() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		t title.Title
	}
	tests := []struct {
		name    string
		args    args
		want    *Base
		wantErr bool
	}{
		{
			name: "пустой заголовок",
			args: args{
				t: title.Title{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "непустой заголовок",
			args: args{
				t: titleFixture(t),
			},
			want: &Base{
				title: titleFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				assert.NotEqual(t, got.ID(), uuid.Nil)
				assert.Equal(t, got.Title(), tt.want.Title())
			}
		})
	}
}

func TestRestore(t *testing.T) {
	type args struct {
		id uuid.UUID
		t  title.Title
	}
	tests := []struct {
		name    string
		args    args
		want    *Base
		wantErr bool
	}{
		{
			name: "пустой идентификатор, не пустой заголовок",
			args: args{
				id: uuid.Nil,
				t:  titleFixture(t),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "пустой заголовок, не пустой идентификатор",
			args: args{
				id: idFixture,
				t:  title.Title{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "пустой заголовок и пустой идентификатор",
			args: args{
				id: uuid.Nil,
				t:  title.Title{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "непустой заголовок и непустой идентификатор",
			args: args{
				id: idFixture,
				t:  titleFixture(t),
			},
			want: &Base{
				id:    idFixture,
				title: titleFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Restore(tt.args.id, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Restore() got = %v, want %v", got, tt.want)
			}
		})
	}
}
