package parser

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"unicode"
)

// Function to remove BOM and non-ASCII characters
func removeBOMAndNonASCII(s string) string {
	isNonASCII := func(r rune) bool {
		return r > unicode.MaxASCII
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isNonASCII), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return s // return the original string in case of error
	}
	return result
}
