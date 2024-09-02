package sexpr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSexprParseKinds(t *testing.T) {
	req := require.New(t)

	r, err := Parse(`(a "b")`)

	req.NoError(err)
	req.Equal(2, len(r))
	req.Equal(Identifier("a"), r[0])
	req.Equal(QuotedString("b"), r[1])
}
