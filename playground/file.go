package playground

import (
	"os"
	"path/filepath"
)

func GetFKitDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".fkit"), nil
}

func GetComposeFilepath() (string, error) {
	dir, err := GetFKitDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "docker-compose.yml"), nil
}
