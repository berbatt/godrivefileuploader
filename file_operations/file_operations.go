package file_operations

import (
	"godrivefileuploader/file_uploader"
	"godrivefileuploader/path_manager"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

var pathManager = path_manager.NewManager()

func ReadFileFromPath(path string) ([]byte, error) {
	// Read the file contents
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return fileContents, err
}

func TraverseThroughDirectoryAndUploadToDrive(dirPath string) error {
	return filepath.Walk(dirPath, HandleFileOrFolder)
}

func HandleFileOrFolder(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	fileUploader, err := file_uploader.GetUploader()
	if err != nil {
		return err
	}
	// Check if the current path points to a regular file
	parentFolderID := pathManager.GetParentFolderID(path)
	if isDir := info.IsDir(); isDir {
		var folderID string
		folderName := path_manager.GetFileOrFolderNameFromPath(path)
		if folderID, err = fileUploader.CreateOrUpdateFolder(folderName, parentFolderID); err != nil {
			return err
		}
		pathManager.SetParentID(path, folderID)
	} else if !isDir {
		input, err := ReadFileFromPath(path)
		if err != nil {
			return err
		}
		return fileUploader.CreateOrUpdateFile(input, info.Name(), parentFolderID)
	}
	return nil
}
