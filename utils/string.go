package utils

// StringFromIndex 从i截取字符串
func StringFromIndex(str string, i int) string {
	w := []rune(str)
	return RuneFromIndex(w, i)
}

// StringInRange 从from到to截取字符串
func StringInRange(str string, from int, to int) string {
	w := []rune(str)
	return RuneInRange(w, from, to)
}

// StringInIndex 获取i位置字符串
func StringInIndex(str string, i int) string {
	w := []rune(str)
	return RuneInIndex(w, i)
}

// StringInSlice 获取[]string中i位置字符串
func StringInSlice(arr []string, i int) string {
	l := len(arr)
	return arr[idxInSlice(l, i)]
}

// StringSliceInRange 获取[]string中from-to位置字符串
func StringSliceInRange(arr []string, from int, to int) []string {
	l := len(arr)
	if from <= -1*l && to <= -1*l {
		return []string{}
	}
	if from >= l {
		return []string{}
	}
	return arr[idxInSlice(l, from) : idxInSlice(l, to-1)+1]
}

// RuneFromIndex 从i截取字符串
func RuneFromIndex(w []rune, i int) string {
	l := len(w)
	return string(w[idxInSlice(l, i):])
}

// RuneInRange 从from到to截取字符串
func RuneInRange(w []rune, from int, to int) string {
	l := len(w)
	if from <= -1*l && to <= -1*l {
		return ""
	}
	if from >= l {
		return ""
	}
	return string(w[idxInSlice(l, from) : idxInSlice(l, to-1)+1])
}

// RuneInIndex 获取i位置字符串
func RuneInIndex(w []rune, i int) string {
	l := len(w)
	return string(w[idxInSlice(l, i)])
}

// StringSplit []rune转[]string
func StringSplit(str string) []string {
	w := []rune(str)
	var ret = make([]string, len(w))
	for idx, r := range w {
		ret[idx] = string(r)
	}
	return ret
}
