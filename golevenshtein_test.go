package golevenshtein

import (
	"testing"
	"fmt"
)

var testCases = []struct {
	source      string
	destination string
	distance    Score
}{
	{"alessandro", "lessandro", 1},
	{"alessandro", "alesasndro", 2},
	{"zzz", "az", 2},
	{"--|", "---", 1},
}

func Test_compareSlicesRecursive(t *testing.T) {
	for _, v := range testCases {
		if cmp := compareSlicesRecursive([]byte(v.source), []byte(v.destination)); cmp != v.distance {
			t.Errorf("Comparing '%s' and '%s' failed: result was %v but expected was %v", v.source, v.destination, cmp, v.distance)
		}
	}
}

func Test_trieFindsWords(t *testing.T) {
	rootTrie := NewTrie()

	// manually appending two short
	// words to make sure test is
	// exclusively run against
	// the find function
	//
	// words added: cat, dog

	// dog
	tr := NewTrie()
	tr.isRoot = false
	tr.Parent = rootTrie
	tr.C = 'd'

	rootTrie.Children['d'] = tr

	tr1 := NewTrie()
	tr1.isRoot = false
	tr1.Parent = tr
	tr1.C = 'o'
	tr.Children['o'] = tr1

	tr2 := NewTrie()
	tr2.isRoot = false
	tr2.Parent = tr1
	tr2.C = 'g'
	tr2.IsWord = true
	tr1.Children['g'] = tr2

	// cat
	tr3 := NewTrie()
	tr3.isRoot = false
	tr3.Parent = rootTrie
	tr3.C = 'c'

	rootTrie.Children['c'] = tr3

	tr4 := NewTrie()
	tr4.isRoot = false
	tr4.Parent = tr3
	tr4.C = 'a'
	tr3.Children['a'] = tr4

	tr5 := NewTrie()
	tr5.isRoot = false
	tr5.Parent = tr4
	tr5.C = 't'
	tr5.IsWord = true
	tr4.Children['t'] = tr5

	if rootTrie.HasWord([]rune("dog")) == false {
		t.Errorf("Finding word 'dog' in trie fails")
	}

	if rootTrie.HasWord([]rune("cat")) == false {
		t.Errorf("Finding word 'cat' in trie fails")
	}

	if rootTrie.HasWord([]rune("foo")) == true {
		t.Errorf("Finding word 'foo' in trie unexpectedly succeeds")
	}

	var i int = 0
	countTrieNodes(rootTrie, &i)
	if i != 7 {
		t.Fatalf("Expected 7 nodes, got %d", i)
	}
}

/*
 * A utility function to make sure
 * node append workd properly for our trie
 */
func countTrieNodes(trie *TrieNode, i *int) {
	if len(trie.Children) == 0 {
		*i = *i + 1
		return
	}
	for _, v := range trie.Children {
		countTrieNodes(v, i)
	}

	*i = *i + 1
}

func Test_trieAppendsWords(t *testing.T) {
	rootTrie := NewTrie()

	const (
		w1 = "testWord1"
	    w2 = "testWord2"
	)

	rootTrie.Append([]rune(w1), true)
	rootTrie.Append([]rune(w2), true)

	if rootTrie.HasWord([]rune(w1)) == false {
		t.Errorf("Finding word '%s' in trie fails", w1)
	}

	if rootTrie.HasWord([]rune(w2)) == false {
		t.Errorf("Finding word '%s' in trie fails", w2)
	}
	var i int = 0
	countTrieNodes(rootTrie, &i)
	if i != 11 {
		t.Fatalf("Expected 11 nodes, got %d", i)
	}

	words := rootTrie.Words()
	for _, word := range words {
		if word != w1 && word != w2 {
			t.Fatal("Cannot find expected words in the list of words in the trie: %v", words)
		}
	}
}

/*
 * A few helper functions
 */
/*func printTrie(trie *TrieNode) {
	for _, v := range trie.Children {
		runes := make([]rune, 0)
		printTrie_(v, append(runes, v.C))
	}
}*/

func printTrie_(trie *TrieNode, runes []rune) {
	for _,v := range trie.Children {
		if v.IsWord {
			fmt.Println(string(append(runes, v.C)))
		}
		printTrie_(v, append(runes, v.C))
	}
}

