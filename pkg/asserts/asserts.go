package asserts

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, name string, got, expect any) {
	t.Helper()

	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("assert equal %q, expect=%v, got=%v", name, expect, got)
	}
}

func NotEqual(t *testing.T, name string, got, expect any) {
	t.Helper()

	if reflect.DeepEqual(got, expect) {
		t.Fatalf("assert not equal %q, expect=%v, got=%v", name, expect, got)
	}
}

func EqualError(t *testing.T, got, expect error) {
	t.Helper()

	if got == nil && expect == nil { // case (nil, nil)
		return
	}

	if (got == nil) != (expect == nil) { // case (nil, notnil) or (notnil, nil)
		t.Fatalf("assert equal error, expect=%v, got=%v", expect, got)
	}

	if got.Error() != expect.Error() { // case (notnil, notnil)
		t.Fatalf("assert equal error, expect=%v, got=%v", expect, got)
	}
}
