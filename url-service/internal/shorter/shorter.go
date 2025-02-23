package shorter

import (
	"strings"
)

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9DOo"

var alphabetLen = uint64(len(alphabet))

func Shorten(id uint64) string {
	var (
		digits  []uint64
		num     = id
		builder strings.Builder
	)

	for num > 0 {
		digits = append(digits, num%alphabetLen)
		num /= alphabetLen
	}

	Reverse(digits)

	for _, digit := range digits {
		builder.WriteString(string(alphabet[digit]))
	}

	return builder.String()
}

func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
