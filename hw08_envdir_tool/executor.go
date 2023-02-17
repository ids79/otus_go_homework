package main

import (
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
	err := comand.Start()
	if err != nil {
		return 1
	}
	err = comand.Wait()
	if err != nil {
		return 1
	}

	return comand.ProcessState.ExitCode()
}
