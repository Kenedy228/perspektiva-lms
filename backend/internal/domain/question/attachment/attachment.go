package attachment

import (
	media2 "gitflic.ru/lms/backend/internal/domain/shared/media"
)

var (
	AllowedMediaTypes = []media2.Type{media2.TypeImage, media2.TypeAudio}
)

type Attachment struct {
	media media2.Media
}

func New(m media2.Media) (Attachment, error) {
	if err := validateMedia(m); err != nil {
		return Attachment{}, err
	}

	return Attachment{
		media: m,
	}, nil
}

func (a Attachment) Media() media2.Media {
	return a.media
}

func (a Attachment) IsZero() bool {
	return a.media.IsIncomplete()
}
