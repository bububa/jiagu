package textrank

type Score struct {
	Idx   int
	Value float64
}

type ScoreSlice []Score

func NewScoreSlice(scores []float64) ScoreSlice {
	s := make(ScoreSlice, 0, len(scores))
	for idx, score := range scores {
		s = append(s, Score{
			Idx:   idx,
			Value: score,
		})
	}
	return s
}

func (s ScoreSlice) Len() int      { return len(s) }
func (s ScoreSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ScoreSlice) Less(i, j int) bool {
	if s[i].Value == s[j].Value {
		return s[i].Idx < s[j].Idx
	}
	return s[i].Value < s[j].Value
}
