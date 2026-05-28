package downloadfile

import (
	"fmt"

	element2 "gitflic.ru/lms/backend/internal/domain/course/element"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
)

var allowedExtensions = map[string]struct{}{
	".docx": {},
	".pptx": {},
	".xlsx": {},
	".pdf":  {},
	".txt":  {},
	".md":   {},
	".csv":  {},
	".png":  {},
	".jpg":  {},
	".jpeg": {},
	".gif":  {},
	".webp": {},
	".bmp":  {},
	".svg":  {},
}

type Content struct {
	f file.File
}

func New(f file.File) (Content, error) {
	if _, ok := allowedExtensions[f.Extension()]; !ok {
		return Content{}, fmt.Errorf("%w: invalid value", element2.ErrInvalid)
	}
	return Content{f: f}, nil
}

func (c Content) File() file.File {
	return c.f
}

func (c Content) ContentType() element2.ContentType {
	return element2.ContentTypeDownloadFile
}

func (c Content) IsInteractive() bool {
	return false
}

func (c Content) Clone() element2.Content {
	return c
}
