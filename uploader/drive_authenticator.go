package uploader

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	pathTokenFile = "token.json"
)

var (
	authenticator       *DriveAuthenticator
	PathCredentialsFile string
)

func Get() (*DriveAuthenticator, error) {
	if authenticator == nil {
		a := NewDriveAuthenticator()
		err := a.ExecuteFlow(PathCredentialsFile, pathTokenFile)
		if err != nil {
			return nil, err
		}
		authenticator = &a
	}
	return authenticator, nil
}

type DriveAuthenticator struct {
	tokenStorage TokenStorage
	config       *oauth2.Config
}

func NewDriveAuthenticator() DriveAuthenticator {
	return DriveAuthenticator{}
}

var server *http.Server

func (a *DriveAuthenticator) ExecuteFlow(pathCredentialsFile, pathTokenFile string) (err error) {
	var credentialFile []byte
	credentialFile, err = ReadFileFromPath(pathCredentialsFile)
	if err != nil {
		return errors.Wrap(err, "Unable to read client secret file")
	}
	a.config, err = google.ConfigFromJSON(credentialFile, drive.DriveFileScope)
	if err != nil {
		return errors.Wrap(err, "Unable to parse client secret file to config")
	}

	a.tokenStorage.Token, err = a.tokenStorage.loadToken(pathTokenFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		err = a.exchangeAuthorizationToken()
		if err != nil {
			return err
		}
	}
	err = a.refreshToken()
	if err != nil {
		return errors.Wrap(err, "Unable to refresh token")
	}
	err = a.tokenStorage.saveToken(pathTokenFile)
	if err != nil {
		return err
	}
	return nil
}

func (a *DriveAuthenticator) GetToken() *oauth2.Token {
	return a.tokenStorage.Token
}

func (a *DriveAuthenticator) GetConfig() *oauth2.Config {
	return a.config
}

func (a *DriveAuthenticator) GetDriveClient() *http.Client {
	return a.config.Client(context.Background(), a.tokenStorage.Token)
}

func (a *DriveAuthenticator) refreshToken() error {
	if !a.tokenStorage.Token.Valid() {
		tokenSource := a.config.TokenSource(context.Background(), a.tokenStorage.Token)
		newToken, err := tokenSource.Token()
		if err != nil {
			return err
		}
		a.tokenStorage.Token = newToken
	}
	fmt.Println("Refresh Token saved successfully.")
	return nil
}

func (a *DriveAuthenticator) exchangeAuthorizationToken() error {
	url := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following URL to authorize the application: \n%v\n", url)

	// Create a channel to receive an interrupt or termination signal
	stop := make(chan struct{}, 1)

	// Set up a simple HTTP server to handle the callback
	server = &http.Server{
		ReadTimeout: 400 * time.Millisecond,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "Authorization code not found", http.StatusBadRequest)
				return
			}
			// Exchange the authorization code for a token
			token, err := a.config.Exchange(context.Background(), code)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to exchange authorization code: %v", err), http.StatusInternalServerError)
				return
			}
			// Save the token
			a.tokenStorage.Token = token
			_, err = w.Write([]byte("Authorization successful. You can close this window."))
			if err != nil {
				return
			}

			stop <- struct{}{}
		}),
	}

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		log.Println("Server gracefully shutdown")
	}()

	// Wait for a signal to stop the server
	<-stop

	return nil
}
