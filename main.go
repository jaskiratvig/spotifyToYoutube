package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	spotifyBearerToken = "BQBQLQsQDANZAUBh2P5M3z4VukxEFirvhYuWuefBYtZncy_yC3jBsklvLFTQfaVZBuRmuA5S2IT4moORBJ71K6aRF-g_YSWf8r42CrzVmbSCgW-E68oUHtVgPOPsqfaSzT4OVTAWwYf4xgiSq00pRv57"
	youtubeAccessToken = "ya29.a0AfH6SMD62w8IpSE2EWN6_Z_K4Ld3CL5hYRyqlW1WLxg4IQCmWNHdaPZzkJSHTC7gsmaGvrJ5gpfLrUV93X_YcvtsecAQcXryE_OneaJa4RDdz0g2jYxMAyvq_h5ruRypL2x6wkBNWIf3EEhghv7x0EE8IYTH1FKXGmk"
	youtubeAPIKey      = "AIzaSyCykjV9NVjOHWZ81iXRjz_oN5T9_Gjwe3o"
	youtubeQuery       = flag.Int64("max-results", 1, "Max YouTube results")
	spotifyPlaylistID  = "0vWNFwmqNyrDL44wGC4lR9"
)

/*func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Spotify and Youtube bearer tokens from authorization headers or JSON body or path parameter?
	spotifyPlaylistID := request.PathParameters["spotifyPlaylistID"]
	return events.APIGatewayProxyResponse{Body: youtubePlaylistID, StatusCode: 200}, nil
}*/

/*func main() {
	lambda.Start(Handler)
}*/

func main() {
	songs, title := getSongsAndTitleFromPlaylist(spotifyPlaylistID)

	youtubePlaylistID := createPlaylist(title)
	var youtubeMusicVideos []string

	for song, artist := range songs {
		youtubeMusicVideoID := getMusicVideoFromSong(song, artist)
		youtubeMusicVideos = append(youtubeMusicVideos, youtubeMusicVideoID)
	}

	for _, video := range youtubeMusicVideos {
		addVideoToPlaylist(video, youtubePlaylistID)
	}

	// Output: youtubePlaylistID
}

// Get list of songs from spotify playlist
func getSongsAndTitleFromPlaylist(spotifyplaylistID string) (map[string]string, string) {
	songToArtist := make(map[string]string)

	url := "https://api.spotify.com/v1/playlists/" + spotifyplaylistID

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	bearer := "Bearer " + spotifyBearerToken

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer)

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var dataJSON map[string]interface{}
	err = json.Unmarshal(body, &dataJSON)
	if err != nil {
		log.Fatalln(err)
	}

	tracks := dataJSON["tracks"].(map[string]interface{})
	playlistName := dataJSON["name"].(string)
	items := tracks["items"]

	for _, value := range items.([]interface{}) {
		item := value.(map[string]interface{})
		track := item["track"].(map[string]interface{})
		song := track["name"].(string)

		artists := track["artists"].([]interface{})
		artist := artists[0].(map[string]interface{})
		artistName := artist["name"].(string)

		songToArtist[song] = artistName
	}

	fmt.Println(songToArtist)

	return songToArtist, playlistName
}

// Create a new Youtube playlist
func createPlaylist(title string) string {
	requestBody, err := json.Marshal(map[string]map[string]string{
		"snippet": {
			"title":           title,
			"defaultLanguage": "en",
		},
		"status": {
			"privacyStatus": "public",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	url := "https://www.googleapis.com/youtube/v3/playlists?part=snippet%2Cstatus"

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	var bearer = "Bearer " + youtubeAccessToken

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer)

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var dataJSON map[string]interface{}
	err = json.Unmarshal(body, &dataJSON)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(dataJSON)

	return dataJSON["id"].(string)
}

// Look up music video for each song on Youtube
func getMusicVideoFromSong(song string, artist string) string {
	// Search query = song artist "official music video"
	client := &http.Client{
		Transport: &transport.APIKey{Key: youtubeAPIKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(song + " " + artist + " official music video").
		MaxResults(*youtubeQuery)
	response, err := call.Do()
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(response.Items[0].Id.VideoId)

	return response.Items[0].Id.VideoId
}

// Populate playlist with music videos
func addVideoToPlaylist(video string, youtubePlaylistID string) {
	// Use youtubePlaylistID
	requestBody, err := json.Marshal(map[string]map[string]interface{}{
		"snippet": {
			"playlistId": youtubePlaylistID,
			"position":   0,
			"resourceId": map[string]string{
				"kind":    "youtube#video",
				"videoId": video,
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	url := "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet"

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	var bearer = "Bearer " + youtubeAccessToken

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bearer)

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var dataJSON map[string]interface{}
	err = json.Unmarshal(body, &dataJSON)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(dataJSON)
}
