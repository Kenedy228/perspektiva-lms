package content

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		cType ContentType
		value string
		err   error
	}{
		{
			name:  "invalid cType",
			cType: "invalid",
			value: "invalid",
			err:   ErrInvalidContentType,
		},
		{
			name:  "valid cType",
			cType: ContentTypeText,
			value: "valid",
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := New(tt.cType, tt.value)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if rc.Value() != tt.value {
					t.Errorf("expected value %v, got %v", tt.value, rc.Value())
				}

				if rc.ContentType() != tt.cType {
					t.Errorf("expected cType %v, got %v", tt.cType, rc.ContentType())
				}
			}
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name        string
		firstCType  ContentType
		firstVal    string
		secondCType ContentType
		secondVal   string
		eq          bool
	}{
		{
			name:        "equal",
			firstCType:  ContentTypeImage,
			firstVal:    "val",
			secondCType: ContentTypeImage,
			secondVal:   "val",
			eq:          true,
		},
		{
			name:        "not cTypes equal",
			firstCType:  ContentTypeImage,
			firstVal:    "val",
			secondCType: ContentTypeText,
			secondVal:   "val",
			eq:          false,
		},
		{
			name:        "not val equal",
			firstCType:  ContentTypeImage,
			firstVal:    "val",
			secondCType: ContentTypeImage,
			secondVal:   "another val",
			eq:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(tt.firstCType, tt.firstVal)
			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			s, err := New(tt.secondCType, tt.secondVal)
			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			if eq := f.Equal(s); eq != tt.eq {
				t.Errorf("expected f Equal s %v, got %v", tt.eq, eq)
			}
		})
	}
}
