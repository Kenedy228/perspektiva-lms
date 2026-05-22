package option

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"github.com/google/uuid"
)

func mustText(t *testing.T, v string) text.Text {
	t.Helper()

	got, err := text.New(v)
	if err != nil {
		t.Fatalf("text.New() error = %v", err)
	}

	return got
}

func TestNew(t *testing.T) {
	validText := mustText(t, "option")
	longText := mustText(t, strings.Repeat("а", TextCharsLimit+1))

	type args struct {
		t         text.Text
		isCorrect bool
	}
	tests := []struct {
		name        string
		args        args
		wantText    text.Text
		wantCorrect bool
		wantErr     bool
	}{
		{
			name: "creates correct option",
			args: args{
				t:         validText,
				isCorrect: true,
			},
			wantText:    validText,
			wantCorrect: true,
			wantErr:     false,
		},
		{
			name: "creates incorrect option",
			args: args{
				t:         validText,
				isCorrect: false,
			},
			wantText:    validText,
			wantCorrect: false,
			wantErr:     false,
		},
		{
			name: "returns error for too long text",
			args: args{
				t:         longText,
				isCorrect: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.t, tt.args.isCorrect)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Errorf("New() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}
			if got.ID() == uuid.Nil {
				t.Errorf("New() ID() = %v, want non-nil uuid", got.ID())
			}
			if !reflect.DeepEqual(got.Text(), tt.wantText) {
				t.Errorf("New() Text() = %v, want %v", got.Text(), tt.wantText)
			}
			if got.IsCorrect() != tt.wantCorrect {
				t.Errorf("New() IsCorrect() = %v, want %v", got.IsCorrect(), tt.wantCorrect)
			}
		})
	}
}

func TestOption_ID(t *testing.T) {
	id := uuid.New()
	txt := mustText(t, "option")

	type fields struct {
		id        uuid.UUID
		text      text.Text
		isCorrect bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns id",
			fields: fields{
				id:        id,
				text:      txt,
				isCorrect: true,
			},
			want: id,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Option{
				id:        tt.fields.id,
				text:      tt.fields.text,
				isCorrect: tt.fields.isCorrect,
			}
			if got := o.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOption_IsCorrect(t *testing.T) {
	txt := mustText(t, "option")

	type fields struct {
		id        uuid.UUID
		text      text.Text
		isCorrect bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "correct option",
			fields: fields{
				id:        uuid.New(),
				text:      txt,
				isCorrect: true,
			},
			want: true,
		},
		{
			name: "incorrect option",
			fields: fields{
				id:        uuid.New(),
				text:      txt,
				isCorrect: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Option{
				id:        tt.fields.id,
				text:      tt.fields.text,
				isCorrect: tt.fields.isCorrect,
			}
			if got := o.IsCorrect(); got != tt.want {
				t.Errorf("IsCorrect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOption_Text(t *testing.T) {
	txt := mustText(t, "option")

	type fields struct {
		id        uuid.UUID
		text      text.Text
		isCorrect bool
	}
	tests := []struct {
		name   string
		fields fields
		want   text.Text
	}{
		{
			name: "returns text",
			fields: fields{
				id:        uuid.New(),
				text:      txt,
				isCorrect: true,
			},
			want: txt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Option{
				id:        tt.fields.id,
				text:      tt.fields.text,
				isCorrect: tt.fields.isCorrect,
			}
			if got := o.Text(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateOptionText(t *testing.T) {
	type args struct {
		t text.Text
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid text",
			args: args{
				t: mustText(t, "option"),
			},
			wantErr: false,
		},
		{
			name: "too long text",
			args: args{
				t: mustText(t, strings.Repeat("а", TextCharsLimit+1)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionText(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("validateOptionText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateOptionTextCharsLimit(t *testing.T) {
	type args struct {
		t text.Text
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "text within limit",
			args: args{
				t: mustText(t, strings.Repeat("а", TextCharsLimit)),
			},
			wantErr: false,
		},
		{
			name: "text exceeds limit",
			args: args{
				t: mustText(t, strings.Repeat("а", TextCharsLimit+1)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionTextCharsLimit(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("validateOptionTextCharsLimit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
