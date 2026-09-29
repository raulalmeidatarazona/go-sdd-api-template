package domain

import "testing"

func TestNormalizeName(t *testing.T) {
	got, err := NormalizeName("  hello  ")
	if err != nil || got != "hello" {
		t.Fatalf("got %q, %v", got, err)
	}
	for _, input := range []string{"", "  ", string(make([]byte, 121))} {
		if _, err := NormalizeName(input); err == nil {
			t.Fatalf("expected invalid name for %q", input)
		}
	}
}
