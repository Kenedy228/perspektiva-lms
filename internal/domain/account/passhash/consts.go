package passhash

import "regexp"

var (
	bcryptFormat = regexp.MustCompile(`^\$2[ayb]\$\d{2}\$[./A-Za-z0-9]{53}$`)
)
