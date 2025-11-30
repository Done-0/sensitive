// Package sensitive provides high-performance sensitive word detection using AC automaton
// Creator: Done-0
// Created: 2025-01-15
package sensitive

type Level int

const (
	LevelLow    Level = 1
	LevelMedium Level = 2
	LevelHigh   Level = 3
)

func (l Level) String() string {
	switch l {
	case LevelLow:
		return "Low"
	case LevelMedium:
		return "Medium"
	case LevelHigh:
		return "High"
	default:
		return "Unknown"
	}
}

func (l Level) IsValid() bool {
	return l >= LevelLow && l <= LevelHigh
}

type Match struct {
	Word  string
	Start int
	End   int
	Level Level
}

type Result struct {
	HasSensitive bool
	Matches      []Match
	FilteredText string
}

type Stats struct {
	TotalWords int
	TreeDepth  int
	MemorySize int64
}

type FilterStrategy int

const (
	StrategyMask FilterStrategy = iota
	StrategyRemove
	StrategyReplace
)

type Options struct {
	FilterStrategy FilterStrategy
	ReplaceChar    rune
	SkipWhitespace bool
	EnableVariant  bool
	CaseSensitive  bool
}

type Option func(*Options)

func WithFilterStrategy(s FilterStrategy) Option {
	return func(o *Options) { o.FilterStrategy = s }
}

func WithReplaceChar(c rune) Option {
	return func(o *Options) { o.ReplaceChar = c }
}

func WithSkipWhitespace(skip bool) Option {
	return func(o *Options) { o.SkipWhitespace = skip }
}

func WithVariant(enable bool) Option {
	return func(o *Options) { o.EnableVariant = enable }
}

func WithCaseSensitive(sensitive bool) Option {
	return func(o *Options) { o.CaseSensitive = sensitive }
}
