package file_uploader

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"os"
)

const pathCredentialFile = "credentials.json"
const folderColorCode = "#00FF00"

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
	CreateFile(input []byte, fileName string, parentID string) error
	CreateFolder(folderName, parentFolderID string) (string, error)
	FindFolderOrFile(name, parentFolderID string) (*drive.File, error)
	UpdateFile(input []byte, fileName string) error
	CreateOrUpdateFile(input []byte, name, parentFolderID string) error
	CreateOrUpdateFolder(name, parentFolderID string) (string, error)
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

func (d *DriveUploader) CreateFile(input []byte, fileName string, parentID string) error {
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
		FolderColorRgb: folderColorCode,
	}

	createdFolder, err := d.client.Files.Create(folder).Do()
	if err != nil {
		return "", errors.Wrapf(err, "unable to create folder with name: %s error: %v", folderName, err)
	}

	return createdFolder.Id, nil
}

func (d *DriveUploader) FindFolderOrFile(name, parentFolderID string) (*drive.File, error) {
	var query string
	if len(parentFolderID) == 0 {
		query = fmt.Sprintf("name = '%s' and trashed = false", name)
	} else if len(parentFolderID) > 0 {
		query = fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", name, parentFolderID)
	}
	files, err := d.client.Files.List().Q(query).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to search for folder or file: %w", err)
	}

	if len(files.Files) > 0 {
		return files.Files[0], nil
	}

	return nil, nil
}

func (d *DriveUploader) UpdateFile(input []byte, fileID string) error {
	// Create a media reader with the new content
	reader := bytes.NewReader(input)

	// Set the media metadata for the update
	update := d.client.Files.Update(fileID, nil)
	update.Media(reader, googleapi.ContentType("text/plain"))

	// Perform the update
	_, err := update.Do()
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	return nil
}

func (d *DriveUploader) CreateOrUpdateFile(input []byte, name, parentFolderID string) error {
	file, err := d.FindFolderOrFile(name, parentFolderID)
	if err != nil {
		return err
	}

	if file == nil {
		// Folder or file not found, create it
		if err = d.CreateFile(input, name, parentFolderID); err != nil {
			return err
		}
	} else {
		if err = d.UpdateFile(input, file.Id); err != nil {
			return err
		}
	}
	return nil
}

func (d *DriveUploader) CreateOrUpdateFolder(name, parentFolderID string) (string, error) {
	folder, err := d.FindFolderOrFile(name, parentFolderID)
	if err != nil {
		return "", err
	}

	if folder == nil {
		var folderID string
		if folderID, err = d.CreateFolder(name, parentFolderID); err != nil {
			return "", err
		}
		return folderID, nil
	}

	return folder.Id, nil
}
