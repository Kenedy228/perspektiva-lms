package person

import (
	"errors"
	"testing"
	"time"

	"gitflic.ru/lms/internal/domain/person/dob"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		firstName      string
		lastName       string
		middleName     string
		jobTitle       string
		snils          string
		dateOfBirth    time.Time
		education      string
		organizationID uuid.UUID
		expectedErr    error
	}{
		{
			name:           "valid: full data",
			firstName:      "Иван",
			lastName:       "Иванов",
			middleName:     "Иванович",
			jobTitle:       "Разработчик",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.New(),
			expectedErr:    nil,
		},
		{
			name:           "valid: empty middleName (allowed by validateMiddleName)",
			firstName:      "Анна",
			lastName:       "Смирнова",
			middleName:     "",
			jobTitle:       "Аналитик",
			snils:          "987-654-321 09",
			dateOfBirth:    time.Now(),
			education:      "Среднее",
			organizationID: uuid.Nil,
			expectedErr:    nil,
		},
		{
			name:           "valid: names with hyphens and apostrophes",
			firstName:      "Жанна",
			lastName:       "Д'Арк-Смирнова",
			middleName:     "Ахмед-оглы",
			jobTitle:       "Менеджер",
			snils:          "111-222-333 44",
			dateOfBirth:    time.Now(),
			education:      "Бакалавр",
			organizationID: uuid.New(),
			expectedErr:    nil,
		},
		{
			name:           "invalid: empty firstName",
			firstName:      "",
			lastName:       "Иванов",
			middleName:     "Иванович",
			jobTitle:       "Разработчик",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrEmptyFirstName,
		},
		{
			name:           "invalid: firstName with spaces only",
			firstName:      "   ",
			lastName:       "Иванов",
			middleName:     "",
			jobTitle:       "Разработчик",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrEmptyFirstName,
		},
		{
			name:           "invalid: firstName with english letters",
			firstName:      "John",
			lastName:       "Иванов",
			middleName:     "",
			jobTitle:       "Разработчик",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrInvalidFirstName,
		},
		{
			name:           "invalid: firstName with numbers",
			firstName:      "Иван123",
			lastName:       "Иванов",
			middleName:     "",
			jobTitle:       "Разработчик",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrInvalidFirstName,
		},
		{
			name:           "invalid: empty lastName",
			firstName:      "Иван",
			lastName:       "",
			middleName:     "",
			jobTitle:       "Разработчик",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrEmptyLastName,
		},
		{
			name:           "invalid: lastName with special chars",
			firstName:      "Иван",
			lastName:       "Иванов@!",
			middleName:     "",
			jobTitle:       "Разработчик",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrInvalidLastName,
		},
		{
			name:           "invalid: middleName with numbers",
			firstName:      "Иван",
			lastName:       "Иванов",
			middleName:     "Иванович1",
			jobTitle:       "Разработчик",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrInvalidMiddleName,
		},
		{
			name:           "invalid: empty jobTitle",
			firstName:      "Иван",
			lastName:       "Иванов",
			middleName:     "Иванович",
			jobTitle:       "",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrEmptyJobTitle,
		},
		{
			name:           "invalid: empty snils",
			firstName:      "Иван",
			lastName:       "Иванов",
			middleName:     "Иванович",
			jobTitle:       "Разработчик",
			snils:          "",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrInvalidSnils,
		},
		{
			name:           "invalid: snils wrong format (no hyphens)",
			firstName:      "Иван",
			lastName:       "Иванов",
			middleName:     "Иванович",
			jobTitle:       "Разработчик",
			snils:          "123456789 12",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrInvalidSnils,
		},
		{
			name:           "invalid: snils wrong length",
			firstName:      "Иван",
			lastName:       "Иванов",
			middleName:     "Иванович",
			jobTitle:       "Разработчик",
			snils:          "12-345-678 90",
			dateOfBirth:    time.Now(),
			education:      "Высшее",
			organizationID: uuid.Nil,
			expectedErr:    ErrInvalidSnils,
		},
		{
			name:           "invalid: empty education",
			firstName:      "Иван",
			lastName:       "Иванов",
			middleName:     "Иванович",
			jobTitle:       "Разработчик",
			snils:          "123-456-789 12",
			dateOfBirth:    time.Now(),
			education:      "   ",
			organizationID: uuid.Nil,
			expectedErr:    ErrEmptyEducation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := createParams(tt.firstName, tt.lastName, tt.middleName, tt.jobTitle, tt.snils, tt.education, tt.organizationID, tt.dateOfBirth)

			p, err := New(params)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected err %v, got %v", tt.expectedErr, err)
			}

			if err == nil {
				if p.FirstName() != tt.firstName {
					t.Errorf("expected firstName %v, got %v", tt.firstName, p.FirstName())
				}

				if p.LastName() != tt.lastName {
					t.Errorf("expected lastName %v, got %v", tt.lastName, p.LastName())
				}

				if p.MiddleName() != tt.middleName {
					t.Errorf("expected middlename %v, got %v", tt.middleName, p.MiddleName())
				}
			
				if p.JobTitle() != tt.jobTitle {
					t.Errorf("expected job title %v, got %v", tt.jobTitle, p.JobTitle())
				}

				if p.Snils() != tt.snils {
					t.Errorf("expected snils %v, got %v", tt.snils, p.Snils())
				}

				if p.Education() != tt.education {
					t.Errorf("expected education %v, got %v", tt.education, p.Education())
				}

				if p.OrganizationID() != tt.organizationID {
					t.Errorf("expected organizationID %v, got %v", tt.organizationID, p.OrganizationID())
				}
			}
		})
	}
}

func TestRename(t *testing.T) {
	tests := []struct {
		name           string
		firstName      string
		lastName       string
		middleName     string
		wantFirstName  string
		wantLastName   string
		wantMiddleName string
		wantErr        error
	}{
		{
			name:           "success full name",
			firstName:      "Петр",
			lastName:       "Иванов",
			middleName:     "Алексеевич",
			wantFirstName:  "Петр",
			wantLastName:   "Иванов",
			wantMiddleName: "Алексеевич",
			wantErr:        nil,
		},
		{
			name:           "success empty middleName",
			firstName:      "Петр",
			lastName:       "Иванов",
			middleName:     "",
			wantFirstName:  "Петр",
			wantLastName:   "Иванов",
			wantMiddleName: "",
			wantErr:        nil,
		},
		{
			name:           "success hyphenated names",
			firstName:      "Анна-Мария",
			lastName:       "Смирнова-Петрова",
			middleName:     "Али-оглы",
			wantFirstName:  "Анна-Мария",
			wantLastName:   "Смирнова-Петрова",
			wantMiddleName: "Али-оглы",
			wantErr:        nil,
		},
		{
			name:       "error empty firstName",
			firstName:  "",
			lastName:   "Иванов",
			middleName: "Петрович",
			wantErr:    ErrEmptyFirstName,
		},
		{
			name:       "error invalid firstName",
			firstName:  "John",
			lastName:   "Иванов",
			middleName: "Петрович",
			wantErr:    ErrInvalidFirstName,
		},
		{
			name:       "error empty lastName",
			firstName:  "Иван",
			lastName:   "",
			middleName: "Петрович",
			wantErr:    ErrEmptyLastName,
		},
		{
			name:       "error invalid lastName",
			firstName:  "Иван",
			lastName:   "Иванов123",
			middleName: "Петрович",
			wantErr:    ErrInvalidLastName,
		},
		{
			name:       "error invalid middleName",
			firstName:  "Иван",
			lastName:   "Иванов",
			middleName: "Петрович123",
			wantErr:    ErrInvalidMiddleName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := basePerson()

			err := p.Rename(tt.firstName, tt.lastName, tt.middleName)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}

			if err == nil {
				if p.FirstName() != tt.wantFirstName {
					t.Fatalf("expected firstName %v, got %v", tt.wantFirstName, p.FirstName())
				}

				if p.LastName() != tt.wantLastName {
					t.Fatalf("expected lastName %v, got %v", tt.wantLastName, p.LastName())
				}

				if p.MiddleName() != tt.wantMiddleName {
					t.Fatalf("expected middleName %v, got %v", tt.wantMiddleName, p.MiddleName())
				}
			}
		})
	}
}

func TestChangeJobTitle(t *testing.T) {
	tests := []struct {
		name     string
		jobTitle string
		want     string
		wantErr  error
	}{
		{
			name:     "success",
			jobTitle: "Методист",
			want:     "Методист",
			wantErr:  nil,
		},
		{
			name:     "error empty",
			jobTitle: "",
			wantErr:  ErrEmptyJobTitle,
		},
		{
			name:     "error spaces only",
			jobTitle: "   ",
			wantErr:  ErrEmptyJobTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := basePerson()

			err := p.ChangeJobTitle(tt.jobTitle)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}

			if err == nil {
				if p.JobTitle() != tt.want {
					t.Fatalf("expected jobTitle %v, got %v", tt.want, p.JobTitle())
				}
			}
		})
	}
}

func TestChangeEducation(t *testing.T) {
	tests := []struct {
		name      string
		education string
		want      string
		wantErr   error
	}{
		{
			name:      "success",
			education: "Среднее профессиональное",
			want:      "Среднее профессиональное",
			wantErr:   nil,
		},
		{
			name:      "error empty",
			education: "",
			wantErr:   ErrEmptyEducation,
		},
		{
			name:      "error spaces only",
			education: "   ",
			wantErr:   ErrEmptyEducation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := basePerson()

			err := p.ChangeEducation(tt.education)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}

			if err == nil {
				if p.Education() != tt.want {
					t.Fatalf("expected education %v, got %v", tt.want, p.Education())
				}
			}
		})
	}
}

func TestChangeSnils(t *testing.T) {
	tests := []struct {
		name    string
		snils   string
		want    string
		wantErr error
	}{
		{
			name:    "success valid snils",
			snils:   "123-456-789 12",
			want:    "123-456-789 12",
			wantErr: nil,
		},
		{
			name:    "error empty",
			snils:   "",
			wantErr: ErrInvalidSnils,
		},
		{
			name:    "error wrong format without hyphens",
			snils:   "12345678912",
			wantErr: ErrInvalidSnils,
		},
		{
			name:    "error wrong format short",
			snils:   "123-456-789 1",
			wantErr: ErrInvalidSnils,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := basePerson()

			err := p.ChangeSnils(tt.snils)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}

			if err == nil {
				if p.Snils() != tt.want {
					t.Fatalf("expected snils %v, got %v", tt.want, p.Snils())
				}
			}
		})
	}
}

func TestChangeOrganization(t *testing.T) {
	tests := []struct {
		name         string
		initialOrgID uuid.UUID
		newOrgID     uuid.UUID
		wantOrgID    uuid.UUID
		wantHasOrg   bool
	}{
		{
			name:       "change to another organization",
			newOrgID:   uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			wantOrgID:  uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			wantHasOrg: true,
		},
		{
			name:       "assign organization from nil",
			newOrgID:   uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			wantOrgID:  uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			wantHasOrg: true,
		},
		{
			name:       "change to nil removes organization",
			newOrgID:   uuid.Nil,
			wantOrgID:  uuid.Nil,
			wantHasOrg: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := basePerson()

			p.ChangeOrganization(tt.newOrgID)

			if p.OrganizationID() != tt.wantOrgID {
				t.Errorf("expected organization id %v, got %v", tt.wantOrgID, p.OrganizationID())

				if p.HasOrganization() != tt.wantHasOrg {
					t.Errorf("expected hasOrganization %v, got %v", tt.wantHasOrg, p.HasOrganization())
				}
			}
		})
	}
}

func TestRemoveOrganization(t *testing.T) {
	tests := []struct {
		name         string
		initialOrgID uuid.UUID
		wantOrgID    uuid.UUID
		wantHasOrg   bool
	}{
		{
			name:         "remove existing organization",
			initialOrgID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			wantOrgID:    uuid.Nil,
			wantHasOrg:   false,
		},
		{
			name:         "remove when already nil",
			initialOrgID: uuid.Nil,
			wantOrgID:    uuid.Nil,
			wantHasOrg:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := basePerson()

			p.RemoveOrganization()

			if p.OrganizationID() != tt.wantOrgID {
				t.Errorf("expected organization id %v, got %v", tt.wantOrgID, p.OrganizationID())

				if p.HasOrganization() != tt.wantHasOrg {
					t.Errorf("expected hasOrganization %v, got %v", tt.wantHasOrg, p.HasOrganization())
				}
			}
		})
	}
}

func createParams(firstName, lastName, middleName, jobTitle, snils, education string, organization uuid.UUID, d time.Time) Params {
	db, err := dob.NewDateOfBirth(d)
	if err != nil {
		panic("err dob")
	}

	return Params{
		FirstName:      firstName,
		LastName:       lastName,
		MiddleName:     middleName,
		JobTitle:       jobTitle,
		Snils:          snils,
		Education:      education,
		OrganizationID: organization,
		DateOfBirth:    db,
	}
}

func basePerson() *Person {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	return &Person{
		id:             uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		firstName:      "Иван",
		lastName:       "Петров",
		middleName:     "Сергеевич",
		jobTitle:       "Преподаватель",
		snils:          "123-456-789 12",
		education:      "Высшее",
		organizationID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		createdAt:      now,
		updatedAt:      now,
	}
}
