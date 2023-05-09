package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFromPathEmpty         = errors.New("from file path is empty")
	ErrToPathEmpty           = errors.New("to file path is empty")
	ErrFromAndToPathsEqual   = errors.New("from and to paths are equal")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validateFileParams(fromPath, toPath); err != nil {
		return err
	}

	fromStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = fromFile.Close()
	}()

	fromSize := fromStat.Size()
	if err = validateOffset(fromSize, offset); err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = toFile.Close()
	}()

	return execCopy(fromFile, toFile, offset, limit, fromSize)
}

func execCopy(source, dest *os.File, offset, limit, size int64) error {
	var copySize int64
	switch {
	case limit == 0 || limit > size-offset:
		copySize = size - offset
	case limit <= size-offset:
		copySize = limit
	}

	if _, err := source.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	bar := pb.New(int(copySize)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond)
	bar.ShowSpeed = true
	bar.Start()
	defer bar.Finish()

	reader := bar.NewProxyReader(source)

	_, err := io.CopyN(dest, reader, copySize)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func validateFileParams(fromPath, toPath string) error {
	if len(fromPath) == 0 {
		return ErrFromPathEmpty
	}
	if len(toPath) == 0 {
		return ErrToPathEmpty
	}
	if fromPath == toPath {
		return ErrFromAndToPathsEqual
	}
	return nil
}

func validateOffset(fileSize, offset int64) error {
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}
	return nil
}
