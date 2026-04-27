package organization

import "regexp"

var (
	individualRegexp = regexp.MustCompile(`^\d{12}$`)
	companyRegexp    = regexp.MustCompile(`^\d{10}$`)
)
