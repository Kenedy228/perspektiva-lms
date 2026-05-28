package element

type ContentType string

const (
	ContentTypeTest            ContentType = "test"
	ContentTypeLectureMaterial ContentType = "lecture_material"
	ContentTypeDownloadFile    ContentType = "download_file"
)

func (t ContentType) Title() string {
	switch t {
	case ContentTypeTest:
		return "тест"
	case ContentTypeLectureMaterial:
		return "лекционный материал"
	case ContentTypeDownloadFile:
		return "файл для скачивания"
	default:
		return ""
	}
}

func (t ContentType) String() string {
	return string(t)
}
