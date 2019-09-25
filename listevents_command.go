package main

import (
	"html/template"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/shouji-kazuo/cliopts"
	"github.com/shouji-kazuo/gocal/google-cal"
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

		optEnsure := cliopts.NewEnsure().
			With(argCredentialJSONPath, cliopts.StdInteract("Enter credential.json path: ").After(func(s string) error {
				credentialPath = s
				return nil
			})).
			With(argTokenJSONPath, cliopts.StdInteract("Enter start date(ex. \"2006/01/02 15:04:05\") (hour,min,sec is optional): ").After(func(s string) error {
				tokenJSONPath = s
				return nil
			})).
			With("start", cliopts.StdInteract("Enter token json path: ").After(func(s string) error {
				tokenJSONPath = s
				return nil
			})).
			With("calendar-name", cliopts.StdInteract("Enter calendar name: ").After(func(s string) error {
				calendarName = s
				return nil
			}))

		if err := optEnsure.Do(ctx); err != nil {
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
		sort.Slice(events, func(i, j int) bool {
			return events[i].Start.Before(events[j].Start)
		})

		templateDataMap := map[string][]*googlecalendar.Event{
			"Events": events,
		}
		eventOutputTemplate, err := template.New("event_summary.tpl").Funcs(template.FuncMap{
			"formatDate": func(date time.Time, format string) string {
				return date.Format(format)
			},
		}).ParseFiles("./event_summary.tpl")
		
		if err != nil {
			return err
		}

		err = eventOutputTemplate.Execute(os.Stdout, templateDataMap)
		if err != nil {
			return err
		}

		return nil
	},
}