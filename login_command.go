package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shouji-kazuo/gocal/cliutil"
	"github.com/shouji-kazuo/gocal/pkg/gocal"

	"github.com/pkg/errors"

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

		if !ctx.IsSet(argCredentialJSONPath) {
			return errors.New("Unable to parse \"credential\" flag.\n" +
				"1.please create google cloud project from https://console.developers.google.com \n" +
				"2.enable Google Calendar API from https://console.developers.google.com/apis/library \n" +
				"3.get OAuth client ID as JSON file from https://console.developers.google.com/apis/credentials \n" +
				"4.specify credential file path to \"credential\" flag.")
		}
		oauth2Token, err := gocal.Auth(ctx.String(argCredentialJSONPath), os.Stdin, os.Stdout)
		if err != nil {
			return errors.Wrap(err, "Unable to authorizate.")
		}

		// tokenファイルを保存するパスの選定．
		// まず引数 -c に指定されたパスを試す
		// →それがダメなら，デフォルトパス($HOMEDIR/.gocal/)の中を試す
		// →それがダメなら，カレントディレクトリへの保存を試す
		// →それもダメなら，諦める
		var tokenFile *os.File = nil
		tryOpenInArgPath := func() error {
			if ctx.IsSet(argTokenJSONPath) {
				if tokenFile, err = os.OpenFile(ctx.String(argTokenJSONPath), os.O_RDWR|os.O_CREATE, 0600); err != nil {
					return err
				}
				return nil
			}
			return errors.New("There is no flag to open file to save token JSON file.")
		}
		tryOpenInDefaultPath := func() error {
			fmt.Fprintln(os.Stderr, "Unable to open file to save oauth JSON in argument.")
			fmt.Fprintln(os.Stderr, "Try to open file in default path...")
			defaultTokenPath, err := cliutil.GetDefaultTokenPathToSave()
			if err != nil {
				return err
			}
			if tokenFile, err = os.OpenFile(defaultTokenPath, os.O_RDWR|os.O_CREATE, 0600); err != nil {
				return err
			}
			return nil
		}
		tryOpenInCurrentDir := func() error {
			fmt.Fprintln(os.Stderr, "Unable to opne file to save oauth JSON in default path.")
			fmt.Fprintln(os.Stderr, "Try to open file in current directory...")
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				return err
			}
			if tokenFile, err = os.OpenFile(filepath.Join(dir, "token.json"), os.O_RDWR|os.O_CREATE, 0600); err != nil {
				return err
			}
			return nil
		}

		if err = tryFuncSeq([]func() error{tryOpenInArgPath, tryOpenInDefaultPath, tryOpenInCurrentDir}); err != nil {
			return errors.Wrap(err, "Unable to open all candidates of path to save token.")
		}
		defer func() {
			if tokenFile != nil {
				tokenFile.Close()
			}
		}()

		if err = gocal.SaveToken(oauth2Token, tokenFile); err != nil {
			return errors.Wrap(err, "Unable to save oauth2 token.")
		}
		return nil
	},
}

func tryFuncSeq(funcs []func() error) error {
	var errorStack error = nil
	for _, f := range funcs {
		if err := f(); err != nil {
			errorStack = errors.Wrap(err, err.Error())
			continue
		}
		return nil
	}
	return errorStack
}
