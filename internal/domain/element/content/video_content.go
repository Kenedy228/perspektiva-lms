package content

import "gitflic.ru/lms/internal/domain/element"

type VideoContent struct {
	key       string
	sizeBytes int64
}

func NewVideoContent(key string, sizeBytes int64) (VideoContent, error) {
	if err := validateVideoFile(key, sizeBytes); err != nil {
		return VideoContent{}, err
	}

	return VideoContent{
		key:       key,
		sizeBytes: sizeBytes,
	}, nil
}

func (c VideoContent) Key() string {
	return c.key
}

func (c VideoContent) SizeBytes() int64 {
	return c.sizeBytes
}

func (c VideoContent) Type() element.Type {
	return element.TypeVideo
}

func (c VideoContent) IsInteractive() bool {
	return false
}

func (c VideoContent) Clone() element.Content {
	return VideoContent{
		key:       c.key,
		sizeBytes: c.sizeBytes,
	}
}
