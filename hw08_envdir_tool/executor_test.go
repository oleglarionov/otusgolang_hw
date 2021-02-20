package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildCmd(t *testing.T) {
	c := buildCmd([]string{"./testdata/echo.sh", "qwe", "123"}, []string{
		"HELLO=HELLO1",
		"BAR=BAR1",
		"FOO=FOO1",
		"UNSET=",
	})

	require.Equal(t, "./testdata/echo.sh", c.Path)
	require.Equal(t, []string{"./testdata/echo.sh", "qwe", "123"}, c.Args)
	require.Subset(t, c.Env, []string{"HELLO=HELLO1", "BAR=BAR1", "FOO=FOO1", "UNSET="})
}

func TestRunCmd(t *testing.T) {
	code := RunCmd([]string{"./testdata/echo.sh", "qwe", "123"}, []string{
		"HELLO=HELLO1",
		"BAR=BAR1",
		"FOO=FOO1",
		"UNSET=",
	})

	require.Equal(t, 0, code)
}
