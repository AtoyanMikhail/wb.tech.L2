package unpackage

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Unpackage осуществляет примитивную распаковку строки, содержащей повторяющиеся символы/руны.
func Unpackage(s string) (string, error) {
	escape := false
	builder := &strings.Builder{}
	runes := []rune(s)

	for i, r := range runes {
		if escape {
			builder.WriteRune(r)
			escape = false
			continue
		}

		if r == '\\' {
			escape = true
			continue
		}

		if unicode.IsDigit(r) {
			if i != len(runes)-1 && unicode.IsDigit(runes[i+1]) {
				return "", fmt.Errorf("Two consecutive digits")
			}
			if i == 0 {
				continue
			}

			intR, _ := strconv.Atoi(string(r))
			builder.WriteString(strings.Repeat(string(runes[i-1]), intR-1))
		} else {
			builder.WriteRune(r)
		}
	}

	res := builder.String()

	return res, nil
}
