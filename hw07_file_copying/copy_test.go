package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Test ErrUnsupportedFile", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp(".", "tmpFile.txt")
		outFileName := "outFile.txt"

		defer func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}()

		err := Copy(tmpFile.Name(), outFileName, 0, 0)

		require.Equal(t, err, ErrUnsupportedFile)
	})
}
