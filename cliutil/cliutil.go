package cliutil

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

func GetDefaultTokenPathToSave() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "Users home directory is not found")
	}
	return filepath.Join(dir, ".gocal-cli-go", "token.json"), nil
}

func GetDefaultCredentialTokenPath() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "Users home directory is not found")
	}
	return filepath.Join(dir, ".gocal-cli-go", "credential.json"), nil
}
