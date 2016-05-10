package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type Element interface {
	Level() int
	Name() string
	SetName(string)
	Header() string
	Body() []string
	WriteLine(string)
}

type MarkdownElement struct {
	level int
	name  string
	body  []string
}

func NewMarkdownElement(level int, name string) Element {
	return &MarkdownElement{
		level: level,
		name:  name,
		body:  make([]string, 0),
	}
}

func (e MarkdownElement) Name() string {
	return e.name
}

func (e *MarkdownElement) SetName(name string) {
	e.name = name
}

func (e MarkdownElement) Level() int {
	return e.level
}

func (e MarkdownElement) Header() string {
	header := strings.Repeat("#", e.level)
	return fmt.Sprintf("%s %s", header, e.name)
}

func (e MarkdownElement) Body() []string {
	return e.body
}

func (e *MarkdownElement) WriteLine(line string) {
	e.body = append(e.body, line)
}

type ElementTreeWalker func(Element) bool
type ElementTree interface {
	ParseFile(string) error
	Parse(io.Reader) error

	WriteFile(string) error
	Write(io.Writer) error

	Len() int
	Walk(ElementTreeWalker)
	Find(string, int) Element
	Pop(string, int) Element
	Push(Element)
}

type MarkdownTree struct {
	elements []Element
}

func NewMarkdownTree() ElementTree {
	return &MarkdownTree{
		elements: make([]Element, 0),
	}
}

func (m MarkdownTree) Len() int {
	return len(m.elements)
}

func (m MarkdownTree) WriteFile(filepath string) error {
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		return err
	}

	return m.Write(file)
}

func (m MarkdownTree) Write(writer io.Writer) error {
	for _, element := range m.elements {
		io.WriteString(writer, element.Header())
		io.WriteString(writer, "\n")

		for _, line := range element.Body() {
			io.WriteString(writer, line)
			io.WriteString(writer, "\n")
		}
	}

	return nil
}

func (m *MarkdownTree) Parse(reader io.Reader) error {
	headerRe := regexp.MustCompile(`^(?P<level>#{1,2})\s(?P<header>.+)$`)
	scanner := bufio.NewScanner(reader)

	var current Element
	for scanner.Scan() {
		line := scanner.Text()
		res := headerRe.FindStringSubmatch(line)

		// if the current element doesn't appear to be a header
		// element, then we write it to the current element's body
		if res == nil {
			if current == nil {
				return fmt.Errorf("error parsing markdown file, no header found")
			}
			current.WriteLine(line)
			continue
		}

		level, header := res[1], res[2]
		if current != nil {
			m.elements = append(m.elements, current)
		}
		current = NewMarkdownElement(len(level), header)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if current != nil {
		m.elements = append(m.elements, current)
	}
	return nil
}

func (m *MarkdownTree) ParseFile(filepath string) error {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return err
	}

	return m.Parse(file)
}

func (m MarkdownTree) Walk(cb ElementTreeWalker) {
	for _, element := range m.elements {
		if !cb(element) {
			break
		}
	}
}

func (m MarkdownTree) Find(name string, level int) Element {
	for _, element := range m.elements {
		if element.Name() == name && element.Level() == level {
			return element
		}
	}

	return nil
}

func (m *MarkdownTree) Push(element Element) {
	m.elements = append(m.elements, element)
}

func (m *MarkdownTree) Pop(name string, level int) Element {
	for index, element := range m.elements {
		if element.Name() == name && element.Level() == level {
			m.elements = append(m.elements[:index], m.elements[index+1:]...)
			return element
		}
	}
	return nil
}
