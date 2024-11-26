package v1

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const bracketOrdinary = "("
const bracketPosLookahead = "(?="
const bracketNegLookahead = "(?!"

type bracket struct {
	t   string
	idx int
}

type lookahead struct {
	t   string
	idx []int
}

type syntaxTree struct {
	idxLookaheads []lookahead
	brackets      []bracket
}

type ErrInvalidSyntax struct{}

func (e ErrInvalidSyntax) Error() string {
	return "invalid syntax"
}

func (tree *syntaxTree) push(t string, idx int) {
	tree.brackets = append(
		tree.brackets,
		bracket{t: t, idx: idx},
	)
}

func (tree *syntaxTree) pop(idx int) error {
	if len(tree.brackets) == 0 {
		return ErrInvalidSyntax{}
	}

	b := tree.brackets[len(tree.brackets)-1]
	if b.t == bracketOrdinary {
		tree.brackets = tree.brackets[:len(tree.brackets)-1]
		return nil
	}
	lookahead := lookahead{
		t:   b.t,
		idx: []int{b.idx, idx + 1},
	}
	tree.idxLookaheads = append(tree.idxLookaheads, lookahead)
	tree.brackets = tree.brackets[:len(tree.brackets)-1]
	return nil
}

func splitRegex(str string) ([]lookahead, error) {
	if len(str) <= 3 {
		return []lookahead{}, nil
	}

	var tree syntaxTree

	i := 0
	for i < len(str) {
		switch str[i] {
		case '(':
			if str[i:i+3] == bracketNegLookahead {
				tree.push(bracketNegLookahead, i)
				i += 2
			} else if str[i:i+3] == bracketPosLookahead {
				tree.push(bracketPosLookahead, i)
				i += 2
			} else {
				tree.push(bracketOrdinary, i)
			}
		case ')':
			if err := tree.pop(i); err != nil {
				return nil, err
			}
		}
		i += 1
	}

	if len(tree.brackets) > 0 {
		panic(ErrInvalidSyntax{}.Error())
	}
	return tree.idxLookaheads, nil
}

func MatchString(pattern string, s string) (matched bool, err error) {
	lookaheads, err := splitRegex(pattern)
	fmt.Println(lookaheads)
	if err != nil {
		return false, err
	}
	return matchString(pattern, lookaheads, 0, 0, s)
}

func matchString(pattern string, lookaheads []lookahead, offset int, patternOffset int, s string) (matched bool, err error) {
	if len(lookaheads) == 0 {
		fmt.Println("----")
		fmt.Printf("try to match %v by %v\n", s, pattern)
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return false, err
		}
		return reg.MatchString(s), nil
	}
	start, end, t := lookaheads[0].idx[0]-patternOffset, lookaheads[0].idx[1]-patternOffset, lookaheads[0].t
	patternPre := pattern[:start]
	lookahead := pattern[start+3 : end-1]

	reg, err := regexp.Compile(patternPre)
	if err != nil {
		return false, err
	}
	fmt.Println("----")
	fmt.Printf("try to match %v by %v\n", s, patternPre)

	idxMatched := reg.FindAllStringIndex(s, -1)
	patternOffset = end + patternOffset
	for _, idx := range idxMatched {
		strOffset := idx[1] + offset
		pattern = "^" + lookahead + pattern[end:]
		var newS string
		newS = strings.Clone(s)
		newS = newS[strOffset:]
		newLookaheads := slices.Clone(lookaheads[1:])
		matched, err := matchString(pattern, newLookaheads, strOffset, patternOffset, newS)
		if err != nil {
			return false, err
		}
		if t == bracketNegLookahead {
			matched = !matched
		}
		if matched {
			return matched, nil
		}
	}
	return matched, nil
}
