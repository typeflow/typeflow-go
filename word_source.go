package typeflow

import (
	. "github.com/typeflow/triego"
)

// Represents a filter to apply to the
// source setted via the SetSource method.
// A filter will be passed the currently processed
// word and will return either a modified version
// for that word or true for the skip param which will
// make the WordSource ignore the current word.
type WordFilter func (w string) (word string, skip bool)

// The actual WordSource type which
// will hold the source you set through
// SetSource
type WordSource struct {
    trie *Trie
	wc   int
}

// Represents a match
// found
type Match struct {
	// The actual word which matched
	// the requested string
	Word       string

	// The similarity value for the
	// current word.
	// Similarity is 1 when
	// the two words are equal.
	// See github.com/typeflow/typeflow-go for
	// more details on how the similarity is computed.
	Similarity float32
}

func computeSimilarity(lenw1, lenw2, ld int) (float32) {
	den, err := maximum(lenw1, lenw2)

	if err != nil {
		panic(err)
	}

	return 1.0 - float32(ld)/float32(den)
}

// Initializes a new empty WordSource
func NewWordSource() (ws *WordSource) {
	ws = new(WordSource)
	ws.trie = NewTrie()
    ws.wc = 0
	return
}

// Sets the given strings slice as the
// current source after applying the given
// filters in order
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

// Finds a match among the current words
// in the source.
// minSimilarity is the minimum accepted similarity
// to use when filling the matches slice.
// substr is the string to match against.
func (ws* WordSource) FindMatch(substr string, minSimilarity float32) (matches []Match, err error) {
	matches = make([]Match, 0, ws.wc)
	word  := make([]rune, 0)

	matrix := InitLState()
	matrix.UpdateState([]rune{}, []rune(substr))

	err = nil
	last_prefix := 0

    ws.trie.EachPrefix(func (prefix string, is_word bool) (skipsubtree, halt bool) {
		if !is_word {
			return false, false
		}
		if len(prefix) < len(last_prefix) {
			matrix.RollbackBy(len(last_prefix) - len(prefix), 0)
		}
		last_prefix = prefix

		if similarity := computeSimilarity(len(word), len(substr), matrix.Distance()); similarity >= minSimilarity {
			matches = append(matches, Match{string(word), similarity})
		} else {
			// the computed similarity is too low
			// so there's no need to proceed further
			// with this subtree
			return true, false
		}

		return false, false
	})

	return
}

