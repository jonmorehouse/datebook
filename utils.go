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


