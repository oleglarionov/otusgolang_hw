package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode"

	"github.com/pkg/errors"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "env dir reading error")
	}

	env := make(Environment, len(files))
	for _, f := range files {
		fileName := f.Name()

		value, err := readValue(dir, fileName)
		if err != nil {
			return nil, err
		}

		env[fileName] = value
	}

	return env, nil
}

func readValue(dir string, fileName string) (string, error) {
	f, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		return "", errors.Wrap(err, "file opening error")
	}
	defer f.Close()

	r := bufio.NewReader(f)
	b, err := r.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", errors.Wrap(err, "file reading error")
	}

	b = bytes.ReplaceAll(b, []byte("\u0000"), []byte("\n"))
	b = bytes.TrimRightFunc(b, unicode.IsSpace)

	return string(b), nil
}
