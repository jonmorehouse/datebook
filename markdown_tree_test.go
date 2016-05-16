package main

import (
	"bytes"
	"testing"
)

var SupportedMarkdown = `# Header 1
text line 1

## Header 2

* this is more text
* this is more text

## Header 3

### Header 3

* this is not a top level header, at this point
`

func TestSupportMarkdownParsing(t *testing.T) {
	parser := NewMarkdownParser()
	parser.Parse(bytes.NewBufferString(SupportedMarkdown))

	iter := 0
	expected := []string{
		"# Header 1",
		"## Header 2",
		"## Header 3",
	}

	cb := func(e Element) bool {
		if e.Header() != expected[iter] {
			t.Fatalf("Did not parse markdown headers correctly")
		}
		iter++
		return true
	}
	parser.Walk(cb)
}
