package main

import (
	"fmt"
	"time"

	"github.com/shouji-kazuo/cliopts"
	"github.com/shouji-kazuo/gocal/google-cal"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	argCalendarName = "calendar"
	argTitle        = "title"
	argStart        = "start"
	argLocation     = "location"
	argEnd          = "end"

	timeLayout = "2006/01/02 15:04:05"
)

var addCommand = &cli.Command{
	Name:        "add",
	Usage:       "",
	Description: "add event",
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
			Aliases: []string{"cal", "cn"},
			Usage:   "set calendar name",
		},
		&cli.StringFlag{
			Name:    argTitle,
			Aliases: []string{"t"},
			Usage:   "set event title",
		},
		&cli.StringFlag{
			Name:    argLocation,
			Aliases: []string{"l"},
			Usage:   "set event location",
		},
		&cli.StringFlag{
			Name:    argStart,
			Aliases: []string{"s"},
			Usage:   "set schedule start added date in 'yyyy/MM/dd hh:mm:ss'",
		},
		&cli.StringFlag{
			Name:    argEnd,
			Aliases: []string{"e"},
			Usage:   "set schedule end added date in 'yyyy/MM/dd hh:mm:ss'",
		},
	},
	Action: func(ctx *cli.Context) error {

		credentialPath := ctx.String(argCredentialJSONPath)
		tokenJSONPath := ctx.String(argTokenJSONPath)
		calendarName := ctx.String(argCalendarName)
		eventTitle := ctx.String(argTitle)
		eventLocation := ctx.String(argLocation)
		startTimeStr := ctx.String(argStart)
		endTimeStr := ctx.String(argEnd)

		optEnsure := cliopts.NewEnsure().
			With(argCredentialJSONPath, cliopts.StdInteract("Enter credential.json path: ").After(func(s string) error {
				credentialPath = s
				return nil
			})).
			With(argTokenJSONPath, cliopts.StdInteract("Enter token json path: ").After(func(s string) error {
				tokenJSONPath = s
				return nil
			})).
			With(argCalendarName, cliopts.StdInteract("Enter calendar name: ").After(func(s string) error {
				calendarName = s
				return nil
			})).
			With(argTitle, cliopts.StdInteract("Enter event title: ").After(func(s string) error {
				eventTitle = s
				return nil
			})).
			With(argStart, cliopts.StdInteract("Enter event start time in 'yyyy/MM/dd hh:mm:ss': ").After(func(s string) error {
				startTimeStr = s
				return nil
			})).
			With(argEnd, cliopts.StdInteract("Enter event end time in 'yyyy/MM/dd hh:mm:ss': ").After(func(s string) error {
				endTimeStr = s
				return nil
			})).
			With(argLocation, cliopts.StdInteract("Enter event location: ").After(func(s string) error {
				eventLocation = s
				return nil
			}))

		if err := optEnsure.Do(ctx); err != nil {
			return err
		}

		cal, err := googlecalendar.New(tokenJSONPath, credentialPath)
		if err != nil {
			return err
		}

		startTime, err := time.Parse(timeLayout, startTimeStr)
		if err != nil {
			return err
		}

		endTime, err := time.Parse(timeLayout, endTimeStr)
		if err != nil {
			return err
		}

		event := googlecalendar.CreateEvent(eventTitle, eventLocation, startTime, endTime)
		added, err := cal.AddEvents(calendarName, event)
		if err != nil {
			return err
		}

		for _, addedEvent := range added {
			fmt.Println(addedEvent)
		}

		return nil
	},
}
