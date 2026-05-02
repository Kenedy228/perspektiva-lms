package inn

type INN struct {
	code    string
	innType Type
}

func New(code string, innType Type) (INN, error) {
	code = normalizeCode(code)

	if err := validateChecksum(code, innType); err != nil {
		return INN{}, err
	}

	return INN{
		code:    code,
		innType: innType,
	}, nil
}

func (inn INN) Code() string {
	return inn.code
}

func (inn INN) Type() Type {
	return inn.innType
}
