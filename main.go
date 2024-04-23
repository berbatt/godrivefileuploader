package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"godrivefileuploader/file_operations"
	"log"
	"os"
)

const (
	appName = "gofileuploader"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErr
)

type options struct {
	Period string `short:"P" long:"period" description:"Period of uploading files, such as '1h', '2m', '3d'. Default is 1 hour" default:"1h"`
	Path   string `short:"p" long:"path" description:"Absolute path of the directory"`
}

func main() {
	code, err := run(os.Args[1:])
	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			os.Exit(int(exitCodeErr))
		}
	}
	os.Exit(int(code))
}

func run(args []string) (exitCode, error) {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = "[OPTIONS] QUERY..."
	_, err := parser.ParseArgs(args)
	if err != nil {
		if flags.WroteHelp(err) {
			return exitCodeOK, nil
		}
		return exitCodeErr, fmt.Errorf("argument parsing failed: %w", err)
	}
	if len(opts.Path) == 0 {
		return exitCodeErr, errors.New("must enter a path")
	}

	err = file_operations.TraverseThroughDirectoryAndUploadToDrive(opts.Path)
	if err != nil {
		log.Fatalf("error while uploading the path with name: %s\nerr: %v", opts.Path, err)
		return exitCodeErr, err
	}

	return exitCodeOK, nil
}
