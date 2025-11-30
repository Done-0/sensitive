// Package main demonstrates dependency injection pattern for production applications
// Creator: Done-0
// Created: 2025-01-15
package main

import (
	"fmt"
	"log"

	"github.com/Done-0/sensitive"
)

type ContentModerator interface {
	Validate(text string) bool
	Filter(text string) string
	Detect(text string) *sensitive.Result
}

type UserService struct {
	moderator ContentModerator
}

func NewUserService(moderator ContentModerator) *UserService {
	return &UserService{moderator: moderator}
}

func (s *UserService) CreatePost(userID int, content string) error {
	if s.moderator.Validate(content) {
		return fmt.Errorf("content contains sensitive words")
	}

	fmt.Printf("[Post Created] User %d: %s\n", userID, content)
	return nil
}

func (s *UserService) CreateComment(postID int, content string) (string, error) {
	result := s.moderator.Detect(content)

	if result.HasSensitive {
		for _, m := range result.Matches {
			if m.Level == sensitive.LevelHigh {
				return "", fmt.Errorf("comment rejected: prohibited content")
			}
		}
		filtered := result.FilteredText
		fmt.Printf("[Comment Created] Post %d: %s (filtered)\n", postID, filtered)
		return filtered, nil
	}

	fmt.Printf("[Comment Created] Post %d: %s\n", postID, content)
	return content, nil
}

type App struct {
	userService *UserService
}

func NewApp(moderator ContentModerator) *App {
	return &App{
		userService: NewUserService(moderator),
	}
}

func main() {
	detector := sensitive.New(
		sensitive.WithFilterStrategy(sensitive.StrategyMask),
	)

	detector.AddWords(map[string]sensitive.Level{
		"badword": sensitive.LevelHigh,
		"spam":    sensitive.LevelMedium,
	})
	detector.Build()

	app := NewApp(detector)

	fmt.Println("=== Dependency Injection Example ===")

	if err := app.userService.CreatePost(100, "This is a normal post"); err != nil {
		log.Printf("Error: %v", err)
	}

	if err := app.userService.CreatePost(101, "This contains badword"); err != nil {
		log.Printf("Error: %v\n", err)
	}

	if _, err := app.userService.CreateComment(1, "Great article!"); err != nil {
		log.Printf("Error: %v", err)
	}

	if _, err := app.userService.CreateComment(2, "This is spam content"); err != nil {
		log.Printf("Error: %v", err)
	}

	if _, err := app.userService.CreateComment(3, "Contains badword here"); err != nil {
		log.Printf("Error: %v\n", err)
	}

	fmt.Println("\n✓ All components use dependency injection")
	fmt.Println("✓ Easy to test with mock implementations")
	fmt.Println("✓ Loose coupling between services")
}
