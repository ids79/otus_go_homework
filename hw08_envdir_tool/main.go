package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		return
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	returnCode := RunCmd(os.Args[2:], env)
	os.Exit(returnCode)
}
