// Package main demonstrates the simplest production usage
// Creator: Done-0
// Created: 2025-01-15
package main

import (
	"fmt"
	"log"

	"github.com/Done-0/sensitive"
)

func main() {
	detector := sensitive.New()

	detector.AddWord("badword", sensitive.LevelHigh)
	detector.AddWord("spam", sensitive.LevelMedium)
	detector.Build()

	userInput := "This is a badword example"
	if detector.Validate(userInput) {
		log.Printf("Content rejected: contains sensitive words")
		return
	}

	fmt.Println("Content approved")
}
