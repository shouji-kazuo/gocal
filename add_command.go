package main

import (
	cli "gopkg.in/urfave/cli.v2"
)

const (
	argCalendarName = "calendar"
	argTitle        = "title"
	argDescription  = "description"
	argWhen         = "when"
	argDuration     = "duration"
)

var addCommand = &cli.Command{
	Name:        "add",
	Usage:       "",
	Description: "add schedule",
	ArgsUsage:   "",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    argCredentialJSONPath,
			Aliases: []string{"credential", "cred", "c"},
			Value:   "",
			Usage:   "set 'credentials.json' path",
		},
		&cli.StringFlag{
			Name:    argTokenJSONPath,
			Aliases: []string{"token"},
			Usage:   "set calendar name",
		},
		&cli.StringFlag{
			Name:    argCalendarName,
			Aliases: []string{"c"},
			Usage:   "set calendar name",
		},
		&cli.StringFlag{
			Name:    argTitle,
			Aliases: []string{"t"},
			Usage:   "set schedule title",
		},
		&cli.StringFlag{
			Name:    argDescription,
			Aliases: []string{"de", "desc", "dc"},
			Usage:   "set schedule description",
		},
		&cli.StringFlag{
			Name:    argWhen,
			Aliases: []string{"w", "wh"},
			Usage:   "set schedule added date in 'yyyy/MM/dd hh:mm:ss'", //TODO タイムゾーン
		},
		&cli.IntFlag{
			Name:    argDuration,
			Aliases: []string{"du", "dr"},
			Usage:   "set schedule duration in minites",
		},
	},
	Action: func(ctx *cli.Context) error {
		// if err := cliutil.IsAllFlagSpecified(ctx, argCalendarName, argWhen); err != nil {
		// 	return err
		// }
		// jsonPaths, err := cliutil.GetJSONPaths(ctx, defaultContextArgKeys)
		// if err != nil {
		// 	return errors.Wrap(err, "cannot get some json path.")
		// }
		// credentialJSONPath := jsonPaths.CredentialJSONPath
		// tokenJSONPath := jsonPaths.TokenJSONPath

		// credentialBytes, err := ioutil.ReadFile(credentialJSONPath)
		// if err != nil {
		// 	return errors.Wrap(err, "Cannot read credential json file.")
		// }
		// config, err := google.ConfigFromJSON(credentialBytes, calendar.CalendarReadonlyScope)
		// if err != nil {
		// 	return errors.Wrap(err, "Unable to parse client secret file to config.")
		// }
		// client := nil
		return nil
	},
}
