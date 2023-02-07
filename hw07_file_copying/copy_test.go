package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var tests = []struct {
	name   string
	offset int64
	limit  int64
	model  string
}{
	{
		name:   "offset 0 limit 0",
		offset: 0,
		limit:  0,
		model:  "testdata/out_offset0_limit0.txt",
	},
	{
		name:   "offset 0 limit 10",
		offset: 0,
		limit:  10,
		model:  "testdata/out_offset0_limit10.txt",
	},
	{
		name:   "offset 0 limit 1000",
		offset: 0,
		limit:  1000,
		model:  "testdata/out_offset0_limit1000.txt",
	},
	{
		name:   "offset 0 limit 10000",
		offset: 0,
		limit:  10000,
		model:  "testdata/out_offset0_limit10000.txt",
	},
	{
		name:   "offset 100 limit 1000",
		offset: 100,
		limit:  1000,
		model:  "testdata/out_offset100_limit1000.txt",
	},
	{
		name:   "offset 6000 limit 1000",
		offset: 6000,
		limit:  1000,
		model:  "testdata/out_offset6000_limit1000.txt",
	},
}

var testsError = []struct {
	name   string
	offset int64
	limit  int64
	from   string
	err    error
}{
	{
		name:   "Error Offset",
		offset: 100000,
		limit:  100,
		from:   "testdata/input.txt",
		err:    ErrOffsetExceedsFileSize,
	},
	{
		name: "Unsupported file",
		from: "/dev/urandom",
		err:  ErrUnsupportedFile,
	},
	{
		name: "Undefined file",
		from: "testdata/non.txt",
		err:  ErrUnsupportedFile,
	},
}

func TestCopy(t *testing.T) {
	from := "testdata/input.txt"
	to := "testdata/out.txt"
	defer os.Remove(to)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Copy(from, to, test.offset, test.limit)
			require.Nil(t, result)
			require.True(t, compareFilesByLine(to, test.model))
		})
	}
	for _, test := range testsError {
		t.Run(test.name, func(t *testing.T) {
			result := Copy(test.from, to, test.offset, test.limit)
			require.ErrorIs(t, result, test.err)
		})
	}
}

func compareFilesByLine(fromPath, toPeath string) bool {
	fSrc, err := os.Open(fromPath)
	if err != nil {
		return false
	}
	defer fSrc.Close()
	fDst, err := os.Open(toPeath)
	if err != nil {
		return false
	}
	defer fDst.Close()
	scan1 := bufio.NewScanner(fSrc)
	scan2 := bufio.NewScanner(fDst)
	for {
		scrEOF := scan1.Scan()
		dstEOF := scan2.Scan()
		if !scrEOF && !dstEOF {
			break
		}
		if scan1.Text() != scan2.Text() {
			return false
		}
	}
	return true
}
