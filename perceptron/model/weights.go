package model

// Each feature gets its own weight vector, so weights is a dict-of-dicts
type Weights map[string]map[string]float64

// NewWeights init Weights
func NewWeights() Weights {
	return make(map[string]map[string]float64)
}

// SetWeight set a weight value for feat->clas
func (w Weights) SetWeight(feat string, clas string, value float64) {
	if _, found := w[feat]; !found {
		w[feat] = make(map[string]float64)
	}
	w[feat][clas] = value
}

// GetWeight get a weight value of feat->clas
func (w Weights) GetWeight(feat string, clas string) float64 {
	if _, found := w[feat]; !found {
		return 0
	}
	if val, found := w[feat][clas]; found {
		return val
	}
	return 0
}

// GetFeatureWeights get a feat classes weights
func (w Weights) GetFeatureWeights(feat string) map[string]float64 {
	return w[feat]
}
