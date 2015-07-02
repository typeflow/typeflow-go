package typeflow

import (
	"errors"
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

func maximum(values ...int) (int, error) {
	if len(values) == 0 {
		return -1, errors.New("Cannot find maximum of empty list")
	}
	var max int = values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	return max, nil
}

func abs(value int) (v int) {
	if value > 0 {
		v = value
		return
	}

	v = -value
	return
}
