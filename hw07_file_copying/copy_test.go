package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	fromFile := "testdata/input.txt"
	toFile := "out.txt"

	t.Run("wrong fromPath", func(t *testing.T) {
		err := Copy("", toFile, 0, 0)
		require.ErrorIs(t, err, ErrFromPathEmpty)
	})
	t.Run("wrong toPath", func(t *testing.T) {
		err := Copy(fromFile, "", 0, 0)
		require.ErrorIs(t, err, ErrToPathEmpty)
	})
	t.Run("wrong offset", func(t *testing.T) {
		err := Copy(fromFile, toFile, 10000000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
	t.Run("from file not exists", func(t *testing.T) {
		err := Copy("testdata/input1.txt", toFile, 0, 0)
		require.ErrorIs(t, err, os.ErrNotExist)
	})
	t.Run("from and to paths are equal", func(t *testing.T) {
		err := Copy(fromFile, fromFile, 0, 0)
		require.ErrorIs(t, err, ErrFromAndToPathsEqual)
	})
	t.Run("ok", func(t *testing.T) {
		err := Copy(fromFile, toFile, 0, 0)
		require.ErrorIs(t, err, nil)
		ff, err := os.ReadFile(fromFile)
		if err != nil {
			fmt.Println(err)
		}
		tf, err := os.ReadFile(toFile)
		if err != nil {
			fmt.Println(err)
		}
		require.EqualValues(t, tf, ff)
		err = os.Remove(toFile)
		if err != nil {
			fmt.Println(err)
		}
	})
	t.Run("ok with limit", func(t *testing.T) {
		err := Copy(fromFile, toFile, 0, 1000)
		require.ErrorIs(t, err, nil)
		ff, err := os.ReadFile(fromFile)
		if err != nil {
			fmt.Println(err)
		}
		tf, err := os.ReadFile(toFile)
		if err != nil {
			fmt.Println(err)
		}
		require.EqualValues(t, tf, ff[:1000])
		err = os.Remove(toFile)
		if err != nil {
			fmt.Println(err)
		}
	})
	t.Run("ok with limit + offset < file size", func(t *testing.T) {
		err := Copy(fromFile, toFile, 1000, 1000)
		require.ErrorIs(t, err, nil)
		ff, err := os.ReadFile(fromFile)
		if err != nil {
			fmt.Println(err)
		}
		tf, err := os.ReadFile(toFile)
		if err != nil {
			fmt.Println(err)
		}
		require.EqualValues(t, tf, ff[1000:2000])
		err = os.Remove(toFile)
		if err != nil {
			fmt.Println(err)
		}
	})
	t.Run("ok with limit + offset > file size", func(t *testing.T) {
		err := Copy(fromFile, toFile, 1000, 10000)
		require.ErrorIs(t, err, nil)
		ff, err := os.ReadFile(fromFile)
		if err != nil {
			fmt.Println(err)
		}
		tf, err := os.ReadFile(toFile)
		if err != nil {
			fmt.Println(err)
		}
		require.EqualValues(t, tf, ff[1000:])
		err = os.Remove(toFile)
		if err != nil {
			fmt.Println(err)
		}
	})
}
