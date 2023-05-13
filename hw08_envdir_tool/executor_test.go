package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("with cmd is empty", func(t *testing.T) {
		code := RunCmd([]string{""}, nil)
		require.Equal(t, -1, code)
	})
	t.Run("with env is empty", func(t *testing.T) {
		code := RunCmd([]string{"/bin/bash"}, nil)
		require.Equal(t, 0, code)
	})
	t.Run("ok", func(t *testing.T) {
		err := os.Setenv("EMPTY", "empty")
		if err != nil {
			fmt.Println(err)
		}
		err = os.Setenv("BAR", "123")
		if err != nil {
			fmt.Println(err)
		}
		envs := Environment{
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"BAR":   EnvValue{Value: "bar"},
		}
		code := RunCmd([]string{"/bin/bash"}, envs)
		require.EqualValues(t, 0, code)

		empty := os.Getenv("EMPTY")
		require.EqualValues(t, "", empty)
		bar := os.Getenv("BAR")
		require.EqualValues(t, "bar", bar)
	})
}
