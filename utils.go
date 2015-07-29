package typeflow

func minimum(values ...int) (int) {
	if len(values) == 0 {
		panic("Cannot find minimum of empty list")
	}
	var min int = values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func maximum(values ...int) (int) {
	if len(values) == 0 {
		panic("Cannot find maximum of empty list")
	}
	var max int = values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	return max
}

func abs(value int) (v int) {
	if value > 0 {
		v = value
		return
	}

	v = -value
	return
}
