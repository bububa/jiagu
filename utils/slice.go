package utils

func idxInSlice(l int, i int) int {
	if i <= -1*l {
		return 0
	} else if i < 0 {
		return l + i
	} else if i >= l {
		return l - 1
	}
	return i
}
