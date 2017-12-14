package models

type Cluster struct {
	Name        string
	Status      string
	Personality string
	Crds        *[]Crd
	Nodes       *[]Node
	Namespace   *[]Namespace
	Events      *[]Events
}
