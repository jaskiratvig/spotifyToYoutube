# SpotifyToYoutube App

This project uses the Spotify API and YouTube API to convert a Spotify songs playlist into a Youtube music videos playlist.

## Dependancies

The following dependancies are required to run the project:

* Install Golang at ``` https://golang.org/doc/install ```
* Create a Spotify OAuth bearer token:
  * Visit ``` https://developer.spotify.com/console/get-current-user-top-artists-and-tracks/ ```
  * Under OAuth token, click the green box that says "Get Token"
  * Under Required scopes for this endpoint, tick "user-top-read" and click "Request Token"
  * Copy the generated OAuth token and set the "spotifyBearerToken" variable in main.go
* Generate YouTube API Key and YouTube access token
  * Visit ``` https://console.developers.google.com/apis/credentials ```
  * Create an API Key and set the "youtubeAPIKey" variable in main.go
* Generate YouTube access token
  * Visit ``` https://developers.google.com/oauthplayground/ ```
  * Select the YouTube Data API v3 scope and click the blue button that says "Authorize APIs"
  * Authenticate with your Google credentials
  * Click the blue botton that says "Exchange authorization code for tokens"
  * Copy the "access_token" in the JSON request response and set the "youtubeAccessToken" variable in main.go
* Pick a playlist from spotify and set the "spotifyPlaylistID" variable in main.go

## Run Project

To run this project, navigate to the spotifyToYoutube root directory and run ``` go run main.go ```. A YouTube playlist should be populated with the music videos from your spotify playlist.
