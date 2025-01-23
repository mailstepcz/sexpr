package sexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	req := require.New(t)

	sc := NewScanner(`(a1 "b2" c (1 2 3) "x\\y" y z)`)
	var toks []string
	for tok := sc.Scan(); tok != EOF; tok = sc.Scan() {
		req.False(tok == Error)
		toks = append(toks, sc.TokenText())
	}
	req.Equal([]string{"(", "a1", "b2", "c", "(", "1", "2", "3", ")", "x\\y", "y", "z", ")"}, toks)
}
