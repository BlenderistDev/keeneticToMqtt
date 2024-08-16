package keeneticdto

type Policy struct {
	Description string         `json:"description"`
	Permit      []PolicyPermit `json:"permit"`
}

type PolicyPermit struct {
	Interface string `json:"interface"`
	Enabled   bool   `json:"enabled,omitempty"`
}
