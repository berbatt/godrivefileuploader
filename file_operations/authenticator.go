package file_operations

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

type DriveAuthenticator struct {
	Token  *oauth2.Token
	Config *oauth2.Config
}

func (a *DriveAuthenticator) RefreshToken() error {
	if !a.Token.Valid() {
		tokenSource := a.Config.TokenSource(context.Background(), a.Token)
		newToken, err := tokenSource.Token()
		if err != nil {
			return err
		}
		a.Token = newToken
	}
	fmt.Println("Refresh Token saved successfully.")
	return nil
}

func (a *DriveAuthenticator) InitiateAuthenticationFlow() error {
	url := a.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following URL to authorize the application: \n%v\n", url)

	// Set up a simple HTTP server to handle the callback
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			log.Fatal("Authorization code not found")
		}
		// Exchange the authorization code for a token
		token, err := a.Config.Exchange(context.Background(), code)
		if err != nil {
			log.Fatalf("Unable to exchange authorization code: %v", err)
		}
		a.Token = token
		fmt.Println("Authorization successful.")
		return
	})

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
	// TODO close the server after finish
	return nil
}
