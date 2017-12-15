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

	if tmpCluster, ok := m.Clusters.Get(clusterName); ok {
		existingCluster := tmpCluster.(Cluster)
		hasNameChanged := existingCluster.Name != cluster.Name
		existingCluster.Name = cluster.Name
		existingCluster.Status = cluster.Status
		existingCluster.Personality = cluster.Personality
		m.Clusters.Set(cluster.Name, existingCluster)
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

	if tmpCluster, ok := m.Clusters.Get(clusterName); ok {
		existingCluster := tmpCluster.(Cluster)
		if existingCluster.Crds.Has(crd.Name) {
			return errors.New("crd already exists in cluster")
		}
		existingCluster.Crds.Set(crd.Name, crd)
		m.Clusters.Set(clusterName, existingCluster)
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}

func (m *RootRouter) UpdateCustomResource(clusterName, crdName string, crd Crd) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	if tmpCluster, ok := m.Clusters.Get(clusterName); ok {
		existingCluster := tmpCluster.(Cluster)
		if tmpCrd, ok := existingCluster.Crds.Get(crdName); ok {
			existingCrd := tmpCrd.(Crd)
			hasNameChanged := existingCrd.Name != crd.Name
			existingCrd.Name = crd.Name
			existingCrd.Version = crd.Version
			existingCluster.Crds.Set(crd.Name, existingCrd)
			if hasNameChanged {
				existingCluster.Crds.Remove(crdName)
			}
			m.Clusters.Set(clusterName, existingCluster)
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

func (m *RootRouter) UpdateNode(clusterName, nodeName string, node Node) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	if existing, ok := m.Clusters.Get(clusterName); ok {
		cluster := existing.(Cluster)
		if tmpNode, ok := m.Clusters.Get(clusterName); ok {
			existingNode := tmpNode.(Node)
			hasNameChanged := existingNode.Name != node.Name
			existingNode.Name = node.Name
			existingNode.Version = node.Version
			cluster.Nodes.Set(existingNode.Name, existingNode)
			if hasNameChanged {
				cluster.Nodes.Remove(nodeName)
			}
			m.Clusters.Set(clusterName, cluster)
		} else {
			return errors.New("node by that name does not exist")
		}
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}
