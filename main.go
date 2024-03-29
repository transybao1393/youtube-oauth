// Sample Go code for user authorization

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"youtube-oauth/services"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var config *oauth2.Config
var client *http.Client

func youtubeOAuth() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("./client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err = google.ConfigFromJSON(b, 
		youtube.YoutubeReadonlyScope,
		youtube.YoutubeUploadScope,
		youtube.YoutubeScope,
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client = services.GetClient(ctx, config)
	service, err := youtube.New(client)

	services.HandleError(err, "Error creating YouTube client")

	services.ChannelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")
}

func main() {

	http.HandleFunc("/oauth", func (w http.ResponseWriter, 
		r *http.Request) {
		fmt.Println("OAuth login start...")
		// when receive the request, print the greeting meassage
		youtubeOAuth()
	
	})

	http.HandleFunc("/auth/callback", func (w http.ResponseWriter, 
		r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")
		fmt.Printf("\n>>> CODE %s\n", code)
		defer r.Body.Close()

		//- code exchange for token
		// conf := &oauth2.Config{
		// 	ClientID: "",
		// 	ClientSecret: "",
		// 	RedirectURL: "", //- can get from config.json
		// 	Scopes:       []string{youtube.YoutubeReadonlyScope},
		// 	Endpoint: oauth2.Endpoint{
		// 		AuthURL:  "", //- can get from config.json
		// 		TokenURL: "", //- can get from config.json
		// 	},
		// }

		tok, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Printf("Unable to retrieve token from web %v", err)	
		}
		
		//- convert to json
		b, err := json.Marshal(tok)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("tokens %s\n", string(b))
	})

	http.HandleFunc("/video/upload", func (w http.ResponseWriter, 
		r *http.Request) {

		b, err := ioutil.ReadFile("./client_secret.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}

		services.YoutubeVideoUpload(b)
	})

	// print out the server is going to start and show the time
	log.Printf("Starting server on port %d....", 8000)
	
	// create server at localhost:8080 and using tcp as the network
	listener, err := net.Listen("tcp", ":8000")
	
	// if recieve error, record it and exit the program
	if err != nil {
		log.Fatal(err)
	}
	
	// setup HTTP connection for the listener of the server
	http.Serve(listener, nil)
}

