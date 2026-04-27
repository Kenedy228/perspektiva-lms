package content

import "gitflic.ru/lms/internal/domain/element"

type SlidesContent struct {
	key       string
	sizeBytes int64
}

func NewSlidesContent(key string, sizeBytes int64) (SlidesContent, error) {
	if err := validateSlidersFile(key, sizeBytes); err != nil {
		return SlidesContent{}, err
	}

	return SlidesContent{
		key:       key,
		sizeBytes: sizeBytes,
	}, nil
}

func (c SlidesContent) Key() string {
	return c.key
}

func (c SlidesContent) SizeBytes() int64 {
	return c.sizeBytes
}

func (c SlidesContent) Type() element.Type {
	return element.TypeSlides
}

func (c SlidesContent) IsInteractive() bool {
	return false
}

func (c SlidesContent) Clone() element.Content {
	return SlidesContent{
		key:       c.key,
		sizeBytes: c.sizeBytes,
	}
}
