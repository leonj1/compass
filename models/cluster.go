package models

import "github.com/orcaman/concurrent-map"

type Cluster struct {
	Name        string
	Status      string
	Personality string
	Crds        *cmap.ConcurrentMap
	Nodes       *[]Node
	Namespace   *[]Namespace
	Events      *[]Events
}
