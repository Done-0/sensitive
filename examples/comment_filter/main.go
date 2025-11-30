// Package main demonstrates user-generated content filtering system
// Creator: Done-0
// Created: 2025-01-15
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Done-0/sensitive"
)

type Comment struct {
	ID        int
	UserID    int
	Content   string
	Timestamp time.Time
}

type ModerationResult struct {
	Comment  *Comment
	Approved bool
	Reason   string
	Filtered string
}

type CommentFilter struct {
	detector *sensitive.Detector
}

func NewCommentFilter(dictPath string) (*CommentFilter, error) {
	detector := sensitive.New(
		sensitive.WithFilterStrategy(sensitive.StrategyMask),
	)

	if err := detector.LoadDict(dictPath); err != nil {
		detector.AddWords(map[string]sensitive.Level{
			"badword": sensitive.LevelHigh,
			"spam":    sensitive.LevelMedium,
			"abuse":   sensitive.LevelMedium,
		})
	}

	detector.Build()

	return &CommentFilter{detector: detector}, nil
}

func (cf *CommentFilter) Moderate(comment *Comment) ModerationResult {
	result := ModerationResult{
		Comment:  comment,
		Approved: true,
	}

	detection := cf.detector.Detect(comment.Content)

	if detection.HasSensitive {
		highSeverity := false
		for _, m := range detection.Matches {
			if m.Level == sensitive.LevelHigh {
				highSeverity = true
				break
			}
		}

		if highSeverity {
			result.Approved = false
			result.Reason = "Contains prohibited content"
		} else {
			result.Approved = true
			result.Filtered = detection.FilteredText
			result.Reason = "Auto-filtered"
		}
	}

	return result
}

func main() {
	filter, err := NewCommentFilter("dict/sensitive.txt")
	if err != nil {
		log.Printf("Warning: %v, using default words", err)
	}

	comments := []Comment{
		{ID: 1, UserID: 100, Content: "Great product!", Timestamp: time.Now()},
		{ID: 2, UserID: 101, Content: "This is spam content", Timestamp: time.Now()},
		{ID: 3, UserID: 102, Content: "Contains badword here", Timestamp: time.Now()},
	}

	for _, comment := range comments {
		result := filter.Moderate(&comment)

		fmt.Printf("\n[Comment #%d from User %d]\n", comment.ID, comment.UserID)
		fmt.Printf("Original: %s\n", comment.Content)

		if result.Approved {
			if result.Filtered != "" {
				fmt.Printf("Status: APPROVED (filtered)\n")
				fmt.Printf("Display: %s\n", result.Filtered)
			} else {
				fmt.Printf("Status: APPROVED\n")
			}
		} else {
			fmt.Printf("Status: REJECTED - %s\n", result.Reason)
		}
	}
}
