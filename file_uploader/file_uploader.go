package file_uploader

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"godrivefileuploader/authentication"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"net/http"
)

const folderColorCode = "#00FF00"

var uploader *FileUploader

func GetUploader() (FileUploader, error) {
	if uploader == nil {
		u, err := NewDriveUploader()
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
	service DriveService
}

func NewDriveUploader() (FileUploader, error) {
	authenticator, err := authentication.Get()
	if err != nil {
		return nil, err
	}
	var driveClient *http.Client
	driveClient = authenticator.GetDriveClient()
	if err != nil {
		return nil, err
	}
	var s *drive.Service
	s, err = drive.NewService(context.Background(), option.WithHTTPClient(driveClient))
	if err != nil {
		return nil, err
	}
	return &DriveUploader{service: NewDriveService(s)}, nil
}

func (d *DriveUploader) CreateFile(input []byte, fileName string, parentID string) error {
	// Create a new file on Google Drive
	file := &drive.File{
		Name:    fileName,
		Parents: []string{parentID},
	}
	_, err := d.service.Create(file, input)
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

	createdFolder, err := d.service.Create(folder, nil)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create folder with name: %s error: %v", folderName, err)
	}

	return createdFolder.Id, nil
}

func (d *DriveUploader) FindFolderOrFile(name, parentFolderID string) (*drive.File, error) {
	files, err := d.service.List(name, parentFolderID)
	if err != nil {
		return nil, fmt.Errorf("unable to search for folder or file: %w", err)
	}

	if len(files.Files) > 0 {
		return files.Files[0], nil
	}

	return nil, nil
}

func (d *DriveUploader) UpdateFile(input []byte, fileID string) error {
	_, err := d.service.Update(fileID, input)
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
