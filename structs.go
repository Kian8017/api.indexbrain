package main

import (
	"strings"
)

type Country struct {
	Abbr string `json:"abbr"`
	Name string `json:"name"`
}

func (c Country) Folder() string {
	return c.Abbr + " " + c.Name
}

func NewCountry(dir string) Country {
	parts := strings.SplitN(dir, " ", 2) // ABBR NAME
	if len(parts) != 2 {
		panic("NewCountry length was not 2: " + strings.Join(parts, " "))
	}
	return Country{Abbr: parts[0], Name: parts[1]}
}

type Entry struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Country string `json:"country"`
}
