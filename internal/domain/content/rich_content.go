package content

type RichContent struct {
	contentType ContentType
	value       string
}

func New(params Params) (RichContent, error) {
	if err := validateContentType(params.Type); err != nil {
		return RichContent{}, err
	}

	return RichContent{
		contentType: params.Type,
		value:       params.Value,
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
