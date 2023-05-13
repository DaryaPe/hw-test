package main

import (
	"fmt"
	"log"
	"os"
)

var ErrNotEnoughArgs = fmt.Errorf("not enough arguments to run")

func main() {
	if len(os.Args) < 3 {
		log.Fatal(ErrNotEnoughArgs)
	}
	args := os.Args[1:]

	envs, err := ReadDir(args[0])
	if err != nil {
		log.Fatalf("ReadDir error: %v\n", err)
	}
	os.Exit(RunCmd(args[1:], envs))
}
