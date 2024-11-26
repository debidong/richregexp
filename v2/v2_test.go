package v2

import (
	"fmt"
	"testing"
)

func TestSplitRegexV2(t *testing.T) {
	str := "foo(?!bar)"
	b := splitRegex(str)
	fmt.Println(b)
}
