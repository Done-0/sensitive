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
	"sync/atomic"

	"github.com/Done-0/sensitive/internal/normalizer"
	"github.com/Done-0/sensitive/internal/pool"
	"github.com/Done-0/sensitive/internal/trie"
)

type Detector struct {
	tree       *trie.Tree
	mu         sync.RWMutex
	normalizer *normalizer.Normalizer
	opts       *Options
	built      atomic.Bool
	count      int
	runePool   sync.Pool
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
		runePool: sync.Pool{
			New: func() any {
				buf := make([]rune, 0, 1024)
				return &buf
			},
		},
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
	d.built.Store(false)
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

	d.tree.Build()
	d.built.Store(true)
	return nil
}

func (d *Detector) Detect(text string) *Result {
	result := &Result{FilteredText: text}
	if text == "" {
		return result
	}

	bufPtr := d.runePool.Get().(*[]rune)
	if cap(*bufPtr) < len(text) {
		*bufPtr = make([]rune, 0, len(text))
	}
	runes := d.normalizer.ToRunes(text, *bufPtr)

	d.mu.RLock()
	if !d.built.Load() {
		d.mu.RUnlock()
		d.runePool.Put(bufPtr)
		return result
	}
	matches := d.tree.SearchDAT(runes)
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

		textRunes := runes
		n := len(textRunes)

		mask := pool.GetBools(n)
		defer pool.PutBools(mask)

		for _, m := range result.Matches {
			for i := m.Start; i < m.End && i < n; i++ {
				(*mask)[i] = true
			}
		}

		filtered := pool.GetRunes(n)
		defer pool.PutRunes(filtered)

		replaceChar := d.opts.ReplaceChar
		if d.opts.FilterStrategy == StrategyMask {
			replaceChar = '*'
		}

		for i, r := range textRunes {
			if (*mask)[i] {
				if d.opts.FilterStrategy != StrategyRemove {
					*filtered = append(*filtered, replaceChar)
				}
			} else {
				*filtered = append(*filtered, r)
			}
		}

		result.FilteredText = string(*filtered)
	}

	*bufPtr = (*bufPtr)[:0]
	d.runePool.Put(bufPtr)
	return result
}

func (d *Detector) Filter(text string) string {
	return d.Detect(text).FilteredText
}

func (d *Detector) Contains(text string) bool {
	if text == "" {
		return false
	}

	bufPtr := d.runePool.Get().(*[]rune)
	if cap(*bufPtr) < len(text) {
		*bufPtr = make([]rune, 0, len(text))
	}
	runes := d.normalizer.ToRunes(text, *bufPtr)

	d.mu.RLock()
	if !d.built.Load() {
		d.mu.RUnlock()
		d.runePool.Put(bufPtr)
		return false
	}
	has := d.tree.Contains(runes)
	d.mu.RUnlock()

	d.runePool.Put(bufPtr)
	return has
}

func (d *Detector) FindFirst(text string) *Match {
	if text == "" {
		return nil
	}

	bufPtr := d.runePool.Get().(*[]rune)
	if cap(*bufPtr) < len(text) {
		*bufPtr = make([]rune, 0, len(text))
	}
	runes := d.normalizer.ToRunes(text, *bufPtr)

	d.mu.RLock()
	if !d.built.Load() {
		d.mu.RUnlock()
		d.runePool.Put(bufPtr)
		return nil
	}
	m := d.tree.FindFirst(runes)
	d.mu.RUnlock()

	d.runePool.Put(bufPtr)
	if m == nil {
		return nil
	}
	return &Match{Word: m.Word, Start: m.Start, End: m.End, Level: Level(m.Level)}
}

func (d *Detector) FindAll(text string) []string {
	result := d.Detect(text)
	if !result.HasSensitive {
		return nil
	}

	seen := make(map[string]struct{}, len(result.Matches))
	words := make([]string, 0, len(result.Matches))
	for _, m := range result.Matches {
		if _, ok := seen[m.Word]; !ok {
			seen[m.Word] = struct{}{}
			words = append(words, m.Word)
		}
	}
	return words
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

	return &Stats{
		TotalWords: d.count,
		TreeDepth:  d.tree.Size(),
		MemorySize: d.tree.MemoryUsage(),
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
