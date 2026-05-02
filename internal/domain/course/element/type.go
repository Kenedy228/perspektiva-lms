package element

type ContentType string

const (
	ContentTypeText       ContentType = "text"
	ContentTypeSlides     ContentType = "slides"
	ContentTypeAttachment ContentType = "document"
	ContentTypeVideo      ContentType = "video"
	ContentTypeQuiz       ContentType = "quiz"
)

func (t ContentType) Title() string {
	switch t {
	case ContentTypeText:
		return "текст"
	case ContentTypeSlides:
		return "слайды"
	case ContentTypeAttachment:
		return "вложение"
	case ContentTypeVideo:
		return "видео"
	case ContentTypeQuiz:
		return "тест"
	default:
		return ""
	}
}

func (t ContentType) String() string {
	return string(t)
}
