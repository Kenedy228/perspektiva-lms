package option

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		isCorrect bool
		wantValue string
		wantErr   bool
	}{
		{name: "creates correct option", value: "option", isCorrect: true, wantValue: "option"},
		{name: "creates incorrect option", value: "option", wantValue: "option"},
		{name: "normalizes option value", value: "  option   text  ", wantValue: "option text"},
		{name: "returns error for empty value", value: "", wantErr: true},
		{name: "returns error for blank value", value: "   ", wantErr: true},
		{name: "returns error for too long value", value: strings.Repeat("а", ValueCharsLimit+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.value, tt.isCorrect)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("New() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}
			if got.ID() == uuid.Nil {
				t.Fatalf("New() ID() = nil, want generated id")
			}
			if got.Value() != tt.wantValue {
				t.Fatalf("New() Value() = %q, want %q", got.Value(), tt.wantValue)
			}
			if got.IsCorrect() != tt.isCorrect {
				t.Fatalf("New() IsCorrect() = %v, want %v", got.IsCorrect(), tt.isCorrect)
			}
		})
	}
}

func TestRestore(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		value     string
		isCorrect bool
		want      Option
		wantErr   bool
	}{
		{name: "restores option", id: id, value: " option ", isCorrect: true, want: Option{id: id, value: "option", isCorrect: true}},
		{name: "returns error for empty id", id: uuid.Nil, value: "option", wantErr: true},
		{name: "returns error for empty value", id: id, value: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Restore(tt.id, tt.value, tt.isCorrect)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Restore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptionAccessors(t *testing.T) {
	id := uuid.New()
	o := Option{id: id, value: "option", isCorrect: true}

	if o.ID() != id {
		t.Fatalf("ID() = %v, want %v", o.ID(), id)
	}
	if o.Value() != "option" {
		t.Fatalf("Value() = %q, want %q", o.Value(), "option")
	}
	if !o.IsCorrect() {
		t.Fatalf("IsCorrect() = false, want true")
	}
	if o.IsZero() {
		t.Fatalf("IsZero() = true, want false")
	}
	if !(Option{}).IsZero() {
		t.Fatalf("zero Option must be zero")
	}
}

func Test_validateValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid value", value: "option"},
		{name: "empty value", value: "", wantErr: true},
		{name: "too long value", value: strings.Repeat("а", ValueCharsLimit+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateValue(tt.value); (err != nil) != tt.wantErr {
				t.Fatalf("validateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
