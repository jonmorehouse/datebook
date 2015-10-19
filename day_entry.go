package main

import (
	"time"
	"fmt"
	"path"
	"os"
	"os/exec"
	"io/ioutil"
	"strings"
	"errors"
)

type DayEntry struct {
	date time.Time
	week *WeekEntry
}

func NewDayEntry(date time.Time) *DayEntry {
	week := NewWeekEntry(date)
	entry := &DayEntry{
		date: date,
		week: week,
	}

	return entry
}

func (e DayEntry) filepath() string {
	filename := fmt.Sprintf("%s_%s_%d_%d.md", e.date.Weekday(), e.date.Month(), e.date.Day(), e.date.Year())
	filename = strings.ToLower(filename)
	return path.Join(AppConfig.directory, DateDirectory(e.date), filename)
}

func (e DayEntry) createIfNotExists() error {
	directory := path.Join(AppConfig.directory, DateDirectory(e.date))

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(e.filepath()); err == nil {
		return nil
	}

	file, err := os.Create(e.filepath())
	defer file.Close()
	if err != nil {
		return err
	}

	// TODO: make this its own function?
	templateData, err := ioutil.ReadFile(AppConfig.template)
	if err != nil {
		return err
	}
	template := string(templateData)
	template = strings.Replace(template, "%year%", fmt.Sprintf("%s", e.date.Year()), -1)
	template = strings.Replace(template, "%month%", fmt.Sprintf("%s", e.date.Month().String()), -1)
	template = strings.Replace(template, "%weekday%", fmt.Sprintf("%s", e.date.Weekday().String()), -1)
	template = strings.Replace(template, "%day%", fmt.Sprintf("%d", e.date.Day()), -1)

	file.Write([]byte(template))

	return nil
}

func (e DayEntry) symlink() error {
	todayYear, todayMonth, todayDay := time.Now().Date()
	year, month, day := e.date.Date()

	// if the datebook entry isn't today's then NOOP
	if year != todayYear || month != todayMonth || day != todayDay {
		return nil
	}

	// we use a subprocess to symlink, in order to avoid having to delete the file and then re-symlink
	command := exec.Command("ln", "-sf", e.filepath(), path.Join(AppConfig.directory, "README.md"))
	if err := command.Run(); err != nil {
		return err
	}

	return nil
}

func (e DayEntry) commit() error {
	// if the datebook directory isn't a git repo, NOOP
	if _, err := os.Stat(path.Join(AppConfig.directory, ".git")); os.IsNotExist(err) {
		return nil
	}

	commitMessage := fmt.Sprintf("%s %s %d", e.date.Weekday(), e.date.Month().String(), e.date.Year())

	// change directories to the datebook repo
	os.Chdir(AppConfig.directory)
	command := exec.Command("git", "add", "--all")
	if err := command.Run(); err != nil {
		return err
	}

	// check to see if any changes are staged
	command = exec.Command("git", "diff", "HEAD", "--exit-code")
	// this command will return an error code if there _are_ changes...
	if err := command.Run(); err == nil {
		return nil
	}

	command = exec.Command("git", "commit", "-m", fmt.Sprintf("\"%s\"", commitMessage))
	if err := command.Run(); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Committed changes for %s %s %d", e.date.Weekday(), e.date.Month().String(), e.date.Year()))
	return nil
}

func (e DayEntry) Print() error {
	if err := e.Prepare(); err != nil {
		return err
	}

	cleanupFilepath := path.Join(AppConfig.directory, ".cleanup")
	if _, err := os.Stat(cleanupFilepath); !os.IsNotExist(err) {
		return errors.New("Cleanup failed. Please try again")
	}

	// create a cleanupFile and template out the timeStamp into it
	cleanupFile, err := os.Create(cleanupFilepath)
	defer cleanupFile.Close()
	if err != nil {
		return err
	}

	// we just write the unix epoch timestamp into the cleanup file to make things easier to parse later
	cleanupFile.Write([]byte(fmt.Sprintf("%d", e.date.Unix())))
	os.Stdout.Write([]byte(fmt.Sprintf("%s\n", e.filepath())))

	return nil
}

func (e DayEntry) Open() error {
	if err := e.Prepare(); err != nil {
		return err
	}

	// actually open the file editor
	command := exec.Command(os.Getenv("EDITOR"), e.filepath())
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		return err
	}

	return e.Save()
}

func (e DayEntry) Prepare() error {
	if err := e.createIfNotExists(); err != nil {
		return err
	}

	if err := e.week.WriteToDay(e.filepath()); err != nil {
		return err
	}

	if err := e.symlink(); err != nil {
		return err
	}

	return nil
}

func (e DayEntry) Save() error {
	if err := e.week.SaveFromDay(e.filepath()); err != nil {
		return err
	}

	if err := e.commit(); err != nil{
		return err
	}

	return nil
}
