package main

import (
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/robfig/cron"
	"godrivefileuploader/authentication"
	"godrivefileuploader/file_operations"
	"godrivefileuploader/file_uploader"
	"os"
	"os/signal"
	"syscall"
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

var ErrDuration = errors.New("must enter a valid period such as '1h', '2m', '3d'")

type options struct {
	Period          string `short:"P" long:"period" description:"Period of uploading files, such as '1m', '2h'. Default is 1 hour" default:"1h"`
	Path            string `short:"p" long:"path" description:"Absolute path of the directory"`
	Duration        string `short:"d" long:"duration" description:"Total duration of the uploader" default:"1h"`
	CredentialsPath string `short:"c" long:"credentials" description:"Path to the credentials file" default:"credentials.json"`
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
		return exitCodeErr, ErrDuration
	}
	authentication.PathCredentialsFile = opts.CredentialsPath
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
	// Channel to listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	duration, err := time.ParseDuration(opts.Duration)
	if err != nil {
		return exitCodeErr, ErrDuration
	}
	// Channel to signal after the given duration
	timeChan := time.After(duration)
	// Block the main function until a signal or given duration
	select {
	case <-timeChan:
		fmt.Println("Main goroutine stopping...")
	case sig := <-sigChan:
		fmt.Printf("Main goroutine received signal: %s. Exiting...\n", sig)
	}
	// Stop the cron scheduler
	c.Stop()
	fmt.Println("Cron scheduler stopped. Exiting main goroutine.")
	return exitCodeOK, nil
}
