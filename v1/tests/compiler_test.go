package v1_test

import (
	"testing"

	v1 "github.com/debidong/re3/v1"
)

func TestCompile(t *testing.T) {
	expr := "this is a(?!regex)expression(?![1-3].*)with(?=[a-z]{1,10})multiple lookaheads."
	_, err := v1.Compile(expr)
	if err != nil {
		t.Fatal(err)
	}
}
