package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ErrNotValidEnvName = fmt.Errorf("not valid environment name")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirContent, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(dirContent))
	for _, file := range dirContent {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		if err = validateEnvName(info.Name()); err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		value, err := os.Open(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		env[file.Name()] = readEnvValue(value)
		err = value.Close()
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

func validateEnvName(name string) error {
	if strings.Count(name, "=") != 0 {
		return ErrNotValidEnvName
	}
	return nil
}

func readEnvValue(file io.Reader) EnvValue {
	fileScanner := bufio.NewScanner(file)

	if !fileScanner.Scan() {
		return EnvValue{NeedRemove: true}
	}
	str := fileScanner.Text()
	str = strings.TrimRight(str, "\t")
	str = strings.TrimRight(str, " ")
	str = strings.Replace(str, "\x00", "\n", 1)
	return EnvValue{Value: str, NeedRemove: false}
}
