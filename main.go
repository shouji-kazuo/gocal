package main

import (
	"fmt"
	"os"

	"github.com/shouji-kazuo/gocal-cli-go/cliutil"
	"gopkg.in/urfave/cli.v2"
)

const (
	argCredentialJSONPath = "credential-json"
	argTokenJSONPath      = "token-json"
)

var defaultContextArgKeys = &cliutil.ContextKeys{
	CredentialJSONPathKey: argCredentialJSONPath,
	TokenJSONPathKey:      argTokenJSONPath,
}

func main() {

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
