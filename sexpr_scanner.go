package pg

import (
	"errors"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// SexprScanner is a scanner for s-expressions.
type SexprScanner struct {
	in  string
	out []byte
	acc string
	err error
}

// NewSexprScanner creates a new s-expression scanner.
func NewSexprScanner[T string | []byte](x T) *SexprScanner {
	out := make([]byte, len(x))
	switch x := interface{}(x).(type) {
	case string:
		return &SexprScanner{in: x, out: out}
	case []byte:
		return &SexprScanner{in: unsafe.String(unsafe.SliceData(x), len(x)), out: out}
	}
	panic("bad input")
}

// Scan scans an s-expression.
func (sc *SexprScanner) Scan() rune {
	if len(sc.in) == 0 {
		return EOF
	}
	r, sz := utf8.DecodeRuneInString(sc.in)
	if r == utf8.RuneError {
		sc.err = errors.New("invalid UTF8 encoding")
		return Error
	}
	for unicode.IsSpace(r) {
		sc.in = sc.in[sz:]
		if len(sc.in) == 0 {
			return EOF
		}
		r, sz = utf8.DecodeRuneInString(sc.in)
		if r == utf8.RuneError {
			sc.err = errors.New("invalid UTF8 encoding")
			return Error
		}
	}
	outStart := unsafe.SliceData(sc.out)
	var outSize int
	switch r {
	case '(', ')':
		sc.in = sc.in[sz:]
		utf8.EncodeRune(sc.out, r)
		sc.acc = unsafe.String(outStart, sz)
		sc.out = sc.out[sz:]
		return r
	case '"':
		sc.in = sc.in[sz:]
		r, sz = utf8.DecodeRuneInString(sc.in)
		if r == utf8.RuneError {
			sc.err = errors.New("invalid UTF8 encoding")
			return Error
		}
		for {
			if r == '"' {
				sc.in = sc.in[sz:]
				sc.acc = unsafe.String(outStart, outSize)
				return String
			} else if r == '\\' {
				sc.in = sc.in[sz:]
				r, sz = utf8.DecodeRuneInString(sc.in)
				if r == utf8.RuneError {
					sc.err = errors.New("invalid UTF8 encoding")
					return Error
				}
				utf8.EncodeRune(sc.out, r)
				outSize += sz
				sc.out = sc.out[sz:]
			} else {
				utf8.EncodeRune(sc.out, r)
				outSize += sz
				sc.out = sc.out[sz:]
			}
			sc.in = sc.in[sz:]
			r, sz = utf8.DecodeRuneInString(sc.in)
			if r == utf8.RuneError {
				sc.err = errors.New("invalid UTF8 encoding")
				return Error
			}
		}
	}
	for {
		if r == ')' || unicode.IsSpace(r) {
			sc.acc = unsafe.String(outStart, outSize)
			if r != ')' {
				sc.in = sc.in[sz:]
			}
			return Ident
		}
		utf8.EncodeRune(sc.out, r)
		outSize += sz
		sc.out = sc.out[sz:]
		sc.in = sc.in[sz:]
		r, sz = utf8.DecodeRuneInString(sc.in)
		if r == utf8.RuneError {
			sc.err = errors.New("invalid UTF8 encoding")
			return Error
		}
	}
}

// TokenText returns the text of the token.
func (sc *SexprScanner) TokenText() string { return sc.acc }

func (sc *SexprScanner) Error() error { return sc.err }
