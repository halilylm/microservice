package test

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

// AssertDeepEqual compares whether two values are semantically equal.
func AssertDeepEqual(t *testing.T, wanted, got any) {
	if diff := cmp.Diff(wanted, got); diff != "" {
		t.Errorf("unexpected diff (-want, +have):\n%s", diff)
	}
}
