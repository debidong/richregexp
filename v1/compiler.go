package v1

type Regexp struct {
}

func Compile(expr string) (regexp *Regexp, err error) {
	lookaheads, err := splitRegex(expr)
	if err != nil {
		return nil, err
	}

}
