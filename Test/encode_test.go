package Test

import (
	"common/tools"
	"testing"
)

func TestEncode(t *testing.T) {
	_, pwd := tools.Encode("123456", nil)
	t.Error(pwd)
	if tools.Verify("123456", "eb59ce35a2e73ec455c5102857d448f76c4dd04dc1788b93cfb4900a07e14f8372112058d4ac2af18d610eed55c16513014797d4b44473d0ef78ce3d7912ac5eb0c00fd089eb7eb4b2eeb7f3e901712032c0288bd77e64b4c93ece79c85aeac8db99269d2ea8a37d87fe61d0a6a092d6b33cce9cae9fb3ed56e24b7891ce56fd", "gNBSkiAkweRT69qLlDiSBQjaoPK7yctE94Qag35OJrr9m5qR2vvhSWIjdS13maU4", nil) {
		t.Error("yes")
		// t.Errorf("%s\n", pwd)
	} else {
		t.Error("no")
	}
}
