package re3

import (
	"fmt"
	"regexp"
	"testing"
)

func TestSplitNegLookAhead(t *testing.T) {

	reg := "foo(?!()bar) and (?!baz) something else"
	re := regexp.MustCompile(`\(\?!([^()]|\([^\)]*\))\)`)
	matches := re.FindAllStringIndex(reg, -1)
	for _, match := range matches {
		fmt.Printf("%v from %v to %v \n", reg[match[0]:match[1]], match[0], match[1])
	}
}
