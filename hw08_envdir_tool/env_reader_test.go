package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Success test case", func(t *testing.T) {
		envVars, err := ReadDir("testdata/env")
		require.Equal(t, len(envVars), 5)

		require.Equal(t, err, nil)
		require.Equal(t, envVars["BAR"].Value, "bar")
		require.Equal(t, envVars["BAR"].NeedRemove, false)
		require.Equal(t, envVars["EMPTY"].Value, "")
		require.Equal(t, envVars["EMPTY"].NeedRemove, false)
		require.Equal(t, envVars["FOO"].Value, "   foo\nwith new line")
		require.Equal(t, envVars["FOO"].NeedRemove, false)
		require.Equal(t, envVars["HELLO"].Value, `"hello"`)
		require.Equal(t, envVars["HELLO"].NeedRemove, false)
		require.Equal(t, envVars["UNSET"].Value, "")
		require.Equal(t, envVars["UNSET"].NeedRemove, true)
	})

	t.Run("Error test read not dir", func(t *testing.T) {
		envVars, err := ReadDir("testdata/env/BAR")
		fmt.Println(err.Error())
		require.Contains(t, err.Error(), "BAR: not a directory")

		require.Equal(t, len(envVars), 0)
	})

	t.Run("Error test read  wrong path", func(t *testing.T) {
		envVars, err := ReadDir("testdata/env/wrong_path")

		fmt.Println(err.Error())
		require.Contains(t, err.Error(), "no such file or directory")

		require.Equal(t, len(envVars), 0)
	})
}
