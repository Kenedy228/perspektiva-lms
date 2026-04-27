package snils

func validate(value string) error {
	if err := validateLength(value); err != nil {
		return err
	}

	if err := validateContent(value); err != nil {
		return err
	}

	if err := validateChecksum(value); err != nil {
		return err
	}

	return nil
}
