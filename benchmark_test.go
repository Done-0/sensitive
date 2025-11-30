// Package sensitive provides high-performance sensitive word detection using AC automaton
// Creator: Done-0
// Created: 2025-01-15
package sensitive

import (
	"strings"
	"testing"
)

func BenchmarkDetector_Detect_SmallDict(b *testing.B) {
	detector := setupDetector(100)
	text := generateText(1000)

	for b.Loop() {
		detector.Detect(text)
	}
}

func BenchmarkDetector_Detect_MediumDict(b *testing.B) {
	detector := setupDetector(1000)
	text := generateText(1000)

	for b.Loop() {
		detector.Detect(text)
	}
}

func BenchmarkDetector_Detect_LargeDict(b *testing.B) {
	detector := setupDetector(10000)
	text := generateText(1000)

	for b.Loop() {
		detector.Detect(text)
	}
}

func BenchmarkDetector_Detect_ShortText(b *testing.B) {
	detector := setupDetector(1000)
	text := generateText(100)

	for b.Loop() {
		detector.Detect(text)
	}
}

func BenchmarkDetector_Detect_LongText(b *testing.B) {
	detector := setupDetector(1000)
	text := generateText(10000)

	for b.Loop() {
		detector.Detect(text)
	}
}

func BenchmarkDetector_AddWord(b *testing.B) {
	detector := New()

	for b.Loop() {
		detector.AddWord("testword", LevelMedium)
	}
}

func BenchmarkDetector_Build(b *testing.B) {
	detector := New()
	for range 10000 {
		detector.AddWord(generateWord(5), LevelMedium)
	}

	for b.Loop() {
		detector.Build()
	}
}

func BenchmarkDetector_Parallel(b *testing.B) {
	detector := setupDetector(1000)
	text := generateText(1000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			detector.Detect(text)
		}
	})
}

func setupDetector(wordCount int) *Detector {
	detector := New()

	for range wordCount {
		word := generateWord(5)
		detector.AddWord(word, LevelMedium)
	}

	detector.Build()
	return detector
}

func generateText(length int) string {
	words := []string{
		"这是", "一段", "测试", "文本", "用于", "性能",
		"基准", "测试", "包含", "一些", "常用", "词汇",
		"and", "some", "english", "words", "for", "testing",
	}

	var builder strings.Builder
	for i := range length {
		builder.WriteString(words[i%len(words)])
		if i%5 == 0 {
			builder.WriteString(" ")
		}
	}

	return builder.String()
}

func generateWord(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	word := make([]byte, length)
	for i := range length {
		word[i] = chars[i%len(chars)]
	}
	return string(word)
}
