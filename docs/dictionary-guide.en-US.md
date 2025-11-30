# Dictionary File Guide

## Built-in Embedded Dictionaries

This library uses Go embed to compile dictionaries into the binary, requiring no external file dependencies.

### Load All Built-in Dictionaries (Recommended)

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().
    MustBuild()
```

### Load Specific Built-in Dictionaries

```go
detector := sensitive.NewBuilder().
    LoadEmbeddedDict(sensitive.DictHighPolitics, sensitive.LevelHigh).
    LoadEmbeddedDict(sensitive.DictMediumGeneral, sensitive.LevelMedium).
    MustBuild()
```

**Available constants:**
- `sensitive.DictHighPolitics` - Political content
- `sensitive.DictHighPornography` - Pornographic content
- `sensitive.DictHighViolence` - Violence/weapons/explosives
- `sensitive.DictMediumGeneral` - General sensitive words
- `sensitive.DictLowAd` - Advertising
- `sensitive.DictLowURL` - URL blacklist

## Custom Dictionaries

### Load Local Dictionary Files

```go
// Load from your project directory
detector := sensitive.NewBuilder().
    LoadDict("configs/dict/my_words.txt").
    LoadDict("custom/company_terms.txt").
    MustBuild()

// Auto-detect level by filename prefix
detector.LoadDict("high_illegal.txt")    // Auto as LevelHigh
detector.LoadDict("medium_abuse.txt")    // Auto as LevelMedium
detector.LoadDict("low_spam.txt")        // Auto as LevelLow

// Explicit level (overrides filename prefix)
detector.LoadDictWithLevel("any_name.txt", sensitive.LevelHigh)
```

### File Format Specification

**File naming (auto-level detection):**
- `high_*.txt` → LevelHigh
- `medium_*.txt` → LevelMedium
- `low_*.txt` → LevelLow
- Others → LevelMedium (default)

**File content (UTF-8, one word per line):**
```
# Comments start with #
word1
word2,
word3
```

### Load Remote Dictionaries from URL

```go
// Single URL
detector := sensitive.NewBuilder().
    LoadDictFromURL("https://example.com/dict/sensitive.txt").
    MustBuild()

// Multiple URLs
urls := []string{
    "https://example.com/dict/high_politics.txt",
    "https://example.com/dict/medium_general.txt",
}
detector := sensitive.NewBuilder().
    LoadDictFromURLs(urls).
    MustBuild()
```

### Combined Usage

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().                       // Built-in dictionaries (optional)
    LoadDict("custom/my_words.txt").         // Local dictionary
    LoadDictFromURL("https://example.com/dict.txt"). // Remote dictionary
    MustBuild()
```

## Built-in Dictionary Details

This library uses `//go:embed` to embed 6 dictionary files:

| Constant | File | Level | Words | Description |
|----------|------|-------|-------|-------------|
| `DictHighPolitics` | high_politics.txt | High | ~325 | Political content |
| `DictHighPornography` | high_pornography.txt | High | ~303 | Pornographic content |
| `DictHighViolence` | high_violence.txt | High | ~436 | Violence/weapons/explosives |
| `DictMediumGeneral` | medium_general.txt | Medium | ~48K | General sensitive words (Tencent database) |
| `DictLowAd` | low_ad.txt | Low | ~122 | Advertising |
| `DictLowURL` | low_url.txt | Low | ~14K | URL blacklist |

**⚠️ No Default Loading**: This library does NOT load any dictionaries by default. You must explicitly call loading methods. Reason: Different apps need different dictionaries, legal/compliance varies by region, prevents accidental blocking.

## Traditional/Simplified Chinese

```go
detector := sensitive.NewBuilder().
    WithVariant(true).
    LoadVariantMap("variant_map.txt").
    MustBuild()
```

Mapping file format (Traditional[TAB]Simplified):

```
體	体
國	国
```

Source: [OpenCC Project](https://github.com/BYVoid/OpenCC)

## Embed Mechanism

This library uses Go 1.16+ `//go:embed` directive to embed dictionaries:

```go
//go:embed configs/dict/*.txt
var dictFS embed.FS
```

**Advantages:**
- Dictionaries are compiled into the binary at build time
- No external file dependencies after `go get` installation
- Dictionaries distributed with the library, ready to use out-of-the-box
- Simple deployment without external files

## Git Management Strategy

⚠️ **Never commit custom sensitive dictionaries to public repositories**

- Built-in dictionaries (`high_*.txt`, `medium_*.txt`, `low_*.txt`) are committed to Git
- User custom dictionaries (`custom_*.txt`, `local_*.txt`, `user_*.txt`) are excluded by `.gitignore`
- Review and validate dictionary contents before production use

## Dictionary Sources

- [Tencent Sensitive Words](https://github.com/cjhnim/tencent-sensitive-words) - Offline dictionary
- [houbb/sensitive-word](https://github.com/fwwdn/sensitive-stop-words) - Internet commonly used sensitive words and stopword database
