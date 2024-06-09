package uploader

import (
	"bytes"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type DriveService interface {
	Create(file *drive.File, input []byte) (*drive.File, error)
	Update(fileID string, input []byte) (*drive.File, error)
	List(name, parentFolderID string) (*drive.FileList, error)
}

type driveService struct {
	client *drive.Service
}

func NewDriveService(service *drive.Service) DriveService {
	return &driveService{client: service}
}

func (d *driveService) Create(file *drive.File, input []byte) (*drive.File, error) {
	if input == nil {
		// Create folder
		return d.client.Files.Create(file).Do()
	}
	return d.client.Files.Create(file).Media(bytes.NewReader(input)).Do()
}

func (d *driveService) Update(fileID string, input []byte) (*drive.File, error) {
	// Create a media reader with the new content
	reader := bytes.NewReader(input)
	return d.client.Files.Update(fileID, nil).Media(reader, googleapi.ContentType("text/plain")).Do()
}

func (d *driveService) List(name, parentFolderID string) (*drive.FileList, error) {
	var query string
	if len(parentFolderID) == 0 {
		query = fmt.Sprintf("name = '%s' and trashed = false", name)
	} else if len(parentFolderID) > 0 {
		query = fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", name, parentFolderID)
	}
	return d.client.Files.List().Q(query).Do()
}
