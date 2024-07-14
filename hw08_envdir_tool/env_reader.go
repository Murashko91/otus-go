package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func newEnvValue(value string, needRemove bool) EnvValue {
	return EnvValue{value, needRemove}
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make(Environment, len(dirEntries))

	files := make([]*os.File, 0, len(dirEntries))

	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()

	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() && !strings.Contains(dirEntry.Name(), "=") {
			file, err := os.Open(dir + "/" + dirEntry.Name())

			if file != nil {
				files = append(files, file)
			}

			if err != nil {
				continue
			}
			fileInfo, err := file.Stat()
			if err != nil {
				continue
			}
			if fileInfo.Size() == 0 {
				result[dirEntry.Name()] = newEnvValue("", true)
				continue
			}
			sc := bufio.NewScanner(file)
			// Read only first line
			sc.Scan()
			line := sc.Text()
			value := strings.TrimRight(line, "\x09\x20")
			value = string(bytes.ReplaceAll(
				[]byte(value),
				[]byte("\x00"), []byte("\n")))
			result[dirEntry.Name()] = newEnvValue(value, false)
		}
	}
	return result, nil
}
