package v1

type lookahead struct {
	t   int   // type
	idx []int // idx in regex expr
}

// MatchString checks whether a string matches the given regular expression, providing additional syntax supports for negative & positive lookaheads compared with regexp.MatchString(). The algorithm in this version uses stacks to extract lookahead expressions, hence, it does not support regular expressions with NESTED lookaheads.
func (r *Regexp) MatchString(s string) bool {
	return r.matchString(s, 0, 0)
}

// the implementation divides regular expr into numbers of pieces, each of which contains two parts:
// TYP | ordinary expr | lookahead expr | ...
// LEN | >=0           | >=0            | ...
func (r *Regexp) matchString(s string, offsetStr int, idxRegexp int) bool {
	if idxRegexp == len(r.regexpMixed) { // end of recursion
		// if the recursion reaches here, the given string will match the regexp
		return true
	}
	// step-1: try to match string with ordinary expr
	expLkahead := r.regexpMixed[idxRegexp]
	expOrd := r.regexpOrd[idxRegexp]

	idxMatched := expOrd.FindAllStringIndex(s, -1)

	// step-2 try to check whether the latter part of the string meets the lookahead assertion
	ret := false
	for _, idx := range idxMatched {
		strSuf := s[idx[1]:]
		matched := expLkahead.r.MatchString(strSuf)
		if expLkahead.t == typeRegexNegLookahead {
			matched = !matched
		}
		if !matched {
			continue
		}
		// step-3 recursively handle the next piece of regular expr
		newOffsetStr := offsetStr + idx[1]
		newIdxRegexp := idxRegexp + 1
		_matched := r.matchString(strSuf, newOffsetStr, newIdxRegexp)
		if _matched {
			ret = _matched
			break
		}
	}
	return ret
}
