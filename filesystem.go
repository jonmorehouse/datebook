package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type FileSystem struct {
	config *Config
}

func NewFileSystem(config *Config) *FileSystem {
	return &FileSystem{
		config: config,
	}
}

func (f *FileSystem) bootstrap(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	return nil
}

func (f *FileSystem) TemplateFilepath() string {
	current, err := user.Current()
	if err != nil {
		return ""
	}

	path := path.Join(current.HomeDir, ".datebook.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ""
	}

	return path
}

func (f FileSystem) DateFilepath(date time.Time) string {
	filename := fmt.Sprintf("%s_%s_%d_%d.md", date.Weekday(), date.Month().String(), date.Day(), date.Year())
	path := path.Join(f.MonthDirectory(date), strings.ToLower(filename))
	f.bootstrap(path)

	return path
}

func (f FileSystem) MonthDirectory(date time.Time) string {
	year, month, _ := date.Date()
	monthDir := strings.ToLower(fmt.Sprintf("%d_%s", month, month.String()))
	yearDir := fmt.Sprintf("%d", year)
	path := path.Join(f.config.Root, yearDir, monthDir)
	f.bootstrap(path)
	return path
}

func (f FileSystem) WeekFilepath(date time.Time) string {
	start := date.AddDate(0, 0, -1*int(date.Weekday()))
	year, month, day := start.Date()
	filename := strings.ToLower(fmt.Sprintf("week_%s_%d_%d.md", month.String(), day, year))
	path := path.Join(f.MonthDirectory(start), filename)
	f.bootstrap(path)
	return path
}

func (f FileSystem) LongtermFilepath() string {
	path := path.Join(f.config.Root, "longterm.md")
	f.bootstrap(path)
	return path
}

func (f FileSystem) UpdateReadme(srcPath string) error {
	// we want this symlink to work with Github so it needs to be relative, as of the root of the datebook dir
	if err := os.Chdir(f.config.Root); err != nil {
		return err
	}
	readmePath := "README.md"

	// remove the current symlinking, ignoring errors if it doesn't exist
	os.Remove(readmePath)

	// now we chop off the full Root path from src path so it is relative now...
	srcPath = strings.Replace(srcPath, f.config.Root, ".", 1)
	if err := os.Symlink(srcPath, readmePath); err != nil {
		return err
	}

	return nil
}

func (f FileSystem) Commit(date time.Time) error {
	if _, err := os.Stat(path.Join(f.config.Root, ".git")); os.IsNotExist(err) {
		fmt.Println("Datebook filesystem is not a git repository")
		return nil
	}

	commitMsg := fmt.Sprintf("%s %s %d", date.Weekday(), date.Month().String(), date.Year())
	os.Chdir(f.config.Root)

	if err := exec.Command("git", "add", "--all").Run(); err != nil {
		return err
	}
	exec.Command("git", "commit", "-m", commitMsg).Run()

	fmt.Println(fmt.Sprintf("Committed changes to datebook filesystem for %s %s %d %d",
		date.Weekday(), date.Month().String(), date.Day(), date.Year()))

	return nil
}

func (f FileSystem) EditFile(path string) error {
	command := exec.Command(os.Getenv("EDITOR"), path)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return err
	}

	return nil
}
