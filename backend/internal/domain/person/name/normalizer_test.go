package name

import "testing"

func Test_normalizeFirstName(t *testing.T) {
	type args struct {
		firstName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "простое имя с пробелами по краям",
			args: args{
				firstName: " иван ",
			},
			want: "Иван",
		},
		{
			name: "составное имя с пробелами по краям",
			args: args{
				firstName: " иван иван ",
			},
			want: "Иван Иван",
		},
		{
			name: "составное имя с несколькими пробелами внутри",
			args: args{
				firstName: "иван    иван",
			},
			want: "Иван Иван",
		},
		{
			name: "составное имя с разделителем '-'",
			args: args{
				firstName: " иван-иван ",
			},
			want: "Иван-Иван",
		},
		{
			name: "составное имя с разделителем '",
			args: args{
				firstName: " иван'иван",
			},
			want: "Иван'Иван",
		},
		{
			name: "составное имя с разделителем ' и пробелами между разделителем",
			args: args{
				firstName: "иван ' иван",
			},
			want: "Иван ' Иван",
		},
		{
			name: "составное имя с разделителем - и пробелами между разделителем",
			args: args{
				firstName: "иван -  иван",
			},
			want: "Иван - Иван",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeFirstName(tt.args.firstName); got != tt.want {
				t.Errorf("normalizeFirstName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeLastName(t *testing.T) {
	type args struct {
		lastName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "простая фамилия с пробелами по краям",
			args: args{
				lastName: " иванов ",
			},
			want: "Иванов",
		},
		{
			name: "составная фамилия с пробелами по краям",
			args: args{
				lastName: " иванов петров ",
			},
			want: "Иванов Петров",
		},
		{
			name: "составная фамилия с несколькими пробелами внутри",
			args: args{
				lastName: "иванов    петров",
			},
			want: "Иванов Петров",
		},
		{
			name: "составная фамилия с разделителем '-'",
			args: args{
				lastName: " иванов-петров ",
			},
			want: "Иванов-Петров",
		},
		{
			name: "составная фамилия с разделителем '",
			args: args{
				lastName: " иванов'петров",
			},
			want: "Иванов'Петров",
		},
		{
			name: "составная фамилия с разделителем ' и пробелами между разделителем",
			args: args{
				lastName: "иванов ' петров",
			},
			want: "Иванов ' Петров",
		},
		{
			name: "составная фамилия с разделителем - и пробелами между разделителем",
			args: args{
				lastName: "иванов -  петров",
			},
			want: "Иванов - Петров",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeLastName(tt.args.lastName); got != tt.want {
				t.Errorf("normalizeLastName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeMiddleName(t *testing.T) {
	type args struct {
		middleName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "простое отчество с пробелами по краям",
			args: args{
				middleName: " иванович ",
			},
			want: "Иванович",
		},
		{
			name: "составное отчество с пробелами по краям",
			args: args{
				middleName: " иванович петрович ",
			},
			want: "Иванович Петрович",
		},
		{
			name: "составное отчество с несколькими пробелами внутри",
			args: args{
				middleName: "иванович    петрович",
			},
			want: "Иванович Петрович",
		},
		{
			name: "составное отчество с разделителем '-'",
			args: args{
				middleName: " иванович-петрович ",
			},
			want: "Иванович-Петрович",
		},
		{
			name: "составное отчество с разделителем '",
			args: args{
				middleName: " иванович'петрович",
			},
			want: "Иванович'Петрович",
		},
		{
			name: "составное отчество с разделителем ' и пробелами между разделителем",
			args: args{
				middleName: "иванович ' петрович",
			},
			want: "Иванович ' Петрович",
		},
		{
			name: "составное отчество с разделителем - и пробелами между разделителем",
			args: args{
				middleName: "иванович -  петрович",
			},
			want: "Иванович - Петрович",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeMiddleName(tt.args.middleName); got != tt.want {
				t.Errorf("normalizeMiddleName() = %v, want %v", got, tt.want)
			}
		})
	}
}
