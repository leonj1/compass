package models

import (
	"github.com/kataras/go-errors"
	"github.com/orcaman/concurrent-map"
)

type Cluster struct {
	Name        string
	Status      string
	Personality string
	Crds        *cmap.ConcurrentMap
	Nodes       *cmap.ConcurrentMap
	Namespace   *[]Namespace
	Events      *[]Events
}

func (m *RootRouter) AddCluster(cluster Cluster) error {
	if m.Clusters.Has(cluster.Name) {
		return errors.New("cluster already exists")
	}
	m.Clusters.Set(cluster.Name, cluster)
	return nil
}

func (m *RootRouter) UpdateCluster(clusterName string, cluster Cluster) error {
	if !m.Clusters.Has(cluster.Name) {
		return errors.New("cluster does not exist")
	}

	if existing, ok := m.Clusters.Get(clusterName); ok {
		current := existing.(Cluster)
		hasNameChanged := current.Name != cluster.Name
		current.Name = cluster.Name
		current.Status = cluster.Status
		current.Personality = cluster.Personality
		m.Clusters.Set(cluster.Name, current)
		if hasNameChanged {
			m.Clusters.Remove(clusterName)
		}
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}

func (m *RootRouter) AddCustomResource(clusterName string, crd Crd) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	if existing, ok := m.Clusters.Get(clusterName); ok {
		current := existing.(Cluster)
		if current.Crds.Has(crd.Name) {
			return errors.New("crd already exists in cluster")
		}
		current.Crds.Set(crd.Name, crd)
		m.Clusters.Set(clusterName, current)
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}

func (m *RootRouter) UpdateCustomResource(clusterName, crdName string, crd Crd) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	if existing, ok := m.Clusters.Get(clusterName); ok {
		current := existing.(Cluster)
		if !current.Crds.Has(crd.Name) {
			return errors.New("crd does not exist in cluster")
		}
		if existingCrd, ok := current.Crds.Get(crdName); ok {
			currentCrd := existingCrd.(Crd)
			hasNameChanged := currentCrd.Name != crd.Name
			currentCrd.Name = crd.Name
			currentCrd.Version = crd.Version
			current.Crds.Set(crd.Name, currentCrd)
			if hasNameChanged {
				current.Crds.Remove(crdName)
			}
		} else {
			return errors.New("unable to fetch crd from cluster some reason")
		}
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}

func (m *RootRouter) AddNode(clusterName string, node Node) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	if existing, ok := m.Clusters.Get(clusterName); ok {
		current := existing.(Cluster)
		if current.Nodes.Has(node.Name) {
			return errors.New("node by that name already exists")
		}
		current.Nodes.Set(node.Name, node)
		m.Clusters.Set(current.Name, current)
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}
