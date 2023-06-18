package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"log"
	"net/http"
	"os"
)

const (
	pathname = "temp"
)

func main() {
	/*
		err := file_operations.TraverseThroughDirectoryAndUploadToDrive(pathname)
		if err != nil {
			log.Fatalf("error while uploading the path with name: %s\nerr: %v", pathname, err)
			return
		}
	*/
	err := authHelper()
	if err != nil {
		return
	}
}

func authHelper() error {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return errors.Wrap(err, "Unable to read client secret file")
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return errors.Wrap(err, "Unable to parse client secret file to config")
	}
	// Create a token storage instance
	//tokenStorage := &TokenStorage{}

	// Load the token from storage
	token := loadToken()
	if token != nil {
		// Check if the token is still valid
		if !token.Valid() {
			// If the token is expired, refresh it
			tokenSource := config.TokenSource(context.Background(), token)
			newToken, err := tokenSource.Token()
			if err != nil {
				log.Fatalf("Unable to refresh token: %v", err)
			}
			token = newToken

			// Save the refreshed token
			storeToken(token)
		}
	}

	// If no token is available, initiate the authorization flow
	if token == nil {
		// Start the authorization flow
		url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		fmt.Printf("Go to the following URL to authorize the application: \n%v\n", url)

		// Set up a simple HTTP server to handle the callback
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				log.Fatal("Authorization code not found")
			}

			// Exchange the authorization code for a token
			token, err := config.Exchange(context.Background(), code)
			if err != nil {
				log.Fatalf("Unable to exchange authorization code: %v", err)
			}

			// Save the token
			storeToken(token)

			fmt.Println("Authorization successful.")
		})

		// Start the HTTP server
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	// Create an OAuth2 HTTP client
	//httpClient := config.Client(context.Background(), token)

	// Use the HTTP client to make authenticated API requests
	// ...
	return nil
}

type TokenStorage struct {
	Token *oauth2.Token
}

func (ts *TokenStorage) SaveToken(token *oauth2.Token) {
	ts.Token = token
}

func (ts *TokenStorage) LoadToken() *oauth2.Token {
	return ts.Token
}

func storeToken(token *oauth2.Token) {
	// Create or open the token file
	file, err := os.Create("token.json")
	if err != nil {
		log.Fatalf("Unable to create token file: %v", err)
	}
	defer file.Close()

	// Serialize the token to JSON
	encoder := json.NewEncoder(file)
	err = encoder.Encode(token)
	if err != nil {
		log.Fatalf("Unable to store token: %v", err)
	}
}

func loadToken() *oauth2.Token {
	// Open the token file
	file, err := os.Open("token.json")
	if err != nil {
		return nil // Token file doesn't exist or cannot be opened
	}
	defer file.Close()

	// Deserialize the token from JSON
	token := &oauth2.Token{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(token)
	if err != nil {
		log.Fatalf("Unable to load token: %v", err)
	}

	return token
}
