package models

type Namespace struct {
	Name     string
	Crds     map[string]int
	PodCount int
}
