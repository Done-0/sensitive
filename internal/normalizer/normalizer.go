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
	variant bool
	lower   bool
}

func New(variant, caseSensitive bool) *Normalizer {
	return &Normalizer{variant: variant, lower: !caseSensitive}
}

func (n *Normalizer) Normalize(text string) string {
	runes := []rune(text)
	for i, r := range runes {
		if n.variant && variantMap != nil {
			if s, ok := variantMap[r]; ok {
				r = s
			}
		}
		if n.lower {
			r = unicode.ToLower(r)
		}
		if r >= 0xFF01 && r <= 0xFF5E {
			r -= 0xFEE0
		} else if r == 0x3000 {
			r = ' '
		}
		runes[i] = r
	}
	return string(runes)
}

func (n *Normalizer) ToRunes(text string, buf []rune) []rune {
	buf = buf[:0]
	if !n.variant && !n.lower {
		for _, r := range text {
			if r >= 0xFF01 && r <= 0xFF5E {
				r -= 0xFEE0
			} else if r == 0x3000 {
				r = ' '
			}
			buf = append(buf, r)
		}
		return buf
	}
	if n.lower && !n.variant {
		for _, r := range text {
			if r >= 'A' && r <= 'Z' {
				r += 32
			} else if r > 127 {
				r = unicode.ToLower(r)
			}
			if r >= 0xFF01 && r <= 0xFF5E {
				r -= 0xFEE0
			} else if r == 0x3000 {
				r = ' '
			}
			buf = append(buf, r)
		}
		return buf
	}
	for _, r := range text {
		if variantMap != nil {
			if s, ok := variantMap[r]; ok {
				r = s
			}
		}
		if n.lower {
			if r >= 'A' && r <= 'Z' {
				r += 32
			} else if r > 127 {
				r = unicode.ToLower(r)
			}
		}
		if r >= 0xFF01 && r <= 0xFF5E {
			r -= 0xFEE0
		} else if r == 0x3000 {
			r = ' '
		}
		buf = append(buf, r)
	}
	return buf
}

func LoadVariantMap(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	variantMap = make(map[rune]rune, 8000)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		t := []rune(strings.TrimSpace(parts[0]))
		s := []rune(strings.TrimSpace(parts[1]))
		if len(t) == 1 && len(s) == 1 {
			variantMap[t[0]] = s[0]
		}
	}
	return sc.Err()
}

func IsVariantLoaded() bool {
	return len(variantMap) > 0
}
