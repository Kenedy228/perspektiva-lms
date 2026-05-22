package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFile_IsIncomplete(t *testing.T) {
	f := File{}

	assert.True(t, f.IsIncomplete())
}

func TestFile_Extension(t *testing.T) {
	type args struct {
		name      string
		sizeBytes int64
	}

	type testCase struct {
		name string
		args args
		want string
	}

	testCases := []testCase{
		{
			name: "with extension",
			args: args{
				name:      "content/image/picture.jpg",
				sizeBytes: 1024,
			},
			want: ".jpg",
		},
		{
			name: "double extension takes last suffix",
			args: args{
				name:      "content/archive/video.mp4.backup",
				sizeBytes: 1024,
			},
			want: ".backup",
		},
		{
			name: "no extension",
			args: args{
				name:      "content/image/picture",
				sizeBytes: 1024,
			},
			want: "",
		},
		{
			name: "dotfile extension behavior",
			args: args{
				name:      "content/image/.gitignore",
				sizeBytes: 1024,
			},
			want: ".gitignore",
		},
		{
			name: "dot is last character",
			args: args{
				name:      ".",
				sizeBytes: 1024,
			},
			want: ".",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := New(tc.args.name, tc.args.sizeBytes)
			require.NoError(t, err)

			ext := f.Extension()

			assert.Equal(t, tc.want, ext)
		})
	}
}

func TestFile_Name(t *testing.T) {
	f, err := New("path/to/path/to/file", 1024)
	require.NoError(t, err)

	assert.Equal(t, "path/to/path/to/file", f.Name())
}

func TestFile_SizeBytes(t *testing.T) {
	f, err := New("file", 1024)
	require.NoError(t, err)

	assert.Equal(t, int64(1024), f.SizeBytes())
}

func TestNew(t *testing.T) {
	type args struct {
		fileName  string
		sizeBytes int64
	}

	type testCase struct {
		name    string
		args    args
		want    File
		wantErr bool
	}

	testCases := []testCase{
		{
			name: "valid file",
			args: args{
				fileName:  "content/image/picture.jpg",
				sizeBytes: 1024,
			},
			want: File{
				name:      "content/image/picture.jpg",
				sizeBytes: 1024,
			},
			wantErr: false,
		},
		{
			name: "invalid file name",
			args: args{
				fileName:  "/..../",
				sizeBytes: 1024,
			},
			want:    File{},
			wantErr: true,
		},
		{
			name: "invalid file size in bytes",
			args: args{
				fileName:  "content",
				sizeBytes: 0,
			},
			want:    File{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := New(tc.args.fileName, tc.args.sizeBytes)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, f)
			}
		})
	}
}

func TestValidateFile(t *testing.T) {
	type testCase struct {
		name     string
		fileName string
		size     int64
		wantErr  bool
	}

	testCases := []testCase{
		{
			name:     "valid file",
			fileName: "image",
			size:     100,
			wantErr:  false,
		},
		{
			name:     "invalid file name",
			fileName: "image/",
			size:     100,
			wantErr:  true,
		},
		{
			name:     "invalid file size",
			fileName: "image/.",
			size:     0,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateFile(tc.fileName, tc.size)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFileName(t *testing.T) {
	type testCase struct {
		name     string
		fileName string
		wantErr  bool
	}

	testCases := []testCase{
		{
			name:     "starts with slash",
			fileName: "/picture.jpg",
			wantErr:  true,
		},
		{
			name:     "fileName without extension",
			fileName: "/picture",
			wantErr:  true,
		},
		{
			name:     "fileName without base",
			fileName: "/.jpg",
			wantErr:  true,
		},
		{
			name:     "fileName without leading slash",
			fileName: "image.png",
			wantErr:  false,
		},
		{
			name:     "fileName without leading slash and base",
			fileName: ".png",
			wantErr:  false,
		},
		{
			name:     "fileName without leading slash and extension, but with base",
			fileName: "image",
			wantErr:  false,
		},
		{
			name:     "ends on slash",
			fileName: "image/",
			wantErr:  true,
		},
		{
			name:     "empty fileName",
			fileName: "",
			wantErr:  true,
		},
		{
			name:     "fileName with .",
			fileName: ".",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateFileName(tc.fileName)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	type testCase struct {
		name    string
		size    int64
		wantErr bool
	}

	testCases := []testCase{
		{
			name:    "positive file size",
			size:    1000,
			wantErr: false,
		},
		{
			name:    "negative file size",
			size:    -1,
			wantErr: true,
		},
		{
			name:    "zero file size",
			size:    0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateFileSize(tc.size)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
