package sexpr

import (
	"errors"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// Scanner is a scanner for s-expressions.
type Scanner struct {
	in  string
	out []byte
	acc string
	err error
}

// NewScanner creates a new s-expression scanner.
func NewScanner[T string | []byte](x T) *Scanner {
	out := make([]byte, len(x))
	switch x := interface{}(x).(type) {
	case string:
		return &Scanner{in: x, out: out}
	case []byte:
		return &Scanner{in: unsafe.String(unsafe.SliceData(x), len(x)), out: out}
	}
	panic("bad input")
}

// Scan scans an s-expression.
func (sc *Scanner) Scan() rune {
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
func (sc *Scanner) TokenText() string { return sc.acc }

func (sc *Scanner) Error() error { return sc.err }
