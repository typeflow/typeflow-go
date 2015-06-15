package golevenshtein

import (
	"testing"
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



