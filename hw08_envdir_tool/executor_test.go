package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	_ = RunCmd([]string{"./testdata/echo.sh", "qwe", "123"}, Environment{
		"HELLO": "HELLO1",
		"BAR":   "BAR1",
		"FOO":   "FOO1",
		"UNSET": "UNSET1",
	})

	_ = w.Close()
	out, _ := ioutil.ReadAll(r)

	require.Equal(
		t,
		"HELLO is (HELLO1)\n"+
			"BAR is (BAR1)\n"+
			"FOO is (FOO1)\n"+
			"UNSET is (UNSET1)\n"+
			"arguments are qwe 123\n",
		string(out),
	)
}
