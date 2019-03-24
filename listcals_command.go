package main

import (
	"fmt"
	"sort"

	reqflags "github.com/shouji-kazuo/cli-reqflags"
	"github.com/shouji-kazuo/gocal-cli-go/google-cal"
	cli "gopkg.in/urfave/cli.v2"
)

var listCalendarsCommand = &cli.Command{
	Name:        "calendars",
	Usage:       "",
	Description: "list calendar IDs",
	ArgsUsage:   "",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    argCredentialJSONPath,
			Aliases: []string{"credential", "cred", "c"},
			Usage:   "set 'credentials.json' path",
		},
		&cli.StringFlag{
			Name:    argTokenJSONPath,
			Aliases: []string{"tok"},
			Usage:   "set token json path to save",
		},
	},
	Action: func(ctx *cli.Context) error {
		credentialPath := ctx.String(argCredentialJSONPath)
		tokenJSONPath := ctx.String(argTokenJSONPath)
		onMissingCredentialJSONPath := func() error {
			fmt.Print("Enter credential.json path: ")
			for nScanned, err := fmt.Scan(&credentialPath); nScanned != 1 || err != nil; {
				continue
			}
			return nil
		}
		onMissingTokenJSONPath := func() error {
			fmt.Print("Enter token json path: ")
			for nScanned, err := fmt.Scan(&tokenJSONPath); nScanned != 1 || err != nil; {
				continue
			}
			return nil
		}
		err := reqflags.Recover(ctx,
			map[string]func() error{
				argCredentialJSONPath: onMissingCredentialJSONPath,
				argTokenJSONPath:      onMissingTokenJSONPath,
			})
		if err != nil {
			return err
		}

		cal, err := googlecalendar.New(tokenJSONPath, credentialPath)
		if err != nil {
			return err
		}

		cals, err := cal.ListCalendars()
		if err != nil {
			return err
		}
		sort.Slice(cals.Items, func(i, j int) bool {
			return cals.Items[i].Id < cals.Items[j].Id
		})
		for _, item := range cals.Items {
			fmt.Println(item.Id)
		}
		return nil
	},
}
