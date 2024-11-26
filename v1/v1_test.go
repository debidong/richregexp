package v1

import (
	"fmt"
	"regexp"
	"testing"
)

func TestSplitNegLookAhead(t *testing.T) {

	s := "ok computer not ok computer"
	re := regexp.MustCompile(`computer`)
	matches := re.FindAllStringIndex(s, -1)
	for _, match := range matches {
		fmt.Printf("%v from %v to %v \n", s[match[0]:match[1]], match[0], match[1])
	}
}

func TestMustSplitRegex(t *testing.T) {
	// reg := "foo(?!()bar) and (?=baz) something else"
	// _ = MustSplitRegex(reg)

	reg := "foo(?!barbazboo(?=okcomputer[0-9].*[a-z]{1,3}anotherNestedLookahead(?=notokcomputer)}}))"
	lookaheads, err := splitRegex(reg)
	if err != nil {
		t.Fatal(err)
	}
	for _, lookahead := range lookaheads {
		fmt.Println(reg[lookahead.idx[0]:lookahead.idx[1]])
	}
}

func TestMatchString(t *testing.T) {
	for i, testcase := range testcases {
		matched, err := MatchString(testcase.reg, testcase.s)
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

	// multiple negative lookahead

	{reg: "a(?!b)c(?!d)", s: "ac", matched: true}, // TODO: this fails
}
