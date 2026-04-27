package content

import "gitflic.ru/lms/internal/domain/element"

type AttachmentContent struct {
	key       string
	sizeBytes int64
}

func NewAttachmentContent(key string, sizeBytes int64) (AttachmentContent, error) {
	if err := validateAttachmentFile(key, sizeBytes); err != nil {
		return AttachmentContent{}, err
	}

	return AttachmentContent{
		key:       key,
		sizeBytes: sizeBytes,
	}, nil
}

func (c AttachmentContent) Key() string {
	return c.key
}

func (c AttachmentContent) SizeBytes() int64 {
	return c.sizeBytes
}

func (c AttachmentContent) Type() element.Type {
	return element.TypeAttachment
}

func (c AttachmentContent) IsInteractive() bool {
	return false
}

func (c AttachmentContent) Clone() element.Content {
	return AttachmentContent{
		key:       c.key,
		sizeBytes: c.sizeBytes,
	}
}
