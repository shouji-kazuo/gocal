package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"

	"gopkg.in/urfave/cli.v2"
)

var (
	defaultTokenPath = ""
)

func main() {
	if defaultTokenPath, err := getDefaultTokenPath(); err != nil {
		// abort? → いやダメでしょ．Specificなパスを指定された時には動いてほしいはずじゃん
	}

	app := &cli.App{
		Name:      "gocal-cli-go",
		Usage:     "GoogleCalendar CLI in Golang",
		ArgsUsage: " ",
		Version:   "v1.0.5",
		Flags:     []cli.Flag{},
		Commands: []*cli.Command{
			loginCommand,
			addCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

}

func getDefaultTokenPath() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "Users home directory is not found")
	}
}
