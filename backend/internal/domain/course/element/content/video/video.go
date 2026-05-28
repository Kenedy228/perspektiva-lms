package video

import (
	element2 "gitflic.ru/lms/backend/internal/domain/course/element"
	"gitflic.ru/lms/backend/internal/domain/shared/file"
)

type Content struct {
	f file.File
}

func New(f file.File) (Content, error) {
	if err := validateFile(f); err != nil {
		return Content{}, err
	}

	return Content{
		f: f,
	}, nil
}

func (c Content) File() file.File {
	return c.f
}

func (c Content) ContentType() element2.ContentType {
	return element2.ContentTypeLectureMaterial
}

func (c Content) IsInteractive() bool {
	return false
}

func (c Content) Clone() element2.Content {
	return c
}
