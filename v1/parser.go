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

// MatchString checks whether a string matches the given regular expression, with additional supports for negative & positive lookaheads compared with regexp.MatchString(). The algorithm in this version uses stacks to extract lookahead expressions, hence, it does not support regular expressions with NESTED lookaheads.
func MatchString(pattern string, s string) (matched bool, err error) {
	lookaheads, err := splitRegex(pattern)
	if err != nil {
		return false, err
	}
	return matchString(pattern, s, lookaheads, 0, 0, true)
}

func matchString(pattern string, s string, lookaheads []lookahead, strOffset int, exprOffset int, isFirstRound bool) (matched bool, err error) {
	if len(lookaheads) == 0 {
		fmt.Printf("try to match %v by %v\n", s, pattern)
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return false, err
		}
		return reg.MatchString(s), nil
	}
	// the implementation divides regular expr into numbers of pieces, each of which contains two parts:
	// TYP | ordinary expr | lookahead expr | ...
	// LEN | >=0           | >=0            | ...

	// step-1: try to match string with ordinary expr
	curLookahead := lookaheads[0]
	start, end := curLookahead.idx[0]-exprOffset, curLookahead.idx[1]-exprOffset

	exprOrd := pattern[:start]
	if !isFirstRound {
		exprOrd = "^" + exprOrd
	}
	exprLkahead := pattern[start+3 : end-1]

	reg, err := regexp.Compile(exprOrd)
	if err != nil {
		return false, err
	}
	fmt.Printf("try to match %v by %v\n", s, exprOrd)

	idxMatched := reg.FindAllStringIndex(s, -1)
	exprOffset = end + exprOffset

	// step-2 try to check whether the latter part of the string meets the lookahead assertion
	ret := false
	for _, idx := range idxMatched {
		strSuf := s[idx[1]:]
		reg, err := regexp.Compile("^" + exprLkahead)
		fmt.Printf("try to match %v by %v\n", strSuf, "^"+exprLkahead)
		if err != nil {
			return false, err
		}
		matched := reg.MatchString(strSuf)
		if curLookahead.t == bracketNegLookahead {
			matched = !matched
		}
		if !matched {
			continue
		}
		newStrOffset := strOffset + idx[1]
		newPattern := strings.Clone(pattern[end:])
		newLookaheads := slices.Clone(lookaheads[1:])

		// step-3 recursively handle the next piece of regular expr
		_matched, err := matchString(newPattern, strSuf, newLookaheads, newStrOffset, exprOffset, false)
		if err != nil {
			return false, err
		}
		if _matched {
			ret = _matched
			break
		}

	}
	return ret, nil
}
