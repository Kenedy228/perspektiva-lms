package file

import (
	"errors"
	"path"
)

var (
	ErrInvalid = errors.New("некорректный файл")
)

type File struct {
	name      string
	sizeBytes int64
}

func New(name string, sizeBytes int64) (File, error) {
	if err := validateFile(name, sizeBytes); err != nil {
		return File{}, err
	}

	return File{
		name:      name,
		sizeBytes: sizeBytes,
	}, nil
}

func (f File) Name() string {
	return f.name
}

func (f File) SizeBytes() int64 {
	return f.sizeBytes
}

func (f File) Extension() string {
	return path.Ext(f.name)
}

func (f File) IsIncomplete() bool {
	return f.name == "" || f.sizeBytes == 0
}
