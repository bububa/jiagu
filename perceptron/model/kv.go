package model

type KV struct {
	Label string
	Value float64
}

func (k KV) IsZero() bool {
	return k.Value < 1e-15
}

type KVSlice []KV

func NewKVSlice(cap int) KVSlice {
	return KVSlice(make([]KV, 0, cap))
}

func (s KVSlice) Len() int      { return len(s) }
func (s KVSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s KVSlice) Less(i, j int) bool {
	if s[i].Value == s[j].Value {
		return s[i].Label < s[j].Label
	}
	return s[i].Value < s[j].Value
}
