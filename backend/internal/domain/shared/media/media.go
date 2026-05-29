package media

import (
	"errors"

	"gitflic.ru/lms/backend/internal/domain/shared/file"
)

var (
	ErrInvalid = errors.New("некорректный медиафайл")
)

type Media struct {
	mType Type
	file  file.File
}

func New(mType Type, file file.File) (Media, error) {
	if err := validateType(mType); err != nil {
		return Media{}, err
	}

	if err := validateFileForType(mType, file); err != nil {
		return Media{}, err
	}

	return Media{
		mType: mType,
		file:  file,
	}, nil
}

func (m Media) File() file.File {
	return m.file
}

func (m Media) Type() Type {
	return m.mType
}

func (m Media) IsIncomplete() bool {
	return !m.mType.IsValid() || m.file.IsIncomplete()
}
