package content

type Content struct {
	cType Type
	value string
}

func New(cType Type, value string) (Content, error) {
	if err := validateValue(value); err != nil {
		return Content{}, err
	}

	return Content{
		cType: cType,
		value: value,
	}, nil
}

func (o Content) Type() Type {
	return o.cType
}

func (o Content) Value() string {
	return o.value
}

func (o Content) IsText() bool {
	return o.cType == TypeText
}
