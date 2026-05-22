package name

import "testing"

func Test_validateFirstName(t *testing.T) {
	type args struct {
		firstName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустое значение",
			args: args{
				firstName: "",
			},
			wantErr: true,
		},
		{
			name: "значение из пробелов",
			args: args{
				firstName: " ",
			},
			wantErr: true,
		},
		{
			name: "начало специальный символ",
			args: args{
				firstName: ".иван",
			},
			wantErr: true,
		},
		{
			name: "конец специальный символ",
			args: args{
				firstName: "иван-",
			},
			wantErr: true,
		},
		{
			name: "конец .",
			args: args{
				firstName: "иван.",
			},
			wantErr: false,
		},
		{
			name: "двойной специальный символ",
			args: args{
				firstName: "иван--петр",
			},
			wantErr: true,
		},
		{
			name: "подряд разные специальные символы без пробела",
			args: args{
				firstName: "иван-,петр",
			},
			wantErr: true,
		},
		{
			name: "цифры",
			args: args{
				firstName: "иван123",
			},
			wantErr: true,
		},
		{
			name: "строчные римские буквы",
			args: args{
				firstName: "иванiv",
			},
			wantErr: true,
		},
		{
			name: "только специальный символ",
			args: args{
				firstName: ".",
			},
			wantErr: true,
		},
		{
			name: "нет кириллицы",
			args: args{
				firstName: "ivan",
			},
			wantErr: true,
		},
		{
			name: "валидное значение",
			args: args{
				firstName: "Иван-Петр",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateFirstName(tt.args.firstName); (err != nil) != tt.wantErr {
				t.Errorf("validateFirstName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateLastName(t *testing.T) {
	type args struct {
		lastName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустое значение",
			args: args{
				lastName: "",
			},
			wantErr: true,
		},
		{
			name: "значение из пробелов",
			args: args{
				lastName: " ",
			},
			wantErr: true,
		},
		{
			name: "начало специальный символ",
			args: args{
				lastName: ".иванов",
			},
			wantErr: true,
		},
		{
			name: "конец специальный символ",
			args: args{
				lastName: "иванов.",
			},
			wantErr: true,
		},
		{
			name: "двойной специальный символ",
			args: args{
				lastName: "иванов--петров",
			},
			wantErr: true,
		},
		{
			name: "подряд разные специальные символы без пробела",
			args: args{
				lastName: "иванов-,петров",
			},
			wantErr: true,
		},
		{
			name: "цифры",
			args: args{
				lastName: "иванов123",
			},
			wantErr: true,
		},
		{
			name: "строчные римские буквы",
			args: args{
				lastName: "ивановiv",
			},
			wantErr: true,
		},
		{
			name: "только специальный символ",
			args: args{
				lastName: ".",
			},
			wantErr: true,
		},
		{
			name: "нет кириллицы",
			args: args{
				lastName: "ivanov",
			},
			wantErr: true,
		},
		{
			name: "валидное значение",
			args: args{
				lastName: "Иванов-Петров",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateLastName(tt.args.lastName); (err != nil) != tt.wantErr {
				t.Errorf("validateLastName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateMiddleName(t *testing.T) {
	type args struct {
		middleName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустое значение",
			args: args{
				middleName: "",
			},
			wantErr: false,
		},
		{
			name: "значение из пробелов",
			args: args{
				middleName: " ",
			},
			wantErr: true,
		},
		{
			name: "начало специальный символ",
			args: args{
				middleName: ".иванович",
			},
			wantErr: true,
		},
		{
			name: "конец специальный символ",
			args: args{
				middleName: "иванович-",
			},
			wantErr: true,
		},
		{
			name: "конец .",
			args: args{
				middleName: "иванович.",
			},
			wantErr: false,
		},
		{
			name: "двойной специальный символ",
			args: args{
				middleName: "иванович--петрович",
			},
			wantErr: true,
		},
		{
			name: "подряд разные специальные символы без пробела",
			args: args{
				middleName: "иванович-,петрович",
			},
			wantErr: true,
		},
		{
			name: "цифры",
			args: args{
				middleName: "иванович123",
			},
			wantErr: true,
		},
		{
			name: "строчные римские буквы",
			args: args{
				middleName: "ивановичiv",
			},
			wantErr: true,
		},
		{
			name: "только специальный символ",
			args: args{
				middleName: ".",
			},
			wantErr: true,
		},
		{
			name: "нет кириллицы",
			args: args{
				middleName: "ivanovich",
			},
			wantErr: true,
		},
		{
			name: "валидное значение",
			args: args{
				middleName: "Иванович-Петрович",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMiddleName(tt.args.middleName); (err != nil) != tt.wantErr {
				t.Errorf("validateMiddleName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
