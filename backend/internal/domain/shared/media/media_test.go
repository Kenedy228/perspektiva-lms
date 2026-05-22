package media

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/shared/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustFile(t *testing.T, key string, size int64) file.File {
	t.Helper()

	f, err := file.New(key, size)
	if err != nil {
		t.Fatalf("file.New() error = %v", err)
	}

	return f
}

func TestMedia_Incomplete(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		m, err := New(TypeVideo, mustFile(t, "content/video/lesson-1.mp4", 1024))
		require.NoError(t, err)

		assert.False(t, m.IsIncomplete())
	})

	t.Run("incomplete", func(t *testing.T) {
		m := Media{}

		assert.True(t, m.IsIncomplete())
	})
}

func TestMedia_File(t *testing.T) {
	f := mustFile(t, "content/video/lesson-1.mp4", 1024)
	m, err := New(TypeVideo, f)
	require.NoError(t, err)

	assert.Equal(t, f, m.File())
}

func TestMedia_Type(t *testing.T) {
	m, err := New(TypeVideo, mustFile(t, "content/video/lesson-1.mp4", 1024))
	require.NoError(t, err)

	assert.Equal(t, TypeVideo, m.Type())
}

func TestNew(t *testing.T) {
	type args struct {
		t Type
		f file.File
	}
	tests := []struct {
		name    string
		args    args
		want    Media
		wantErr bool
	}{
		{
			name: "valid slides",
			args: args{
				t: TypeSlides,
				f: mustFile(t, "content/slides/topic-1.pptx", 1024),
			},
			want: Media{
				mType: TypeSlides,
				file:  mustFile(t, "content/slides/topic-1.pptx", 1024),
			},
			wantErr: false,
		},
		{
			name: "invalid unknown type",
			args: args{
				t: Type("unknown"),
				f: mustFile(t, "content/audio/task.mp3", 1024),
			},
			want:    Media{},
			wantErr: true,
		},
		{
			name: "invalid extension for slides",
			args: args{
				t: TypeSlides,
				f: mustFile(t, "content/slides/topic-1.pdf", 1024),
			},
			want:    Media{},
			wantErr: true,
		},
		{
			name: "size too large for image",
			args: args{
				t: TypeImage,
				f: mustFile(t, "content/image/picture.jpg", TypeImage.MaxSizeInBytes()+1),
			},
			want:    Media{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.t, tt.args.f)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestValidateFileExtension(t *testing.T) {
	tests := []struct {
		name    string
		allowed []string
		actual  string
		wantErr bool
	}{
		{
			name:    "allowed contains actual",
			allowed: []string{".pdf", ".pptx"},
			actual:  ".pdf",
			wantErr: false,
		},
		{
			name:    "allowed contains part of actual",
			allowed: []string{".pdf", ".pptx"},
			actual:  "pdf",
			wantErr: true,
		},
		{
			name:    "allowed empty",
			allowed: []string{},
			actual:  ".pdf",
			wantErr: true,
		},
		{
			name:    "actual empty",
			allowed: []string{".pdf", ".pptx"},
			actual:  "",
			wantErr: true,
		},
		{
			name:    "allowed doesn't contain actual",
			allowed: []string{".pdf", ".pptx"},
			actual:  ".jpeg",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFileExtension(tt.actual, tt.allowed)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFileForType(t *testing.T) {
	type args struct {
		t Type
		f file.File
	}

	type testCase struct {
		name    string
		args    args
		wantErr bool
	}

	testCases := []testCase{
		{
			name: "valid video",
			args: args{
				t: TypeVideo,
				f: mustFile(t, "content/video/lesson-1.mp4", 1024),
			},
			wantErr: false,
		},
		{
			name: "valid image boundary size",
			args: args{
				t: TypeImage,
				f: mustFile(t, "content/image/picture.jpg", TypeImage.MaxSizeInBytes()),
			},
			wantErr: false,
		},
		{
			name: "valid pdf document",
			args: args{
				t: TypeDocument,
				f: mustFile(t, "content/docs/topic-1.pdf", 1024),
			},
			wantErr: false,
		},
		{
			name: "invalid document extension",
			args: args{
				t: TypeDocument,
				f: mustFile(t, "content/docs/topic-1.docx", 1024),
			},
			wantErr: true,
		},
		{
			name: "invalid size",
			args: args{
				t: TypeAudio,
				f: mustFile(t, "content/audio/task.wav", TypeAudio.MaxSizeInBytes()+1),
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateFileForType(tc.args.t, tc.args.f)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFileComplete(t *testing.T) {
	t.Run("incomplete file", func(t *testing.T) {
		err := validateFileComplete(file.File{})

		assert.Error(t, err)
	})

	t.Run("complete file", func(t *testing.T) {
		err := validateFileComplete(mustFile(t, "content", 100))

		assert.NoError(t, err)
	})
}

func TestValidateFileSize(t *testing.T) {
	type testCase struct {
		name    string
		actual  int64
		max     int64
		wantErr bool
	}

	testCases := []testCase{
		{
			name:    "actual greater than max",
			actual:  100,
			max:     10,
			wantErr: true,
		},
		{
			name:    "actual is negative",
			actual:  -100,
			max:     100,
			wantErr: false,
		},
		{
			name:    "max is negative",
			actual:  100,
			max:     -100,
			wantErr: true,
		},
		{
			name:    "actual equal to max",
			actual:  100,
			max:     100,
			wantErr: false,
		},
		{
			name:    "zero actual size",
			actual:  0,
			max:     100,
			wantErr: false,
		},
		{
			name:    "both positive",
			actual:  100,
			max:     200,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateFileSize(tc.actual, tc.max)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateType(t *testing.T) {
	type testCase struct {
		name    string
		mType   Type
		wantErr bool
	}

	testCases := []testCase{
		{
			name:    "slides valid",
			mType:   TypeSlides,
			wantErr: false,
		},
		{
			name:    "document valid",
			mType:   TypeDocument,
			wantErr: false,
		},
		{
			name:    "video valid",
			mType:   TypeVideo,
			wantErr: false,
		},
		{
			name:    "image valid",
			mType:   TypeImage,
			wantErr: false,
		},
		{
			name:    "audio valid",
			mType:   TypeAudio,
			wantErr: false,
		},
		{
			name:    "unknown invalid",
			mType:   Type("unknown"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateType(tc.mType)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
