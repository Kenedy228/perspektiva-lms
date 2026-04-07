package permission

import "testing"

func BenchmarkReviewActions(b *testing.B) {
	given := []Action{ActionRead, ActionWrite}

	for b.Loop() {
		reviewActions(given)
	}
}
