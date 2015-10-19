package main

import (
	"time"
	"path"
	"fmt"
	"strings"
)

func DateDirectory(date time.Time) string {
	year, month, _ := date.Date()
	name := fmt.Sprintf("%d_%s", month, strings.ToLower(month.String()))

	return path.Join(fmt.Sprintf("%d", year), name)
}

func StripNewLines(buf []byte) []byte {
	var startIndex, rightOffset int

	for index, word := range buf {
		if word == byte('\n') {
			continue
		}

		startIndex = index
		break
	}

	for index := len(buf)- 1; index > startIndex; index-- {
		if buf[index] != byte('\n') {
			break
		}

		rightOffset++
	}

	return buf[startIndex:len(buf) - rightOffset]
}
