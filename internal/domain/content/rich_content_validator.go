package content

import "errors"

var (
	ErrInvalidContentType = errors.New("invalid content type")
)

func validateContentType(cType ContentType) error {
	switch cType {
	case ContentTypeText, ContentTypeImage:
		return nil
	default:
		return ErrInvalidContentType
	}
}
