package pic

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		s1 := fmt.Sprintf("%v", a)
		s2 := fmt.Sprintf("%v", b)
		if s1 == s2 {
			t.Fatalf("Type mismatch: %T != %T", a, b)
		} else {
			t.Fatalf("'%s' != '%s'", s1, s2)
		}
	}
}
