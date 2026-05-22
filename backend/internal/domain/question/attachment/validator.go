package attachment

import (
	"errors"
	"fmt"
	"slices"

	media2 "gitflic.ru/lms/backend/internal/domain/shared/media"
)

func validateMedia(m media2.Media) error {
	if err := validateMediaComplete(m); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	if err := validateAllowedType(m.Type()); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	return nil
}

func validateMediaComplete(m media2.Media) error {
	if m.IsIncomplete() {
		return errors.New("invalid value")
	}

	return nil
}

func validateAllowedType(t media2.Type) error {
	if !slices.Contains(AllowedMediaTypes, t) {
		return fmt.Errorf("invalid value (%q)", t.Title())
	}

	return nil
}
