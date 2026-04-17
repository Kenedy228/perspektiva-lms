package option

type ContentOption struct {
	cType ContentType
	value string
}

func NewContentOption(cType ContentType, value string) (ContentOption, error) {
	if err := validateContentType(cType); err != nil {
		return ContentOption{}, err
	}

	if err := validateValue(cType, value); err != nil {
		return ContentOption{}, err
	}

	return ContentOption{
		cType: cType,
		value: value,
	}, nil
}

func (o ContentOption) ContentType() ContentType {
	return o.cType
}

func (o ContentOption) Value() string {
	return o.value
}

func (o ContentOption) Equal(other ContentOption) bool {
	return o.cType == other.cType && o.value == other.value
}
