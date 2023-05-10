package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	name := cmd[0]
	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}

	err := setEnvs(env)
	if err != nil {
		log.Printf("setEnvs error: %v\n", err)
	}

	command := exec.Command(name, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		log.Printf("command.Run error: %v\n", err)
	}
	return command.ProcessState.ExitCode()
}

func setEnvs(env Environment) error {
	for key, val := range env {
		if val.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return err
			}
			delete(env, key)
			continue
		}
		err := os.Setenv(key, val.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
