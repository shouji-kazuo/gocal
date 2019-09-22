package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v2"
)

const (
	argCredentialJSONPath = "credential-json"
	argTokenJSONPath      = "token-json"
)

func main() {

	app := &cli.App{
		Name:      "gocal",
		Usage:     "GoogleCalendar CLI in Golang",
		ArgsUsage: " ",
		Version:   "v1.0.5",
		Flags:     []cli.Flag{},
		Commands: []*cli.Command{
			loginCommand,
			addCommand,
			listCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

}
