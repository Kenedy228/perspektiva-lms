package content

import "errors"

var (
	ErrInvalidContentType = errors.New("invalid content type")
)

func validateContentType(cType ContentType) error {
	if !cType.IsValid() {
		return ErrInvalidContentType
	}

	return nil
}
