package profile

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
)

func Test_validateDateOfBirth(t *testing.T) {
	type args struct {
		dob dob.DateOfBirth
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "нулевое значение",
			args: args{
				dob: dob.DateOfBirth{},
			},
			wantErr: true,
		},
		{
			name: "ненулевое значение",
			args: args{
				dob: dateOfBirthFixture(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDateOfBirth(tt.args.dob); (err != nil) != tt.wantErr {
				t.Errorf("validateDateOfBirth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateSNILS(t *testing.T) {
	type args struct {
		s snils.SNILS
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "нулевое значение",
			args: args{
				s: snils.SNILS{},
			},
			wantErr: true,
		},
		{
			name: "ненулевое значение",
			args: args{
				s: snilsFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSNILS(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("validateSNILS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
