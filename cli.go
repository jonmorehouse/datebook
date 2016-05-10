package main

import (
	"log"
	"strings"
	"time"

	"github.com/jonmorehouse/datelp"
)

type CLI struct {
	config *Config
}

func NewCLI(config *Config) *CLI {
	return &CLI{
		config: config,
	}
}

func (c CLI) ParseDate(args []string) (time.Time, error) {
	if len(args) == 0 {
		return time.Now(), nil
	}

	date, err := datelp.Parse(strings.Join(args, " "))
	if err != nil {
		log.Fatalf("Passed in a date string which could not be parsed")
	}
	return date, nil
}

func (c CLI) Open(args []string) error {
	date, err := c.ParseDate(args)
	if err != nil {
		return err
	}

	fs := NewFileSystem(c.config)
	entry := NewEntry(fs, date)
	if err := entry.Prepare(); err != nil {
		return err
	}

	if err := fs.EditFile(entry.Filepath()); err != nil {
		return err
	}

	if err := entry.Save(); err != nil {
		return err
	}

	today := time.Now()
	if date.Unix()-date.Unix()%(3600*24) == today.Unix()-today.Unix()%(3600*24) {
		if err := fs.UpdateReadme(entry.Filepath()); err != nil {
			return err
		}
	}

	if err := fs.Commit(date); err != nil {
		return err
	}
	return nil
}

func (c CLI) Cleanup() error {
	// finds any files that need to be cleaned up and then parses, saves and commits them
	panic("Not implemented yet!")
	return nil
}

func (c CLI) Search(args []string) error {
	panic("Not implemented yet!")
	return nil
}
