package bayes

type Model struct {
	Total float64          `json:"total,omitempty"`
	Data  map[string]Probe `json:"d,omitempty"`
}
