package model

import (
	"strings"
)

// Each feature gets its own weight vector, so weights is a dict-of-dicts
type Weights struct {
	values   [][]float64
	classes  [][]string
	weights  map[string]int
	features map[string]int
}

// NewWeights init Weights
func NewWeights(cap int) *Weights {
	return &Weights{
		values:   make([][]float64, 0, cap),
		classes:  make([][]string, 0, cap),
		weights:  make(map[string]int, cap),
		features: make(map[string]int, cap),
	}
}

func NewWeightsFromMap(mp map[string]map[string]float64) *Weights {
	l := len(mp)
	w := NewWeights(l)
	for feat, classes := range mp {
		for clas, val := range classes {
			w.SetWeight(feat, clas, val)
		}
	}
	return w
}

//GetWeight get class weight value
func (w Weights) GetWeight(feat string, clas string) float64 {
	key := FeatureClassKey(feat, clas)
	if featIdx, found := w.features[feat]; found {
		if weightIdx, found := w.weights[key]; found {
			return w.values[featIdx][weightIdx]
		}
	}
	return 0
}

// SetWeight set a weight value for feat->clas
func (w *Weights) SetWeight(feat string, clas string, value float64) {
	key := FeatureClassKey(feat, clas)
	if featIdx, found := w.features[feat]; found {
		if weightIdx, found := w.weights[key]; found {
			w.values[featIdx][weightIdx] = value
		} else {
			lastWeightIdx := len(w.weights)
			w.weights[key] = lastWeightIdx
			w.classes[featIdx] = append(w.classes[featIdx], clas)
			w.values[featIdx] = append(w.values[featIdx], value)
		}
	} else {
		lastFeatIdx := len(w.features)
		w.features[feat] = lastFeatIdx
		w.weights[key] = 0
		w.classes = append(w.classes, []string{clas})
		w.values = append(w.values, []float64{value})
	}
}

func (w Weights) GetFeatureLength(feat string) int {
	if featIdx, found := w.features[feat]; found {
		return len(w.values[featIdx])
	}
	return 0
}

// GetFeatureWeights get a feat classes weights
func (w Weights) GetFeatureWeights(feat string) <-chan KV {
	ch := make(chan KV, w.GetFeatureLength(feat))
	var (
		values  []float64
		classes []string
	)
	if featIdx, found := w.features[feat]; found {
		values = w.values[featIdx]
		classes = w.classes[featIdx]
	}
	go func(classes []string, values []float64) {
		for idx, clas := range classes {
			ch <- KV{
				Label: clas,
				Value: values[idx],
			}
		}
		close(ch)
	}(classes, values)
	return ch
}

// Map convert weights to map
func (w Weights) Map() map[string]map[string]float64 {
	l := len(w.values)
	mp := make(map[string]map[string]float64, l)
	for feat, featIdx := range w.features {
		mp[feat] = make(map[string]float64, len(w.values[featIdx]))
		for clasIdx, clas := range w.classes[featIdx] {
			mp[feat][clas] = w.values[featIdx][clasIdx]
		}
	}
	return mp
}

// Features
func (w Weights) Features() <-chan string {
	ch := make(chan string, len(w.features))
	go func(ch chan string) {
		for feat, _ := range w.features {
			ch <- feat
		}
		close(ch)
	}(ch)
	return ch
}

func FeatureClassKey(feat string, clas string) string {
	var buf strings.Builder
	buf.WriteString(feat)
	buf.WriteString("-")
	buf.WriteString(clas)
	return buf.String()
}
