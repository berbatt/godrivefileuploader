package main

import (
	"godrivefileuploader/file_operations"
	"log"
)

const (
	pathname = "temp"
)

func main() {
	err := file_operations.TraverseThroughDirectoryAndUploadToDrive(pathname)
	if err != nil {
		log.Fatalf("error while uploading the path with name: %s", pathname)
		return
	}
}
