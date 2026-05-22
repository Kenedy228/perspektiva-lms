package passhash

// Hash объект-значение над хешем пароля аккаунта.
type Hash struct {
	hash string
}

// New создает новый объект Hash с нормализацией и валидацией хеша пароля.
func New(hash string) (Hash, error) {
	hash = normalizeHash(hash)

	if err := validateHash(hash); err != nil {
		return Hash{}, err
	}

	return Hash{
		hash: hash,
	}, nil
}

// Value возвращает хеш пароля.
func (ph Hash) Value() string {
	return ph.hash
}

// IsZero сигнализирует, был ли проинициализирован объект.
func (ph Hash) IsZero() bool {
	return ph.hash == ""
}
