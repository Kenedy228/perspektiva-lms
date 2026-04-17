package content

type ContentType string

const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image_url"
)

func (c ContentType) IsValid() bool {
	switch c {
	case ContentTypeText, ContentTypeImage:
		return true
	default:
		return false
	}
}

func (c ContentType) String() string {
	return string(c)
}
