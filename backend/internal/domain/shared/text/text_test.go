package text

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Text
		wantErr bool
	}{
		{
			name: "valid text",
			args: args{
				value: "Привет, мир",
			},
			want: Text{
				value: "Привет, мир",
			},
			wantErr: false,
		},
		{
			name: "valid text with spaces around",
			args: args{
				value: "  Привет, мир  ",
			},
			want: Text{
				value: "  Привет, мир  ",
			},
			wantErr: false,
		},
		{
			name: "empty text",
			args: args{
				value: "",
			},
			want:    Text{},
			wantErr: true,
		},
		{
			name: "spaces only",
			args: args{
				value: "   \n\t   ",
			},
			want:    Text{},
			wantErr: true,
		},
		{
			name: "boundary chars limit",
			args: args{
				value: strings.Repeat("а", ValueCharsLimit),
			},
			want: Text{
				value: strings.Repeat("а", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "exceeds chars limit",
			args: args{
				value: strings.Repeat("а", ValueCharsLimit+1),
			},
			want:    Text{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestText_Value(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "returns value",
			fields: fields{
				value: "контент вопроса",
			},
			want: "контент вопроса",
		},
		{
			name: "returns value with spaces",
			fields: fields{
				value: "  контент вопроса  ",
			},
			want: "  контент вопроса  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Text{
				value: tt.fields.value,
			}

			assert.Equal(t, tt.want, got.Value())
		})
	}
}

func TestText_IsIncomplete(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		got, err := New("text")
		require.NoError(t, err)

		assert.False(t, got.IsIncomplete())
	})

	t.Run("incomplete", func(t *testing.T) {
		got := Text{}

		assert.True(t, got.IsIncomplete())
	})
}

func Test_validateCharsLimit(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "short text",
			args: args{
				value: "привет",
			},
			wantErr: false,
		},
		{
			name: "boundary limit",
			args: args{
				value: strings.Repeat("а", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "exceeds limit",
			args: args{
				value: strings.Repeat("а", ValueCharsLimit+1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCharsLimit(tt.args.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_validateNotEmptyValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "regular text",
			args: args{
				value: "привет",
			},
			wantErr: false,
		},
		{
			name: "text with surrounding spaces",
			args: args{
				value: "  привет  ",
			},
			wantErr: false,
		},
		{
			name: "empty string",
			args: args{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "spaces only",
			args: args{
				value: "   ",
			},
			wantErr: true,
		},
		{
			name: "tabs and newlines only",
			args: args{
				value: "\n\t  ",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNotEmptyValue(tt.args.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_validateValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid text",
			args: args{
				value: "контент вопроса",
			},
			wantErr: false,
		},
		{
			name: "valid text with spaces around",
			args: args{
				value: "  контент вопроса  ",
			},
			wantErr: false,
		},
		{
			name: "empty text",
			args: args{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "spaces only",
			args: args{
				value: "   \n\t  ",
			},
			wantErr: true,
		},
		{
			name: "boundary limit",
			args: args{
				value: strings.Repeat("а", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "exceeds limit",
			args: args{
				value: strings.Repeat("а", ValueCharsLimit+1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateValue(tt.args.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
