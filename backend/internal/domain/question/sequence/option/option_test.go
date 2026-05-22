package option

import (
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
	type args struct {
		t text.Text
	}
	tests := []struct {
		name    string
		args    args
		want    Option
		wantErr bool
	}{
		{
			name: "creates option with valid text",
			args: args{
				t: mustText(t, "option"),
			},
			wantErr: false,
		},
		{
			name: "error when text exceeds limit",
			args: args{
				t: mustText(t, strings.Repeat("x", TextCharsLimit+1)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.ID() == uuid.Nil {
				t.Errorf("New() ID() = %v, want non-nil uuid", got.ID())
			}
			if !reflect.DeepEqual(got.Text(), tt.args.t) {
				t.Errorf("New() Text() = %v, want %v", got.Text(), tt.args.t)
			}
		})
	}
}

func TestOption_ID(t *testing.T) {
	id := uuid.New()
	txt := mustText(t, "opt")

	type fields struct {
		id   uuid.UUID
		text text.Text
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns id",
			fields: fields{
				id:   id,
				text: txt,
			},
			want: id,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Option{
				id:   tt.fields.id,
				text: tt.fields.text,
			}
			if got := o.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOption_Text(t *testing.T) {
	id := uuid.New()
	txt := mustText(t, "opt")

	type fields struct {
		id   uuid.UUID
		text text.Text
	}
	tests := []struct {
		name   string
		fields fields
		want   text.Text
	}{
		{
			name: "returns text",
			fields: fields{
				id:   id,
				text: txt,
			},
			want: txt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Option{
				id:   tt.fields.id,
				text: tt.fields.text,
			}
			if got := o.Text(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateText(t *testing.T) {
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
				t: mustText(t, "opt"),
			},
			wantErr: false,
		},
		{
			name: "too long",
			args: args{
				t: mustText(t, strings.Repeat("x", TextCharsLimit+1)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateText(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("validateText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateTextCharsLimit(t *testing.T) {
	type args struct {
		t text.Text
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "within limit",
			args: args{
				t: mustText(t, strings.Repeat("x", TextCharsLimit)),
			},
			wantErr: false,
		},
		{
			name: "over limit",
			args: args{
				t: mustText(t, strings.Repeat("x", TextCharsLimit+1)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateTextCharsLimit(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("validateTextCharsLimit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
