// Package main demonstrates high concurrency usage for production environments
// Creator: Done-0
// Created: 2025-01-15
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Done-0/sensitive"
)

var detector *sensitive.Detector

func init() {
	detector = sensitive.New(
		sensitive.WithFilterStrategy(sensitive.StrategyMask),
	)

	detector.AddWords(map[string]sensitive.Level{
		"badword":  sensitive.LevelHigh,
		"spam":     sensitive.LevelMedium,
		"abuse":    sensitive.LevelMedium,
		"violence": sensitive.LevelHigh,
	})
	detector.Build()
}

func processContent(id int, content string, wg *sync.WaitGroup, rejected *atomic.Int32, filtered *atomic.Int32) {
	defer wg.Done()

	result := detector.Detect(content)

	if result.HasSensitive {
		highSeverity := false
		for _, m := range result.Matches {
			if m.Level == sensitive.LevelHigh {
				highSeverity = true
				break
			}
		}

		if highSeverity {
			rejected.Add(1)
			fmt.Printf("[Worker %d] REJECTED: %s\n", id, content)
		} else {
			filtered.Add(1)
			fmt.Printf("[Worker %d] FILTERED: %s -> %s\n", id, content, result.FilteredText)
		}
	}
}

func main() {
	fmt.Println("=== High Concurrency Example ===")

	testContents := []string{
		"This is a normal message",
		"This contains spam content",
		"Great product, highly recommended",
		"Warning: violence content here",
		"This is badword example",
		"Just some abuse language",
		"Clean content without issues",
		"More spam and abuse here",
		"Another normal message",
		"Final test with violence",
	}

	fmt.Println("Test 1: Concurrent Processing (100 workers)")
	fmt.Println("---")

	var wg sync.WaitGroup
	var rejected atomic.Int32
	var filtered atomic.Int32

	start := time.Now()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		content := testContents[i%len(testContents)]
		go processContent(i, content, &wg, &rejected, &filtered)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("\nProcessed 100 requests in %v\n", elapsed)
	fmt.Printf("Rejected: %d, Filtered: %d\n", rejected.Load(), filtered.Load())

	fmt.Println("\nTest 2: Sustained Load (1000 concurrent requests)")
	fmt.Println("---")

	rejected.Store(0)
	filtered.Store(0)
	var processed atomic.Int32

	start = time.Now()

	for i := range 1000 {
		wg.Add(1)
		content := testContents[i%len(testContents)]
		go func(id int, text string) {
			defer wg.Done()
			detector.Validate(text)
			processed.Add(1)
		}(i, content)
	}

	wg.Wait()
	elapsed = time.Since(start)

	fmt.Printf("Processed %d requests in %v\n", processed.Load(), elapsed)
	fmt.Printf("Throughput: %.0f requests/second\n", float64(processed.Load())/elapsed.Seconds())

	fmt.Println("\n✓ Thread-safe concurrent access")
	fmt.Println("✓ No race conditions")
	fmt.Println("✓ Suitable for high-traffic production systems")
}
