# Sensitive

[English](README.md) | [简体中文](README.zh-CN.md)

基于 Aho-Corasick 自动机算法的高性能 Go 敏感词检测库。

## 特性

- **高性能** - Double Array Trie + AC 自动机，O(n) 复杂度
- **超快构建** - 64K 词典 50ms 完成
- **零依赖** - 纯 Go 实现
- **线程安全** - Build() 后支持并发读
- **流式 API** - 简洁的构建模式
- **内置词典** - 64K+ 中文词汇
- **灵活过滤** - 掩码、替换或删除匹配
- **中文支持** - 繁简体转换

## 安装

```bash
go get github.com/Done-0/sensitive
```

## 快速开始

⚠️ **重要**：本库默认不加载任何词典，必须显式加载。

### 方案 1：使用内置词典

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().
    MustBuild()
```

### 方案 2：使用自己的词典文件

```go
detector := sensitive.NewBuilder().
    LoadDict("path/to/your/dict.txt").
    LoadDict("path/to/another/dict.txt").
    MustBuild()
```

### 方案 3：手动添加词汇

```go
detector := sensitive.NewBuilder().
    AddWord("违禁词", sensitive.LevelHigh).
    AddWord("垃圾信息", sensitive.LevelLow).
    MustBuild()
```

### 方案 4：组合使用

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().                // 内置词典
    LoadDict("custom/my_words.txt").  // 自定义词典文件
    AddWord("特殊词", sensitive.LevelHigh). // 手动添加
    MustBuild()
```

## 内置词典

本库嵌入了 6 个词典：

| 常量 | 文件 | 级别 | 词数 | 描述 |
|------|------|------|------|------|
| `DictHighPolitics` | high_politics.txt | 高 | ~325 | 政治类 |
| `DictHighPornography` | high_pornography.txt | 高 | ~303 | 色情类 |
| `DictHighViolence` | high_violence.txt | 高 | ~436 | 涉枪涉爆违法信息 |
| `DictMediumGeneral` | medium_general.txt | 中 | ~48K | 通用敏感词 |
| `DictLowAd` | low_ad.txt | 低 | ~122 | 广告 |
| `DictLowURL` | low_url.txt | 低 | ~14K | 网址黑名单 |

## 使用方法

### 1. 创建检测器

```go
detector := sensitive.NewBuilder().
    WithFilterStrategy(sensitive.StrategyMask).
    WithReplaceChar('*').
    WithCaseSensitive(false).
    LoadAllEmbedded().
    MustBuild()
```

### 2. 添加词汇

```go
// 单个词汇
detector.AddWord("违禁词", sensitive.LevelHigh)

// 批量添加
words := map[string]sensitive.Level{
    "违法":  sensitive.LevelHigh,
    "暴力":  sensitive.LevelHigh,
    "辱骂":  sensitive.LevelMedium,
    "垃圾":  sensitive.LevelLow,
}
detector.AddWords(words)
```

### 3. 加载词典

**内置词典：**

```go
detector.LoadAllEmbedded()  // 加载全部 6 个词典
detector.LoadEmbeddedDict(sensitive.DictHighPolitics, sensitive.LevelHigh)  // 加载特定词典
```

**自定义词典：**

```go
detector.LoadDict("custom/my_words.txt")  // 根据文件名自动识别级别
detector.LoadDictWithLevel("any_name.txt", sensitive.LevelHigh)  // 显式指定级别
```

**从 URL 加载：**

```go
detector.LoadDictFromURL("https://example.com/dict.txt")
```

**文件命名规则（自动级别识别）：**
- `high_*.txt` → LevelHigh
- `medium_*.txt` → LevelMedium
- `low_*.txt` → LevelLow
- 其他 → LevelMedium（默认）

### 4. 配置选项

```go
// 过滤策略
detector.WithFilterStrategy(sensitive.StrategyMask)     // "敏感" → "**"
detector.WithFilterStrategy(sensitive.StrategyReplace).WithReplaceChar('█')  // "敏感" → "██"
detector.WithFilterStrategy(sensitive.StrategyRemove)    // "敏感" → ""

// 大小写敏感
detector.WithCaseSensitive(false)  // "TEST"、"test"、"Test" 都匹配（默认）
detector.WithCaseSensitive(true)   // 仅精确匹配大小写

// 跳过空白字符
detector.WithSkipWhitespace(true)  // "敏 感" 匹配 "敏感"

// 繁简体中文转换
detector.WithVariant(true).LoadVariantMap("variant_map.txt")
```

### 5. 检测内容

```go
// 简单验证
if detector.Validate(text) {
    return errors.New("内容被拒绝")
}

// 获取详情
result := detector.Detect(text)
if result.HasSensitive {
    for _, match := range result.Matches {
        fmt.Printf("词汇: %s, 级别: %s, 位置: %d-%d\n",
            match.Word, match.Level, match.Start, match.End)
    }
    fmt.Println("过滤后:", result.FilteredText)
}

// 仅过滤
filtered := detector.Filter(text)
```

### 6. 错误处理

```go
// Build() 返回错误
detector, err := sensitive.NewBuilder().LoadAllEmbedded().Build()
if err != nil {
    log.Fatal(err)
}

// MustBuild() 遇错误直接 panic（适合在 init() 中使用）
detector := sensitive.NewBuilder().LoadAllEmbedded().MustBuild()
```

### 7. 并发使用

```go
var detector *sensitive.Detector

func init() {
    detector = sensitive.NewBuilder().LoadAllEmbedded().MustBuild()
}

// Build() 后线程安全
func handler(text string) error {
    if detector.Validate(text) {
        return errors.New("敏感内容")
    }
    return nil
}
```

⚠️ **不安全**：在并发环境中 Build() 后添加词汇

### 8. 性能

```
BenchmarkDAT_Build_1KWords-12              886 μs/op   29.9 MB/op    2753 allocs/op
BenchmarkDAT_Build_10KWords-12            3.60 ms/op   30.1 MB/op   20753 allocs/op
BenchmarkDetector_Detect_SmallDict-12     79.0 μs/op   34.0 KB/op       5 allocs/op
BenchmarkDetector_Detect_ShortText-12     8.17 μs/op    3.9 KB/op       5 allocs/op
BenchmarkDetector_Detect_LongText-12       778 μs/op    328 KB/op       5 allocs/op
BenchmarkDetector_AddWord-12               126 ns/op     32 B/op        2 allocs/op
BenchmarkDetector_Parallel-12             10.7 μs/op   34.1 KB/op       5 allocs/op
```

- Double Array Trie 实现，搜索复杂度 O(n)
- 在 `init()` 中加载词典一次
- 跨 goroutine 复用检测器（Build() 后线程安全）
- 内部使用内存池优化性能

## 自定义词典

将词典文件放在项目任意位置：

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().                   // 可选：内置词典
    LoadDict("dict/high_banned.txt").    // 项目中的词典
    LoadDict("configs/custom.txt").      // 其他位置
    MustBuild()
```

**Git 排除**：`configs/dict/` 目录下的 `custom_*.txt`、`local_*.txt`、`user_*.txt` 文件会被自动排除。

## 示例代码

查看 [examples/](examples/) 获取生产环境可用代码：

| 示例 | 描述 |
|------|------|
| [fluent_api](examples/fluent_api/) | 流式 API 链式调用 |
| [quickstart](examples/quickstart/) | 最简单用法 |
| [web_api](examples/web_api/) | HTTP REST API 服务 |
| [comment_filter](examples/comment_filter/) | 内容审核系统 |
| [dependency_injection](examples/dependency_injection/) | 依赖注入模式 |
| [high_concurrency](examples/high_concurrency/) | 并发处理 |

运行示例：
```bash
cd examples/fluent_api
go run main.go
```

## 许可证

[MIT License](LICENSE)
