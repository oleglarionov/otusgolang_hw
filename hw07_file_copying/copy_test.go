package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	inputFilePath := "./testdata/input.txt"
	tests := []struct {
		offset             int64
		limit              int64
		inputFile          string
		expectedResultFile string
	}{
		{
			offset:             0,
			limit:              0,
			inputFile:          "./testdata/input.txt",
			expectedResultFile: "./testdata/out_offset0_limit0.txt",
		},
		{
			offset:             0,
			limit:              10,
			inputFile:          "./testdata/input.txt",
			expectedResultFile: "./testdata/out_offset0_limit10.txt",
		},
		{
			offset:             0,
			limit:              1000,
			inputFile:          "./testdata/input.txt",
			expectedResultFile: "./testdata/out_offset0_limit1000.txt",
		},
		{
			offset:             0,
			limit:              10000,
			inputFile:          "./testdata/input.txt",
			expectedResultFile: "./testdata/out_offset0_limit10000.txt",
		},
		{
			offset:             100,
			limit:              1000,
			inputFile:          "./testdata/input.txt",
			expectedResultFile: "./testdata/out_offset100_limit1000.txt",
		},
		{
			offset:             6000,
			limit:              1000,
			inputFile:          "./testdata/input.txt",
			expectedResultFile: "./testdata/out_offset6000_limit1000.txt",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run("", func(t *testing.T) {
			resultFile, _ := ioutil.TempFile("", "")

			err := Copy(inputFilePath, resultFile.Name(), tc.offset, tc.limit)
			require.NoError(t, err)

			expected, _ := ioutil.ReadFile(tc.expectedResultFile)
			actual, _ := ioutil.ReadAll(resultFile)
			require.Equal(t, expected, actual, "The contents of the files are not equal")

			_ = os.Remove(resultFile.Name())
		})
	}
}

func TestCopyWithWrongArguments(t *testing.T) {
	inputFilePath := "./testdata/input.txt"
	file, _ := os.Open(inputFilePath)
	stat, _ := file.Stat()
	incorrectOffset := stat.Size() + 1

	err := Copy(inputFilePath, "/dev/null", incorrectOffset, 0)
	require.EqualError(t, err, "offset exceeds file size")
}
