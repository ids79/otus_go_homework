package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 2 && os.Args[1] == "to-file" {
		file, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		os.Stdout = file
	}
	bar, ok := os.LookupEnv("BAR")
	if ok {
		fmt.Println("BAR:", bar)
	}
	empty, ok := os.LookupEnv("EMPTY")
	if ok {
		fmt.Println("EMPTY:", empty)
	}
	foo, ok := os.LookupEnv("FOO")
	if ok {
		fmt.Println("FOO:", foo)
	}
	hello, ok := os.LookupEnv("HELLO")
	if ok {
		fmt.Println("HELLO:", hello)
	}
	_, ok = os.LookupEnv("UNSET")
	if !ok {
		fmt.Println("UNSET: is unsetted")
	}
}
