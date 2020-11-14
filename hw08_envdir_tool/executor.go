package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		log.Fatal("invalid cmd")
	}

	program := cmd[0]
	args := cmd[1:]
	c := exec.Command(program, args...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout

	if env != nil {
		cEnv := os.Environ()
		for key, value := range env {
			cEnv = append(cEnv, fmt.Sprintf("%s=%s", key, value))
		}
		c.Env = cEnv
	}

	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}

	return c.ProcessState.ExitCode()
}
