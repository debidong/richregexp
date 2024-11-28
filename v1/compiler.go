package v1

import (
	"fmt"
	"regexp"
)

const (
	// external flags for identification of a regexp
	TypeRegexpOrd  = 0 // ordinary regexp with all syntax supported by re2
	TypeRegexMixed = 1 // regexp with looaheads

	// internal flags
	typeRegexNegLookahead = 2
	typeRegexPosLookahead = 3

	// types of left brackets
	typeBracketOrd          = "("
	typeBracketPosLookahead = "(?="
	typeBracketNegLookahead = "(?!"
)

// Regexp stores compiled regular expressions.
type Regexp struct {
	T           int
	regexpOrd   []*regexp.Regexp
	regexpMixed []*regexpMixed
}

type regexpMixed struct {
	t int
	r *regexp.Regexp
}
type syntaxStack struct {
	idxLookaheads []lookahead
	brackets      []bracket
}

type bracket struct {
	t   string // type
	idx int    // idx in regex expr
}

func Compile(expr string) (*Regexp, error) {
	lookaheads, err := splitRegex(expr)
	if err != nil {
		return nil, err
	}

	r := new(Regexp)
	offset := 0
	for i, l := range lookaheads {
		start, end := l.idx[0], l.idx[1]

		_regOrd := expr[offset:start]
		if i >= 1 {
			_regOrd = "^" + _regOrd
		}
		regOrd, err := regexp.Compile(_regOrd)
		if err != nil {
			return nil, fmt.Errorf("v1.Compile: %w", err)
		}
		r.regexpOrd = append(r.regexpOrd, regOrd)

		regLookahead, err := regexp.Compile("^" + expr[start+3:end-1])
		if err != nil {
			return nil, fmt.Errorf("v1.Compile: %w", err)
		}
		r.regexpMixed = append(r.regexpMixed,
			&regexpMixed{t: l.t, r: regLookahead},
		)
		offset = end
	}
	return r, nil
}

func splitRegex(str string) ([]lookahead, error) {
	if len(str) <= 3 {
		return []lookahead{}, nil
	}

	var tree syntaxStack

	i := 0
	for i < len(str) {
		switch str[i] {
		case '(':
			if str[i:i+3] == typeBracketNegLookahead {
				tree.push(typeBracketNegLookahead, i)
				i += 2
			} else if str[i:i+3] == typeBracketPosLookahead {
				tree.push(typeBracketPosLookahead, i)
				i += 2
			} else {
				tree.push(typeBracketOrd, i)
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

func (s *syntaxStack) push(t string, idx int) {
	s.brackets = append(
		s.brackets,
		bracket{t: t, idx: idx},
	)
}

func (s *syntaxStack) pop(idx int) error {
	if len(s.brackets) == 0 {
		return ErrInvalidSyntax{}
	}

	b := s.brackets[len(s.brackets)-1]
	var _t int
	switch b.t {
	case typeBracketOrd:
		s.brackets = s.brackets[:len(s.brackets)-1]
		return nil
	case typeBracketNegLookahead:
		_t = typeRegexNegLookahead
	case typeBracketPosLookahead:
		_t = typeRegexPosLookahead
	}
	lookahead := lookahead{
		t:   _t,
		idx: []int{b.idx, idx + 1},
	}
	s.idxLookaheads = append(s.idxLookaheads, lookahead)
	s.brackets = s.brackets[:len(s.brackets)-1]
	return nil
}
