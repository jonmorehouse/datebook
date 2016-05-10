package main

import (
	"log"
	"os"
	"os/user"
	"path"
)

type Config struct {
	Root string
}

func NewDefaultConfig() *Config {
	current, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		Root: path.Join(current.HomeDir, ".datebook"),
	}
}

func main() {
	cli := NewCLI(NewDefaultConfig())

	args := os.Args[1:]
	if err := cli.Open(args); err != nil {
		log.Fatal(err)
	}
}
