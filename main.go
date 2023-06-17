package main

import (
	"fmt"
	"godrivefileuploader/file_uploader"
	"log"
)

const (
	filename = "temp.txt"
)

func main() {
	fileUploader, err := file_uploader.NewDefaultDriveUploader()
	if err != nil {
		log.Fatal("error while creating new default file uploader")
		return
	}
	if err = fileUploader.UploadFileFrom(filename); err != nil {
		log.Fatalf("error while uploading the file with name: %s", filename)
		return
	}

	fmt.Printf("File %s uploaded to Google Drive\n", filename)
}
