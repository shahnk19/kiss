package baseEnc

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

var testCasesVariousKnownRanges = []struct {
	enc int
	dec string
}{
	{0,"0"},
	{9, "9"},
	{9, "9"},
	{10, "a"},
	{35, "z"},
	{36, "A"},
	{61, "Z"},
	{62, "10"},
	{99, "1B"},
	{3844, "100"},
	{3860, "10g"},
	{4815162342, "5fRVGK"},
}

var encoder16 = Base16Encoding()
var encoder62 = Base62Encoding()
// encode or decode known numbers and its corresponding values
func TestEncodeBase62VariousKnownRanges(t *testing.T) {
	for _, ttc := range testCasesVariousKnownRanges {
		got := encoder62.BaseEncode(ttc.enc)
		if ttc.dec != got {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%s)", ttc.enc, ttc.dec, got))
		}
	}
}

func TestDecodeBase62VariousKnownRanges(t *testing.T) {
	for _, ttc := range testCasesVariousKnownRanges {
		got,err := encoder62.BaseDecode(ttc.dec)
		if err != nil {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,error:%v)", ttc.dec, ttc.enc, err))
		} else if ttc.enc != got {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%s)", ttc.dec, ttc.enc, got))
		}
	}
}

// encode/Decode out of range or error cases
func TestEncodeBase62OutOfRanges(t *testing.T) {
	testCases := []struct {
		what int
		want string
	}{
		{-1,""},
	}
	for _, ttc := range testCases {
		got := encoder62.BaseEncode(ttc.what)
		if ttc.want != got {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%s)", ttc.what, ttc.want, got))
		}
	}
}

func TestDecodeBase62OutOfRanges(t *testing.T) {
	testCases := []struct {
		what string
		want error
	}{
		{"0", nil}, //control test
		{"", errors.New("Encoded string is empty")},
		{"-", errors.New("Encoded string is not properly formated for the choosen base")},
		{"1-0", errors.New("Encoded string is not properly formated for the choosen base")},
	}
	for _, ttc := range testCases {
		_,err := encoder62.BaseDecode(ttc.what)
		if err != nil && err.Error() != ttc.want.Error() {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%s)", ttc.what, ttc.want, err))
		}
	}
}

// Can this encode match against golang base16 encoder for the range of decimal 0-100?
func TestEncodeAgainstGolangBase16(t *testing.T) {
	for what := 0; what < 100; what++ {
		got := encoder16.BaseEncode(what)
		want := strconv.FormatInt(int64(what), 16)
		if want != got {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%s)", what, want, got))
		}
	}
}

// since this encode case range over any bases, put it to test various bases for a range of decimal 0-20
func TestEncodeAnyBaseRanges(t *testing.T) {
	testCases := []struct {
		base string
		want []string
	}{
		{"0123456789", []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"}},
		{"abcdefghij", []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "ba", "bb", "bc", "bd", "be", "bf", "bg", "bh", "bi", "bj", "ca"}},
		{"abc", []string{"a", "b", "c", "ba", "bb", "bc", "ca", "cb", "cc", "baa", "bab", "bac", "bba", "bbb", "bbc", "bca", "bcb", "bcc", "caa", "cab", "cac"}},
	}
	for _, ttc := range testCases {
		encoder := NewEncoding(ttc.base)
		for what := 0; what < 20; what++ {
			got := encoder.BaseEncode(what)
			if ttc.want[what] != got {
				t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%s)", what, ttc.want[what], got))
			}
		}
	}
}
