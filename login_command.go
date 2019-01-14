package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shouji-kazuo/gocal-cli-go/cliutil"

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
			Name:    argCredentialJSONPath,
			Aliases: []string{"credential", "cred", "c"},
			Value:   "",
			Usage:   "set 'credentials.json' path",
		},
		&cli.StringFlag{
			Name:    argTokenJSONPath,
			Aliases: []string{"o"},
			Value:   "",
			Usage:   "set token json path to save",
		},
	},
	Action: func(ctx *cli.Context) error {
		jsonPaths, err := cliutil.GetJSONPaths(ctx, defaultContextArgKeys)
		if err != nil {
			return errors.Wrap(err, "cannot get some json path.")
		}
		credentialJSONPath := jsonPaths.CredentialJSONPath
		tokenJSONPath := jsonPaths.TokenJSONPath
		b, err := ioutil.ReadFile(credentialJSONPath)
		if err != nil {
			return errors.Wrap(err, "Unable to read client secret file from path: "+credentialJSONPath)
		}

		// If modifying these scopes, delete your previously saved token.json.
		config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
		if err != nil {
			return errors.Wrap(err, "Unable to parse client secret file to config")
		}

		token, err := getTokenFromWeb(config)
		if err != nil {
			return errors.Wrap(err, "Unable to get token from web")
		}

		if err = saveToken(tokenJSONPath, token); err != nil {
			return errors.Wrap(err, "Unable to save token to path: "+tokenJSONPath)
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
	dirPath := filepath.Dir(path)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.Mkdir(dirPath, 0700); err != nil {
			return errors.Wrap(err, "cannot create directory to save token.")
		}
	}
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return errors.Wrap(err, "something wrong during get directory stat.")
	}
	if dirInfo.Mode().Perm() != 0700 {
		if err = os.Chmod(dirPath, 0700); err != nil {
			return errors.Wrap(err, "cannot change directory permission")
		}
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}
