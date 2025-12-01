# Sensitive

[English](README.md) | [简体中文](README.zh-CN.md)

High-performance sensitive word detection library for Go using Aho-Corasick automaton.

## Features

- **High Performance** - Double Array Trie with AC automaton, O(n) complexity
- **Ultra-Fast Build** - 64K dictionary in 50ms
- **Zero Dependencies** - Pure Go implementation
- **Thread-Safe** - Concurrent reads after Build()
- **Fluent API** - Clean builder pattern
- **Built-in Dictionaries** - 64K+ Chinese words
- **Flexible Filtering** - Mask, replace, or remove matches
- **Chinese Support** - Traditional/Simplified conversion

## Installation

```bash
go get github.com/Done-0/sensitive
```

## Quick Start

⚠️ **Important**: This library does NOT load any dictionaries by default. You must explicitly load dictionaries.

### Option 1: Use Built-in Dictionaries

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().
    MustBuild()
```

### Option 2: Use Your Own Dictionary Files

```go
detector := sensitive.NewBuilder().
    LoadDict("path/to/your/dict.txt").
    LoadDict("path/to/another/dict.txt").
    MustBuild()
```

### Option 3: Add Words Manually

```go
detector := sensitive.NewBuilder().
    AddWord("badword", sensitive.LevelHigh).
    AddWord("spam", sensitive.LevelLow).
    MustBuild()
```

### Option 4: Combine All

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().                // Built-in dictionaries
    LoadDict("custom/my_words.txt").  // Your dictionary files
    AddWord("special", sensitive.LevelHigh). // Manual words
    MustBuild()
```

## Built-in Dictionaries

The library embeds 6 dictionaries:

| Constant | File | Level | Words | Description |
|----------|------|-------|-------|-------------|
| `DictHighPolitics` | high_politics.txt | High | ~325 | Political content |
| `DictHighPornography` | high_pornography.txt | High | ~303 | Pornographic content |
| `DictHighViolence` | high_violence.txt | High | ~436 | Violence/weapons/explosives |
| `DictMediumGeneral` | medium_general.txt | Medium | ~48K | General sensitive words |
| `DictLowAd` | low_ad.txt | Low | ~122 | Advertising |
| `DictLowURL` | low_url.txt | Low | ~14K | URL blacklist |

## Usage

### 1. Create Detector

```go
detector := sensitive.NewBuilder().
    WithFilterStrategy(sensitive.StrategyMask).
    WithReplaceChar('*').
    WithCaseSensitive(false).
    LoadAllEmbedded().
    MustBuild()
```

### 2. Add Words

```go
// Single word
detector.AddWord("badword", sensitive.LevelHigh)

// Multiple words
words := map[string]sensitive.Level{
    "illegal":  sensitive.LevelHigh,
    "violence": sensitive.LevelHigh,
    "abuse":    sensitive.LevelMedium,
    "spam":     sensitive.LevelLow,
}
detector.AddWords(words)
```

### 3. Load Dictionary

**Built-in dictionaries:**

```go
detector.LoadAllEmbedded()  // All 6 dictionaries
detector.LoadEmbeddedDict(sensitive.DictHighPolitics, sensitive.LevelHigh)  // Specific
```

**Custom dictionaries:**

```go
detector.LoadDict("custom/my_words.txt")  // Auto-detect level from filename
detector.LoadDictWithLevel("any_name.txt", sensitive.LevelHigh)  // Explicit level
```

**From URL:**

```go
detector.LoadDictFromURL("https://example.com/dict.txt")
```

**File naming (auto-level detection):**
- `high_*.txt` → LevelHigh
- `medium_*.txt` → LevelMedium
- `low_*.txt` → LevelLow
- Other → LevelMedium (default)

### 4. Configure Options

```go
// Filter strategy
detector.WithFilterStrategy(sensitive.StrategyMask)     // "bad" → "***"
detector.WithFilterStrategy(sensitive.StrategyReplace).WithReplaceChar('█')  // "bad" → "███"
detector.WithFilterStrategy(sensitive.StrategyRemove)    // "bad" → ""

// Case sensitivity
detector.WithCaseSensitive(false)  // "TEST", "test", "Test" all match (default)
detector.WithCaseSensitive(true)   // Only exact case

// Skip whitespace
detector.WithSkipWhitespace(true)  // "b a d" matches "bad"

// Traditional/Simplified Chinese
detector.WithVariant(true).LoadVariantMap("variant_map.txt")
```

### 5. Detect Content

```go
// Simple validation
if detector.Validate(text) {
    return errors.New("content rejected")
}

// Get details
result := detector.Detect(text)
if result.HasSensitive {
    for _, match := range result.Matches {
        fmt.Printf("Word: %s, Level: %s, Position: %d-%d\n",
            match.Word, match.Level, match.Start, match.End)
    }
    fmt.Println("Filtered:", result.FilteredText)
}

// Filter only
filtered := detector.Filter(text)
```

### 6. Error Handling

```go
// Build() returns error
detector, err := sensitive.NewBuilder().LoadAllEmbedded().Build()
if err != nil {
    log.Fatal(err)
}

// MustBuild() panics on error (use in init())
detector := sensitive.NewBuilder().LoadAllEmbedded().MustBuild()
```

### 7. Concurrent Usage

```go
var detector *sensitive.Detector

func init() {
    detector = sensitive.NewBuilder().LoadAllEmbedded().MustBuild()
}

// Thread-safe after Build()
func handler(text string) error {
    if detector.Validate(text) {
        return errors.New("sensitive content")
    }
    return nil
}
```

⚠️ **Not safe:** Adding words after Build() in concurrent environment

### 8. Performance

```
BenchmarkDAT_Build_1KWords-12              886 μs/op   29.9 MB/op    2753 allocs/op
BenchmarkDAT_Build_10KWords-12            3.60 ms/op   30.1 MB/op   20753 allocs/op
BenchmarkDetector_Detect_SmallDict-12     79.0 μs/op   34.0 KB/op       5 allocs/op
BenchmarkDetector_Detect_ShortText-12     8.17 μs/op    3.9 KB/op       5 allocs/op
BenchmarkDetector_Detect_LongText-12       778 μs/op    328 KB/op       5 allocs/op
BenchmarkDetector_AddWord-12               126 ns/op     32 B/op        2 allocs/op
BenchmarkDetector_Parallel-12             10.7 μs/op   34.1 KB/op       5 allocs/op
```

- Double Array Trie implementation for O(n) search complexity
- Load dictionaries once in `init()`
- Reuse detector across goroutines (thread-safe after Build())
- Memory pool used internally

## Custom Dictionaries

Place your dictionary files anywhere in your project:

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().                   // Optional: built-in dictionaries
    LoadDict("dict/high_banned.txt").    // Your dictionary in project
    LoadDict("configs/custom.txt").      // Another location
    MustBuild()
```

**Git exclusion**: Files named `custom_*.txt`, `local_*.txt`, `user_*.txt` in `configs/dict/` are auto-excluded.

## Examples

See [examples/](examples/) for production-ready code:

| Example | Description |
|---------|-------------|
| [fluent_api](examples/fluent_api/) | Fluent API chain calls |
| [quickstart](examples/quickstart/) | Simplest usage |
| [web_api](examples/web_api/) | HTTP REST API service |
| [comment_filter](examples/comment_filter/) | Content moderation system |
| [dependency_injection](examples/dependency_injection/) | DI pattern |
| [high_concurrency](examples/high_concurrency/) | Concurrent processing |

Run example:
```bash
cd examples/fluent_api
go run main.go
```

## License

[MIT License](LICENSE)
