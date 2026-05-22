package profile

import (
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	type args struct {
		s              snils.SNILS
		dob            dob.DateOfBirth
		jt             jobtitle.JobTitle
		e              education.Education
		organizationID uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    Profile
		wantErr bool
	}{
		{
			name: "без СНИЛС",
			args: args{
				s:              snils.SNILS{},
				dob:            dateOfBirthFixture(t),
				jt:             jobtitle.JobTitle{},
				e:              education.Education{},
				organizationID: uuid.UUID{},
			},
			want: Profile{
				snils:          snils.SNILS{},
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			wantErr: true,
		},
		{
			name: "без даты рождения",
			args: args{
				s:              snilsFixture(t),
				dob:            dob.DateOfBirth{},
				jt:             jobtitle.JobTitle{},
				e:              education.Education{},
				organizationID: uuid.UUID{},
			},
			want: Profile{
				snils:          snils.SNILS{},
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			wantErr: true,
		},
		{
			name: "валидное состояние",
			args: args{
				s:              snilsFixture(t),
				dob:            dateOfBirthFixture(t),
				jt:             jobtitle.JobTitle{},
				e:              education.Education{},
				organizationID: uuid.UUID{},
			},
			want: Profile{
				snils:          snilsFixture(t),
				dateOfBirth:    dateOfBirthFixture(t),
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.s, tt.args.dob, tt.args.jt, tt.args.e, tt.args.organizationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfile_DateOfBirth(t *testing.T) {
	type fields struct {
		snils          snils.SNILS
		dateOfBirth    dob.DateOfBirth
		jobTitle       jobtitle.JobTitle
		education      education.Education
		organizationID uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   dob.DateOfBirth
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				snils:          snils.SNILS{},
				dateOfBirth:    dateOfBirthFixture(t),
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			want: dateOfBirthFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				snils:          tt.fields.snils,
				dateOfBirth:    tt.fields.dateOfBirth,
				jobTitle:       tt.fields.jobTitle,
				education:      tt.fields.education,
				organizationID: tt.fields.organizationID,
			}
			if got := p.DateOfBirth(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateOfBirth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfile_Education(t *testing.T) {
	type fields struct {
		snils          snils.SNILS
		dateOfBirth    dob.DateOfBirth
		jobTitle       jobtitle.JobTitle
		education      education.Education
		organizationID uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   education.Education
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				snils:          snils.SNILS{},
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobtitle.JobTitle{},
				education:      educationFixture(t),
				organizationID: uuid.UUID{},
			},
			want: educationFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				snils:          tt.fields.snils,
				dateOfBirth:    tt.fields.dateOfBirth,
				jobTitle:       tt.fields.jobTitle,
				education:      tt.fields.education,
				organizationID: tt.fields.organizationID,
			}
			if got := p.Education(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Education() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfile_HasOrganization(t *testing.T) {
	type fields struct {
		snils          snils.SNILS
		dateOfBirth    dob.DateOfBirth
		jobTitle       jobtitle.JobTitle
		education      education.Education
		organizationID uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "есть организация",
			fields: fields{
				snils:          snils.SNILS{},
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: orgIDFixture,
			},
			want: true,
		},
		{
			name: "нет организации",
			fields: fields{
				snils:          snils.SNILS{},
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				snils:          tt.fields.snils,
				dateOfBirth:    tt.fields.dateOfBirth,
				jobTitle:       tt.fields.jobTitle,
				education:      tt.fields.education,
				organizationID: tt.fields.organizationID,
			}
			if got := p.HasOrganization(); got != tt.want {
				t.Errorf("HasOrganization() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfile_IsZero(t *testing.T) {
	type fields struct {
		snils          snils.SNILS
		dateOfBirth    dob.DateOfBirth
		jobTitle       jobtitle.JobTitle
		education      education.Education
		organizationID uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "нет СНИЛС",
			fields: fields{
				snils:          snils.SNILS{},
				dateOfBirth:    dateOfBirthFixture(t),
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			want: true,
		},
		{
			name: "нет даты рождения",
			fields: fields{
				snils:          snilsFixture(t),
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			want: true,
		},
		{
			name: "нет СНИЛС и даты рождения",
			fields: fields{
				snils:          snils.SNILS{},
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			want: true,
		},
		{
			name: "есть СНИЛС и дата рождения",
			fields: fields{
				snils:          snilsFixture(t),
				dateOfBirth:    dateOfBirthFixture(t),
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				snils:          tt.fields.snils,
				dateOfBirth:    tt.fields.dateOfBirth,
				jobTitle:       tt.fields.jobTitle,
				education:      tt.fields.education,
				organizationID: tt.fields.organizationID,
			}
			if got := p.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfile_JobTitle(t *testing.T) {
	type fields struct {
		snils          snils.SNILS
		dateOfBirth    dob.DateOfBirth
		jobTitle       jobtitle.JobTitle
		education      education.Education
		organizationID uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   jobtitle.JobTitle
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				snils:          snils.SNILS{},
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobTitleFixture(t),
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			want: jobTitleFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				snils:          tt.fields.snils,
				dateOfBirth:    tt.fields.dateOfBirth,
				jobTitle:       tt.fields.jobTitle,
				education:      tt.fields.education,
				organizationID: tt.fields.organizationID,
			}
			if got := p.JobTitle(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfile_OrganizationID(t *testing.T) {
	type fields struct {
		snils          snils.SNILS
		dateOfBirth    dob.DateOfBirth
		jobTitle       jobtitle.JobTitle
		education      education.Education
		organizationID uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				snils:          snils.SNILS{},
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: orgIDFixture,
			},
			want: orgIDFixture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				snils:          tt.fields.snils,
				dateOfBirth:    tt.fields.dateOfBirth,
				jobTitle:       tt.fields.jobTitle,
				education:      tt.fields.education,
				organizationID: tt.fields.organizationID,
			}
			if got := p.OrganizationID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrganizationID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfile_SNILS(t *testing.T) {
	type fields struct {
		snils          snils.SNILS
		dateOfBirth    dob.DateOfBirth
		jobTitle       jobtitle.JobTitle
		education      education.Education
		organizationID uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   snils.SNILS
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				snils:          snilsFixture(t),
				dateOfBirth:    dob.DateOfBirth{},
				jobTitle:       jobtitle.JobTitle{},
				education:      education.Education{},
				organizationID: uuid.UUID{},
			},
			want: snilsFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				snils:          tt.fields.snils,
				dateOfBirth:    tt.fields.dateOfBirth,
				jobTitle:       tt.fields.jobTitle,
				education:      tt.fields.education,
				organizationID: tt.fields.organizationID,
			}
			if got := p.SNILS(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SNILS() = %v, want %v", got, tt.want)
			}
		})
	}
}
