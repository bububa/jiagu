package model

type PerceptronJSONModel struct {
	Weights map[string]map[string]float64
	Classes map[string]struct{}
}
