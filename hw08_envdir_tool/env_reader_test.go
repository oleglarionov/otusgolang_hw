package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	env, err := ReadDir("testdata/env")
	require.NoError(t, err)
	require.Equal(t, "bar", env["BAR"])
	require.Equal(t, "   foo\nwith new line", env["FOO"])
	require.Equal(t, `"hello"`, env["HELLO"])
	require.Equal(t, "", env["UNSET"])
}
