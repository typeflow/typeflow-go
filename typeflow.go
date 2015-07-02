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

	// w1 = cols, w2 = rows
	w1       []rune
	w2       []rune
}

// Initializes an empty LState.
// This is the first step required
// for working with a matrix-based
// levenshtein computation.
func InitLState() (ls *LState) {
	ls = new(LState)
	ls.matrix = nil

	return
}

// Updates the current state. w1part and w2part
// are respectively the parts to be appended
// to the current state.
//
// E.g. you may want to compare the words
// 'levenshtein' and 'einstein' in two steps:
//
// 1) the first UpdateState will be called
// passing in 'leven' as w1part and 'ein' as w2part.
//
// 2) the second UpdateState will be called passing
// in 'stein' as w1part and 'stein' as w2part.
//
// You will then be able to call Distance() on the
// updated state and get 4 as result.
//
// Complexity: O(MAX(fullw1,fullw2))
// fullw1, fullw2 are actually the full lengths of the
// all characters added by all the UpdateState calls
// until now. Plans are there to reduce complexity
// to just O(MAX(w1part, w2part)) but I've not yet
// investigated about it.
func (state *LState) UpdateState(w1part, w2part []rune) {
    if state.matrix == nil {
        state.initializeMatrix(len(w1part), len(w2part))

		state.w1 = make([]rune, len(w1part))
		copy(state.w1, w1part)

		state.w2 = make([]rune, len(w2part))
		copy(state.w2, w2part)

		state.fillMatrix(0, 0)
	} else {
		incr1 := len(state.w1)
		incr2 := len(state.w2)

		// we need to increase size
		state.w1 = append(state.w1, w1part...)
		state.w2 = append(state.w2, w2part...)

		cols := make([][]int, len(state.w1) + 1)
		copy(cols, state.matrix)
		for i := 0; i < len(cols); i++ {
			col := make([]int, len(state.w2) + 1)
			if i < len(state.matrix) {
				copy(col, state.matrix[i])
			}
			cols[i] = col
		}

		state.matrix = cols

		// initializing the extended part now
		for i := incr1; i < len(state.w1) + 1; i++ {
			state.matrix[i][0] = i
		}

		// initializing the extended part now
		for i := incr2; i < len(state.w2) + 1; i++ {
			state.matrix[0][i] = i
		}

		// refilling matrix
		// TODO this can probably be improved
		// TODO avoiding refilling not needed cells
		state.fillMatrix(0, 0)
	}
}

// Rolls back the current state. cols is the number of
// characters to roll back referencing what was w1part
// in UpdateState. rows will be the number of characters
// to roll back referencing what was w2part in UpdateState.
func (state *LState) RollbackBy(cols, rows int) (error) {
    if len(state.matrix) == 0 {
		return EmptyStateError
	}

    if len(state.matrix) < cols {
		return OutOfRangeRollbackError
	}

	state.matrix = state.matrix[:len(state.matrix) - cols]

    for index, row := range state.matrix {
		if len(row) < rows {
			return OutOfRangeRollbackError
		}

		state.matrix[index] = row[:len(row) - rows]
	}

	state.w2 = state.w2[:len(state.w2) - rows]
	state.w1 = state.w1[:len(state.w1) - cols]

	return nil
}

// Returns the newly computed distance
// or math.MaxInt32 if no distance has been
// computed yet.
// This method has complexity O(1) as
// the distance is computed upon every
// UpdateState call.
func (state *LState) Distance() int {
	if state.matrix != nil {
		return state.matrix[len(state.w1)][len(state.w2)]
	}
	return math.MaxInt32
}

func (state *LState) initializeMatrix(s1, s2 int) {
	state.matrix = make([][]int, s1 + 1)
	for i := 0; i < len(state.matrix); i++ {
		state.matrix[i] = make([]int, s2 + 1)
	}

	// initializing
	for i := 1; i < s1 + 1; i++ {
		state.matrix[i][0] = i
	}

	// initializing
	for i := 1; i < s2 + 1; i++ {
		state.matrix[0][i] = i
	}
}

func (state *LState) fillMatrix(startRow, startCol int) {
	for j := startRow; j < len(state.w2); j++ {
		for i := startCol; i < len(state.w1); i++ {
			if state.w1[i] == state.w2[j] {
				state.matrix[i+1][j+1] = state.matrix[i][j]
			} else {
				v, err := minimum(state.matrix[i][j+1] + 1,
					state.matrix[i + 1][j] + 1,
					state.matrix[i][j] + 1)
				if err != nil {
					panic(err)
				}

				state.matrix[i+1][j+1] = v
			}
		}
	}
}