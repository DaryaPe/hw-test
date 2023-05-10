package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	defer func() { _ = deleteTestData() }()

	err := generateTestData()
	if err != nil {
		fmt.Printf("generateTestData:%v", err)
	}

	t.Run("dir not exists", func(t *testing.T) {
		_, err := ReadDir("testdata/env111")
		require.ErrorIs(t, err, os.ErrNotExist)
	})
	t.Run("not valid env name", func(t *testing.T) {
		_, err := ReadDir("tests")
		require.ErrorIs(t, err, ErrNotValidEnvName)
	})
	t.Run("ok", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.ErrorIs(t, err, nil)
		envExp := Environment{
			"BAR": EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
			"EMPTY": EnvValue{
				Value:      "",
				NeedRemove: false,
			},
			"FOO": EnvValue{
				Value:      "   foo\nwith new line",
				NeedRemove: false,
			},
			"HELLO": EnvValue{
				Value:      "\"hello\"",
				NeedRemove: false,
			},
			"UNSET": EnvValue{
				Value:      "",
				NeedRemove: true,
			},
		}
		require.EqualValues(t, envExp, env)
	})
}

func generateTestData() error {
	err := os.MkdirAll("tests", os.ModePerm)
	if err != nil {
		return fmt.Errorf("os.MkdirAll:%w", err)
	}

	file, err := os.Create("tests/test=test")
	if err != nil {
		return fmt.Errorf("os.Create:%w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.WriteString("testtest")
	if err != nil {
		return fmt.Errorf("file.WriteString:%w", err)
	}
	return nil
}

func deleteTestData() error {
	err := os.RemoveAll("tests")
	if err != nil {
		return fmt.Errorf("os.RemoveAll:%w", err)
	}
	return nil
}
