package typeflow

import (
	"testing"
)

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
		if min := minimum(v.values...); min != v.result {
			t.Errorf("minimum(%v): unexpected result got %d; expected: %d", v.values, min, v.result)
		}
	}
}
