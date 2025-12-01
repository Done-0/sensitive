// Package pool provides memory pool optimization
// Creator: Done-0
// Created: 2025-01-15
package pool

import "sync"

var (
	runePool = sync.Pool{New: func() any { s := make([]rune, 0, 1024); return &s }}
	boolPool = sync.Pool{New: func() any { s := make([]bool, 0, 1024); return &s }}
)

func GetRunes(n int) *[]rune {
	s := runePool.Get().(*[]rune)
	if cap(*s) < n {
		*s = make([]rune, 0, n)
	} else {
		*s = (*s)[:0]
	}
	return s
}

func PutRunes(s *[]rune) {
	if s != nil && cap(*s) <= 65536 {
		runePool.Put(s)
	}
}

func GetBools(n int) *[]bool {
	p := boolPool.Get().(*[]bool)
	if cap(*p) < n {
		*p = make([]bool, n)
	} else {
		*p = (*p)[:n]
		for i := range *p {
			(*p)[i] = false
		}
	}
	return p
}

func PutBools(s *[]bool) {
	if s != nil && cap(*s) <= 65536 {
		boolPool.Put(s)
	}
}
