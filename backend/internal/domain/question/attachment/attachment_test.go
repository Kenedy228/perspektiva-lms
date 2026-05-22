package attachment

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/shared/file"
	media2 "gitflic.ru/lms/backend/internal/domain/shared/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustMedia(t *testing.T, mt media2.Type, key string, size int64) media2.Media {
	t.Helper()

	f, err := file.New(key, size)
	if err != nil {
		t.Fatalf("file.New() error = %v", err)
	}

	m, err := media2.New(mt, f)
	if err != nil {
		t.Fatalf("media.New() error = %v", err)
	}

	return m
}

func TestAttachment_Media(t *testing.T) {
	m := mustMedia(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024)
	a, err := New(m)
	require.NoError(t, err)

	assert.Equal(t, m, a.Media())
}

func TestNew(t *testing.T) {
	type args struct {
		m media2.Media
	}
	tests := []struct {
		name    string
		args    args
		want    Attachment
		wantErr bool
	}{
		{
			name: "valid image attachment",
			args: args{
				m: mustMedia(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024),
			},
			want: Attachment{
				media: mustMedia(t, media2.TypeImage, "content/questions/image/picture.jpg", 1024),
			},
			wantErr: false,
		},
		{
			name: "valid audio attachment",
			args: args{
				m: mustMedia(t, media2.TypeAudio, "content/questions/audio/task.mp3", 2048),
			},
			want: Attachment{
				media: mustMedia(t, media2.TypeAudio, "content/questions/audio/task.mp3", 2048),
			},
			wantErr: false,
		},
		{
			name: "invalid video attachment",
			args: args{
				m: mustMedia(t, media2.TypeVideo, "content/questions/video/lesson.mp4", 10*1024*1024),
			},
			want:    Attachment{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.m)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_validateMediaComplete(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		m := mustMedia(t, media2.TypeSlides, "content/questions/slides/topic.pptx", 1024)

		assert.NoError(t, validateMediaComplete(m))
	})

	t.Run("incomplete", func(t *testing.T) {
		m := media2.Media{}

		assert.Error(t, validateMediaComplete(m))
	})
}

func Test_validateAllowedType(t *testing.T) {
	type args struct {
		t media2.Type
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "image valid",
			args:    args{t: media2.TypeImage},
			wantErr: false,
		},
		{
			name:    "audio valid",
			args:    args{t: media2.TypeAudio},
			wantErr: false,
		},
		{
			name:    "video invalid",
			args:    args{t: media2.TypeVideo},
			wantErr: true,
		},
		{
			name:    "document invalid",
			args:    args{t: media2.TypeDocument},
			wantErr: true,
		},
		{
			name:    "slides invalid",
			args:    args{t: media2.TypeSlides},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAllowedType(tt.args.t)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
