package dob

import "time"

// normalize приводит time.Time к UTC и обнуляет время, оставляя только дату.
func normalize(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)
}
