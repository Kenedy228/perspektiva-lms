package element

type Type string

const (
	TypeText     Type = "text"
	TypeSlides   Type = "slides"
	TypeDocument Type = "document"
	TypeVideo    Type = "video"
	TypeQuiz     Type = "quiz"
)

func (t Type) Title() string {
	switch t {
	case TypeText:
		return "текст"
	case TypeSlides:
		return "слайды"
	case TypeDocument:
		return "документ"
	case TypeVideo:
		return "видео"
	case TypeQuiz:
		return "тест"
	default:
		return ""
	}
}

func (t Type) String() string {
	return string(t)
}
