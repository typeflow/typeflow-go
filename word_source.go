package typeflow

import (
	. "github.com/typeflow/triego"
)

type WordFilter func (w string) (word string, skip bool)

type WordSource struct {
    trie *Trie
	wc   int
}

type Match struct {
	Word     string
	Similarity float32
}

func computeSimilarity(lenw1, lenw2, ld int) (float32) {
	den, err := maximum(lenw1, lenw2)

	if err != nil {
		panic(err)
	}

	return 1.0 - float32(ld)/float32(den)
}

func NewWordSource() (ws *WordSource) {
	ws = new(WordSource)
	ws.trie = NewTrie()
    ws.wc = 0
	return
}

func (ws *WordSource) SetSource(words []string, filters []WordFilter) {

	for _, w := range words {
		for _, filter := range filters {
			w, skip := filter(w)
			if !skip {
				ws.trie.AppendWord(w)
				ws.wc += 1
			}
		}
	}
}

type dirty_range struct {
	low    int
	length int
}

func (ws* WordSource) FindMatch(substr string, minSimilarity float32) (matches []Match, err error) {
	matches = make([]Match, 0, ws.wc)
	word  := make([]rune, 0)

	matrix := InitLState()
	matrix.UpdateState([]rune{}, []rune(substr))

	err = nil
	last_depth := 0
	dirty_range := dirty_range{0, 0}

    ws.trie.EachNode(func (node *TrieNode, halt *bool) {
		rollback_size := 0
		if !node.IsRoot() {
			if node.Depth() <= last_depth {
				rollback_size = (last_depth - node.Depth() + 1)
				word = word[:len(word) - rollback_size]
				err = matrix.RollbackBy(rollback_size, 0)
				if err != nil {
				    *halt = true
					return
				}
				dirty_range.low -= rollback_size
				dirty_range.length = 0
			}

			word = append(word, node.C)
			dirty_range.length += 1
			matrix.UpdateState([]rune(word[dirty_range.low:(dirty_range.low + dirty_range.length)]), []rune{})
			dirty_range.low    += dirty_range.length
			dirty_range.length = 0
		}

		if node.IsWord {
			if similarity := computeSimilarity(len(word), len(substr), matrix.Distance()); similarity >= minSimilarity {
				matches = append(matches, Match{string(word), similarity})
			}
		}

		last_depth = node.Depth()
	})

	return
}

