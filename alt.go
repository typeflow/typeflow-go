package golevenshtein

import (
	"fmt"
	"math"
)

// this file contains just some
// unexported alternative levenshtein
// implementations


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

