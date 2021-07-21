package model

type PerceptronModel struct {
	Weights map[string]map[string]float64
	Classes map[string]struct{}
}
