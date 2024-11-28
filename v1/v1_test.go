package v1

import (
	"fmt"
	"testing"
)

func TestMatchString(t *testing.T) {
	for i, testcase := range testcases {
		fmt.Printf("--- testcase %d ---\n", i)
		fmt.Println(":regex: ", testcase.reg)
		fmt.Println(":str: ", testcase.s)
		matched, err := MatchString(testcase.reg, testcase.s)
		fmt.Println(":result: ", matched)
		if err != nil {
			t.Fatalf("failed at testcase %d: %v", i, err)
		}
		if matched != testcase.matched {
			t.Fatalf("failed at testcase %d, want %v, got %v", i, testcase.matched, matched)
		}
	}
}

type testcase struct {
	reg     string
	s       string
	matched bool
}

var testcases = []testcase{
	// single negative lookahead
	{reg: "foo(?!bar)", s: "foobar foobar", matched: false},
	{reg: "foo(?!bar)", s: "foobaz foobak", matched: true},
	{reg: "foo(?!bar)", s: "foobar foobak", matched: true},
	{reg: "foo(?!bar)", s: "foobaz foobar", matched: true},

	// single positive lookahead
	{reg: "foo(?=bar)", s: "fooboo foobaw", matched: false},
	{reg: "foo(?=bar)", s: "fooboo foobar", matched: true},
	{reg: "foo(?=bar)", s: "foobar fooboo", matched: true},
	{reg: "foo(?=bar)", s: "foobar foobar", matched: true},

	// multiple negative lookaheads
	{reg: "a(?!b)c(?!d)", s: "ac", matched: true},
	{reg: "a(?!b)c(?!d)", s: "abcac", matched: true},
	{reg: "a(?!b)c(?!d)", s: "ad", matched: false},
	{reg: "a(?!b)c(?!d)", s: "dc", matched: false},
	{reg: "a(?!b)c(?!d)", s: "bbb", matched: false},

	// multiple positive lookaheads
	{reg: "a(?=[0-9])1(?=[a-z])e", s: "a1e", matched: true},

	// mixed expr with negative and positive lookaheads
	{reg: "a(?![0-9])c(?=[a-z])", s: "a1cd", matched: false},
	{reg: "a(?![0-9])c(?=[a-z])", s: "acd", matched: true},
	{reg: "a(?=[0-9])3(?![a-z])[0-9]", s: "a34", matched: true},

	// custom testcases
	{reg: "^(?!OK$).*", s: "OK", matched: false},
	{reg: "^(?!OK$).*", s: "NOTOK", matched: true},
	{reg: "^(?!OK$).*", s: "WARNING", matched: true},
	{reg: "^(?!OK$).*", s: "ERROR", matched: true},
}
