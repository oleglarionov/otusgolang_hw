package hw02_unpack_string //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
		{
			input:    "\tabc\n",
			expected: "\tabc\n",
		},
		{
			input:    "\t0a1b2c3\n4",
			expected: "abbccc\n\n\n\n",
		},
		{
			input:    "走吧",
			expected: "走吧",
		},
		{
			input:    "走2吧3",
			expected: "走走吧吧吧",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
		{
			input:    `\`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `\1`,
			expected: `1`,
		},
		{
			input:    `\10`,
			expected: ``,
		},
		{
			input:    `q\we`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `q\\we`,
			expected: `q\we`,
		},
		{
			input:    `qwe\4\5\`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `\1qwe`,
			expected: `1qwe`,
		},
		{
			input:    `走\2吧3`,
			expected: "走2吧吧吧",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}
