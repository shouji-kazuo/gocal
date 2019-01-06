package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	cli "gopkg.in/urfave/cli.v2"
)

var loginCommand = &cli.Command{
	Name:        "login",
	Usage:       "",
	Description: "login to google calendar",
	ArgsUsage:   "",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "credential-json-path",
			Aliases: []string{"credential", "json", "c"},
			Usage:   "set 'credentials.json' path",
		},
		&cli.StringFlag{
			Name:    "saved-token-path",
			Aliases: []string{"o"},
			Value:   "",
			Usage:   "set saved token path",
		},
	},
	Action: func(ctx *cli.Context) error {
		if !ctx.IsSet("credential-json-path") {
			return errors.New("credential-json-path flag is not set.")
		}

		credentialJSONPath := ctx.String("credential-json-path")
		savedTokenPath := ctx.String("saved-token-path")

		b, err := ioutil.ReadFile(credentialJSONPath)
		if err != nil {
			message := "Unable to read client secret file from path: " + credentialJSONPath
			log.Fatal(message)
			return errors.Wrap(err, message)
		}

		// If modifying these scopes, delete your previously saved token.json.
		config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
		if err != nil {
			message := "Unable to parse client secret file to config"
			log.Fatal(message)
			return errors.Wrap(err, message)
		}

		token, err := getTokenFromWeb(config)
		if err != nil {
			message := "Unable to get token from web"
			log.Fatal(message)
			return errors.Wrap(err, message)
		}

		if err = saveToken(savedTokenPath, token); err != nil {
			message := "Unable to save token to path: " + savedTokenPath
			log.Fatal(message)
			return errors.Wrap(err, message)
		}

		return nil
	},
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}
