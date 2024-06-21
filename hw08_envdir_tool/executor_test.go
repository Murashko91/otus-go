package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fileName  = "testFile.txt"
	trueValue = "trueValue"
	initValue = "initValue"
)

func TestRunCmd(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		env := make(Environment, 0)
		os.Setenv("HELLO", initValue)

		env["HELLO"] = EnvValue{trueValue, false}

		cmd := []string{"/bin/bash", "-c", "echo $HELLO"}
		file, _ := os.Create(fileName)
		std := os.Stdout
		os.Stdout = file

		result := RunCmd(cmd, env)
		require.Equal(t, result, 0)
		fileOut, _ := os.Open(fileName)
		testResult, _ := io.ReadAll(fileOut)

		require.Contains(t, string(testResult), trueValue)
		require.NotContains(t, string(testResult), initValue)

		file.Close()
		fileOut.Close()
		os.Remove(fileName)
		os.Stdout = std
	})

	t.Run("Success case remove env var", func(t *testing.T) {
		env := make(Environment, 0)
		os.Setenv("HELLO", initValue)

		env["HELLO"] = EnvValue{trueValue, true}

		cmd := []string{"/bin/bash", "-c", "echo $HELLO"}
		file, _ := os.Create(fileName)
		std := os.Stdout
		os.Stdout = file

		result := RunCmd(cmd, env)
		require.Equal(t, result, 0)
		fileOut, _ := os.Open(fileName)
		testResult, _ := io.ReadAll(fileOut)

		require.NotContains(t, string(testResult), trueValue)
		require.NotContains(t, string(testResult), initValue)

		file.Close()
		fileOut.Close()
		os.Remove(fileName)
		os.Stdout = std
	})

	t.Run("Success case nil env var", func(t *testing.T) {
		env := make(Environment, 0)
		os.Setenv("HELLO", initValue)
		cmd := []string{"/bin/bash", "-c", "echo $HELLO"}
		file, _ := os.Create(fileName)
		std := os.Stdout
		os.Stdout = file

		result := RunCmd(cmd, env)
		require.Equal(t, result, 0)
		fileOut, _ := os.Open(fileName)
		testResult, _ := io.ReadAll(fileOut)

		require.NotContains(t, string(testResult), trueValue)
		require.Contains(t, string(testResult), initValue)

		file.Close()
		fileOut.Close()
		os.Remove(fileName)
		os.Stdout = std
	})

	t.Run("Wrong command", func(t *testing.T) {
		env := make(Environment, 0)
		cmd := []string{"testwrongcommand", "echo", "hello"}

		result := RunCmd(cmd, env)
		require.Equal(t, result, 1)
	})
}
