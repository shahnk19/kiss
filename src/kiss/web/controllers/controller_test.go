package controllers

import (
	"fmt"
	"testing"
)

//test any url validation that is not covered by golang lib url.ParseRequestURI
func TestIsValidUrl(t *testing.T) {
	testCases := []struct {
		what string
		want bool
	}{
		{"", false},
		{"https://www.google.com", true},
	}
	for _, ttc := range testCases {
		got := isValidUrl(ttc.what)
		if ttc.want != got {
			t.Error(fmt.Sprintf("(!!!)Fail (what:%v,want:%v,got:%s)", ttc.what, ttc.want, got))
		}
	}
}
