package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/robfig/cron"
	"godrivefileuploader/file_operations"
	"godrivefileuploader/file_uploader"
	"os"
	"time"
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

var opts options

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

	_, err = time.ParseDuration(opts.Period)
	if err != nil {
		return exitCodeErr, errors.New("must enter a valid period such as '1h', '2m', '3d'")
	}

	_, err = file_uploader.GetUploader()
	if err != nil {
		return exitCodeErr, err
	}

	c := cron.New()
	// Schedule the function to run periodically
	err = c.AddFunc(fmt.Sprintf("@every %s", opts.Period), func() {
		fmt.Printf("Uploading files...\n")
		err = file_operations.TraverseThroughDirectoryAndUploadToDrive(opts.Path)
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		return exitCodeErr, err
	}
	// Start the cron scheduler
	c.Start()

	select {}
}
