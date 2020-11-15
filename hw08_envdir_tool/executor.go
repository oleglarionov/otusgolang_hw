package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd, env []string) (returnCode int) {
	c := buildCmd(cmd, env)

	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}

	return c.ProcessState.ExitCode()
}

func buildCmd(cmd, env []string) *exec.Cmd {
	if len(cmd) == 0 {
		log.Fatal("invalid cmd")
	}

	program := cmd[0]
	args := cmd[1:]
	c := exec.Command(program, args...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Env = env

	return c
}
