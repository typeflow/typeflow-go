package golevenshtein

import (
	"testing"
)

type perm struct {
	source      string
	destination string
	distance    int
}

var testCases = []perm{
	{"alessandro", "lessandro", 1},
	{"alessandro", "alesasndro", 2},
	{"zzz", "az", 2},
	{"--|", "---", 1},
}

var incrementalTestCases = [][]perm{
	{ {"aaa", "abc", 2 } },
	{{ "al", "al",  0 }, { "ex", "es", 1}},
	{{"a", "b", 1}, {"aa", "ba", 2}},
	{{"a", "a", 0}, { "l", "l", 0}, {"e", "e", 0}, {"s", "s", 0}, {"s", "a", 1}, {"a", "s", 2}, {"n", "n", 2} },
}

func Test_compareSlicesRecursive(t *testing.T) {
	for _, v := range testCases {
		if cmp := compareSlicesRecursive([]byte(v.source), []byte(v.destination)); cmp != v.distance {
			t.Errorf("Comparing '%s' and '%s' failed: result was %v but expected was %v", v.source, v.destination, cmp, v.distance)
		}
	}
}

func Test_compareSlicesWithMatrixBaseCase(t *testing.T) {
	for _, v := range incrementalTestCases[0] {
		ls := InitLState()
		ls.UpdateState([]rune(v.source), []rune(v.destination))

		if cmp := ls.Distance(); cmp != v.distance {
			t.Errorf("Comparing '%s' and '%s' failed: result was %v but expected was %v", v.source, v.destination, cmp, v.distance)
		}
	}
}

func Test_compareSlicesWithMatrixIncremental(t* testing.T) {
	for _, v := range incrementalTestCases {
		ls := InitLState()
		for _, increment := range v {
			ls.UpdateState([]rune(increment.source), []rune(increment.destination))

			if cmp := ls.Distance(); cmp != increment.distance {
				t.Errorf("Comparing '%s' and '%s' failed: result was %v but expected was %v", increment.source, increment.destination, cmp, increment.distance)
			}
		}
	}
}

func Benchmark_recursiveLevenshtein(b *testing.B) {
	for _, v := range testCases {
		compareSlicesRecursive([]byte(v.source), []byte(v.destination))
	}
}

func Benchmark_matrixLevenshtein(b *testing.B) {
	for _, v := range testCases {
		ls := InitLState()
		ls.UpdateState([]rune(v.source), []rune(v.destination))
	}
}
