package file_uploader

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"os"
)

const pathCredentialFile = "credentials.json"

var uploader *FileUploader

func GetUploader() (FileUploader, error) {
	if uploader == nil {
		u, err := NewDefaultDriveUploader()
		if err != nil {
			return nil, err
		}
		uploader = &u
	}
	return *uploader, nil
}

type FileUploader interface {
	UploadFileFrom(input []byte, fileName string, parentID string) error
	CreateFolder(folderName, parentFolderID string) (string, error)
}

type DriveUploader struct {
	client *drive.Service
}

func NewDriveUploader(pathCredentialFile string) (FileUploader, error) {
	b, err := os.ReadFile(pathCredentialFile)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read client secret file")
	}
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse client secret file to config")
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

func (d *DriveUploader) UploadFileFrom(input []byte, fileName string, parentID string) error {
	// Create a new file on Google Drive
	file := &drive.File{
		Name:    fileName,
		Parents: []string{parentID},
	}
	_, err := d.client.Files.Create(file).Media(bytes.NewReader(input)).Do()
	if err != nil {
		return err
	}
	return nil
}

// CreateFolder creates a folder with the given name in the specified parent folder (if any) and returns its ID.
func (d *DriveUploader) CreateFolder(folderName, parentFolderID string) (string, error) {
	var parents []string
	if len(parentFolderID) > 0 {
		parents = append(parents, parentFolderID)
	}
	folder := &drive.File{
		Name:           folderName,
		MimeType:       "application/vnd.google-apps.folder",
		Parents:        parents,
		FolderColorRgb: "#00FF00",
	}

	createdFolder, err := d.client.Files.Create(folder).Do()
	if err != nil {
		return "", errors.Wrapf(err, "unable to create folder with name: %s error: %v", folderName, err)
	}

	return createdFolder.Id, nil
}
