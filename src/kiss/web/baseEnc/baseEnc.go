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

func BaseEncode(n int, baseVar string) string {
	strBased := ""
	base := len(baseVar)
	if base <= 0 || n < 0 {
		return ""
	}
	if n == 0 {
		return string(baseVar[0])
	}

	for n > 0 {
		strBased = string(baseVar[n%base]) + strBased
		n = n / base
	}
	return strBased
}

func BaseDecode(s string, baseVar string) (int, error) {
	if s == "" {
		return 0, errors.New("Encoded string is empty")
	}
	n := 0
	for i, e := range s {
		idx := strings.IndexRune(baseVar, e)
		if idx < 0 {
			return 0, errors.New("Encoded string is not properly formated for the choosen base")
		}
		power := len(s) - (i + 1)
		c := int64(idx) * int64(math.Pow(float64(len(baseVar)), float64(power)))
		n = n + int(c)
	}
	return n, nil
}
