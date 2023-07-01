package main

import (
	"github.com/pkg/errors"
	"godrivefileuploader/file_operations"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"log"
)

const (
	pathCredentialsFile = "credentials.json"
	pathname            = "temp"
)

func main() {
	if err := handleAuthenticationFlow(); err != nil {
		log.Fatalf("error while authentication flow, err: %v", err)
		return
	}

	err := file_operations.TraverseThroughDirectoryAndUploadToDrive(pathname)
	if err != nil {
		log.Fatalf("error while uploading the path with name: %s\nerr: %v", pathname, err)
		return
	}

}

func handleAuthenticationFlow() error {
	credentialFile, err := file_operations.ReadFileFromPath(pathCredentialsFile)
	if err != nil {
		return errors.Wrap(err, "Unable to read client secret file")
	}

	config, err := google.ConfigFromJSON(credentialFile, drive.DriveFileScope)
	if err != nil {
		return errors.Wrap(err, "Unable to parse client secret file to config")
	}

	tokenStorage := file_operations.TokenStorage{}
	token, err := tokenStorage.LoadToken()
	if err != nil {
		return err
	}

	driveAuthenticator := file_operations.DriveAuthenticator{
		Token:  token,
		Config: config,
	}
	if token != nil {
		if err = driveAuthenticator.RefreshToken(); err != nil {
			return err
		}
	} else if token == nil {
		if err = driveAuthenticator.InitiateAuthenticationFlow(); err != nil {
			return err
		}
	}
	if err = tokenStorage.SaveToken(driveAuthenticator.Token); err != nil {
		return err
	}
	return nil
}
