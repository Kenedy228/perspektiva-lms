package file

import (
	"path"
	"strings"
)

type File struct {
	key       string
	sizeBytes int64
}

func New(key string, sizeBytes int64) (File, error) {
	if err := validateFile(key, sizeBytes); err != nil {
		return File{}, err
	}

	return File{
		key:       key,
		sizeBytes: sizeBytes,
	}, nil
}

func (f File) Key() string {
	return f.key
}

func (f File) SizeBytes() int64 {
	return f.sizeBytes
}

func (f File) Extension() string {
	filename := path.Base(f.key)
	return strings.ToLower(path.Ext(filename))
}
