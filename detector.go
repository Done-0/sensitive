// Package sensitive provides high-performance sensitive word detection using AC automaton
// Creator: Done-0
// Created: 2025-01-15
package sensitive

import (
	"bufio"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Done-0/sensitive/internal/normalizer"
	"github.com/Done-0/sensitive/internal/pool"
	"github.com/Done-0/sensitive/internal/trie"
)

type Detector struct {
	tree       *trie.Tree
	mu         sync.RWMutex
	normalizer *normalizer.Normalizer
	opts       *Options
	isBuilt    bool
	count      int
}

func New(opts ...Option) *Detector {
	o := &Options{
		FilterStrategy: StrategyMask,
		ReplaceChar:    '*',
		SkipWhitespace: true,
		EnableVariant:  false,
		CaseSensitive:  false,
	}
	for _, opt := range opts {
		opt(o)
	}

	return &Detector{
		tree:       trie.New(),
		normalizer: normalizer.New(o.EnableVariant, o.CaseSensitive),
		opts:       o,
	}
}

func (d *Detector) AddWord(word string, level Level) error {
	if word == "" {
		return errors.New("empty word")
	}
	if !level.IsValid() {
		return errors.New("invalid level")
	}

	d.mu.Lock()
	normalized := d.normalizer.Normalize(word)
	if normalized == "" {
		d.mu.Unlock()
		return errors.New("normalized word is empty")
	}

	d.tree.Insert(normalized, int(level))
	d.count++
	d.isBuilt = false
	d.mu.Unlock()
	return nil
}

func (d *Detector) AddWords(words map[string]Level) error {
	for word, level := range words {
		if err := d.AddWord(word, level); err != nil {
			return err
		}
	}
	return nil
}

func (d *Detector) Build() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	trie.BuildFailureLinks(d.tree.Root())
	d.isBuilt = true
	return nil
}

func (d *Detector) Detect(text string) *Result {
	result := &Result{FilteredText: text}
	if text == "" {
		return result
	}

	d.mu.RLock()
	if !d.isBuilt {
		d.mu.RUnlock()
		return result
	}

	normalized := d.normalizer.Normalize(text)
	runes := []rune(normalized)
	matches := trie.Search(d.tree.Root(), runes)
	d.mu.RUnlock()

	if len(matches) > 0 {
		result.HasSensitive = true
		result.Matches = make([]Match, len(matches))
		for i, m := range matches {
			result.Matches[i] = Match{
				Word:  m.Word,
				Start: m.Start,
				End:   m.End,
				Level: Level(m.Level),
			}
		}

		textRunes := []rune(text)
		filtered := pool.Get(len(textRunes))
		defer pool.Put(filtered)

		mask := make([]bool, len(textRunes))
		for _, m := range result.Matches {
			for i := m.Start; i < m.End && i < len(mask); i++ {
				mask[i] = true
			}
		}

		for i, r := range textRunes {
			if mask[i] {
				switch d.opts.FilterStrategy {
				case StrategyReplace:
					*filtered = append(*filtered, d.opts.ReplaceChar)
				case StrategyMask:
					*filtered = append(*filtered, '*')
				}
			} else {
				*filtered = append(*filtered, r)
			}
		}

		result.FilteredText = string(*filtered)
	}

	return result
}

func (d *Detector) Filter(text string) string {
	return d.Detect(text).FilteredText
}

func (d *Detector) IsVariantEnabled() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.opts.EnableVariant && normalizer.IsVariantLoaded()
}

func (d *Detector) LoadDict(path string) error {
	level := inferLevel(path)
	return d.LoadDictWithLevel(path, level)
}

func (d *Detector) LoadDictWithLevel(path string, level Level) error {
	if !level.IsValid() {
		return errors.New("invalid level")
	}

	words, err := loadFile(path)
	if err != nil {
		return err
	}

	wordMap := make(map[string]Level, len(words))
	for _, word := range words {
		wordMap[word] = level
	}

	return d.AddWords(wordMap)
}

func (d *Detector) LoadDictFromURL(url string) error {
	level := inferLevel(url)
	return d.LoadDictFromURLWithLevel(url, level)
}

func (d *Detector) LoadDictFromURLWithLevel(url string, level Level) error {
	if !level.IsValid() {
		return errors.New("invalid level")
	}

	words, err := loadURL(url)
	if err != nil {
		return err
	}

	wordMap := make(map[string]Level, len(words))
	for _, word := range words {
		wordMap[word] = level
	}

	return d.AddWords(wordMap)
}

func (d *Detector) LoadDictFromURLs(urls []string) error {
	for _, url := range urls {
		if err := d.LoadDictFromURL(url); err != nil {
			return err
		}
	}
	return nil
}

func (d *Detector) LoadVariantMap(path string) error {
	return normalizer.LoadVariantMap(path)
}

func (d *Detector) Stats() *Stats {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var depth func(*trie.Node, int) int
	depth = func(node *trie.Node, current int) int {
		if node == nil {
			return current
		}
		max := current
		for _, child := range node.Children() {
			if d := depth(child, current+1); d > max {
				max = d
			}
		}
		return max
	}

	return &Stats{
		TotalWords: d.count,
		TreeDepth:  depth(d.tree.Root(), 0),
		MemorySize: int64(d.count) * 400,
	}
}

func (d *Detector) Validate(text string) bool {
	return d.Detect(text).HasSensitive
}

func LoadDictDir(dir string) (map[string]Level, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return nil, err
	}

	words := make(map[string]Level)
	for _, file := range files {
		if strings.HasSuffix(file, ".example.txt") {
			continue
		}

		level := inferLevel(file)
		fileWords, err := loadFile(file)
		if err != nil {
			return nil, err
		}

		for _, word := range fileWords {
			words[word] = level
		}
	}

	return words, nil
}

func inferLevel(path string) Level {
	name := strings.ToLower(filepath.Base(path))

	if strings.HasPrefix(name, "low_") {
		return LevelLow
	}
	if strings.HasPrefix(name, "high_") {
		return LevelHigh
	}

	return LevelMedium
}

func loadFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	words := make([]string, 0, 512)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSuffix(line, ",")
		line = strings.TrimSpace(line)
		if line != "" {
			words = append(words, line)
		}
	}

	return words, scanner.Err()
}

func loadURL(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch dictionary: " + resp.Status)
	}

	words := make([]string, 0, 512)
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSuffix(line, ",")
		line = strings.TrimSpace(line)
		if line != "" {
			words = append(words, line)
		}
	}

	return words, scanner.Err()
}
