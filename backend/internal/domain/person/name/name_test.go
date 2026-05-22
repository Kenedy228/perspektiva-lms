package name

import (
	"reflect"
	"testing"
)

func TestName_FirstName(t *testing.T) {
	type fields struct {
		firstName  string
		lastName   string
		middleName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "возвращает имя как есть",
			fields: fields{
				firstName:  "Иван",
				lastName:   "Иванов",
				middleName: "Иванович",
			},
			want: "Иван",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Name{
				firstName:  tt.fields.firstName,
				lastName:   tt.fields.lastName,
				middleName: tt.fields.middleName,
			}
			if got := n.FirstName(); got != tt.want {
				t.Errorf("FirstName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName_IsZero(t *testing.T) {
	type fields struct {
		firstName  string
		lastName   string
		middleName string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "все поля пустые",
			fields: fields{
				firstName:  "",
				lastName:   "",
				middleName: "",
			},
			want: true,
		},
		{
			name: "поля заполнены",
			fields: fields{
				firstName:  "s",
				lastName:   "s",
				middleName: "s",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Name{
				firstName:  tt.fields.firstName,
				lastName:   tt.fields.lastName,
				middleName: tt.fields.middleName,
			}
			if got := n.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName_LastName(t *testing.T) {
	type fields struct {
		firstName  string
		lastName   string
		middleName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "возвращает фамилию как есть",
			fields: fields{
				firstName:  "Иван",
				lastName:   "Иванов",
				middleName: "Иванович",
			},
			want: "Иванов",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Name{
				firstName:  tt.fields.firstName,
				lastName:   tt.fields.lastName,
				middleName: tt.fields.middleName,
			}
			if got := n.LastName(); got != tt.want {
				t.Errorf("LastName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName_MiddleName(t *testing.T) {
	type fields struct {
		firstName  string
		lastName   string
		middleName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "возвращает отчество как есть",
			fields: fields{
				firstName:  "Иван",
				lastName:   "Иванов",
				middleName: "Иванович",
			},
			want: "Иванович",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Name{
				firstName:  tt.fields.firstName,
				lastName:   tt.fields.lastName,
				middleName: tt.fields.middleName,
			}
			if got := n.MiddleName(); got != tt.want {
				t.Errorf("MiddleName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		firstName  string
		lastName   string
		middleName string
	}
	tests := []struct {
		name    string
		args    args
		want    Name
		wantErr bool
	}{
		{
			name: "некорректное имя",
			args: args{
				firstName:  "123",
				lastName:   "Иванов",
				middleName: "Иванович",
			},
			want: Name{
				firstName:  "",
				lastName:   "",
				middleName: "",
			},
			wantErr: true,
		},
		{
			name: "некорректная фамилия",
			args: args{
				firstName:  "Иван",
				lastName:   "123",
				middleName: "Иванович",
			},
			want: Name{
				firstName:  "",
				lastName:   "",
				middleName: "",
			},
			wantErr: true,
		},
		{
			name: "некорректное отчество",
			args: args{
				firstName:  "Иван",
				lastName:   "Иванов",
				middleName: "123",
			},
			want: Name{
				firstName:  "",
				lastName:   "",
				middleName: "",
			},
			wantErr: true,
		},
		{
			name: "пустое отчество",
			args: args{
				firstName:  "Иван",
				lastName:   "Иванов",
				middleName: "",
			},
			want: Name{
				firstName:  "Иван",
				lastName:   "Иванов",
				middleName: "",
			},
			wantErr: false,
		},
		{
			name: "непустое отчество",
			args: args{
				firstName:  "Иван",
				lastName:   "Иванов",
				middleName: "Иванович",
			},
			want: Name{
				firstName:  "Иван",
				lastName:   "Иванов",
				middleName: "Иванович",
			},
			wantErr: false,
		},
		{
			name: "составное имя и фамилия",
			args: args{
				firstName:  "Иван-Петр",
				lastName:   "Иванов-Петров",
				middleName: "Иванович",
			},
			want: Name{
				firstName:  "Иван-Петр",
				lastName:   "Иванов-Петров",
				middleName: "Иванович",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.firstName, tt.args.lastName, tt.args.middleName)
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
