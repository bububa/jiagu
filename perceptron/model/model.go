package model

type PerceptronModel struct {
	Weights Weights
	Classes map[string]struct{}
}
