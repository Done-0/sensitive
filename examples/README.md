# Production Examples

Real-world usage examples for production environments.

## Fluent API

**Elegant chain calls** - Modern fluent interface design.

```bash
cd examples/fluent_api
go run main.go
```

Use case: Clean, readable configuration.

Features:
- Method chaining (Builder pattern)
- Compile-time type safety
- Zero runtime overhead
- Follows Go idioms

```go
detector := sensitive.NewBuilder().
    WithFilterStrategy(sensitive.StrategyMask).
    WithReplaceChar('*').
    AddWord("badword", sensitive.LevelHigh).
    AddWords(words).
    LoadDict("dict/sensitive.txt").
    MustBuild()
```

## Quickstart

**Simplest usage** - Content validation in 20 lines.

```bash
cd examples/quickstart
go run main.go
```

Use case: Basic content approval/rejection.

## Web API

**HTTP content moderation service** - RESTful API for content checking.

```bash
cd examples/web_api
go run main.go
```

Test:
```bash
curl -X POST http://localhost:8080/api/check \
  -H "Content-Type: application/json" \
  -d '{"content":"This has spam in it"}'
```

Response:
```json
{
  "safe": false,
  "reason": "Content contains sensitive words",
  "matches": ["spam"],
  "filtered": "This has **** in it"
}
```

Use case: Content moderation microservice, API gateway integration.

## Comment Filter

**User-generated content filtering** - Comment moderation system with approval workflow.

```bash
cd examples/comment_filter
go run main.go
```

Use case: Forum comments, user reviews, chat messages.

Features:
- High severity = reject
- Medium/low severity = auto-filter
- Clean content = approve

## Dependency Injection

**Dependency injection pattern** - Enterprise application architecture.

```bash
cd examples/dependency_injection
go run main.go
```

Use case: Large-scale applications with multiple services.

Features:
- Interface-based design
- Loose coupling between components
- Easy unit testing with mocks
- Clean architecture

## High Concurrency

**High-traffic production system** - Thread-safe concurrent processing.

```bash
cd examples/high_concurrency
go run main.go
```

Use case: High-throughput services, real-time content moderation.

Features:
- 1000+ concurrent goroutines
- No race conditions
- Thread-safe operations
- Performance benchmarking

## Production Tips

1. **Use Fluent API** - Cleaner than traditional New() + multiple calls
2. **Load dictionaries once** - Use `init()` or singleton pattern
3. **Build before use** - Always call `Build()` or `MustBuild()` after adding words
4. **Thread-safe** - Detector is safe for concurrent use (read-only after Build)
5. **Remote dictionaries** - Use `LoadDictFromURL()` to load from HTTP/HTTPS
6. **Error handling** - Use `Build()` for errors, `MustBuild()` to panic
7. **Dependency injection** - Use interface `ContentModerator` for loose coupling
8. **High performance** - All slices pre-allocated, pool reuses memory
