package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/shouji-kazuo/cliopts"
	"github.com/shouji-kazuo/gocal/pkg/gocal"
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
		&cli.StringFlag{
			Name:  "template",
			Usage: "set event list output template file path.",
			Value: "./daily_event_list.tmpl",
		},
		&cli.BoolFlag{
			Name: "single-event",
			Usage: "set whether show recurring events as single-event." +
				"if recurring event partially removed between \"start\" and \"end\"," +
				"and assing \"true\" to this flag, then partially removed event does not appear.",
			Value: true,
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

		cal, err := gocal.New(tokenJSONPath, credentialPath)
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
		events, err := cal.ListEvents(calendarName, startDate, endDate, ctx.Bool("single-event"))
		if err != nil {
			return err
		}

		yaegiTemplatePath := ctx.String("template")
		src, err := ioutil.ReadFile(yaegiTemplatePath)
		if err != nil {
			return err
		}

		interp := interp.New(interp.Options{})
		interp.Use(stdlib.Symbols)
		interp.Use(map[string]map[string]reflect.Value{
			"github.com/shouji-kazuo/gocal/pkg/gocal": map[string]reflect.Value{
				"GoogleCalendar": reflect.ValueOf((*gocal.GoogleCalendar)(nil)),
				"Event":          reflect.ValueOf((*gocal.Event)(nil)),
			},
		})

		_, err = interp.Eval(string(src))
		if err != nil {
			return err
		}

		v, err := interp.Eval(`tmpl.ListEvents`)
		listEvents := v.Interface().(func([]*gocal.Event) (io.Reader, error))
		reader, err := listEvents(events)
		if err != nil {
			return err
		}

		bufreader := bufio.NewReader(reader)
		for {
			//TODO 手抜き．本当は isPrefix(第2戻り値) == true の時を考慮しないといけない．が，多分すごくめんどくさい．
			line, _, err := bufreader.ReadLine()
			if err == io.EOF {
				fmt.Fprintln(os.Stdout, string(line))
				break
			}
			fmt.Fprintln(os.Stdout, string(line))
		}

		return nil
	},
}
