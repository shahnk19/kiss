package baseEnc

import (
	"errors"
	"math"
	"strings"
)

const (
	Base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Base16 = "0123456789abcdef"
)

type Encoding struct {
	encode string
	base   int
}

func NewEncoding(encoder string) (*Encoding, error) {
	//make sure the encode has unique chars
	if !isUnique(encoder) {
		return nil, errors.New("base chars are not unique!")
	}
	return &Encoding{
		encode: encoder,
		base:   len(encoder),
	}, nil
}

func Base62Encoding() (*Encoding, error) {
	return NewEncoding(Base62)
}

func Base16Encoding() (*Encoding, error) {
	return NewEncoding(Base16)
}

func (e *Encoding) BaseEncode(n int) string {
	strBased := ""
	if e.base <= 0 || n < 0 {
		return ""
	}
	if n == 0 {
		return string(e.encode[0])
	}

	for n > 0 {
		strBased = string(e.encode[n%e.base]) + strBased
		n = n / e.base
	}
	return strBased
}

func (e *Encoding) BaseDecode(s string) (int, error) {
	if s == "" {
		return 0, errors.New("Encoded string is empty")
	}
	n := 0
	for i, v := range s {
		idx := strings.IndexRune(e.encode, v)
		if idx < 0 {
			return 0, errors.New("Encoded string is not properly formated for the choosen base")
		}
		power := len(s) - (i + 1)
		c := int64(idx) * int64(math.Pow(float64(e.base), float64(power)))
		n = n + int(c)
	}
	return n, nil
}

func isUnique(s string) bool {
	if s == "" {
		return false
	}
	isUnique := make(map[string]bool, 0)
	for _, v := range s {
		if isUnique[string(v)] {
			return false
		}
		isUnique[string(v)] = true
	}
	return true
}
