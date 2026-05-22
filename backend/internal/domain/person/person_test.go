package person

import (
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/person/profile"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		n name.Name
	}
	tests := []struct {
		name    string
		args    args
		want    *Person
		wantErr bool
	}{
		{
			name: "пустое имя",
			args: args{
				n: name.Name{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "непустое имя",
			args: args{
				n: nameFixture(t),
			},
			want: &Person{
				id:      uuid.UUID{},
				name:    nameFixture(t),
				profile: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotEqual(t, uuid.Nil, got.id)
				assert.Equal(t, tt.want.name, got.name)
			}
		})
	}
}

func TestPerson_AttachOrReplaceProfile(t *testing.T) {
	type fields struct {
		id      uuid.UUID
		name    name.Name
		profile *profile.Profile
	}
	type args struct {
		prof profile.Profile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "прикрепляем профиль",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: nil,
			},
			args: args{
				prof: profileFixture(t),
			},
			wantErr: false,
		},
		{
			name: "прикрепляем профиль с ошибкой",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: nil,
			},
			args: args{
				prof: profile.Profile{},
			},
			wantErr: true,
		},
		{
			name: "заменяем профиль",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: &profile.Profile{},
			},
			args: args{
				prof: profileFixture(t),
			},
			wantErr: false,
		},
		{
			name: "заменяем профиль с ошибкой",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: &profile.Profile{},
			},
			args: args{
				prof: profile.Profile{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				id:      tt.fields.id,
				name:    tt.fields.name,
				profile: tt.fields.profile,
			}
			if err := p.AttachOrReplaceProfile(tt.args.prof); (err != nil) != tt.wantErr {
				t.Errorf("AttachOrReplaceProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPerson_ChangeName(t *testing.T) {
	type fields struct {
		id      uuid.UUID
		name    name.Name
		profile *profile.Profile
	}
	type args struct {
		n name.Name
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "пустое имя",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: nil,
			},
			args: args{
				n: name.Name{},
			},
			wantErr: true,
		},
		{
			name: "непустое имя",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: nil,
			},
			args: args{
				n: nameFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				id:      tt.fields.id,
				name:    tt.fields.name,
				profile: tt.fields.profile,
			}
			if err := p.ChangeName(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("ChangeName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPerson_DetachProfile(t *testing.T) {
	type fields struct {
		id      uuid.UUID
		name    name.Name
		profile *profile.Profile
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "есть профиль",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: &profile.Profile{},
			},
		},
		{
			name: "нет профиля",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				id:      tt.fields.id,
				name:    tt.fields.name,
				profile: tt.fields.profile,
			}
			p.DetachProfile()
		})
	}
}

func TestPerson_HasProfile(t *testing.T) {
	type fields struct {
		id      uuid.UUID
		name    name.Name
		profile *profile.Profile
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "есть профиль",
			fields: fields{
				id:      idFixture,
				name:    nameFixture(t),
				profile: &profile.Profile{},
			},
			want: true,
		},
		{
			name: "нет профиля",
			fields: fields{
				id:      idFixture,
				name:    nameFixture(t),
				profile: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				id:      tt.fields.id,
				name:    tt.fields.name,
				profile: tt.fields.profile,
			}
			if got := p.HasProfile(); got != tt.want {
				t.Errorf("HasProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPerson_ID(t *testing.T) {
	type fields struct {
		id      uuid.UUID
		name    name.Name
		profile *profile.Profile
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				id:      idFixture,
				name:    name.Name{},
				profile: nil,
			},
			want: idFixture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				id:      tt.fields.id,
				name:    tt.fields.name,
				profile: tt.fields.profile,
			}
			if got := p.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPerson_Name(t *testing.T) {
	type fields struct {
		id      uuid.UUID
		name    name.Name
		profile *profile.Profile
	}
	tests := []struct {
		name   string
		fields fields
		want   name.Name
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				id:      uuid.UUID{},
				name:    nameFixture(t),
				profile: nil,
			},
			want: nameFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				id:      tt.fields.id,
				name:    tt.fields.name,
				profile: tt.fields.profile,
			}
			if got := p.Name(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPerson_Profile(t *testing.T) {
	type fields struct {
		id      uuid.UUID
		name    name.Name
		profile profile.Profile
	}
	tests := []struct {
		name   string
		fields fields
		want   profile.Profile
		ok     bool
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: profileFixture(t),
			},
			want: profileFixture(t),
			ok:   true,
		},
		{
			name: "возвращает пустое значение",
			fields: fields{
				id:      uuid.UUID{},
				name:    name.Name{},
				profile: profile.Profile{},
			},
			want: profile.Profile{},
			ok:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Person{
				id:      tt.fields.id,
				name:    tt.fields.name,
				profile: &tt.fields.profile,
			}

			if tt.fields.profile.IsZero() {
				p.profile = nil
			}

			got, ok := p.Profile()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Profile() got = %v, want %v", got, tt.want)
			}
			if ok != tt.ok {
				t.Errorf("Profile() ok = %v, want %v", ok, tt.ok)
			}
		})
	}
}

func TestRestore(t *testing.T) {
	type args struct {
		id   uuid.UUID
		n    name.Name
		prof *profile.Profile
	}
	tests := []struct {
		name    string
		args    args
		want    *Person
		wantErr bool
	}{
		{
			name: "нет идентификатора",
			args: args{
				id:   uuid.UUID{},
				n:    nameFixture(t),
				prof: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "нет имени",
			args: args{
				id:   idFixture,
				n:    name.Name{},
				prof: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "есть пустой профиль",
			args: args{
				id:   idFixture,
				n:    nameFixture(t),
				prof: &profile.Profile{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "есть все данные",
			args: args{
				id:   idFixture,
				n:    nameFixture(t),
				prof: nil,
			},
			want: &Person{
				id:      idFixture,
				name:    nameFixture(t),
				profile: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Restore(tt.args.id, tt.args.n, tt.args.prof)
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
