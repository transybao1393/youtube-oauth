package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func YoutubeVideoUpload(configByte []byte) {
	//- default params
	filename    := "/Users/transybao/Downloads/cat.mp4" //- upload file path
	title       := ""
	description := ""
	category    := "22"
	keywords    := "video, test"
	privacy     := "unlisted"

	config, err := google.ConfigFromJSON(configByte, 
		youtube.YoutubeReadonlyScope,
		youtube.YoutubeUploadScope,
		youtube.YoutubeScope,
	)

	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	if filename == "" {
			log.Fatalf("You must provide a filename of a video file to upload")
	}

	client := GetClient(context.Background(), config)

	service, err := youtube.New(client)
	if err != nil {
			log.Fatalf("Error creating YouTube client: %v", err)
	}

	upload := &youtube.Video{
			Snippet: &youtube.VideoSnippet{
					Title:       title,
					Description: description,
					CategoryId:  category,
			},
			Status: &youtube.VideoStatus{PrivacyStatus: privacy},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(keywords, "") != "" {
			upload.Snippet.Tags = strings.Split(keywords, ",")
	}
	parts := []string {
		"snippet",
		"status",
	}
	call := service.Videos.Insert(parts, upload)
	

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening %v: %v", filename, err)
	}
	defer file.Close()
	

	response, err := call.Media(file).Do()
	HandleError(err, "")
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
}
