package passhash

import "strings"

func normalizeHash(hash string) string {
	return strings.TrimSpace(hash)
}
