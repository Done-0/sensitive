# 词典文件使用指南

## 内置嵌入词典

本库使用 Go embed 技术将词典文件编译进二进制，无需外部文件依赖。

### 加载所有内置词典（推荐）

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().
    MustBuild()
```

### 加载特定内置词典

```go
detector := sensitive.NewBuilder().
    LoadEmbeddedDict(sensitive.DictHighPolitics, sensitive.LevelHigh).
    LoadEmbeddedDict(sensitive.DictMediumGeneral, sensitive.LevelMedium).
    MustBuild()
```

**可用常量**：
- `sensitive.DictHighPolitics` - 政治类
- `sensitive.DictHighPornography` - 色情类
- `sensitive.DictHighViolence` - 涉枪涉爆违法信息
- `sensitive.DictMediumGeneral` - 通用敏感词
- `sensitive.DictLowAd` - 广告
- `sensitive.DictLowURL` - 网址黑名单

## 自定义词典

### 加载本地词典文件

```go
// 从项目目录加载
detector := sensitive.NewBuilder().
    LoadDict("configs/dict/my_words.txt").
    LoadDict("custom/company_terms.txt").
    MustBuild()

// 根据文件名前缀自动识别级别
detector.LoadDict("high_illegal.txt")    // 自动识别为 LevelHigh
detector.LoadDict("medium_abuse.txt")    // 自动识别为 LevelMedium
detector.LoadDict("low_spam.txt")        // 自动识别为 LevelLow

// 明确指定级别（覆盖文件名前缀）
detector.LoadDictWithLevel("any_name.txt", sensitive.LevelHigh)
```

### 文件格式规范

**文件命名（自动级别识别）**：
- `high_*.txt` → LevelHigh
- `medium_*.txt` → LevelMedium
- `low_*.txt` → LevelLow
- 其他 → LevelMedium（默认）

**文件内容（UTF-8编码，每行一词）**：
```
# 以 # 开头的是注释
词汇1
词汇2,
词汇3
```

### 从 URL 加载远程词典

```go
// 单个 URL
detector := sensitive.NewBuilder().
    LoadDictFromURL("https://example.com/dict/sensitive.txt").
    MustBuild()

// 多个 URL
urls := []string{
    "https://example.com/dict/high_politics.txt",
    "https://example.com/dict/medium_general.txt",
}
detector := sensitive.NewBuilder().
    LoadDictFromURLs(urls).
    MustBuild()
```

### 组合使用

```go
detector := sensitive.NewBuilder().
    LoadAllEmbedded().                        // 内置词典（可选）
    LoadDict("custom/my_words.txt").          // 本地词典
    LoadDictFromURL("https://example.com/dict.txt"). // 远程词典
    MustBuild()
```

## 内置词典详情

本库使用 `//go:embed` 嵌入 6 个词典文件：

| 常量 | 文件名 | 级别 | 词数 | 描述 |
|------|--------|------|------|------|
| `DictHighPolitics` | high_politics.txt | 高 | ~325 | 政治类 |
| `DictHighPornography` | high_pornography.txt | 高 | ~303 | 色情类 |
| `DictHighViolence` | high_violence.txt | 高 | ~436 | 涉枪涉爆违法信息 |
| `DictMediumGeneral` | medium_general.txt | 中 | ~48K | 通用敏感词 |
| `DictLowAd` | low_ad.txt | 低 | ~122 | 广告 |
| `DictLowURL` | low_url.txt | 低 | ~14K | 网址黑名单 |

**⚠️ 无默认加载**：本库默认不加载任何词典，必须显式调用。原因：不同应用需求不同、法律合规因地区而异、避免误拦截。

## 繁简体中文转换

```go
detector := sensitive.NewBuilder().
    WithVariant(true).
    LoadVariantMap("variant_map.txt").
    MustBuild()
```

映射文件格式（繁体[TAB]简体）：

```
體	体
國	国
```

来源：[OpenCC 项目](https://github.com/BYVoid/OpenCC)

## 嵌入机制

本库使用 Go 1.16+ 的 `//go:embed` 指令实现词典嵌入，编译时打包进二进制，用户通过 `go get` 安装后无需额外文件。

## Git 管理

⚠️ **绝不将自定义敏感词提交到公开代码仓库**

内置词典（`high_*.txt`, `medium_*.txt`, `low_*.txt`）已提交到 Git，用户自定义词典（`custom_*.txt`, `local_*.txt`, `user_*.txt`）被 `.gitignore` 排除。

## 词典来源

- [Tencent Sensitive Words](https://github.com/cjhnim/tencent-sensitive-words) - 腾讯离线词库
- [houbb/sensitive-word](https://github.com/fwwdn/sensitive-stop-words) - 互联网常用敏感词、停止词库
