package media

type Type string

const (
	TypeSlides   Type = "slides"
	TypeDocument Type = "document"
	TypeVideo    Type = "video"
	TypeStatic   Type = "static"
	TypeImage    Type = "image"
	TypeAudio    Type = "audio"
)

func (t Type) IsValid() bool {
	switch t {
	case TypeSlides, TypeDocument, TypeVideo, TypeStatic, TypeImage, TypeAudio:
		return true
	default:
		return false
	}
}

func (t Type) Title() string {
	switch t {
	case TypeSlides:
		return "слайды"
	case TypeDocument:
		return "документ"
	case TypeVideo:
		return "видео"
	case TypeStatic:
		return "статический файл"
	case TypeImage:
		return "изображение"
	case TypeAudio:
		return "аудио"
	default:
		panic("unknown value: " + string(t))
	}
}

func (t Type) AllowedExtensions() []string {
	switch t {
	case TypeSlides:
		return []string{".pptx"}
	case TypeDocument:
		return []string{".pdf"}
	case TypeVideo:
		return []string{".mp4", ".webm"}
	case TypeStatic:
		return []string{".pdf", ".docx", ".xlsx", ".pptx", ".zip"}
	case TypeImage:
		return []string{".jpg", ".jpeg", ".webp"}
	case TypeAudio:
		return []string{".mp3", ".ogg", ".wav"}
	default:
		panic("unknown value: " + string(t))
	}
}

func (t Type) MaxSizeInBytes() int64 {
	switch t {
	case TypeSlides:
		return 500 * 1024 * 1024
	case TypeDocument:
		return 700 * 1024 * 1024
	case TypeVideo:
		return 500 * 1024 * 1024
	case TypeStatic:
		return 700 * 1024 * 1024
	case TypeImage:
		return 10 * 1024 * 1024
	case TypeAudio:
		return 50 * 1024 * 1024
	default:
		panic("unknown value: " + string(t))
	}
}

func (t Type) String() string {
	return string(t)
}
