package httpapi

import "testing"

func TestClosureRequiresEvidenceAndNoOutstandingWork(t *testing.T) {
	for _, tc := range []struct {
		missing, reports, reviews, work int
		want                            bool
	}{{0, 1, 1, 0, true}, {1, 1, 1, 0, false}, {0, 0, 1, 0, false}, {0, 1, 0, 0, false}, {0, 1, 1, 1, false}} {
		if got := closureReady(tc.missing, tc.reports, tc.reviews, tc.work); got != tc.want {
			t.Errorf("closureReady(%+v)=%v", tc, got)
		}
	}
}
func TestContactValidation(t *testing.T) {
	for _, phone := range []string{"", "123", "+0123456789", "+919876543210 extra"} {
		if phonePattern.MatchString(phone) {
			t.Errorf("accepted invalid phone %q", phone)
		}
	}
	if !phonePattern.MatchString("+919876543210") {
		t.Fatal("valid E.164 rejected")
	}
}
