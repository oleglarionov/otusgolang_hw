package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const (
	beginState = iota
	letterState
	escapeState
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	builder := &strings.Builder{}

	state := beginState
	var prevSymbol, curSymbol rune
	for _, curSymbol = range input {
		switch state {
		case beginState:
			if unicode.IsDigit(curSymbol) {
				return "", ErrInvalidString
			}
			if curSymbol == '\\' {
				state = escapeState
				break
			}

			prevSymbol = curSymbol
			state = letterState

		case letterState:
			if unicode.IsDigit(curSymbol) {
				n, _ := strconv.Atoi(string(curSymbol))
				value := strings.Repeat(string(prevSymbol), n)
				builder.WriteString(value)
				state = beginState
				break
			}

			builder.WriteRune(prevSymbol)
			if curSymbol == '\\' {
				state = escapeState
			}
			prevSymbol = curSymbol

		case escapeState:
			if curSymbol == '\\' || unicode.IsDigit(curSymbol) {
				prevSymbol = curSymbol
				state = letterState
				break
			}
			return "", ErrInvalidString
		}
	}
	if state == escapeState {
		return "", ErrInvalidString
	}
	if state == letterState {
		builder.WriteRune(curSymbol)
	}

	return builder.String(), nil
}
