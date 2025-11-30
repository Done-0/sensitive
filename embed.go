// Package sensitive provides high-performance sensitive word detection using AC automaton
// Creator: Done-0
// Created: 2025-01-15
package sensitive

import (
	"embed"
	"errors"
	"strings"
)

//go:embed configs/dict/*.txt
var dictFS embed.FS

const (
	DictHighPolitics     = "high_politics.txt"
	DictHighPornography  = "high_pornography.txt"
	DictHighViolence     = "high_violence.txt"
	DictMediumGeneral    = "medium_general.txt"
	DictLowAd            = "low_ad.txt"
	DictLowURL           = "low_url.txt"
)

func LoadAllEmbedded(detector *Detector) error {
	dicts := map[string]Level{
		DictHighPolitics:    LevelHigh,
		DictHighPornography: LevelHigh,
		DictHighViolence:    LevelHigh,
		DictMediumGeneral:   LevelMedium,
		DictLowAd:           LevelLow,
		DictLowURL:          LevelLow,
	}

	for name, level := range dicts {
		if err := LoadEmbeddedDict(detector, name, level); err != nil {
			return err
		}
	}

	return nil
}

func LoadEmbeddedDict(detector *Detector, name string, level Level) error {
	if !level.IsValid() {
		return errors.New("invalid level")
	}

	data, err := dictFS.ReadFile("configs/dict/" + name)
	if err != nil {
		return err
	}

	words := make([]string, 0, 512)
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSuffix(line, ",")
		if line != "" {
			words = append(words, line)
		}
	}

	wordMap := make(map[string]Level, len(words))
	for _, word := range words {
		wordMap[word] = level
	}

	return detector.AddWords(wordMap)
}
