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

var server *http.Server

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

	// Create a channel to receive an interrupt or termination signal
	stop := make(chan struct{}, 1)

	// Set up a simple HTTP server to handle the callback
	server = &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "Authorization code not found", http.StatusBadRequest)
				return
			}
			// Exchange the authorization code for a token
			token, err := a.Config.Exchange(context.Background(), code)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to exchange authorization code: %v", err), http.StatusInternalServerError)
				return
			}
			// Save the token
			a.Token = token
			_, err = w.Write([]byte("Authorization successful. You can close this window."))
			if err != nil {
				return
			}

			stop <- struct{}{}
		}),
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		} else {
			log.Println("Server gracefully shutdown")
		}
	}()

	// Wait for a signal to stop the server
	<-stop

	return nil
}
