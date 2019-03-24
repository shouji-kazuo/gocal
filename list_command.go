package main

import (
	cli "gopkg.in/urfave/cli.v2"
)

var listCommand = &cli.Command{
	Name:        "list",
	Usage:       "",
	Description: "list schedule with duration from google calendar",
	ArgsUsage:   "",
	Subcommands: []*cli.Command{
		listEventsCommand,
		listCalendarsCommand,
	},
}
