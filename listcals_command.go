package main

import (
	"fmt"
	"sort"

	"github.com/shouji-kazuo/cliopts"
	"github.com/shouji-kazuo/gocal/google-cal"
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
		optEnsure := cliopts.NewEnsure().
			With(argCredentialJSONPath, cliopts.StdInteract("Enter credential.json path: ").After(func(s string) error {
				credentialPath = s
				return nil
			})).
			With(argTokenJSONPath, cliopts.StdInteract("Enter token json path: ").After(func(s string) error {
				tokenJSONPath = s
				return nil
			}))

		if err := optEnsure.Do(ctx); err != nil {
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
