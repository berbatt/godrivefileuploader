package uploader

import (
	"io/fs"
	"path/filepath"
)

var pathManager = NewManager()

func TraverseThroughDirectoryAndUploadToDrive(dirPath string) error {
	return filepath.Walk(dirPath, HandleFileOrFolder)
}

func HandleFileOrFolder(path string, info fs.FileInfo, _ error) (err error) {
	var fileUploader FileUploader
	fileUploader, err = GetUploader()
	if err != nil {
		return err
	}
	// Check if the current path points to a regular file
	parentFolderID := pathManager.GetParentFolderID(path)
	var isDirectory, isFile = info.IsDir(), !info.IsDir()
	if isDirectory {
		var folderID string
		folderName := GetFileOrFolderNameFromPath(path)
		if folderID, err = fileUploader.CreateOrUpdateFolder(folderName, parentFolderID); err != nil {
			return err
		}
		pathManager.SetParentID(path, folderID)
	} else if isFile {
		var input []byte
		input, err = ReadFileFromPath(path)
		if err != nil {
			return err
		}
		return fileUploader.CreateOrUpdateFile(input, info.Name(), parentFolderID)
	}
	return nil
}
