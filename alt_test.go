package typeflow

import (
	"math"
)

// This file contains just some
// unexported alternative levenshtein
// implementations


// A slow recursive Levenshtein implementation
// only used for benchmarking purposes
func compare_silces_r(first, second []rune) int {
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

	min := minimum(int(compare_silces_r(first[0:len1-1], second))+1,
		int(compare_silces_r(first, second[0:len2-1]))+1,
		int(compare_silces_r(first[0:len1-1], second[0:len2-1]))+cost)
	return min
}

