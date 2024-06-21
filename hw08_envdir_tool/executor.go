package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	commandName := cmd[0]
	cmdExec := exec.Command(commandName, cmd[1:]...)
	cmdExec.Stdout = os.Stdout

	for envKey := range env {
		os.Unsetenv(envKey)

		_, ok := os.LookupEnv(envKey)
		if ok {
			os.Unsetenv(envKey)
		}

		if !env[envKey].NeedRemove {
			os.Setenv(envKey, env[envKey].Value)
		}
	}

	cmdExec.Env = os.Environ()

	err := cmdExec.Run()
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	return 0
}
