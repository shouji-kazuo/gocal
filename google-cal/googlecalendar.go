package googlecalendar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

type GoogleCalendar struct {
	client *calendar.Service
}

func Auth(credentialJSONPath string, fin io.Reader, fout io.Writer) (*oauth2.Token, error) {
	b, err := ioutil.ReadFile(credentialJSONPath)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read client secret file from path: "+credentialJSONPath)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse client secret file to config")
	}

	token, err := getTokenFromWeb(config, fin, fout)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get token from web")
	}

	return token, nil
}

func SaveToken(oauth2Token *oauth2.Token, destFile io.Writer) error {
	return json.NewEncoder(destFile).Encode(oauth2Token)
}

func LoadToken(tokenJSONPath string) (*oauth2.Token, error) {
	var token oauth2.Token
	var bytes []byte
	var err error
	if bytes, err = ioutil.ReadFile(tokenJSONPath); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func New(tokenJSONPath string, credentialJSONPath string) (*GoogleCalendar, error) {
	// memo:
	//	*oauth2.Config, *oauth2.Token を渡す方式ではダメなのか
	//		→ユーザ側がその方法を知ってなきゃダメ
	//		→googlecalendar.Auth()で生成していることから，「oauth2関連の操作はgooglecalendarパッケージに閉じ込める」
	//		のならば，ユーザ側がoauth2関連の操作を知っていることはおかしい．
	b, err := ioutil.ReadFile(credentialJSONPath)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read client secret file from path: "+credentialJSONPath)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create oauth2.Config from credential json file.")
	}
	client, err := getClient(config, tokenJSONPath)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create http.Client from saved token file.")
	}
	calendar, err := calendar.New(client)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create calendar.Service from http.Client.")
	}
	return &GoogleCalendar{
		client: calendar,
	}, nil
}

func (cal *GoogleCalendar) ListCalendars() (*calendar.CalendarList, error) {
	return cal.client.CalendarList.List().Do()
}

func (cal *GoogleCalendar) ListEvents(calID string, startDate time.Time, endTime time.Time) ([]*calendar.Event, error) {
	startDateRFC3339 := startDate.Format(time.RFC3339)
	endDateRFC3339 := endTime.Format(time.RFC3339)
	events, err := cal.client.Events.List(calID).TimeMin(startDateRFC3339).TimeMax(endDateRFC3339).Do()
	if err != nil {
		return nil, err
	}
	return events.Items, nil
}

func getTokenFromWeb(config *oauth2.Config, fin io.Reader, fout io.Writer) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Fprintf(fout, "Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Fscan(fin, &authCode); err != nil {
		return nil, err
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokenPath string) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
