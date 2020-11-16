package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	dirEnv, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	env := populateOsEnv(dirEnv)

	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}

func populateOsEnv(env Environment) []string {
	resultEnv := os.Environ()
	for key, value := range env {
		resultEnv = append(resultEnv, fmt.Sprintf("%s=%s", key, value))
	}

	return resultEnv
}
