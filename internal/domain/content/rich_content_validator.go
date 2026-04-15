package content

import "errors"

type validator struct{}

var (
	ErrInvalidContentType = errors.New("invalid content type")
)

func (v validator) validateContentType(cType ContentType) error {
	switch cType {
	case ContentTypeText, ContentTypeImage:
		return nil
	default:
		return ErrInvalidContentType
	}
}
