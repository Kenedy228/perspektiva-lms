package matching

func mapPairs(rawPairs []PairParams) ([]Pair, error) {
	pairs := make([]Pair, 0, len(rawPairs))

	for i := range rawPairs {
		p, err := NewPair(rawPairs[i])
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, p)
	}

	return pairs, nil
}
