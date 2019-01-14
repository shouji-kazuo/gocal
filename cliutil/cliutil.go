package cliutil

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v2"
)

type ContextKeys struct {
	CredentialJSONPathKey string
	TokenJSONPathKey      string
}

type JSONPaths struct {
	CredentialJSONPath string
	TokenJSONPath      string
}

func GetJSONPaths(ctx *cli.Context, keys *ContextKeys) (*JSONPaths, error) {
	var err error
	credentialJSONPath := ctx.String(keys.CredentialJSONPathKey)
	if credentialJSONPath == "" {
		if credentialJSONPath, err = getDefaultCredentialTokenPath(); err != nil {
			return nil, errors.Wrap(err, "cannot get credential json file path.")
		}
	}

	tokenJSONPath := ctx.String(keys.TokenJSONPathKey)
	if tokenJSONPath == "" {
		if tokenJSONPath, err = getDefaultTokenPathToSave(); err != nil {
			return nil, errors.Wrap(err, "cannot get token path to save.")
		}
	}
	return &JSONPaths{
		CredentialJSONPath: credentialJSONPath,
		TokenJSONPath:      tokenJSONPath,
	}, nil
}

func getDefaultTokenPathToSave() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "Users home directory is not found")
	}
	return filepath.Join(dir, ".gocal-cli-go", "token.json"), nil
}

func getDefaultCredentialTokenPath() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "Users home directory is not found")
	}
	return filepath.Join(dir, ".gocal-cli-go", "credential.json"), nil
}
