package option

type ContentType string

const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
	ContentTypeAudio ContentType = "audio"
)

func (c ContentType) IsValid() bool {
	switch c {
	case ContentTypeText, ContentTypeImage, ContentTypeAudio:
		return true
	default:
		return false
	}
}

func (c ContentType) String() string {
	return string(c)
}

func (c ContentType) Equal(other ContentType) bool {
	return c.String() == other.String()
}
