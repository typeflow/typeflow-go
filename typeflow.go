// This package contains everything you need
// to work with word-based searching.
// Its internals are founded on the Levenshtein
// distance.
// Check https://github.com/typeflow/typeflow-go
// for an example.
package typeflow

import (
	"math"
	"errors"
)

// Errors
var (
	// This error occurs when trying to rollback more than needed
	OutOfRangeRollbackError = errors.New("Unexpected rollback: out of range")
	EmptyStateError         = errors.New("Unexpected: current state is empty")
)

// An LState encapsulates the current
// state of a levenshtein-based distance
// computation.
// In particular, the matrix-based
// levenshtein computation is used.
// An LState instance can be updated
// or rolled back which is useful
// for iterative word computations. It will
// take care of managing the internal state
// and memory allocation for you.
type LState struct {

	matrix   [][]int

	target       []rune // the target word to compare against
	source       []rune
}

// Initializes an empty LState.
// This is the first step required
// for working with a matrix-based
// levenshtein computation.
func InitLState(target string) (ls *LState) {
	ls = new(LState)
	ls.target = []rune(target)

	return
}

// Updates the current state with param source.
//
// Since LState can be rolled back, source can
// also represent a slice to be appended
// to the previous one passed through a previous
// call to UpdateState.
//
// E.g. you may want to compare the words
// 'levenshtein' and 'einstein' in two steps
// after having initialized the state using
// InitLState('einstein')
//
// 1) the first UpdateState will be called
// passing in 'leven' as source.
//
// 2) the second UpdateState will be called passing
// in 'stein' as source.
//
// You will then be able to call Distance() on the
// updated state and get 4 as result.
//
func (state *LState) UpdateState(source []rune) {
    if state.source == nil {
		state.source = source
		state.initializeMatrix(len(state.target), len(state.source))
		state.fillMatrix(0, 0)
	} else {
		delta := len(source)
		state.source = append(state.source, source...)
		state.extendMatrix(delta)
		state.fillMatrix(0, len(state.source) - delta)
	}
}

// Rolls back the current state. Specify the amount of characters
// to roll back for the source string through the chars param
func (state *LState) RollbackBy(chars int) (error) {
	if chars > len(state.matrix) {
		return OutOfRangeRollbackError
	}
    //state.matrix = state.matrix[:len(state.matrix) - chars]
	state.source = state.source[:len(state.source) - chars]

	return nil
}

// Returns the newly computed distance
// or math.MaxInt32 if no distance has been
// computed yet.
// This method has complexity O(1) as
// the distance is computed upon every
// UpdateState call.
func (state *LState) Distance() int {
	if state.matrix == nil {
		return math.MaxInt32
	}
	return state.matrix[len(state.source)][len(state.target)]
}

// initializes the matrix
// param length represents
// the size of the source string
// while width is the size of the target
// string.
// The target string is the one that doesn't change.
func (state *LState) initializeMatrix(width, length int) {
    state.matrix = make([][]int, length + 1)
    // now initializing each column
	for i := 0; i < len(state.matrix); i++ {
		state.matrix[i] = make([]int, width + 1)
	}

	for i := 1; i < length + 1; i++ {
		state.matrix[i][0] = i
	}

	for i := 1; i < width + 1; i++ {
		state.matrix[0][i] = i
	}
}

// TODO this function needs to be more
// clever about memory allocations
func (state *LState) extendMatrix(delta int) {
	l_now := len(state.matrix)
	m := make([][]int, len(state.matrix) + delta)
	copy(m, state.matrix)
	state.matrix = m
	for i := l_now; i < len(state.matrix); i++ {
		state.matrix[i] = make([]int, len(state.target) + 1)
	}

	for i := l_now; i < len(state.matrix); i++ {
		state.matrix[i][0] = i
	}
}

func (state *LState) fillMatrix(targetStart, sourceStart int) {
	for j := targetStart; j < len(state.target); j++ {
		for i := sourceStart; i < len(state.source); i++ {
			if state.source[i] == state.target[j] {
				state.matrix[i + 1][j + 1] = state.matrix[i][j]
			} else {
				state.matrix[i + 1][j + 1] = minimum(
					state.matrix[i][j + 1] + 1,
					state.matrix[i + 1][j] + 1,
					state.matrix[i][j] + 1,
				)
			}
		}
	}
}

// This implementation takes advantage of the 2-columns
// approach since doesn't expose any incremental update
// functionality
// You are encouraged to use this function when simply
// interested in the levenshtein distance between 2 words
func LevenshteinDistance(source, destination string) int {
	vec1 := make([]int, len(destination) + 1)
	vec2 := make([]int, len(destination) + 1)

	w1 := []rune(source)
	w2 := []rune(destination)

	// initializing vec1
	for i := 0; i < len(vec1); i++ {
		vec1[i] = i
	}

	// initializing the matrix
	for i := 0; i < len(w1); i++ {
		vec2[0] = i + 1;

		for j := 0; j < len(w2); j++ {
			cost := 1
			if (w1[i] == w2[j]) {
				cost = 0
			}
			min := minimum(vec2[j] + 1,
				vec1[j + 1] + 1,
				vec1[j] + cost)
			vec2[j + 1] = min
		}

		for j := 0; j < len(vec1); j++ {
			vec1[j] = vec2[j]
		}
	}

	return vec2[len(w2)]
}
