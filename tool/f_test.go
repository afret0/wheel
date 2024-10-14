package tool

import (
	"testing"
)

func Test_f(t *testing.T) {
	a := "world"
	s := f("Hello, {a}!", a)

	if s != "Hello, world!" {
		t.Errorf("Expected: Hello, world!, got: %s", s)
	}
}
