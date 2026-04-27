package content

import "errors"

var (
	ErrInvalidFormat = errors.New("некорректный формат файла")
	ErrEmptyFileSize = errors.New("пустой размер файла")
	ErrTooLargeFile  = errors.New("слишком большой файл")
	ErrEmptyQuizID   = errors.New("указан несуществующий квиз")
)
