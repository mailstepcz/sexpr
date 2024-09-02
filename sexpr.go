package sexpr

import (
	"errors"
)

// token types
const (
	Ident rune = -(iota + 1)
	String
	EOF
	Error
)

type (
	// Identifier is a symbolic atom.
	Identifier string
	// QuotedString is a quoted string.
	QuotedString string
)

func (x Identifier) String() string { return string(x) }

func (x QuotedString) String() string { return string(x) }

// Parse parses an s-expression.
func Parse(s string) ([]interface{}, error) {
	sc := NewScanner(s)
	return parseSexpr(sc, true)
}

func parseSexpr(sc *Scanner, checkFirstPar bool) ([]interface{}, error) {
	var (
		els    = make([]interface{}, 0, 10)
		inList = !checkFirstPar
	)
	for tok := sc.Scan(); tok != EOF; tok = sc.Scan() {
		if tok == Error {
			return nil, sc.Error()
		}
		if !inList {
			if tok != '(' {
				return nil, errors.New("s-expressions must begin with '('")
			}
			inList = true
		} else {
			if tok == ')' {
				break
			} else if tok == '(' {
				ex, err := parseSexpr(sc, false)
				if err != nil {
					return nil, err
				}
				els = append(els, ex)
			} else {
				switch tok {
				case Ident:
					els = append(els, Identifier(sc.TokenText()))
				case String:
					els = append(els, QuotedString(sc.TokenText()))
				}
			}
		}
	}
	return els, nil
}
