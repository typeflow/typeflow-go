package typeflow

import (
	. "github.com/typeflow/triego"
	"github.com/alediaferia/stackgo"
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
	den := maximum(lenw1, lenw2)

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
// Param substr is the string to match against.
func (ws* WordSource) FindMatch(substr string, minSimilarity float32) (matches []Match, err error) {
	matches = make([]Match, 0, ws.wc)

	matrix := InitLState(substr)

	err = nil
	last_depth  := -1

	// we'll take advantage
	// of a stack for keeping
	// track of the added prefix
	// lengths
	lengths := stackgo.NewStack()

    ws.trie.EachPrefix(func (info PrefixInfo) (skipsubtree, halt bool) {

		if info.Depth <= last_depth {
			// we are going up again
			// in the radix tree so
			// we roll back the previous
			// update
			matrix.RollbackBy(lengths.Pop().(int) - info.SharedLength)
		}
		last_depth  = info.Depth
		matrix.UpdateState([]rune(info.Prefix)[info.SharedLength:])
		lengths.Push(len(info.Prefix))

		// if this prefix
		// is not yet a word we took
		// the advantage of preparing
		// the comparison matrix
		// in advance as for sure the
		// words that are about to come
		// will have this as prefix string
		if !info.IsWord {
			return false, false
		}

		similarity := computeSimilarity(len(matrix.source), len(substr), matrix.Distance())

		if similarity >= minSimilarity {
			matches = append(matches, Match{string(matrix.source), similarity})
			return false, false
		}

		// the computed similarity is too low
		// so there's no need to proceed further
		// with this subtree
		return true, false
	})

	return
}

