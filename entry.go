package main

import (
	"fmt"
	"strings"
	"time"
)

type Week struct {
	start time.Time
	end   time.Time
}

type Entry struct {
	date time.Time
	fs   *FileSystem
}

func NewEntry(fs *FileSystem, date time.Time) *Entry {
	return &Entry{
		date: date,
		fs:   fs,
	}
}

func (e Entry) Filepath() string {
	return e.fs.DateFilepath(e.date)
}

func (e Entry) Template(tree ElementTree) error {
	templateTree := NewMarkdownTree()
	if err := templateTree.ParseFile(e.fs.TemplateFilepath()); err != nil {
		return err
	}

	replacements := make(map[string]string)
	replacements["{{year}}"] = fmt.Sprintf("%s", e.date.Year())
	replacements["{{month}}"] = e.date.Month().String()
	replacements["{{weekday}}"] = e.date.Weekday().String()
	replacements["{{day}}"] = fmt.Sprintf("%d", e.date.Day())

	walker := func(element Element) bool {
		for k, v := range replacements {
			// this template engine only supports templating out name attributes as part of the header
			name := element.Name()
			name = strings.Replace(name, k, v, -1)
			element.SetName(name)
		}

		tree.Push(element)
		return true
	}

	templateTree.Walk(walker)
	return nil
}

func (e Entry) Prepare() error {
	// load the entry tree, templating it out if it doesn't already exist
	dateTree := NewMarkdownTree()
	if err := dateTree.ParseFile(e.fs.DateFilepath(e.date)); err != nil {
		if err := e.Template(dateTree); err != nil {
			return err
		}
	}

	// we pop both of these elements because we want to ensure that they are the latest!
	// we always make sure to save these files
	dateTree.Pop("longterm", 2)
	dateTree.Pop("week", 2)

	longtermTree := NewMarkdownTree()
	if err := longtermTree.ParseFile(e.fs.LongtermFilepath()); err != nil {
		element := NewMarkdownElement(2, "longterm")
		element.WriteLine("\n")
		longtermTree.Push(element)
	}

	// if the week tree is empty, then we go ahead and add in a default header
	weekTree := NewMarkdownTree()
	if err := weekTree.ParseFile(e.fs.WeekFilepath(e.date)); err != nil || weekTree.Len() == 0 {
		element := NewMarkdownElement(2, "week")
		element.WriteLine("\n")
		weekTree.Push(element)
	}

	// walk the longterm and week trees and push any relevant items into the current date entry
	cb := func(element Element) bool {
		if element.Level() == 2 {
			dateTree.Push(element)
		}
		return true
	}
	weekTree.Walk(cb)
	longtermTree.Walk(cb)
	if err := dateTree.WriteFile(e.fs.DateFilepath(e.date)); err != nil || longtermTree.Len() == 0 {
		return err
	}
	return nil
}

func (e Entry) Save() error {
	dateFilepath := e.fs.DateFilepath(e.date)
	dateTree := NewMarkdownTree()
	if err := dateTree.ParseFile(dateFilepath); err != nil {
		return err
	}

	// fetch the longterm elements from the entry and write them to the corresponding file
	longtermTree := NewMarkdownTree()
	if element := dateTree.Find("longterm", 2); element != nil {
		longtermTree.Push(element)
	}
	if err := longtermTree.WriteFile(e.fs.LongtermFilepath()); err != nil {
		return err
	}

	weekTree := NewMarkdownTree()
	if element := dateTree.Find("week", 2); element != nil {
		weekTree.Push(element)
	}
	if err := weekTree.WriteFile(e.fs.WeekFilepath(e.date)); err != nil {
		return err
	}

	return nil
}
