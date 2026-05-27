package base

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	"github.com/google/uuid"
)

func Test_validateID(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "нулевой идентификатор",
			args: args{
				id: uuid.UUID{},
			},
			wantErr: true,
		},
		{
			name: "существующий идентификатор",
			args: args{
				id: uuid.New(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateID(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("validateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateIDRequired(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "нулевой идентификатор",
			args: args{
				id: uuid.UUID{},
			},
			wantErr: true,
		},
		{
			name: "существующий идентификатор",
			args: args{
				id: uuid.New(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateIDRequired(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("validateIDRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateTitle(t *testing.T) {
	type args struct {
		t title.Title
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустой заголовок",
			args: args{
				t: title.Title{},
			},
			wantErr: true,
		},
		{
			name: "непустой заголовок",
			args: args{
				t: titleFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateTitle(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("validateTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateTitleRequired(t *testing.T) {
	type args struct {
		t title.Title
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустой заголовок",
			args: args{
				t: title.Title{},
			},
			wantErr: true,
		},
		{
			name: "непустой заголовок",
			args: args{
				t: titleFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateTitleRequired(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("validateTitleRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
