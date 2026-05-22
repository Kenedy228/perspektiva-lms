package name

import (
	"github.com/Kenedy228/fnsfio"
)

func normalizeFirstName(firstName string) string {
	return fnsfio.NormalizeFirstName(firstName)
}

func normalizeLastName(lastName string) string {
	return fnsfio.NormalizeLastName(lastName)
}

func normalizeMiddleName(middleName string) string {
	return fnsfio.NormalizeMiddleName(middleName)
}
