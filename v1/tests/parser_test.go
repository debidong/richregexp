package v1_test

import (
	v1 "debidong/re3/v1"
	"fmt"
	"testing"
)

func TestMatchString(t *testing.T) {
	for i, testcase := range testcases {
		fmt.Printf("--- testcase %d ---\n", i)
		fmt.Println(":regex: ", testcase.reg)
		fmt.Println(":str: ", testcase.s)
		regexp, err := v1.Compile(testcase.reg)
		if err != nil {
			t.Fatal(err)
		}
		matched := regexp.MatchString(testcase.s)
		fmt.Println(":result: ", matched)
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

	// examples of matching status codes
	// 1. `5xx`
	// Regexp: `^(5[0-9]{2})$`
	{reg: "^(5[0-9]{2})$", s: "499", matched: false},
	{reg: "^(5[0-9]{2})$", s: "500", matched: true},
	{reg: "^(5[0-9]{2})$", s: "501", matched: true},
	// 2. `warn|error|debug`
	// Regexp: `^(warn|error|debug)$`
	{reg: "^(warn|error|debug)$", s: "info", matched: false},
	{reg: "^(warn|error|debug)$", s: "warn", matched: true},
	{reg: "^(warn|error|debug)$", s: "error", matched: true},
	{reg: "^(warn|error|debug)$", s: "debug", matched: true},
	// 3. !0
	// Regexp: `^[^0].*$`
	{reg: "^[^0].*$", s: "0", matched: false},
	{reg: "^[^0].*$", s: "ERROR", matched: true},
	{reg: "^[^0].*$", s: "400", matched: true},
	{reg: "^[^0].*$", s: "500", matched: true},
	// 4. >0
	// Regexp: `^[1-9]\d*$`
	{reg: "^[1-9]\\d*$", s: "0", matched: false},
	{reg: "^[1-9]\\d*$", s: "1", matched: true},
	{reg: "^[1-9]\\d*$", s: "100", matched: true},
	{reg: "^[1-9]\\d*$", s: "200", matched: true},
	{reg: "^[1-9]\\d*$", s: "500", matched: true},
	// 5. 0 not match 00x match
	// Regexp: `^(?!0$)0\d*$`
	{reg: "^(?!0$)0\\d*$", s: "0", matched: false},
	{reg: "^(?!0$)0\\d*$", s: "001", matched: true},
	{reg: "^(?!0$)0\\d*$", s: "00100", matched: true},
}
