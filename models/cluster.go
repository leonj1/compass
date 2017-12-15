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
	Namespace   *cmap.ConcurrentMap
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
		if _, ok := existingCluster.Crds.Get(crdName); ok {
			existingCluster.Crds.Set(crd.Name, crd)
			if crdName != crd.Name {
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
		if _, ok := m.Clusters.Get(clusterName); ok {
			cluster.Nodes.Set(node.Name, node)
			if nodeName != node.Name {
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

func (m *RootRouter) AddNamespace(clusterName string, namespace Namespace) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	if tmpCluster, ok := m.Clusters.Get(clusterName); ok {
		existingCluster := tmpCluster.(Cluster)
		if existingCluster.Namespace.Has(namespace.Name) {
			return errors.New("namespace by that name already exists")
		}
		existingCluster.Namespace.Set(namespace.Name, namespace)
		m.Clusters.Set(existingCluster.Name, existingCluster)
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}

func (m *RootRouter) UpdateNamespace(clusterName, namespaceName string, namespace Namespace) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	if tmpCluster, ok := m.Clusters.Get(clusterName); ok {
		existingCluster := tmpCluster.(Cluster)
		if _, ok := existingCluster.Namespace.Get(namespaceName); ok {
			if namespaceName != namespace.Name {
				existingCluster.Namespace.Remove(namespaceName)
			}
			existingCluster.Namespace.Set(namespace.Name, namespace)
			m.Clusters.Set(clusterName, existingCluster)
		} else {
			return errors.New("node by that name does not exist")
		}
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}
