package organization

import (
	"errors"
	"testing"
)

// New() should return error if name is empty
func TestNew(t *testing.T) {
	_, err := New("")

	if err == nil {
		t.Errorf("expected err not nil")
	}

	if !errors.Is(err, ErrEmptyName) {
		t.Errorf("expected err of type ErrEmptyName, got %T", err)
	}
}

// Rename() should return error if newName is empty
func TestRename(t *testing.T) {
	tests := []struct {
		name        string
		initialName string
		newName     string
		wantName    string
		wantErr     error
	}{
		{
			name:        "rename with valid newName",
			initialName: "old name",
			newName:     "new name",
			wantName:    "new name",
			wantErr:     nil,
		},
		{
			name:        "rename with invalid newName",
			initialName: "old name",
			newName:     "",
			wantName:    "old name",
			wantErr:     ErrEmptyName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			organization, err := New(tt.initialName)
			if err != nil {
				t.Errorf("expected no errors, got %v", err)
			}

			err = organization.Rename(tt.newName)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err %v, got %v", tt.wantErr, err)
			}

			if organization.name != tt.wantName {
				t.Errorf("expected name '%s', got '%s'", tt.wantName, organization.name)
			}
		})
	}
}

// Equal should compare values by id
func TestEqual(t *testing.T) {
	type testCase struct {
		name      string
		first     *Organization
		second    *Organization
		wantEqual bool
	}

	first, _ := New("name")
	second, _ := New("name")

	tests := []testCase{
		{
			name:      "equal values",
			first:     first,
			second:    first,
			wantEqual: true,
		},
		{
			name:      "not equal values",
			first:     first,
			second:    second,
			wantEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equality := tt.first.Equal(tt.second)
			if equality != tt.wantEqual {
				t.Errorf("expected equal %v, got %v", tt.wantEqual, equality)
			}
		})
	}
}
