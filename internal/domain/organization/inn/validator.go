package inn

import (
	"fmt"
)

func validateChecksum(value string, t Type) error {
	digits, ok := getDigits(value)
	if !ok {
		return fmt.Errorf("%w, детали: код ИНН содержит недопустимые символы", ErrInvalid)
	}

	switch t {
	case TypeOrganization:
		if len(digits) != t.CodeLength() {
			return fmt.Errorf("%w, детали: длина кода ИНН не совпадает с длиной кода для ИНН юр. лица", ErrInvalid)
		}

		return validateOrganizationChecksum(digits)
	case TypeIP:
		if len(digits) != t.CodeLength() {
			return fmt.Errorf("%w, детали: длина кода ИНН не совпадает с длиной кода для ИНН ИП", ErrInvalid)
		}

		return validateIPChecksum(digits)
	case TypePhysical:
		if len(digits) != t.CodeLength() {
			return fmt.Errorf("%w, детали: длина кода ИНН не совпадает с длиной кода для ИНН физ. лица", ErrInvalid)
		}

		return validatePhysicalChecksum(digits)
	default:
		return fmt.Errorf("%w, детали: указанный тип ИНН некорректный", ErrInvalid)
	}
}

func validateOrganizationChecksum(digits []int) error {
	sum, ok := calculateSum(digits[:9], TypeOrganization.Coefficients()[0])
	if !ok {
		return fmt.Errorf("%w, детали: ошибка подсчета чексуммы", ErrInvalid)
	}
	remainder := calcRemainder(sum)

	if remainder != digits[9] {
		return fmt.Errorf("%w, детали: код ИНН не прошел валидацию чексуммы", ErrInvalid)
	}

	return nil
}

func validateIPChecksum(digits []int) error {
	firstPartSum, ok := calculateSum(digits[:10], TypeIP.Coefficients()[0])
	if !ok {
		return fmt.Errorf("%w, детали: ошибка подсчета чексуммы", ErrInvalid)
	}
	firstPartRemainder := calcRemainder(firstPartSum)

	if firstPartRemainder != digits[10] {
		return fmt.Errorf("%w, детали: код ИНН не прошел валидацию чексуммы", ErrInvalid)
	}

	secondPartSum, ok := calculateSum(digits[:11], TypeIP.Coefficients()[1])
	if !ok {
		return fmt.Errorf("%w, детали: ошибка подсчета чексуммы", ErrInvalid)
	}
	secondPartRemainder := calcRemainder(secondPartSum)

	if secondPartRemainder != digits[11] {
		return fmt.Errorf("%w, детали: код ИНН не прошел валидацию чексуммы", ErrInvalid)
	}

	return nil
}

func validatePhysicalChecksum(digits []int) error {
	firstPartSum, ok := calculateSum(digits[:10], TypePhysical.Coefficients()[0])
	if !ok {
		return fmt.Errorf("%w, детали: ошибка подсчета чексуммы", ErrInvalid)
	}
	firstPartRemainder := calcRemainder(firstPartSum)

	if firstPartRemainder != digits[10] {
		return fmt.Errorf("%w, детали: код ИНН не прошел валидацию чексуммы", ErrInvalid)
	}

	secondPartSum, ok := calculateSum(digits[:11], TypePhysical.Coefficients()[1])
	if !ok {
		return fmt.Errorf("%w, детали: ошибка подсчета чексуммы", ErrInvalid)
	}
	secondPartRemainder := calcRemainder(secondPartSum)

	if secondPartRemainder != digits[11] {
		return fmt.Errorf("%w, детали: код ИНН не прошел валидацию чексуммы", ErrInvalid)
	}

	return nil
}

func calculateSum(digits []int, coefficients []int) (int, bool) {
	if len(digits) != len(coefficients) {
		return 0, false
	}

	sum := 0
	for i := range digits {
		sum += digits[i] * coefficients[i]
	}

	return sum, true
}

func calcRemainder(value int) int {
	remainder := value % 11
	if remainder > 9 {
		remainder %= 10
	}

	return remainder
}
