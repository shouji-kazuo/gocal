package main

import (
	"errors"

	cli "gopkg.in/urfave/cli.v2"
)

var addCommand = &cli.Command{
	Name:        "add",
	Usage:       "",
	Description: "add schedule",
	ArgsUsage:   "",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "calendar",
			Aliases: []string{"c"},
			Usage:   "set calendar name",
		},
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "set schedule title",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"de", "dc"},
			Usage:   "set schedule description",
		},
		&cli.StringFlag{
			Name:    "when",
			Aliases: []string{"w", "wh"},
			Usage:   "set schedule added date in 'yyyy/MM/dd hh:mm:ss'", //TODO タイムゾーン
		},
		&cli.IntFlag{
			Name:    "duration",
			Aliases: []string{"du", "dr"},
			Usage:   "set schedule duration in minites",
		},
	},
	Action: func(ctx *cli.Context) error {
		if !ctx.IsSet("calendar") {
			return errors.New("calendar flag is not set.")
		}
		if !ctx.IsSet("when") {
			return errors.New("when flag is not set.")
		}
		return nil
	},
}
