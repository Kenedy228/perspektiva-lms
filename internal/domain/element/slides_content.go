package element

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
