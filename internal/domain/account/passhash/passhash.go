package passhash

type Hash struct {
	hash string
}

func New(hash string) (Hash, error) {
	hash = normalizeHash(hash)

	if err := validateHash(hash); err != nil {
		return Hash{}, err
	}

	return Hash{
		hash: hash,
	}, nil
}

func (ph Hash) Hash() string {
	return ph.hash
}
