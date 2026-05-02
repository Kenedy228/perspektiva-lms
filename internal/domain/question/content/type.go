package content

type Type string

const (
	TypeText  Type = "text"
	TypeImage Type = "image"
	TypeAudio Type = "audio"
)

func (c Type) Title() string {
	switch c {
	case TypeText:
		return "текст"
	case TypeImage:
		return "изображение"
	case TypeAudio:
		return "аудио"
	default:
		return ""
	}
}

func (c Type) String() string {
	return string(c)
}
