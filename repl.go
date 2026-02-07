package main

import (
	"strings"
)

func cleanInput(text string) []string {
	var clean []string
	if text == "" {
		return clean
	}
	clean = strings.Fields(strings.ToLower(text))
	return clean
}
