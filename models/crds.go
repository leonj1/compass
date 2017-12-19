package models

type Crd struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type NameVersionContext struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}
