package main

import (
	"time"
	"flag"
	"os"
	"path"
	"os/user"
	"fmt"
	"errors"
	"log"
	"io/ioutil"
)

type Config struct {
	location *time.Location
	directory string
	template string
}
var AppConfig *Config

func ParseConfig() error {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
		return errors.New("Internal error; unable to fetch current user")
	}

	// fetch configuration settings from flags
	var location = flag.String("location", "Local", "location to use")
	var directory = flag.String("directory", fmt.Sprintf("%s/.go-backlog", usr.HomeDir), "Backlog directory. If this directory doesn't exist it will be created")
	var template = flag.String("template", fmt.Sprintf("%s/.go-template.md", usr.HomeDir), "Template file for new backlog entries.")
	flag.Parse()

	loc, err := time.LoadLocation(*location)
	if err != nil {
		log.Println(err)
		return errors.New(fmt.Sprintf("Invalid location %s", *location))
	}

	// attempt to create the backlog directory at boot time, erring out if we are unable to create it
	// NOTE ~ characters are expected to be expanded by the shell before getting to the application
	if _, err := os.Stat(*directory); os.IsNotExist(err) {
		err := os.Mkdir(*directory, 0755)
		if err != nil {
			log.Println(err)
			return errors.New(fmt.Sprintf("Unable to create backlog directory"))
		}
	}

	AppConfig = &Config{
		location: loc,
		directory: *directory,
		template: *template,
	}

	return nil
}

func ParseDate(args []string) (time.Time, error) {
	if len(args) == 0 {
		return time.Now(), nil
	}
	
	// this is temporary while we build out datelp
	date := time.Now()
	//date, err := datelp.Parse(args)
	//if err != nil {
		//return nil, err
	//}

	return date, nil
}

func Cleanup() error {
	return nil
	filepath := path.Join(AppConfig.directory, ".cleanup")
	// NOOP if there is nothing to cleanup
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil
	}

	timestamp, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil
	}

	fmt.Println(timestamp)

	return nil

}

func main() {
	var action = flag.String("action", "open", "Open/Print to either open an entry file or print its full filepath")
	err := ParseConfig()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	date, err := ParseDate(flag.Args())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// TODO: clean up from a previous print command
	if err := Cleanup(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	entry := NewDayEntry(date)
	
	// initialize the method based upon input parameters
	var method func() error
	switch {
	case "open" == *action:
		method = entry.Open
	case "print" == *action:
		method = entry.Print
	default:
		log.Fatal("Unknown action")
		os.Exit(1)
	}

	if err := method(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

