package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	builder := &strings.Builder{}

	var prevSymbol rune
	var curSymbol rune
	var curState state = beginState
	for _, nextSymbol := range input {
		nextState, err := curState(prevSymbol, curSymbol, nextSymbol, builder, false)
		if err != nil {
			return "", err
		}

		prevSymbol = curSymbol
		curSymbol = nextSymbol
		curState = nextState
	}
	_, err := curState(prevSymbol, curSymbol, '0', builder, true)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

type state func(prevSymbol, curSymbol, nextSymbol rune, builder *strings.Builder, isFinal bool) (state, error)

func letterState(_, curSymbol, nextSymbol rune, builder *strings.Builder, isFinal bool) (state, error) {
	if isFinal {
		builder.WriteRune(curSymbol)
		return nil, nil
	}

	if unicode.IsDigit(nextSymbol) {
		return digitState, nil
	}

	builder.WriteRune(curSymbol)

	if nextSymbol == '\\' {
		return backslashState, nil
	}

	return letterState, nil
}

func digitState(prevSymbol, curSymbol, nextSymbol rune, builder *strings.Builder, isFinal bool) (state, error) {
	n, _ := strconv.Atoi(string(curSymbol))
	value := strings.Repeat(string(prevSymbol), n)

	if isFinal {
		builder.WriteString(value)
		return nil, nil
	}

	if unicode.IsDigit(nextSymbol) {
		return nil, ErrInvalidString
	}

	builder.WriteString(value)

	if nextSymbol == '\\' {
		return backslashState, nil
	}

	return letterState, nil
}

func backslashState(_, _, nextSymbol rune, _ *strings.Builder, isFinal bool) (state, error) {
	if isFinal {
		return nil, ErrInvalidString
	}

	if unicode.IsDigit(nextSymbol) {
		return letterState, nil
	}

	if nextSymbol == '\\' {
		return letterState, nil
	}

	return nil, ErrInvalidString
}

func beginState(_, _, nextSymbol rune, _ *strings.Builder, isFinal bool) (state, error) {
	if isFinal {
		return nil, nil
	}

	if unicode.IsDigit(nextSymbol) {
		return nil, ErrInvalidString
	}

	if nextSymbol == '\\' {
		return backslashState, nil
	}

	return letterState, nil
}
