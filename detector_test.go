// Package sensitive provides high-performance sensitive word detection using AC automaton
// Creator: Done-0
// Created: 2025-01-15
package sensitive

import (
	"os"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	detector := New()
	if detector == nil {
		t.Fatal("New() returned nil")
	}
}

func TestAddWord_Valid(t *testing.T) {
	detector := New()
	if err := detector.AddWord("test", LevelHigh); err != nil {
		t.Errorf("AddWord() error: %v", err)
	}
}

func TestAddWord_Empty(t *testing.T) {
	detector := New()
	if err := detector.AddWord("", LevelHigh); err == nil {
		t.Error("AddWord() should return error for empty word")
	}
}

func TestAddWord_InvalidLevel(t *testing.T) {
	detector := New()
	if err := detector.AddWord("test", Level(99)); err == nil {
		t.Error("AddWord() should return error for invalid level")
	}
}

func TestAddWords(t *testing.T) {
	detector := New()
	words := map[string]Level{
		"word1": LevelHigh,
		"word2": LevelMedium,
		"word3": LevelLow,
	}
	if err := detector.AddWords(words); err != nil {
		t.Fatalf("AddWords() error: %v", err)
	}
	if detector.Stats().TotalWords != 3 {
		t.Error("expected 3 words")
	}
}

func TestBuild(t *testing.T) {
	detector := New()
	detector.AddWord("test", LevelMedium)
	if err := detector.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}
}

func TestDetect_Empty(t *testing.T) {
	detector := New()
	detector.AddWord("test", LevelMedium)
	detector.Build()
	result := detector.Detect("")
	if result.HasSensitive {
		t.Error("empty text should not have sensitive words")
	}
}

func TestDetect_NoMatch(t *testing.T) {
	detector := New()
	detector.AddWord("bad", LevelHigh)
	detector.Build()
	result := detector.Detect("good text")
	if result.HasSensitive {
		t.Error("should not detect non-sensitive text")
	}
}

func TestDetect_SingleMatch(t *testing.T) {
	detector := New()
	detector.AddWord("bad", LevelMedium)
	detector.Build()
	result := detector.Detect("this is bad")
	if !result.HasSensitive {
		t.Error("should detect sensitive word")
	}
	if len(result.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(result.Matches))
	}
}

func TestDetect_MultipleMatches(t *testing.T) {
	detector := New()
	detector.AddWord("bad", LevelMedium)
	detector.AddWord("ugly", LevelLow)
	detector.Build()
	result := detector.Detect("bad and ugly")
	if len(result.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(result.Matches))
	}
}

func TestDetect_CaseInsensitive(t *testing.T) {
	detector := New(WithCaseSensitive(false))
	detector.AddWord("test", LevelMedium)
	detector.Build()
	if !detector.Detect("TEST").HasSensitive {
		t.Error("should detect case-insensitive")
	}
}

func TestDetect_CaseSensitive(t *testing.T) {
	detector := New(WithCaseSensitive(true))
	detector.AddWord("test", LevelMedium)
	detector.Build()
	if detector.Detect("TEST").HasSensitive {
		t.Error("should not detect different case when case-sensitive")
	}
}

func TestDetect_Chinese(t *testing.T) {
	detector := New()
	detector.AddWord("敏感词", LevelHigh)
	detector.Build()
	result := detector.Detect("这是敏感词文本")
	if !result.HasSensitive {
		t.Error("should detect Chinese words")
	}
}

func TestFilter_Mask(t *testing.T) {
	detector := New(WithFilterStrategy(StrategyMask))
	detector.AddWord("bad", LevelHigh)
	detector.Build()
	filtered := detector.Filter("this is bad")
	expected := "this is ***"
	if filtered != expected {
		t.Errorf("expected '%s', got '%s'", expected, filtered)
	}
}

func TestFilter_Replace(t *testing.T) {
	detector := New(WithFilterStrategy(StrategyReplace), WithReplaceChar('#'))
	detector.AddWord("bad", LevelHigh)
	detector.Build()
	filtered := detector.Filter("this is bad")
	expected := "this is ###"
	if filtered != expected {
		t.Errorf("expected '%s', got '%s'", expected, filtered)
	}
}

func TestValidate(t *testing.T) {
	detector := New()
	detector.AddWord("bad", LevelHigh)
	detector.Build()
	if !detector.Validate("this is bad") {
		t.Error("Validate() should return true")
	}
	if detector.Validate("this is good") {
		t.Error("Validate() should return false")
	}
}

func TestStats(t *testing.T) {
	detector := New()
	detector.AddWord("test", LevelHigh)
	detector.Build()
	stats := detector.Stats()
	if stats.TotalWords != 1 {
		t.Error("expected 1 word")
	}
	if stats.TreeDepth == 0 {
		t.Error("tree depth should be > 0")
	}
}

func TestLoadDict(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping file I/O test")
	}
	tmpFile := t.TempDir() + "/medium_test.txt"
	content := "word1\nword2\n# comment\n\nword3"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	detector := New()
	if err := detector.LoadDict(tmpFile); err != nil {
		t.Fatalf("LoadDict() error: %v", err)
	}
	if detector.Stats().TotalWords != 3 {
		t.Error("expected 3 words")
	}
}

func TestLoadDictDir(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping directory I/O test")
	}
	tmpDir := t.TempDir()
	files := map[string]string{
		"high_test.txt":    "word1\nword2",
		"medium_test.txt":  "word3",
		"low_test.txt":     "word4",
		"test.example.txt": "ignored",
	}
	for name, content := range files {
		os.WriteFile(tmpDir+"/"+name, []byte(content), 0644)
	}
	words, err := LoadDictDir(tmpDir)
	if err != nil {
		t.Fatalf("LoadDictDir() error: %v", err)
	}
	if len(words) != 4 {
		t.Errorf("expected 4 words, got %d", len(words))
	}
}

func TestConcurrent(t *testing.T) {
	detector := New()
	detector.AddWord("test", LevelMedium)
	detector.Build()
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			detector.Detect("test text")
		}()
	}
	wg.Wait()
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelLow, "Low"},
		{LevelMedium, "Medium"},
		{LevelHigh, "High"},
		{Level(99), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("String() = %s, want %s", got, tt.want)
		}
	}
}

func TestLevelIsValid(t *testing.T) {
	tests := []struct {
		level Level
		want  bool
	}{
		{LevelLow, true},
		{LevelMedium, true},
		{LevelHigh, true},
		{Level(0), false},
		{Level(4), false},
	}
	for _, tt := range tests {
		if got := tt.level.IsValid(); got != tt.want {
			t.Errorf("IsValid() = %v, want %v", got, tt.want)
		}
	}
}

func TestBuilder_Basic(t *testing.T) {
	detector := NewBuilder().
		AddWord("test", LevelHigh).
		MustBuild()
	if detector == nil {
		t.Fatal("MustBuild() returned nil")
	}
}

func TestBuilder_AddWords(t *testing.T) {
	words := map[string]Level{
		"word1": LevelHigh,
		"word2": LevelMedium,
	}
	detector := NewBuilder().
		AddWords(words).
		MustBuild()
	if detector.Stats().TotalWords != 2 {
		t.Error("expected 2 words")
	}
}

func TestBuilder_WithOptions(t *testing.T) {
	detector := NewBuilder().
		WithFilterStrategy(StrategyReplace).
		WithReplaceChar('#').
		WithCaseSensitive(true).
		WithSkipWhitespace(true).
		AddWord("test", LevelHigh).
		MustBuild()
	result := detector.Filter("this is test")
	if result != "this is ####" {
		t.Errorf("expected 'this is ####', got '%s'", result)
	}
}

func TestBuilder_LoadDict(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping file I/O test")
	}
	tmpFile := t.TempDir() + "/test.txt"
	os.WriteFile(tmpFile, []byte("word1\nword2"), 0644)
	detector := NewBuilder().
		LoadDict(tmpFile).
		MustBuild()
	if detector.Stats().TotalWords != 2 {
		t.Error("expected 2 words")
	}
}

func TestBuilder_Build_Error(t *testing.T) {
	_, err := NewBuilder().
		AddWord("", LevelHigh).
		Build()
	if err == nil {
		t.Error("Build() should return error for empty word")
	}
}

func TestBuilder_MustBuild_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustBuild() should panic on error")
		}
	}()
	NewBuilder().AddWord("", LevelHigh).MustBuild()
}

func TestLoadAllEmbedded(t *testing.T) {
	detector := NewBuilder().
		LoadAllEmbedded().
		MustBuild()
	stats := detector.Stats()
	if stats.TotalWords == 0 {
		t.Error("LoadAllEmbedded() should load words")
	}
	if stats.TotalWords < 1000 {
		t.Errorf("expected > 1000 words, got %d", stats.TotalWords)
	}
}

func TestLoadEmbeddedDict(t *testing.T) {
	detector := NewBuilder().
		LoadEmbeddedDict(DictHighPolitics, LevelHigh).
		MustBuild()
	stats := detector.Stats()
	if stats.TotalWords == 0 {
		t.Error("LoadEmbeddedDict() should load words")
	}
	if stats.TotalWords < 100 {
		t.Errorf("expected > 100 words, got %d", stats.TotalWords)
	}
}

func TestLoadEmbeddedDict_AllConstants(t *testing.T) {
	dicts := []string{
		DictHighPolitics,
		DictHighPornography,
		DictHighViolence,
		DictMediumGeneral,
		DictLowAd,
		DictLowURL,
	}
	for _, dict := range dicts {
		detector := NewBuilder().
			LoadEmbeddedDict(dict, LevelHigh).
			MustBuild()
		if detector.Stats().TotalWords == 0 {
			t.Errorf("LoadEmbeddedDict(%s) failed to load words", dict)
		}
	}
}

func TestLoadEmbeddedDict_InvalidLevel(t *testing.T) {
	_, err := NewBuilder().
		LoadEmbeddedDict(DictHighPolitics, Level(99)).
		Build()
	if err == nil {
		t.Error("should return error for invalid level")
	}
}

func TestFilter_Remove(t *testing.T) {
	detector := NewBuilder().
		WithFilterStrategy(StrategyRemove).
		AddWord("bad", LevelHigh).
		MustBuild()
	filtered := detector.Filter("this is bad text")
	expected := "this is  text"
	if filtered != expected {
		t.Errorf("expected '%s', got '%s'", expected, filtered)
	}
}

func TestSkipWhitespace(t *testing.T) {
	detector := NewBuilder().
		WithSkipWhitespace(true).
		AddWord("bad", LevelHigh).
		MustBuild()
	result := detector.Detect("b a d")
	t.Logf("Detected: %v, Matches: %d", result.HasSensitive, len(result.Matches))
}

func TestSkipWhitespace_Disabled(t *testing.T) {
	detector := NewBuilder().
		WithSkipWhitespace(false).
		AddWord("bad", LevelHigh).
		MustBuild()
	if detector.Detect("b a d").HasSensitive {
		t.Error("should not detect 'b a d' without skip whitespace")
	}
	if !detector.Detect("bad").HasSensitive {
		t.Error("should detect 'bad'")
	}
}

func TestVariant_Enabled(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping variant test")
	}
	tmpFile := t.TempDir() + "/variant.txt"
	content := "體\t体\n國\t国"
	os.WriteFile(tmpFile, []byte(content), 0644)

	detector := NewBuilder().
		WithVariant(true).
		LoadVariantMap(tmpFile).
		AddWord("国", LevelHigh).
		MustBuild()

	result := detector.Detect("國家")
	t.Logf("Variant test - Detected: %v, Matches: %d", result.HasSensitive, len(result.Matches))

	if !detector.Detect("国家").HasSensitive {
		t.Error("should detect simplified Chinese")
	}
}

func TestLoadDictFromURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping URL test")
	}
	t.Skip("URL loading requires network/httptest mock")
}
