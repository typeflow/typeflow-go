package typeflow

import (
	"testing"
	"fmt"
	"strings"
	"os"
	"io"
	"bufio"
)

func (s* LState) printMatrix(t* testing.T) {
	row := ""
	for col_index, _ := range s.vec1 {
		row = fmt.Sprintf("%s[ %v ]", row, s.vec1[col_index])
	}
    row = fmt.Sprintf("%s\n", row)
	for col_index, _ := range s.vec2 {
		row = fmt.Sprintf("%s[ %v ]", row, s.vec2[col_index])
	}
	t.Logf("\n%s\n", row)
}

type minimum_testcase struct {
	values []int
	result int
}

var minimum_testcases = []minimum_testcase{
	{ []int{ 3,56,21,45,2,4,1,2 }, 1 },
	{ []int{ 63,12,4,32,0,7,8,5,34,90}, 0 },
	{ []int{ 1,2,3,4}, 1},
	{ []int{ 9,8,7,6,5}, 5},
}

func Test_minimum(t *testing.T) {
	for _, v := range minimum_testcases {
		if min, err := minimum(v.values...); err != nil || min != v.result {
			t.Errorf("minimum(%v): unexpected result got %d; expected: %d", v.values, min, v.result)
		}
	}
}

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

// the following cases represent subsequent
// additions: {'al, 'al', 0}, {'ex', 'es', 1} == 'alex', 'ales', distance: 1
var incrementalTestCases = [][]perm{
	{ {"aaa", "abc", 2 } },
	{{ "al", "al",  0 }, { "ex", "es", 1}},
	{{"a", "b", 1}, {"aa", "ba", 2}, {"le", "a", 4}},
	{{"a", "a", 0}, { "l", "l", 0}, {"e", "e", 0}, {"s", "s", 0}, {"s", "a", 1}, {"a", "s", 2}, {"n", "n", 2} },
	{{"iraq", "rep of ireland", 11}},
}

func Test_compareSlicesRecursive(t *testing.T) {
	for _, v := range testCases {
		if cmp := compareSlicesRecursive([]byte(v.source), []byte(v.destination)); cmp != v.distance {
			t.Errorf("Comparing '%s' and '%s' failed: result was %v but expected was %v", v.source, v.destination, cmp, v.distance)
		}
	}
}

func Test_compareSlicesWithMatrixBaseCase(t *testing.T) {
	for _, v := range incrementalTestCases {
		source := make([]rune, 0)
		dest   := make([]rune, 0)
		for _, test_case := range v {
			source = append(source, []rune(test_case.source)...)
			dest   = append(dest, []rune(test_case.destination)...)
		}

		ls := InitLState()
		ls.UpdateState([]rune(source), []rune(dest))

		if cmp := ls.Distance(); cmp != v[len(v)-1].distance {
			t.Logf("Comparing '%s' and '%s' failed: result was %v but expected was %v\n", string(source), string(dest), cmp, v[len(v)-1].distance)
			t.Logf("Matrix was:\n")
			ls.printMatrix(t)
			t.FailNow()
		}
	}
}

func Test_compareSlicesWithMatrixIncremental(t* testing.T) {
	for _, v := range incrementalTestCases {
		ls := InitLState()

		for _, increment := range v {
			ls.UpdateState([]rune(increment.source), []rune(increment.destination))

			if cmp := ls.Distance(); cmp != increment.distance {
				t.Logf("Comparing '%s' and '%s' incrementally failed: result was %v but expected was %v\n",
					string(ls.w1),
					string(ls.w2),
					cmp,
					increment.distance)
				t.Logf("Matrix was:\n")
				ls.printMatrix(t)
				t.FailNow()
			}
		}

		if len(v) > 1 { // test allows for roll back
			err := ls.RollbackBy(len(v[len(v) - 1].source), len(v[len(v) - 1].destination))
			if err != nil {
				t.Errorf("An error occurred during rollback: %v", err)
			}

			if newDistance := ls.Distance(); newDistance != v[len(v) - 2].distance {
				t.Errorf("Distance after rolling back is %d, expected is %d", newDistance, v[len(v) - 2].distance)
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

type word_source_test struct {
	substr   string
	expected_matches []expected_match
}

var word_source_tests = []word_source_test{
	{ "rep of ireland", []expected_match{{"Ireland (Republic)", similarity_range{0.3, 0.35}}} },
}

func Test_wordSource(t *testing.T) {
    ws := NewWordSource()

	var filter WordFilter = func (w string) (word string, skip bool) {
		word = strings.ToLower(w)
		skip = false

		return
	}

	// building country name
	// source from file
    file, err := os.Open("testdata/countries.txt")
	country_names := make([]string, 0)
	if err != nil {
		t.Log("Cannot open expected file testdata/countries.txt. Skipping this test.")
		t.SkipNow()
		return
	}
    reader := bufio.NewReader(file)
	for  {
		line, err := reader.ReadString('\n');
		if err == io.EOF {
			break
		}
		country_names = append(country_names, line[:len(line)-1])
	}

	ws.SetSource(country_names, []WordFilter{ filter })

OuterLoop:
	for _, test := range word_source_tests {
		t.Logf("Finding matches for substring '%s'", test.substr)
		matches, err := ws.FindMatch(test.substr, 0.32)
		if err != nil {
			t.Logf("An error occurred: %v", err)
			for _, m := range matches {
				t.Logf("%s, %f", m.Word, m.Similarity)
			}
			t.FailNow()
		}

		for _, match := range matches {
			for _, expected := range test.expected_matches {
				if match.Similarity >= expected.similarity_range.low &&
				match.Similarity <= expected.similarity_range.high {
					t.Log("Found!")
					continue OuterLoop
				}
			}
		}
		t.Logf("Couldn't find expected match")
		t.Logf("Found the following matches:")
		for _, m := range matches {
			t.Logf("'%s', '%f'", m.Word, m.Similarity)
		}
		t.FailNow()
	}
}

func Benchmark_recursiveLevenshtein(b *testing.B) {
    for i:= 0; i < b.N; i++ {
    	for _, v := range testCases {
    		compareSlicesRecursive([]byte(v.source), []byte(v.destination))
    	}
    }
}

func Benchmark_matrixLevenshtein(b *testing.B) {
    for i:= 0; i < b.N; i++ {    
    	for _, v := range testCases {
    		ls := InitLState()
    		ls.UpdateState([]rune(v.source), []rune(v.destination))
    	}
    }
}
