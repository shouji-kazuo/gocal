package main

import (
	"fmt"
	"strings"
	"time"

	reqflags "github.com/shouji-kazuo/cli-reqflags"
	"github.com/shouji-kazuo/gocal-cli-go/google-cal"
	cli "gopkg.in/urfave/cli.v2"
)

var listEventsCommand = &cli.Command{
	Name:        "events",
	Usage:       "",
	Description: "list schedule with duration from google calendar",
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
		&cli.StringFlag{
			Name:    "calendar-name",
			Aliases: []string{"cal", "calendar"},
			Usage:   "set calendar name",
		},
		&cli.StringFlag{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "set start date(ex. \"2006/01/02 15:04:05\") json path to save. hour, min, sec is optional",
		},
		&cli.StringFlag{
			Name:    "end",
			Aliases: []string{"e"},
			Usage:   "set end date(ex. \"2006/01/02 15:04:05\") json path to save. hour, min, sec is optional",
		},
	},
	Action: func(ctx *cli.Context) error {
		credentialPath := ctx.String(argCredentialJSONPath)
		tokenJSONPath := ctx.String(argTokenJSONPath)
		startDateRaw := ctx.String("start")
		endDateRaw := ctx.String("end")
		calendarName := ctx.String("calendar-name")

		onMissingCredentialJSONPath := func() error {
			fmt.Print("Enter credential.json path: ")
			for nScanned, err := fmt.Scan(&credentialPath); nScanned != 1 || err != nil; {
				continue
			}
			return nil
		}
		onMissingCalendarName := func() error {
			fmt.Print("Enter calendar name: ")
			for nScanned, err := fmt.Scan(&calendarName); nScanned != 1 || err != nil; {
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
		onMissingStartDate := func() error {
			fmt.Print("Enter start date(ex. \"2006/01/02 15:04:05\") (hour,min,sec is optional): ")
			for nScanned, err := fmt.Scan(&startDateRaw); nScanned != 1 || err != nil; {
				continue
			}
			return nil
		}
		err := reqflags.Recover(ctx,
			map[string]func() error{
				argCredentialJSONPath: onMissingCredentialJSONPath,
				argTokenJSONPath:      onMissingTokenJSONPath,
				"start":               onMissingStartDate,
				"calendar-name":       onMissingCalendarName,
			})
		if err != nil {
			return err
		}

		cal, err := googlecalendar.New(tokenJSONPath, credentialPath)
		if err != nil {
			return err
		}

		location := time.Now().Location()
		startDate, err := time.ParseInLocation("2006/01/02 15:04:05", startDateRaw, location)
		if err != nil {
			startDate, err = time.ParseInLocation("2006/01/02", strings.TrimSpace(startDateRaw), location)
			if err != nil {
				return err
			}
		}
		endDate := time.Now() //TODO 未来永劫の予定を出す or 今日まで?
		if endDateRaw != "" {
			endDate, err = time.ParseInLocation("2006/01/02 15:04:05", endDateRaw, location)
			if err != nil {
				endDate, err = time.ParseInLocation("2006/01/02", strings.TrimSpace(endDateRaw), location)
				if err != nil {
					return err
				}
			}
		}
		events, err := cal.ListEvents(calendarName, startDate, endDate)
		if err != nil {
			return err
		}
		for _, event := range events {
			fmt.Println(event)
		}

		return nil
	},
}
