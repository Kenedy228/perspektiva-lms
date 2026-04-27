package snils

import "fmt"

type Snils struct {
	value string
}

func New(value string) (Snils, error) {
	nValue := normalize(value)

	if err := validate(nValue); err != nil {
		return Snils{}, err
	}

	return Snils{
		value: nValue,
	}, nil
}

func (s Snils) Value() string {
	return s.value
}

func (s Snils) Formatted() string {
	if len(s.value) == 0 {
		return ""
	}

	return fmt.Sprintf("%s-%s-%s %s", s.value[0:3], s.value[3:6], s.value[6:9], s.value[9:11])
}
