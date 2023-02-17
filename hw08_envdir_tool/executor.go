package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 || cmd[0] == "" {
		return 127
	}
	comand := exec.Command(cmd[0], cmd[1:]...)
	for key, value := range env {
		os.Unsetenv(key)
		if !value.NeedRemove {
			os.Setenv(key, value.Value)
		}
	}
	comand.Env = os.Environ()
	comand.Stdin = os.Stdin
	comand.Stdout = os.Stdout
	comand.Stderr = os.Stderr
	err := comand.Run()
	if err != nil {
		var errExitCode *exec.ExitError
		if errors.As(err, &errExitCode) {
			return errExitCode.ProcessState.ExitCode()
		}
		return 1
	}

	return comand.ProcessState.ExitCode()
}
