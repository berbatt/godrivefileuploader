package file_uploader

import (
	"bytes"
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
)

const pathCredentialFile = "credentials.json"

type FileUploader interface {
	UploadFileFrom(path string) error
}

type DriveUploader struct {
	client *drive.Service
}

func NewDriveUploader(pathCredentialFile string) (FileUploader, error) {
	b, err := os.ReadFile(pathCredentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	httpClient, err := getClientFromConfig(config)
	if err != nil {
		return nil, err
	}
	c, err := drive.NewService(context.Background(), option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return &DriveUploader{c}, nil
}

func NewDefaultDriveUploader() (FileUploader, error) {
	return NewDriveUploader(pathCredentialFile)
}

func (d *DriveUploader) UploadFileFrom(path string) error {
	// Read the file contents
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new file on Google Drive
	file := &drive.File{Name: path}
	_, err = d.client.Files.Create(file).Media(bytes.NewReader(fileContents)).Do()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
