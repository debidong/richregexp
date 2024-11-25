package re3

import (
	"regexp"
)

const bracketOrdinary = "("
const bracketPosLookahead = "(?="
const bracketNegLookahead = "(?!"

type bracket struct {
	t   string
	idx int
}

type Lookaheads struct {
	idxNegLookAheads [][]int
	idxPosLookAheads [][]int
	brackets         []bracket
}

type ErrInvalidSyntax struct{}

func (e ErrInvalidSyntax) Error() string {
	return "invalid syntax"
}

func (l *Lookaheads) push(t string, idx int) {
	l.brackets = append(
		l.brackets,
		bracket{t: t, idx: idx},
	)
}

func (l *Lookaheads) pop(idx int) error {
	if len(l.brackets) == 0 {
		return ErrInvalidSyntax{}
	}

	l.brackets = l.brackets[:len(l.brackets)-1]
	b := l.brackets[len(l.brackets)-1]
	switch b.t {
	case bracketPosLookahead:
		l.idxPosLookAheads = append(l.idxPosLookAheads, []int{b.idx, idx})
	case bracketNegLookahead:
		l.idxNegLookAheads = append(l.idxNegLookAheads, []int{b.idx, idx})
	}
	return nil
}

func MustCompile(str string) []*regexp.Regexp {
	if len(str) <= 3 {
		return []*regexp.Regexp{regexp.MustCompile(str)}
	}
	var l Lookaheads
	i := 0
	for i < len(str)-3 {
		i += 1
		switch str[i] {
		case '(':
			if str[i:i+3] == bracketNegLookahead {
				i += 2
				l.push(bracketNegLookahead, i)
			} else if str[i:i+3] == bracketPosLookahead {
				l.push(bracketPosLookahead, i)
			} else {
				l.push(bracketOrdinary, i)
			}
		case ')':
			if err := l.pop(i); err != nil {
				panic(err)
			}
		}
	}

	// TODO: add return logic

	return nil
}
