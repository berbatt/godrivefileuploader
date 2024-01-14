package file_operations

import (
	"godrivefileuploader/file_uploader"
	"godrivefileuploader/path_manager"
	"godrivefileuploader/utils"
	"io/fs"
	"path/filepath"
)

var pathManager = path_manager.NewManager()

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
	var isDirectory, isFile = info.IsDir(), !info.IsDir()
	if isDirectory {
		var folderID string
		folderName := path_manager.GetFileOrFolderNameFromPath(path)
		if folderID, err = fileUploader.CreateOrUpdateFolder(folderName, parentFolderID); err != nil {
			return err
		}
		pathManager.SetParentID(path, folderID)
	} else if isFile {
		var input []byte
		input, err = utils.ReadFileFromPath(path)
		if err != nil {
			return err
		}
		return fileUploader.CreateOrUpdateFile(input, info.Name(), parentFolderID)
	}
	return nil
}
