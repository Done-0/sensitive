# Sensitive

[English](README.md) | [简体中文](README.zh-CN.md)

High-performance sensitive word detection library for Go using Aho-Corasick automaton.

## Features

- **High Performance** - Double Array Trie with AC automaton, O(n) complexity
- **High Concurrency** - sync.RWMutex + sync.Pool, 6x faster than alternatives
- **Zero Allocation** - Hot path (Contains, FindFirst) with 0 allocs
- **Multi-Language** - Full Unicode support (CJK, Cyrillic, Arabic, etc.)
- **Thread-Safe** - Concurrent reads after Build()
- **Fluent API** - Clean builder pattern
- **Built-in Dictionaries** - 64K+ words included
- **Flexible Filtering** - Mask, replace, or remove matches
- **Chinese Support** - Traditional/Simplified conversion, Full-width/Half-width

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

**Benchmark Environment:** Apple M2 Max, Go 1.25, 1000 words dictionary, mixed Chinese/English text

**Comparison with popular Go libraries:**

| Benchmark | Done-0/sensitive | importcjj/sensitive | anknown/ahocorasick |
|-----------|------------------|---------------------|---------------------|
| **Contains** | 36.6 μs, **0B**, 0 allocs | 89.4 μs, 42KB, 15 allocs | 24.1 μs, 0B, 0 allocs |
| **FindAll** | 37.0 μs, 752B, 2 allocs | 21.5 μs, 13KB, 1 alloc | 23.5 μs, 0B, 0 allocs |
| **Filter** | 36.8 μs, 752B, 2 allocs | 37.1 μs, 19KB, 2 allocs | N/A |
| **Parallel (12-core)** | **4.3 μs**, ~0B, 0 allocs | 27.0 μs, 46KB, 15 allocs | 2.7 μs, 0B, 0 allocs |
| **Short Text (100 chars)** | 678 ns, 0B, 0 allocs | 1.59 μs, 461B, 6 allocs | 398 ns, 0B, 0 allocs |
| **Long Text (10K chars)** | **367 μs**, ~0B, 0 allocs | 1.35 ms, 393KB, 22 allocs | 239 μs, 0B, 0 allocs |

**Key Advantages:**

- ✅ **Zero allocation** in hot path (Contains, FindFirst)
- ✅ **High concurrency**: 4.3μs on 12-core parallel, 6x faster than importcjj
- ✅ **26x less memory** than importcjj/sensitive in Filter
- ✅ **3.7x faster** for long text vs importcjj
- ✅ **Full-featured**: Filter, levels, variant support (vs ahocorasick's search-only)
- ✅ **Thread-safe**: sync.RWMutex + sync.Pool optimization

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
