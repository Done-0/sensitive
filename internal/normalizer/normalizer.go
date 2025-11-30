// Package normalizer provides text normalization for sensitive word detection
// Creator: Done-0
// Created: 2025-01-15
package normalizer

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

var variantMap map[rune]rune

type Normalizer struct {
	enableVariant bool
	caseSensitive bool
}

func New(enableVariant, caseSensitive bool) *Normalizer {
	return &Normalizer{
		enableVariant: enableVariant,
		caseSensitive: caseSensitive,
	}
}

func (n *Normalizer) Normalize(text string) string {
	runes := []rune(text)
	for i, r := range runes {
		if n.enableVariant && len(variantMap) > 0 {
			if simplified, ok := variantMap[r]; ok {
				runes[i] = simplified
				r = simplified
			}
		}

		if !n.caseSensitive {
			runes[i] = unicode.ToLower(r)
			r = unicode.ToLower(r)
		}

		if r >= 0xFF01 && r <= 0xFF5E {
			runes[i] = r - 0xFEE0
		} else if r == 0x3000 {
			runes[i] = 0x0020
		}
	}

	return string(runes)
}

func LoadVariantMap(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	variantMap = make(map[rune]rune, 8000)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}

		traditional := []rune(strings.TrimSpace(parts[0]))
		simplified := []rune(strings.TrimSpace(parts[1]))

		if len(traditional) == 1 && len(simplified) == 1 {
			variantMap[traditional[0]] = simplified[0]
		}
	}

	return scanner.Err()
}

func IsVariantLoaded() bool {
	return len(variantMap) > 0
}
