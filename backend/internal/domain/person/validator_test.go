package person

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/person/profile"
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
			name: "uuid.Nil",
			args: args{
				id: uuid.Nil,
			},
			wantErr: true,
		},
		{
			name: "непустой id",
			args: args{
				id: idFixture,
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

func Test_validateName(t *testing.T) {
	type args struct {
		n name.Name
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустое имя",
			args: args{
				n: name.Name{},
			},
			wantErr: true,
		},
		{
			name: "непустое имя",
			args: args{
				n: nameFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateName(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("validateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateProfile(t *testing.T) {
	type args struct {
		prof profile.Profile
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустой профиль",
			args: args{
				prof: profile.Profile{},
			},
			wantErr: true,
		},
		{
			name: "непустой профиль",
			args: args{
				prof: profileFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateProfile(tt.args.prof); (err != nil) != tt.wantErr {
				t.Errorf("validateProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
