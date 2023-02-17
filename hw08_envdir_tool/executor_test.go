package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	out := "out.txt"
	defer os.Remove(out)
	t.Run("check execute program", func(t *testing.T) {
		args := []string{"./cmd/cmd", "to-file", out}
		env := make(Environment, 5)
		env["BAR"] = EnvValue{"bar", false}
		env["EMPTY"] = EnvValue{"", false}
		env["FOO"] = EnvValue{"foo", false}
		env["HELLO"] = EnvValue{`"hello"`, false}
		env["UNSET"] = EnvValue{"", true}
		returnCode := RunCmd(args, env)
		require.Equal(t, 0, returnCode)
		require.True(t, compareFilesByLine("./"+out, "./cmd/model.txt"))
	})
	t.Run("check error exit code", func(t *testing.T) {
		args := []string{""}
		env := make(Environment)
		returnCode := RunCmd(args, env)
		require.Equal(t, 127, returnCode)
		args[0] = "undefined"
		returnCode = RunCmd(args, env)
		require.Equal(t, 1, returnCode)
	})
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
