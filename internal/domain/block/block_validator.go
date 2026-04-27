package block

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"gitflic.ru/lms/internal/domain/element"
	"gitflic.ru/lms/internal/domain/shared/duplicate"
	"github.com/google/uuid"
)

func validateTitle(title string) error {
	if err := validateNotEmptyStr(title); err != nil {
		return err
	}

	if err := validateLimitStr(title, titleCharsLimit); err != nil {
		return err
	}

	return nil
}

func validateElements(elements []*element.Element) error {
	if err := validateElementsLimit(elements, elementsLimit); err != nil {
		return err
	}

	if err := validateElementsNil(elements); err != nil {
		return err
	}

	if err := validateElementsDuplication(elements); err != nil {
		return err
	}

	return nil
}

func validateNotEmptyStr(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w, детали: заголовок блока должен содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateLimitStr(s string, limit int) error {
	if utf8.RuneCountInString(s) > limit {
		return fmt.Errorf("%w, детали: заголовок должен содержать не более %d символов", ErrInvalid, limit)
	}

	return nil
}

func validateElementsDuplication(elements []*element.Element) error {
	ids := make([]uuid.UUID, 0, len(elements))

	for i := range elements {
		ids = append(ids, elements[i].ID())
	}

	if has := duplicate.FindUUID(ids); has {
		return fmt.Errorf("%w, детали: блок не должен содержать дубликаты элементов (разрешаются копии)", ErrInvalid)
	}

	return nil
}

func validateElementsNil(elements []*element.Element) error {
	for i := range elements {
		if elements[i] == nil {
			return fmt.Errorf("%w, детали: блок не должен содержать не существующие элементы", ErrInvalid)
		}
	}

	return nil
}

func validateElementsLimit(elements []*element.Element, limit int) error {
	if len(elements) > limit {
		return fmt.Errorf("%w, детали: количество элементов в блоке не должно превышать %d штук", ErrInvalid, limit)
	}

	return nil
}
