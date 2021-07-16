package knowledge

// Subject model subject
type Subject struct {
	Word  string
	Label string
	Index int
}

// Entity entity result
type Entity struct {
	Subject   string `json:"subject,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	Value     string `json:"value,omitempty"`
}

func (e Entity) String() string {
	return e.Subject + e.Attribute + e.Value
}
