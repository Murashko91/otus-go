package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const outFileName, tmpFileName = "outFile.txt", "tmpFile.txt"

func TestCopy(t *testing.T) {
	t.Run("Test ErrUnsupportedFile", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp(".", tmpFileName)

		defer func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}()

		err := Copy(tmpFile.Name(), outFileName, 0, 0)

		require.Equal(t, err, ErrUnsupportedFile)

		_, outFileErr := os.Stat(outFileName)

		if outFileErr == nil {
			os.Remove(outFileName)
		}

		require.Equal(t, true, os.IsNotExist(outFileErr))
	})

	t.Run("Test ErrOffsetExceedsFileSize", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp(".", tmpFileName)
		tmpFile.WriteString("1")

		defer func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}()

		err := Copy(tmpFile.Name(), outFileName, 2, 0)

		_, outFileErr := os.Stat(outFileName)

		if outFileErr == nil {
			os.Remove(outFileName)
		}

		require.Equal(t, true, os.IsNotExist(outFileErr))

		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("Test Positive test", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp(".", tmpFileName)
		tmpFile.WriteString("12345")

		defer func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}()

		err := Copy(tmpFile.Name(), outFileName, 2, 1)

		_, outFileErr := os.Stat(outFileName)

		outFile, _ := os.Open(outFileName)

		data, _ := io.ReadAll(outFile)

		fmt.Println(string(data) == "3")

		if outFileErr == nil {
			os.Remove(outFileName)
		}

		require.Equal(t, string(data), "3")
		require.Equal(t, err, nil)
	})
}
