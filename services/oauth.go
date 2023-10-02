package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func GetClient(ctx context.Context, config *oauth2.Config) *http.Client {
	fmt.Println("here at GetClient() function")
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
		fmt.Printf("Unable to get path to cached credential file. %v", err)
	}
	fmt.Printf("\ncache file %s\n", cacheFile)
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		fmt.Println(">>> Load token from Cache file  failed")
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	fmt.Println("here at getTokenFromWeb() function")
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
		fmt.Printf("Unable to read authorization code %v", err)
	}
	fmt.Printf("code from web %s", code)

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
		fmt.Printf("Unable to retrieve token from web %v", err)
		
	}
	fmt.Printf("token from web %+v", tok)
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func HandleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message + ": %v", err.Error())
	}
}

func ChannelsListByUsername(service *youtube.Service, part string, forUsername string) {
	var parts []string
	parts = append(parts, part)
	call := service.Channels.List(parts)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	HandleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', " +
				"and it has %d views.",
				response.Items[0].Id,
				response.Items[0].Snippet.Title,
				response.Items[0].Statistics.ViewCount))
}