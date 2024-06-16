package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fr, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	fi, err := fr.Stat()
	if err != nil {
		return err
	}

	fileSize := fi.Size()

	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	if fileSize < 1 {
		return ErrUnsupportedFile
	}
	_, err = fr.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	fw, err := os.Create(toPath)
	if err != nil {
		return err
	}

	fileOutSize := getOutFileSize(fileSize, offset, limit)

	bar := pb.Full.Start64(fileOutSize)

	// create proxy reader
	barReader := bar.NewProxyReader(fr)

	defer func() {
		barReader.Close()
		barReader.Close()
	}()

	if limit > 0 {
		io.CopyN(fw, barReader, limit)
	} else {
		io.Copy(fw, barReader)
	}

	return nil
}

func getOutFileSize(fileSize, offset, limit int64) int64 {
	fileOutSize := fileSize

	if offset > 0 {
		fileOutSize -= offset
	}

	if limit > 0 && fileOutSize > limit {
		fileOutSize = limit
	}

	return fileOutSize
}
