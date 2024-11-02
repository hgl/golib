package base64

// This package automatically picks the right encoding to decode base64 string
// and also ignores characters like newlines and paddings
// https://groups.google.com/g/golang-nuts/c/777TQzaZBXM
//
// Credit: https://github.com/rogpeppe/misc/blob/master/cmd/base64/base64.go

import (
	"encoding/base64"
	"io"
)

func NewDecoder(r io.Reader) io.Reader {
	return base64.NewDecoder(base64.RawStdEncoding, transformer{r})
}

func Decode(dst, src []byte) (n int, err error) {
	transform(src)
	return base64.RawStdEncoding.Decode(dst, src)
}

func DecodeString(s string) ([]byte, error) {
	enc := base64.RawStdEncoding
	dbuf := make([]byte, enc.DecodedLen(len(s)))
	n, err := Decode(dbuf, []byte(s))
	return dbuf[:n], err
}

type transformer struct {
	r io.Reader
}

func (xfrm transformer) Read(p []byte) (int, error) {
	n, err := xfrm.r.Read(p)
	if n == 0 {
		return n, err
	}
	p = p[0:n]
	i := transform(p)
	return i, nil
}

func transform(p []byte) int {
	i := 0
	for _, b := range p {
		switch b {
		case '-':
			p[i] = '+'
		case '_':
			p[i] = '/'
		case ' ', '\n', '\r', '\t', '=':
			continue
		default:
			p[i] = b
		}
		i++
	}
	return i
}
