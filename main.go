package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jonmorehouse/datelp"
)

type Config struct {
	location  *time.Location
	directory string
	template  string
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
	var directory = flag.String("directory", fmt.Sprintf("%s/.datebook", usr.HomeDir), "Datebook directory. If this directory doesn't exist it will be created")
	var template = flag.String("template", fmt.Sprintf("%s/.datebook.md", usr.HomeDir), "Template file for new datebook entries.")
	flag.Parse()

	loc, err := time.LoadLocation(*location)
	if err != nil {
		log.Println(err)
		return errors.New(fmt.Sprintf("Invalid location %s", *location))
	}

	// attempt to create the datebook directory at boot time, erring out if we are unable to create it
	// NOTE ~ characters are expected to be expanded by the shell before getting to the application
	if _, err := os.Stat(*directory); os.IsNotExist(err) {
		err := os.Mkdir(*directory, 0755)
		if err != nil {
			log.Println(err)
			return errors.New(fmt.Sprintf("Unable to create datebook directory"))
		}
	}

	AppConfig = &Config{
		location:  loc,
		directory: *directory,
		template:  *template,
	}

	return nil
}

func ParseDate(args []string) (time.Time, error) {
	if len(args) == 0 {
		return time.Now(), nil
	}

	dateString := strings.Join(args, " ")
	date, err := datelp.Parse(dateString)
	if err != nil {
		log.Fatalf("Passed in a date string that could not be parsed")
	}

	return date, nil
}

func Cleanup() error {
	filepath := path.Join(AppConfig.directory, ".cleanup")
	// NOOP if there is nothing to cleanup
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil
	}

	// delete the cleanup file regardless of the outcome
	defer os.Remove(filepath)

	timestamp, err := ioutil.ReadFile(filepath)
	if err != nil {
		return errors.New("Cleanup file corrupted")
	}

	epoch, err := strconv.ParseInt(string(timestamp), 10, 64)
	if err != nil {
		return errors.New("Cleanup file corrupted")
	}
	date := time.Unix(epoch, 0)
	entry := NewDayEntry(date)
	return entry.Save()
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
