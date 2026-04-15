package content

type RichContent struct {
	contentType ContentType
	value       string
}

func New(cType ContentType, value string) (RichContent, error) {
	validator := validator{}
	if err := validator.validateContentType(cType); err != nil {
		return RichContent{}, err
	}

	return RichContent{
		contentType: cType,
		value:       value,
	}, nil
}

func (c RichContent) ContentType() ContentType {
	return c.contentType
}

func (c RichContent) Value() string {
	return c.value
}

func (c RichContent) Equal(other RichContent) bool {
	return c.contentType == other.contentType && c.value == other.value
}
