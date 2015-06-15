package golevenshtein

import (
	"github.com/alediaferia/stackgo"
	"unsafe"
)

type TrieNode struct {
	IsWord   bool
	Parent   *TrieNode
	C        rune
	Children map[rune]*TrieNode
	isRoot   bool
}

type TrieNodePtr *TrieNode

func NewTrie() *TrieNode {
	t := new(TrieNode)
	t.IsWord = false
	t.Parent = nil
	t.C = 0
	t.isRoot = true
	t.Children = make(map[rune](*TrieNode))

	return t
}

func (t *TrieNode) IsRoot() bool {
	return t.isRoot
}

/*
 * Appends a word.
 * This is a recursive function, so not that
 * efficient.
 */
func (t *TrieNode) Append(suffix []rune, makesWord bool) {
	if len(suffix) == 0 {
		return
	}

	// if there is already a node
	// holding this character we
	// move forward and append
	// the remaining part
	if c,ok := t.Children[suffix[0]]; ok {
		c.Append(suffix[1:], makesWord)
		return
	}

	tc := NewTrie()
	tc.Parent = t
	t.Children[suffix[0]] = tc
	tc.C = suffix[0]
	tc.isRoot = false

	if len(suffix) > 1 {
		tc.Append(suffix[1:], makesWord)
	} else {
		tc.IsWord = makesWord
	}
}

func (t *TrieNode) HasWord(word []rune) bool {
	currentSlice := word
	currentRoot  := t

	for len(currentSlice) > 0 {
		c, ok := currentRoot.Children[currentSlice[0]]
		if len(currentSlice) == 1 && ok == true && c.IsWord {
			return true
		} else if !ok {
			return false
		}
		currentSlice = currentSlice[1:]
		currentRoot  = c
	}

	return false
}

// Returns a list with all the
// words present in the trie
func (t *TrieNode) Words() (words []string) {
    // DFS-based implementation for returning
	// al the words in the trie

	stack := stackgo.NewStack()
	node := t

	words = make([]string, 0)
	word  := make([]rune, 0)

	stack.Push(unsafe.Pointer(node))
	for stack.Size() > 0 {
		node = TrieNodePtr(stack.Pop().(unsafe.Pointer))
		word = append(word, node.C)

		if len(node.Children) == 0 {
			if node.IsWord {
				words = append(words, string(word))
			}
			word = word[:len(word) - 1]
		}
		for c := range node.Children {
			stack.Push(c)
		}
	}

	return
}
