package main

import (
	"os"
)

func main() {
	var path string

	cmdArgs := make([]string, 0, len(os.Args)-1)

	for i, arg := range os.Args {
		switch i {
		case 1:
			path = arg
		case 2:
			cmdArgs = os.Args[2:]
		}
	}

	envVars, err := ReadDir(path)
	if err != nil {
		os.Exit(1)
	}

	result := RunCmd(cmdArgs, envVars)
	os.Exit(result)
}
