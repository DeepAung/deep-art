package utils_test

import (
	"testing"

	"github.com/DeepAung/deep-art/pkg/utils"
)

func TestJoin(t *testing.T) {
	tests := []struct {
		inputs []string
		expect string
	}{
		{[]string{"a", "b", "c"}, "a/b/c"},
		{[]string{"a", "/b/", "/c/"}, "a/b/c/"},
		{[]string{"a", "///b/", "/c/"}, "a/b/c/"},
	}

	for _, tt := range tests {
		got := utils.Join(tt.inputs...)
		if got != tt.expect {
			t.Fatalf("expect=%q, got=%q", tt.expect, got)
		}
	}
}
