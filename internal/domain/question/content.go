package question

type Content struct {
	cType ContentType
	value string
}

func NewContent(cType ContentType, value string) (Content, error) {
	if err := validateContentType(cType); err != nil {
		return Content{}, err
	}

	if err := validateValue(cType, value); err != nil {
		return Content{}, err
	}

	return Content{
		cType: cType,
		value: value,
	}, nil
}

func (o Content) ContentType() ContentType {
	return o.cType
}

func (o Content) Value() string {
	return o.value
}

func (o Content) IsText() bool {
	return o.cType == ContentTypeText
}

func (o Content) Equal(other Content) bool {
	return o.cType == other.cType && o.value == other.value
}
