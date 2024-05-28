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

		_, outFileErr := os.Stat(outFileName)

		require.Equal(t, true, os.IsNotExist(outFileErr))
	})

	t.Run("Test ErrOffsetExceedsFileSize", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp(".", "tmpFile.txt")
		outFileName := "outFile.txt"
		tmpFile.WriteString("1")

		defer func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}()

		err := Copy(tmpFile.Name(), outFileName, 2, 0)

		_, outFileErr := os.Stat(outFileName)

		require.Equal(t, true, os.IsNotExist(outFileErr))

		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})
}
