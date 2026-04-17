package matching

import "gitflic.ru/lms/internal/domain/content"

func mapPairs(rawPairs map[string]content.RichContent) ([]Pair, []Option, error) {
	pairs := make([]Pair, 0, len(rawPairs))
	options := make([]Option, 0, len(rawPairs))

	for prompt, content := range rawPairs {
		option, err := NewOption(content)
		if err != nil {
			return nil, nil, err
		}
		options = append(options, option)

		pair, err := NewPair(prompt, option.id)
		if err != nil {
			return nil, nil, err
		}

		pairs = append(pairs, pair)
	}

	return pairs, options, nil
}
