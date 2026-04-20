package typed

func mapBlanks(rawBlanks []BlankParams) ([]Blank, error) {
	blanks := make([]Blank, 0, len(rawBlanks))
	for i := range rawBlanks {
		b, err := NewBlank(rawBlanks[i])
		if err != nil {
			return nil, err
		}

		blanks = append(blanks, b)
	}

	return blanks, nil
}
