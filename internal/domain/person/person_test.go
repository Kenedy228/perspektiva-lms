package person

import (
	"testing"
	"time"

	"gitflic.ru/lms/internal/domain/person/dob"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name:           "valid: empty middleName",
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
			params := createParams(
				tt.firstName,
				tt.lastName,
				tt.middleName,
				tt.jobTitle,
				tt.snils,
				tt.education,
				tt.organizationID,
				tt.dateOfBirth,
			)

			p, err := New(params)

			assert.ErrorIs(t, err, tt.expectedErr)

			if tt.expectedErr != nil {
				assert.Nil(t, p)
				return
			}

			require.NotNil(t, p)
			assert.Equal(t, tt.firstName, p.FirstName())
			assert.Equal(t, tt.lastName, p.LastName())
			assert.Equal(t, tt.middleName, p.MiddleName())
			assert.Equal(t, tt.jobTitle, p.JobTitle())
			assert.Equal(t, tt.snils, p.Snils())
			assert.Equal(t, tt.education, p.Education())
			assert.Equal(t, tt.organizationID, p.OrganizationID())
			assert.NotEqual(t, uuid.Nil, p.ID())
			assert.False(t, p.CreatedAt().IsZero())
			assert.False(t, p.UpdatedAt().IsZero())
			assert.Equal(t, p.CreatedAt(), p.UpdatedAt())
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
			oldUpdatedAt := p.UpdatedAt()

			err := p.Rename(tt.firstName, tt.lastName, tt.middleName)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr != nil {
				assert.Equal(t, "Иван", p.FirstName())
				assert.Equal(t, "Петров", p.LastName())
				assert.Equal(t, "Сергеевич", p.MiddleName())
				assert.Equal(t, oldUpdatedAt, p.UpdatedAt())
				return
			}

			assert.Equal(t, tt.wantFirstName, p.FirstName())
			assert.Equal(t, tt.wantLastName, p.LastName())
			assert.Equal(t, tt.wantMiddleName, p.MiddleName())
			assert.True(t, p.UpdatedAt().After(oldUpdatedAt) || p.UpdatedAt().Equal(oldUpdatedAt))
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
			oldValue := p.JobTitle()
			oldUpdatedAt := p.UpdatedAt()

			err := p.ChangeJobTitle(tt.jobTitle)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr != nil {
				assert.Equal(t, oldValue, p.JobTitle())
				assert.Equal(t, oldUpdatedAt, p.UpdatedAt())
				return
			}

			assert.Equal(t, tt.want, p.JobTitle())
			assert.True(t, p.UpdatedAt().After(oldUpdatedAt) || p.UpdatedAt().Equal(oldUpdatedAt))
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
			oldValue := p.Education()
			oldUpdatedAt := p.UpdatedAt()

			err := p.ChangeEducation(tt.education)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr != nil {
				assert.Equal(t, oldValue, p.Education())
				assert.Equal(t, oldUpdatedAt, p.UpdatedAt())
				return
			}

			assert.Equal(t, tt.want, p.Education())
			assert.True(t, p.UpdatedAt().After(oldUpdatedAt) || p.UpdatedAt().Equal(oldUpdatedAt))
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
			oldValue := p.Snils()
			oldUpdatedAt := p.UpdatedAt()

			err := p.ChangeSnils(tt.snils)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr != nil {
				assert.Equal(t, oldValue, p.Snils())
				assert.Equal(t, oldUpdatedAt, p.UpdatedAt())
				return
			}

			assert.Equal(t, tt.want, p.Snils())
			assert.True(t, p.UpdatedAt().After(oldUpdatedAt) || p.UpdatedAt().Equal(oldUpdatedAt))
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
			name:         "change to another organization",
			initialOrgID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			newOrgID:     uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			wantOrgID:    uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			wantHasOrg:   true,
		},
		{
			name:         "assign organization from nil",
			initialOrgID: uuid.Nil,
			newOrgID:     uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			wantOrgID:    uuid.MustParse("33333333-3333-3333-3333-333333333333"),
			wantHasOrg:   true,
		},
		{
			name:         "change to nil removes organization",
			initialOrgID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			newOrgID:     uuid.Nil,
			wantOrgID:    uuid.Nil,
			wantHasOrg:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := basePerson()
			p.organizationID = tt.initialOrgID
			oldUpdatedAt := p.UpdatedAt()

			p.ChangeOrganization(tt.newOrgID)

			assert.Equal(t, tt.wantOrgID, p.OrganizationID())
			assert.Equal(t, tt.wantHasOrg, p.HasOrganization())
			assert.True(t, p.UpdatedAt().After(oldUpdatedAt) || p.UpdatedAt().Equal(oldUpdatedAt))
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
			p.organizationID = tt.initialOrgID
			oldUpdatedAt := p.UpdatedAt()

			p.RemoveOrganization()

			assert.Equal(t, tt.wantOrgID, p.OrganizationID())
			assert.Equal(t, tt.wantHasOrg, p.HasOrganization())
			assert.True(t, p.UpdatedAt().After(oldUpdatedAt) || p.UpdatedAt().Equal(oldUpdatedAt))
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

