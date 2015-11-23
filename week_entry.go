package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type WeekEntry struct {
	date time.Time
}

func (w *WeekEntry) weekBounds() (time.Time, time.Time) {
	current := int(w.date.Weekday())

	// calculate the offset into the current week we are and use that as
	// the basis for creating the first/last day of the week
	left := w.date.AddDate(0, 0, int(current)*-1)
	right := w.date.AddDate(0, 0, 7-current)

	return left, right
}

func (w *WeekEntry) filepath() string {
	left, _ := w.weekBounds()
	filename := fmt.Sprintf("week_%s_%d.md", left.Month().String(), left.Day())
	filename = strings.ToLower(filename)

	return path.Join(AppConfig.directory, DateDirectory(left), filename)
}

func (w *WeekEntry) getWeekBounds(buf []byte) (int, int, bool) {
	// look for the week header
	startRegex := regexp.MustCompile(`## week.*`)
	start := startRegex.FindIndex(buf)
	if start == nil {
		return 0, 0, false
	}

	// chop off the beginning of the day's entry before looking for the ending
	remainder := buf[start[1]:]

	// find the end of the week data, which is denoted by a new header or EOF
	endRegex := regexp.MustCompile(`(## |# |$)`)
	end := endRegex.FindIndex(remainder)

	// if for some reason the end fails... which it shouldn't ... then just
	// return the end of the buffer
	if end == nil {
		return start[0], len(buf) - 1, true
	}

	return start[0], start[1] + end[0], true
}

func (w *WeekEntry) SaveFromDay(filepath string) error {
	day, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	left, right, found := w.getWeekBounds(day)
	if !found {
		return nil
	}

	// grab the week data, and write it out to the week file
	week := StripNewLines(day[left:right])

	if err := ioutil.WriteFile(w.filepath(), week, 0644); err != nil {
		return err
	}

	return nil
}

func (w *WeekEntry) WriteToDay(filepath string) error {
	// if the week file doesn't exist then NOOP
	if _, err := os.Stat(w.filepath()); os.IsNotExist(err) {
		return nil
	}

	week, err := ioutil.ReadFile(w.filepath())
	if err != nil {
		return err
	}
	day, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// fetch the bounds of the existing week block
	left, right, found := w.getWeekBounds(day)
	if !found {
		left = len(day) - 1
	}

	newDay := append(day[:left], week...)
	newDay = append(newDay, '\n', '\n')
	if found {
		newDay = append(newDay, day[right:]...)
	}

	if err := ioutil.WriteFile(filepath, newDay, 0644); err != nil {
		return err
	}

	return nil
}

func NewWeekEntry(date time.Time) *WeekEntry {
	return &WeekEntry{
		date: date,
	}
}
