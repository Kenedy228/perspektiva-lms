package question

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

func (c ContentType) Title() string {
	switch c {
	case ContentTypeText:
		return "текст"
	case ContentTypeImage:
		return "изображение"
	case ContentTypeAudio:
		return "аудио"
	default:
		return ""
	}
}

func (c ContentType) String() string {
	return string(c)
}
