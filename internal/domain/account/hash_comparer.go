package account

type PasswordComparer interface {
	Compare(hash, plain string) bool
}
