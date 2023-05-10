package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
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
	if err := validateFilePaths(fromPath, toPath); err != nil {
		return err
	}

	fromStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	fromSize := fromStat.Size()
	if offset > fromSize {
		return ErrOffsetExceedsFileSize
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = fromFile.Close()
	}()

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

	bar := pb.New(int(copySize)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond)
	bar.ShowSpeed = true
	bar.Start()
	defer bar.Finish()

	reader := bar.NewProxyReader(source)
	if _, err := source.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	_, err := io.CopyN(dest, reader, copySize)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func validateFilePaths(fromPath, toPath string) error {
	if len(fromPath) == 0 {
		return ErrFromPathEmpty
	}
	if len(toPath) == 0 {
		return ErrToPathEmpty
	}

	fromAbs, err := filepath.Abs(fromPath)
	if err != nil {
		return err
	}

	toAbs, err := filepath.Abs(toPath)
	if err != nil {
		return err
	}

	if fromAbs == toAbs {
		return ErrFromAndToPathsEqual
	}
	return nil
}
