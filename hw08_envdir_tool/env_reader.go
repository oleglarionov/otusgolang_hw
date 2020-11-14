package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	env := Environment{}
	for _, f := range files {
		fileName := f.Name()

		value, err := readValue(dir, fileName)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		env[fileName] = value
	}

	return env, nil
}

func readValue(dir string, fileName string) (string, error) {
	f, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	s, err := r.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", errors.WithStack(err)
	}

	s = strings.ReplaceAll(s, "\u0000", "\n")
	s = strings.TrimRightFunc(s, unicode.IsSpace)

	return s, nil
}
