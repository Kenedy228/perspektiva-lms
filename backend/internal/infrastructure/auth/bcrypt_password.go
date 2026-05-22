package auth

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"golang.org/x/crypto/bcrypt"
)

// DefaultBcryptCost is the bcrypt work factor used for newly generated hashes.
const DefaultBcryptCost = 12

// BcryptPasswordComparer compares plaintext passwords with bcrypt hashes.
type BcryptPasswordComparer struct {
	cost int
}

// NewBcryptPasswordComparer creates a bcrypt password comparer using the
// default work factor for generated hashes.
func NewBcryptPasswordComparer() BcryptPasswordComparer {
	return BcryptPasswordComparer{
		cost: DefaultBcryptCost,
	}
}

// NewBcryptPasswordComparerWithCost creates a bcrypt password comparer with a
// custom work factor for generated hashes.
func NewBcryptPasswordComparerWithCost(cost int) (BcryptPasswordComparer, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return BcryptPasswordComparer{}, fmt.Errorf("invalid bcrypt cost: %d", cost)
	}

	return BcryptPasswordComparer{
		cost: cost,
	}, nil
}

// Hash creates a bcrypt hash for a plaintext password.
func (c BcryptPasswordComparer) Hash(plain string) (passhash.Hash, error) {
	cost := c.cost
	if cost == 0 {
		cost = DefaultBcryptCost
	}

	raw, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return passhash.Hash{}, fmt.Errorf("create bcrypt password hash: %w", err)
	}

	h, err := passhash.New(string(raw))
	if err != nil {
		return passhash.Hash{}, fmt.Errorf("create password hash value: %w", err)
	}

	return h, nil
}

// Compare reports whether plain matches hash.
func (c BcryptPasswordComparer) Compare(hash passhash.Hash, plain string) bool {
	if hash.IsZero() {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(hash.Value()), []byte(plain)) == nil
}
