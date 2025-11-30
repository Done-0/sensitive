// Package main demonstrates Fluent API pattern for elegant chain calls
// Creator: Done-0
// Created: 2025-01-15
package main

import (
	"fmt"
	"log"

	"github.com/Done-0/sensitive"
)

func main() {
	fmt.Println("=== Fluent API Example ===")
	fmt.Println()

	fmt.Println("1. Simple Chain")
	fmt.Println("---")

	detector := sensitive.NewBuilder().
		AddWord("badword", sensitive.LevelHigh).
		AddWord("spam", sensitive.LevelMedium).
		AddWord("abuse", sensitive.LevelMedium).
		MustBuild()

	text := "This contains badword and spam"
	result := detector.Detect(text)
	fmt.Printf("Original: %s\n", text)
	fmt.Printf("Filtered: %s\n", result.FilteredText)
	fmt.Println()

	fmt.Println("2. Complex Configuration Chain")
	fmt.Println("---")

	detector2 := sensitive.NewBuilder().
		WithFilterStrategy(sensitive.StrategyReplace).
		WithReplaceChar('█').
		WithCaseSensitive(false).
		WithSkipWhitespace(true).
		AddWords(map[string]sensitive.Level{
			"test":    sensitive.LevelHigh,
			"example": sensitive.LevelMedium,
			"sample":  sensitive.LevelLow,
		}).
		MustBuild()

	text2 := "This is a TEST example with SAMPLE content"
	fmt.Printf("Text: %s\n", text2)
	fmt.Printf("Filtered: %s\n", detector2.Filter(text2))
	fmt.Println()

	fmt.Println("3. Load Dictionary with Chain")
	fmt.Println("---")

	detector3, err := sensitive.NewBuilder().
		WithFilterStrategy(sensitive.StrategyMask).
		LoadDict("dict/high_illegal.txt").
		LoadDict("dict/medium_abuse.txt").
		Build()

	if err != nil {
		log.Printf("Load failed (expected): %v\n", err)
		detector3 = sensitive.NewBuilder().
			AddWord("fallback", sensitive.LevelHigh).
			MustBuild()
	}

	fmt.Printf("Detector initialized: %v words\n", detector3.Stats().TotalWords)
	fmt.Println()

	fmt.Println("4. All Features Combined")
	fmt.Println("---")

	detector4 := sensitive.NewBuilder().
		WithFilterStrategy(sensitive.StrategyMask).
		WithReplaceChar('*').
		WithCaseSensitive(false).
		WithVariant(true).
		AddWord("illegal", sensitive.LevelHigh).
		AddWord("violence", sensitive.LevelHigh).
		AddWords(map[string]sensitive.Level{
			"abuse":     sensitive.LevelMedium,
			"offensive": sensitive.LevelMedium,
		}).
		MustBuild()

	testCases := []string{
		"Clean content here",
		"Contains illegal content",
		"Some abuse and offensive language",
	}

	for i, tc := range testCases {
		result := detector4.Detect(tc)
		fmt.Printf("[%d] %s\n", i+1, tc)
		if result.HasSensitive {
			fmt.Printf("    → %s (%d matches)\n", result.FilteredText, len(result.Matches))
		} else {
			fmt.Printf("    → OK\n")
		}
	}

	fmt.Println("\n✓ Fluent API enables elegant chain calls")
	fmt.Println("✓ Methods return *Builder for chaining")
	fmt.Println("✓ Build() for error handling, MustBuild() for panic on error")
	fmt.Println("✓ Clean and readable configuration")
}
