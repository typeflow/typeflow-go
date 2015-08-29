package typeflow

import (
	"testing"
	"fmt"
)

func (s* LState) print_matrix(t* testing.T) {
	rows := ""
	for i := 0; i < len(s.matrix[0]); i++ {
		row := ""
		for j := 0; j < len(s.matrix); j++ {
			row = fmt.Sprintf("%s[ %v ]", row, s.matrix[j][i])
		}
		rows = fmt.Sprintf("%s%s\n", rows, row)
	}
	t.Logf("\n%s\n", rows)
}

type permutation_t struct {
	source      string
	destination string
	distance    int
}

var basic_test_cases = []permutation_t{
	{"alessandro", "lessandro", 1},
	{"alessandro", "alesasndro", 2},
	{"zzz", "az", 2},
	{"--|", "---", 1},
}

func TestLevenshteinDistance(t *testing.T) {
	for _, v := range basic_test_cases {
		if cmp := LevenshteinDistance(v.source, v.destination); cmp != v.distance {
			t.Errorf("Comparing '%s' and '%s' failed: result was %v but expected was %v",
				v.source,
				v.destination,
				cmp,
				v.distance)
		}
	}
}

// represents a single increment
// within a test case with its
// associated distance value
type increment_t struct {
	source   string
	distance int
}

// represents a generic
// incremental test case
type incremental_test_case_t struct {
	target    			string
	source_increments []increment_t
}

var incremental_test_cases = []incremental_test_case_t{
	{"abc", []increment_t{{"aaa", 2}}},
	{ "alex", []increment_t{{"al",  2}, { "es", 1}}}, // alex, ales
	{ "aaale", []increment_t{{"b", 5}, {"ba", 4}, {"a", 4}}}, // aaale, bbaa
	{"iraq", []increment_t{{"rep of ireland", 11}}},
}

func Test_compare_silces_r(t *testing.T) {
	for _, v := range incremental_test_cases {
		var source []rune
		for _, increment := range v.source_increments {
			source = append(source, []rune(increment.source)...)

			if cmp := compare_silces_r([]rune(source), []rune(v.target)); cmp != increment.distance {
				t.Errorf("Comparing '%s' and '%s' failed: result was %v but expected was %v",
					source,
					v.target,
					cmp,
					increment.distance)
			}
		}
	}
}

func TestDistanceBaseCase(t *testing.T) {
	for _, v := range incremental_test_cases {
		var source   []rune
		var target   []rune = []rune(v.target)
		for _, tc := range v.source_increments {
			source = append(source, []rune(tc.source)...)
		}

		ls := InitLState(string(target))
		ls.UpdateState([]rune(source))

		if cmp := ls.Distance(); cmp != v.source_increments[len(v.source_increments)-1].distance {
			t.Logf("Comparing '%s' and '%s' failed: result was %v but expected was %v\n",
				string(source),
				string(target), cmp, v.source_increments[len(v.source_increments)-1].distance)
			t.Logf("Matrix was:\n")
			ls.print_matrix(t)
			t.FailNow()
		}
	}
}

func TestDistanceIncremental(t* testing.T) {
	for _, v := range incremental_test_cases {
		ls := InitLState(v.target)

		for _, increment := range v.source_increments {
			ls.UpdateState([]rune(increment.source))

			if cmp := ls.Distance(); cmp != increment.distance {
				t.Logf("Comparing '%s' and '%s' incrementally failed: result was %v but expected was %v\n",
					string(ls.source),
					string(ls.target),
					cmp,
					increment.distance)
				t.Logf("Matrix was:\n")
				ls.print_matrix(t)
				t.FailNow()
			}
		}

		if len(v.source_increments) > 1 { // test allows for roll back
			err := ls.RollbackBy(len(v.source_increments[len(v.source_increments) - 1].source))
			if err != nil {
				t.Errorf("An error occurred during rollback: %v", err)
			}

			if newDistance := ls.Distance(); newDistance != v.source_increments[len(v.source_increments) - 2].distance {
				t.Errorf("Distance after rolling back is %d, expected is %d",
					newDistance,
					v.source_increments[len(v.source_increments) - 2].distance)
			}
		}
	}
}

type similarity_range struct {
	low  float32
	high float32
}

type expected_match struct {
	word string
	similarity_range similarity_range
}

func Benchmark_recursiveLevenshtein(b *testing.B) {
    for i := 0; i < b.N; i++ {
    	for _, v := range basic_test_cases {
    		compare_silces_r([]rune(v.source), []rune(v.destination))
    	}
    }
}

func Benchmark_matrixLevenshtein(b *testing.B) {
    for i:= 0; i < b.N; i++ {    
    	for _, v := range basic_test_cases {
    		ls := InitLState(v.destination)
    		ls.UpdateState([]rune(v.source))
    	}
    }
}

func Benchmark_TwoRowsLevenshtein(b *testing.B) {
    for i := 0; i < b.N; i++ {
        for _, v := range basic_test_cases {
            LevenshteinDistance(v.destination, v.source)
        }
    }
}
