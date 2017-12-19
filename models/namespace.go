package models

type Namespace struct {
	Name     string         `json:"name,omitempty"`
	Crds     map[string]int `json:"crds,omitempty"`
	PodCount int            `json:"pod_count,omitempty"`
}
