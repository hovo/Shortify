package base62

import (
	"fmt"
	"strings"
)

const (
	base         = 62
	characterSet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Encode integer to base 62 string
func Encode(n uint64) string {
	if n == 0 {
		return string(characterSet[0])
	}

	s := ""
	for n > 0 {
		s = string(characterSet[n%base]) + s
		n = n / base
	}

	return string(s)
}

// Decode a base62 encoded string into an integer
func Decode(s string) (uint64, error) {
	var r uint64
	for _, c := range []byte(s) {
		i := strings.IndexByte(characterSet, c)
		if i < 0 {
			return 0, fmt.Errorf("Unexpected character %c in base62 literal", c)
		}
		r = base*r + uint64(i)
	}
	return r, nil
}
