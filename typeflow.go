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
	vec1     []int
    vec2     []int

	// w1 = cols, w2 = rows
	w1       []rune
	w2       []rune
    
    // init flag
    setup    bool
}

// Initializes an empty LState.
// This is the first step required
// for working with a matrix-based
// levenshtein computation.
func InitLState() (ls *LState) {
	ls = new(LState)
	ls.setup = false

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
func (state *LState) UpdateState(w1part, w2part []rune) {
    if state.setup == false {
		state.w1 = make([]rune, len(w1part))
		copy(state.w1, w1part)

		state.w2 = make([]rune, len(w2part))
		copy(state.w2, w2part)

        state.setup = true
        state.initializeMatrix(len(state.w2))
        
	} else {
		state.w1 = append(state.w1, w1part...)
		state.w2 = append(state.w2, w2part...)
        
        new_vec1 := make([]int, len(state.vec1) + len(w2part))
        new_vec2 := make([]int, len(state.vec2) + len(w2part))
        
        copy(new_vec1, state.vec1)
        copy(new_vec2, state.vec2)
        
        state.vec1 = new_vec1
        state.vec2 = new_vec2
        
    	for i := 0; i < len(state.vec1); i++ {
    		state.vec1[i] = i
    	}
    }
    state.fillMatrix(0, 0)    
}

// Rolls back the current state. cols is the number of
// characters to roll back referencing what was w1part
// in UpdateState. rows will be the number of characters
// to roll back referencing what was w2part in UpdateState.
func (state *LState) RollbackBy(cols, rows int) (error) {
    if state.setup == false {
		return EmptyStateError
	}

    if len(state.vec1) < rows {
		return OutOfRangeRollbackError
	}

	state.w2 = state.w2[:len(state.w2) - rows]
	state.w1 = state.w1[:len(state.w1) - cols]
    
    state.vec1 = state.vec1[:len(state.vec1) - rows]
    state.vec2 = state.vec2[:len(state.vec2) - rows]
    
	for i := 0; i < len(state.vec1); i++ {
		state.vec1[i] = i
	}
    state.fillMatrix(0, 0)

	return nil
}

// Returns the newly computed distance
// or math.MaxInt32 if no distance has been
// computed yet.
// This method has complexity O(1) as
// the distance is computed upon every
// UpdateState call.
func (state *LState) Distance() int {
	if state.setup != false {
		return state.vec2[len(state.w2)]
	}
	return math.MaxInt32
}

func (state *LState) initializeMatrix(s2 int) {
    state.vec1 = make([]int, s2 + 1)
    state.vec2 = make([]int, s2 + 1)
    
    // initializing vec1
	for i := 0; i < len(state.vec1); i++ {
		state.vec1[i] = i
	}
}

func (state *LState) fillMatrix(startRow, startCol int) {
	for i := startRow; i < len(state.w1); i++ {
	    state.vec2[0] = i + 1;
        
        for j := startCol; j < len(state.w2); j++ {
            cost := 1
            if (state.w1[i] == state.w2[j]) {
                cost = 0
            }
            min, err := minimum(state.vec2[j] + 1,
                                state.vec1[j + 1] + 1,
                                state.vec1[j] + cost)
            if err != nil {
                panic(err)
            }
            state.vec2[j + 1] = min
        }
        
        for j := startCol; j < len(state.vec1); j++ {
            state.vec1[j] = state.vec2[j]
        }
	}
}
