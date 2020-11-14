package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadDir(t *testing.T) {
	env, _ := ReadDir("testdata/env")
	require.Equal(t, "bar", env["BAR"])
	require.Equal(t, "   foo\nwith new line", env["FOO"])
	require.Equal(t, `"hello"`, env["HELLO"])
	require.Equal(t, "", env["UNSET"])
}
