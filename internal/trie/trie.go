// Package trie implements Double Array Trie and AC automaton for high-performance sensitive word detection
// Creator: Done-0
// Created: 2025-01-15
package trie

import "sort"

const (
	initialSize = 524288
)

type Match struct {
	Word  string
	Start int
	End   int
	Level int
}

type output struct {
	word  *string
	level int
	len   int
}

type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
	word     *string
	level    int
}

type Tree struct {
	base         []int
	check        []int
	fail         []int
	output       []*[]output
	children     [][]int
	used         []bool
	size         int
	nextCheckPos int
	root         *trieNode
}

func New() *Tree {
	return &Tree{
		root:         &trieNode{children: make(map[rune]*trieNode, 8)},
		nextCheckPos: 1,
	}
}

func (t *Tree) Insert(word string, level int) {
	current := t.root
	for _, r := range word {
		if _, exists := current.children[r]; !exists {
			current.children[r] = &trieNode{children: make(map[rune]*trieNode, 4)}
		}
		current = current.children[r]
	}
	current.isEnd = true
	current.word = &word
	current.level = level
}

func (t *Tree) Build() {
	if t.root == nil {
		return
	}

	t.base = make([]int, initialSize)
	t.check = make([]int, initialSize)
	t.fail = make([]int, initialSize)
	t.output = make([]*[]output, initialSize)
	t.children = make([][]int, initialSize)
	t.used = make([]bool, initialSize)
	t.size = 1

	t.used[0] = true

	chars := make([]int, 0, len(t.root.children))
	for r := range t.root.children {
		chars = append(chars, int(r))
	}
	sort.Ints(chars)

	for _, c := range chars {
		child := t.root.children[rune(c)]
		next := c
		if next >= len(t.base) {
			newSize := len(t.base) * 2
			if next+1 > newSize {
				newSize = next + 1
			}
			newBase := make([]int, newSize)
			copy(newBase, t.base)
			t.base = newBase
			newCheck := make([]int, newSize)
			copy(newCheck, t.check)
			t.check = newCheck
			newFail := make([]int, newSize)
			copy(newFail, t.fail)
			t.fail = newFail
			newOutput := make([]*[]output, newSize)
			copy(newOutput, t.output)
			t.output = newOutput
			newChildren := make([][]int, newSize)
			copy(newChildren, t.children)
			t.children = newChildren
			newUsed := make([]bool, newSize)
			copy(newUsed, t.used)
			t.used = newUsed
		}
		t.check[next] = 0
		t.used[next] = true
		t.children[0] = append(t.children[0], c)

		if child.isEnd {
			wordLen := len([]rune(*child.word))
			out := make([]output, 0, 1)
			t.output[next] = &out
			*t.output[next] = append(*t.output[next], output{
				word:  child.word,
				level: child.level,
				len:   wordLen,
			})
		}

		if next >= t.size {
			t.size = next + 1
		}

		t.buildDATRecursive(child, next)
	}

	queue := make([]int, 0, 8192)
	head := 0

	for _, c := range t.children[0] {
		t.fail[c] = 0
		queue = append(queue, c)
	}

	for head < len(queue) {
		state := queue[head]
		head++

		for _, c := range t.children[state] {
			next := t.base[state] + c
			if next >= len(t.check) || t.check[next] != state || !t.used[next] {
				continue
			}
			queue = append(queue, next)

			failState := t.fail[state]
			for {
				if failState == 0 {
					rootNext := c
					if rootNext < len(t.check) && t.check[rootNext] == 0 && t.used[rootNext] && rootNext != next {
						t.fail[next] = rootNext
					} else {
						t.fail[next] = 0
					}
					break
				}

				failNext := t.base[failState] + c
				if failNext < len(t.check) && t.check[failNext] == failState && t.used[failNext] {
					t.fail[next] = failNext
					break
				}

				failState = t.fail[failState]
			}
		}
	}

	t.root = nil
	t.children = nil
}

func (t *Tree) buildDATRecursive(node *trieNode, state int) {
	if node == nil || len(node.children) == 0 {
		return
	}

	chars := make([]int, 0, len(node.children))
	for r := range node.children {
		chars = append(chars, int(r))
	}
	sort.Ints(chars)

	pos := t.nextCheckPos
	if pos < chars[0]+1 {
		pos = chars[0] + 1
	}
	base := pos - chars[0]

	for {
		collision := false
		for _, c := range chars {
			next := base + c
			if next >= len(t.base) {
				newSize := len(t.base) * 2
				if next+1 > newSize {
					newSize = next + 1
				}
				newBase := make([]int, newSize)
				copy(newBase, t.base)
				t.base = newBase
				newCheck := make([]int, newSize)
				copy(newCheck, t.check)
				t.check = newCheck
				newFail := make([]int, newSize)
				copy(newFail, t.fail)
				t.fail = newFail
				newOutput := make([]*[]output, newSize)
				copy(newOutput, t.output)
				t.output = newOutput
				newChildren := make([][]int, newSize)
				copy(newChildren, t.children)
				t.children = newChildren
				newUsed := make([]bool, newSize)
				copy(newUsed, t.used)
				t.used = newUsed
			}
			if t.used[next] {
				collision = true
				break
			}
		}
		if !collision {
			break
		}
		base++
	}

	t.base[state] = base
	if base > t.nextCheckPos {
		t.nextCheckPos = base
	}

	for _, c := range chars {
		next := base + c
		t.check[next] = state
		t.used[next] = true
		t.children[state] = append(t.children[state], c)

		child := node.children[rune(c)]
		if child.isEnd {
			wordLen := len([]rune(*child.word))
			if t.output[next] == nil {
				out := make([]output, 0, 1)
				t.output[next] = &out
			}
			*t.output[next] = append(*t.output[next], output{
				word:  child.word,
				level: child.level,
				len:   wordLen,
			})
		}

		if next >= t.size {
			t.size = next + 1
		}

		t.buildDATRecursive(child, next)
	}
}

func (t *Tree) SearchDAT(text []rune) []Match {
	matches := make([]Match, 0, 16)
	state := 0

	for i, r := range text {
		c := int(r)
		for {
			if state >= len(t.base) {
				state = 0
				break
			}

			next := t.base[state] + c
			if next < len(t.check) && t.check[next] == state && t.used[next] {
				state = next
				break
			}

			if state == 0 {
				break
			}
			state = t.fail[state]
		}

		temp := state
		for temp > 0 {
			if temp < len(t.output) && t.output[temp] != nil {
				for _, out := range *t.output[temp] {
					matches = append(matches, Match{
						Word:  *out.word,
						Start: i - out.len + 1,
						End:   i + 1,
						Level: out.level,
					})
				}
			}
			temp = t.fail[temp]
		}
	}

	return matches
}

func (t *Tree) Size() int {
	return t.size
}

func (t *Tree) MemoryUsage() int64 {
	return int64(len(t.base)*8 + len(t.check)*8 + len(t.fail)*8 + len(t.used))
}
