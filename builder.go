// Package sensitive provides high-performance sensitive word detection using AC automaton
// Creator: Done-0
// Created: 2025-01-15
package sensitive

import "errors"

type Builder struct {
	detector *Detector
	errors   []error
}

func NewBuilder(opts ...Option) *Builder {
	return &Builder{
		detector: New(opts...),
		errors:   make([]error, 0, 4),
	}
}

func (b *Builder) AddWord(word string, level Level) *Builder {
	if err := b.detector.AddWord(word, level); err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *Builder) AddWords(words map[string]Level) *Builder {
	if err := b.detector.AddWords(words); err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *Builder) LoadDict(path string) *Builder {
	if err := b.detector.LoadDict(path); err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *Builder) LoadDictWithLevel(path string, level Level) *Builder {
	if err := b.detector.LoadDictWithLevel(path, level); err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *Builder) LoadDictFromURL(url string) *Builder {
	if err := b.detector.LoadDictFromURL(url); err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *Builder) LoadDictFromURLWithLevel(url string, level Level) *Builder {
	if err := b.detector.LoadDictFromURLWithLevel(url, level); err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *Builder) LoadDictFromURLs(urls []string) *Builder {
	for _, url := range urls {
		if err := b.detector.LoadDictFromURL(url); err != nil {
			b.errors = append(b.errors, err)
		}
	}
	return b
}

func (b *Builder) LoadVariantMap(path string) *Builder {
	if err := b.detector.LoadVariantMap(path); err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *Builder) LoadEmbeddedDict(name string, level Level) *Builder {
	if err := LoadEmbeddedDict(b.detector, name, level); err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *Builder) LoadAllEmbedded() *Builder {
	if err := LoadAllEmbedded(b.detector); err != nil {
		b.errors = append(b.errors, err)
	}
	return b
}

func (b *Builder) WithFilterStrategy(strategy FilterStrategy) *Builder {
	b.detector.opts.FilterStrategy = strategy
	return b
}

func (b *Builder) WithReplaceChar(char rune) *Builder {
	b.detector.opts.ReplaceChar = char
	return b
}

func (b *Builder) WithSkipWhitespace(skip bool) *Builder {
	b.detector.opts.SkipWhitespace = skip
	return b
}

func (b *Builder) WithVariant(enable bool) *Builder {
	b.detector.opts.EnableVariant = enable
	return b
}

func (b *Builder) WithCaseSensitive(sensitive bool) *Builder {
	b.detector.opts.CaseSensitive = sensitive
	return b
}

func (b *Builder) Build() (*Detector, error) {
	if len(b.errors) > 0 {
		return nil, errors.Join(b.errors...)
	}

	b.detector.Build()
	return b.detector, nil
}

func (b *Builder) MustBuild() *Detector {
	detector, err := b.Build()
	if err != nil {
		panic(err)
	}
	return detector
}
