package v1

import (
	"regexp"
	"slices"
	"strings"
)

type lookahead struct {
	t   int   // type
	idx []int // idx in regex expr
}

type ErrInvalidSyntax struct{}

func (e ErrInvalidSyntax) Error() string {
	return "invalid syntax"
}

// MatchString checks whether a string matches the given regular expression, providing additional syntax supports for negative & positive lookaheads compared with regexp.MatchString(). The algorithm in this version uses stacks to extract lookahead expressions, hence, it does not support regular expressions with NESTED lookaheads.
func MatchString(pattern string, s string) (matched bool, err error) {
	lookaheads, err := splitRegex(pattern)
	if err != nil {
		return false, err
	}
	return matchString(pattern, s, lookaheads, 0, 0, true)
}

func matchString(pattern string, s string, lookaheads []lookahead, strOffset int, exprOffset int, isFirstRound bool) (matched bool, err error) {
	if len(lookaheads) == 0 {
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

	idxMatched := reg.FindAllStringIndex(s, -1)
	exprOffset = end + exprOffset

	// step-2 try to check whether the latter part of the string meets the lookahead assertion
	ret := false
	for _, idx := range idxMatched {
		strSuf := s[idx[1]:]
		reg, err := regexp.Compile("^" + exprLkahead)
		if err != nil {
			return false, err
		}
		matched := reg.MatchString(strSuf)
		if curLookahead.t == typeRegexNegLookahead {
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
