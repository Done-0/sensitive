// Package trie implements Trie tree and AC automaton for sensitive word detection
// Creator: Done-0
// Created: 2025-01-15
package trie

type Node struct {
	children map[rune]*Node
	fail     *Node
	isEnd    bool
	word     string
	level    int
}

func (n *Node) Children() map[rune]*Node {
	return n.children
}

type Tree struct {
	root *Node
}

func New() *Tree {
	return &Tree{
		root: &Node{children: make(map[rune]*Node)},
	}
}

func (t *Tree) Insert(word string, level int) {
	current := t.root
	for _, r := range word {
		if _, exists := current.children[r]; !exists {
			current.children[r] = &Node{children: make(map[rune]*Node)}
		}
		current = current.children[r]
	}
	current.isEnd = true
	current.word = word
	current.level = level
}

func (t *Tree) Root() *Node {
	return t.root
}

func BuildFailureLinks(root *Node) {
	if root == nil {
		return
	}

	queue := make([]*Node, 0, 256)
	root.fail = root

	for _, child := range root.children {
		child.fail = root
		queue = append(queue, child)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for char, child := range current.children {
			queue = append(queue, child)

			fail := current.fail
			for fail != root && fail.children[char] == nil {
				fail = fail.fail
			}

			if fail != root && fail.children[char] != nil {
				child.fail = fail.children[char]
			} else if root.children[char] != nil && child != root.children[char] {
				child.fail = root.children[char]
			} else {
				child.fail = root
			}
		}
	}
}

type Match struct {
	Word  string
	Start int
	End   int
	Level int
}

func Search(root *Node, text []rune) []Match {
	matches := make([]Match, 0, 16)
	current := root

	for i, r := range text {
		for current != root && current.children[r] == nil {
			current = current.fail
		}

		if current.children[r] != nil {
			current = current.children[r]
		} else {
			current = root
		}

		temp := current
		for temp != root {
			if temp.isEnd {
				matches = append(matches, Match{
					Word:  temp.word,
					Start: i - len([]rune(temp.word)) + 1,
					End:   i + 1,
					Level: temp.level,
				})
			}
			temp = temp.fail
		}
	}

	return matches
}
