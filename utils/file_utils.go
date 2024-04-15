package utils

import (
	"bufio"
	"errors"
	"io"
	"os"
)

func ReadFileFromPath(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	reader := bufio.NewReader(f)
	fileContents := make([]byte, 1024)
	var size int
	size, err = reader.Read(fileContents)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return fileContents[:size], err
}
