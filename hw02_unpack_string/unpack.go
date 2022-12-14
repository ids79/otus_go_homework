package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var bilder strings.Builder
	lastSumbol := ' '
	screening := false
	for _, r := range str {
		if lastSumbol != '\\' && !screening && (lastSumbol == ' ' || unicode.IsDigit(lastSumbol)) && unicode.IsDigit(r) {
			return "", ErrInvalidString
		}
		if unicode.IsLetter(lastSumbol) || screening {
			if unicode.IsDigit(r) {
				c, _ := strconv.Atoi(string(r))
				bilder.Grow(c)
				bilder.WriteString(strings.Repeat(string(lastSumbol), c))
			} else {
				bilder.WriteRune(lastSumbol)
			}
			screening = false
			lastSumbol = r
			continue
		}
		if lastSumbol == '\\' {
			screening = true
		}
		lastSumbol = r
	}
	if lastSumbol == '\\' {
		return "", ErrInvalidString
	}
	if unicode.IsLetter(lastSumbol) || screening {
		bilder.WriteRune(lastSumbol)
	}
	return bilder.String(), nil
}
