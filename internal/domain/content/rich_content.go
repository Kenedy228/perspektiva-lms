package content

type RichContent struct {
	cType ContentType
	value string
}

func New(params Params) (RichContent, error) {
	if err := validateContentType(params.Type); err != nil {
		return RichContent{}, err
	}

	return RichContent{
		cType: params.Type,
		value: params.Value,
	}, nil
}

func (c RichContent) ContentType() ContentType {
	return c.cType
}

func (c RichContent) Value() string {
	return c.value
}

func (c RichContent) Equal(other RichContent) bool {
	return c.cType == other.cType && c.value == other.value
}
