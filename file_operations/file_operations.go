package file_operations

import (
	"godrivefileuploader/file_uploader"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var pathToParentIDMap = map[string]string{}

func ReadFileFromPath(path string) ([]byte, error) {
	// Read the file contents
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return fileContents, err
}

func TraverseThroughDirectoryAndUploadToDrive(dirPath string) error {
	return filepath.Walk(dirPath, WalkFunction)
}

func WalkFunction(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	fileUploader, err := file_uploader.GetUploader()
	if err != nil {
		return err
	}
	// Check if the current path points to a regular file
	if isDir := info.IsDir(); isDir {
		folderID, err := fileUploader.CreateFolder(path, pathToParentIDMap[path])
		pathToParentIDMap[path] = folderID
		if err != nil {
			return err
		}
	} else if !isDir {
		input, err := ReadFileFromPath(path)
		if err != nil {
			return err
		}
		return fileUploader.UploadFileFrom(input, info.Name(), pathToParentIDMap[getRootFolderName(path)])
	}

	return nil
}

func isRootFolder(path string) bool {
	return strings.Contains(path, "/")
}

func getRootFolderName(path string) (root string) {
	paths := strings.Split(path, "/")
	root = paths[0]
	return root
}
