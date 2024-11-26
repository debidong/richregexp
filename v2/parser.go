package v2

const bracketOrdinary = "("
const bracketPosLookahead = "(?="
const bracketNegLookahead = "(?!"

const (
	typeBlockNormal       = 0
	typeBlockMixed        = 1
	typeBlockPosLookahead = 2
	typeBlockNegLookahead = 3
	typeBlockGroup        = -1
)

type syntaxTree struct {
	root *block
}

type block struct {
	expr  string
	t     int
	root  *block
	child *block
	next  *block
}

type ErrInvalidSyntax struct{}

func (e ErrInvalidSyntax) Error() string {
	return "invalid syntax"
}

func splitRegex(str string) *block {
	root := new(block)
	splitRegexRecur(root, str)
	return root
}

func splitRegexRecur(b *block, str string) *block {
	if len(str) == 0 {
		return b
	}

	pos := -1
	for i := range str {
		if str[i] == '(' || str[i] == ')' {
			pos = i
			break
		}
	}
	if pos == -1 {
		b.expr = str
		b.t = typeBlockNormal
		b.next = nil
		return b
	}
	b.expr = str[:pos]
	b.t = typeBlockMixed

	switch str[pos] {
	case '(':
		var t int
		isComplex := true
		if len(str)-pos >= 3 {
			switch str[:3] {
			case bracketPosLookahead:
				t = typeBlockPosLookahead
			case bracketNegLookahead:
				t = typeBlockNegLookahead
			default:
				isComplex = false
				t = typeBlockGroup
			}
		} else {
			isComplex = false
			t = typeBlockGroup
		}
		newB := new(block)
		newB.t = t
		b.child = newB
		if isComplex {
			str = str[3:]
		} else {
			str = str[1:]
		}
		ret := splitRegexRecur(newB, str)
		newB.next = ret
	case ')':
		b.next = nil
		return nil
	default:
		b.next = splitRegexRecur(new(block), str)
	}
	return b
}
