// Package pool provides memory pool optimization
// Creator: Done-0
// Created: 2025-01-15
package pool

import "sync"

var RunePool = sync.Pool{
	New: func() any {
		s := make([]rune, 0, 1024)
		return &s
	},
}

func Get(capacity int) *[]rune {
	s := RunePool.Get().(*[]rune)
	if cap(*s) < capacity {
		*s = make([]rune, 0, capacity)
	}
	return s
}

func Put(s *[]rune) {
	if cap(*s) > 65536 {
		return
	}
	*s = (*s)[:0]
	RunePool.Put(s)
}
