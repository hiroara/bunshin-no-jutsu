package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

var configFileName = ".makimono.yml"

func findConfig() (string, error) {
	dir := "."
	for {
		path := filepath.Join(dir, configFileName)
		stat, err := os.Stat(path)
		if os.IsNotExist(err) {
			dir, err = getParentDir(dir)
			if err != nil {
				return "", err
			}
			if dir == "" {
				return "", fmt.Errorf("Cannot find any configuration file named `%s`.", configFileName)
			}
		} else {
			if err != nil {
				return "", err
			}
			if !stat.Mode().IsRegular() {
				return "", fmt.Errorf("Cannot find any configuration file named `%s`.", configFileName)
			}
			return filepath.Abs(dir)
		}
	}
}

func getParentDir(path string) (string, error) {
	current, err := filepath.Abs(path)
	if current == "/" {
		return "", err
	}
	return filepath.Abs(filepath.Join(path, ".."))
}
