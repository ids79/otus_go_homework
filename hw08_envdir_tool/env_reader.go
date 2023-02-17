package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment, len(files))
	for _, entry := range files {
		name := entry.Name()
		if strings.Contains(name, "=") {
			return nil, errors.New("name contains =")
		}
		file, err := os.Open(dir + "/" + name)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		info, err := file.Stat()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			env[name] = EnvValue{"", true}
			continue
		}
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		str := strings.TrimRight(string(scanner.Bytes()), "\u0020\u0009")
		str = strings.TrimLeft(str, "\u0020\u0009")
		str = strings.ReplaceAll(str, "\x00", "\n")
		str = strings.Split(str, "\n")[0]
		env[name] = EnvValue{str, false}
	}
	return env, nil
}
