package golevenshtein

import (
	"errors"
	"fmt"
	"math"
)

func minimum(values ...int) (int, error) {
	if len(values) == 0 {
		return -1, errors.New("Cannot find minimum of empty list")
	}
	var min int = values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min, nil
}

/*
 * A slow recursive levenshtein implementation
 * only used for benchmarking purposes
 */
func compareSlicesRecursive(first, second []byte) int {
	if len(first) == 0 || len(second) == 0 {
		return  int(math.Max(float64(len(first)), float64(len(second))))
	}
	len1 := len(first)
	len2 := len(second)
	var cost int
	if first[len1-1] == second[len2-1] {
		cost = 0
	} else {
		cost = 1
	}

	min, err := minimum(int(compareSlicesRecursive(first[0:len1-1], second))+1,
		int(compareSlicesRecursive(first, second[0:len2-1]))+1,
		int(compareSlicesRecursive(first[0:len1-1], second[0:len2-1]))+cost)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return min
}

type LState struct {
	matrix   [][]int

	// w1 = cols, w2 = rows
	w1       []rune
	w2       []rune
}

func InitLState() (ls *LState) {
	ls = new(LState)
	ls.matrix = nil

	return
}

func (state *LState) UpdateState(w1part, w2part []rune) {
	if len(w1part) == 0 || len(w2part) == 0 {
		panic("Unexpected 0 length argument")
	}
    if state.matrix == nil {
        state.initializeMatrix(len(w1part), len(w2part))
		state.w1 = w1part
		state.w2 = w2part

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
		for i := incr1; i < len(state.w1); i++ {
			if cap(state.matrix[i]) > i {
				state.matrix[i][0] = i
			} else {
				state.matrix[i] = append(state.matrix[i], i)
			}
		}

		// initializing
		for i := incr2; i < len(state.w2); i++ {
			if cap(state.matrix[0]) > i {
				state.matrix[0][i] = i
			} else {
				state.matrix[0] = append(state.matrix[0], i)
			}
		}

		state.fillMatrix(incr2 - 1, incr1 - 1)
	}
}

func (state *LState) initializeMatrix(s1, s2 int) {
	state.matrix = make([][]int, s1 + 1)
	for i := 0; i < len(state.matrix); i++ {
		state.matrix[i] = make([]int, s2 + 1)
	}

	// initializing
	for i := 0; i < s1; i++ {
		state.matrix[i][0] = i
	}

	// initializing
	for i := 0; i < s2; i++ {
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

func (state *LState) Distance() int {
	if state.matrix != nil {
		return state.matrix[len(state.w1)][len(state.w2)]
	}
	return math.MaxInt32
}
