package models

import (
	"github.com/kataras/go-errors"
	"github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"
)

type Cluster struct {
	Name        string              `json:"name,omitempty"`
	Status      string              `json:"status,omitempty"`
	Personality string              `json:"personality,omitempty"`
	Crds        *cmap.ConcurrentMap `json:"crds,omitempty"`
	Nodes       *cmap.ConcurrentMap `json:"nodes,omitempty"`
	Namespace   *cmap.ConcurrentMap `json:"namespace,omitempty"`
	Events      string              `json:"events,omitempty"`
}

type ClusterContext struct {
	Name        string                        `json:"name,omitempty"`
	Status      string                        `json:"status,omitempty"`
	Personality string                        `json:"personality,omitempty"`
	Crds        map[string]NameVersionContext `json:"crds,omitempty"`
	Nodes       map[string]NameVersionContext `json:"nodes,omitempty"`
	Namespace   map[string]Namespace          `json:"namespace,omitempty"`
	Events      string                        `json:"events,omitempty"`
}

func (m *RootRouter) AddCluster(cluster ClusterContext) error {
	if m.Clusters.Has(cluster.Name) {
		log.Print("cluster already exists")
		return errors.New("cluster already exists")
	}

	// iterate over Crds, Namespace, Nodes
	crds := cmap.New()
	for k, v := range cluster.Crds {
		crds.Set(k, v)
	}
	namespaces := cmap.New()
	for k, v := range cluster.Namespace {
		namespaces.Set(k, v)
	}
	nodes := cmap.New()
	for k, v := range cluster.Nodes {
		nodes.Set(k, v)
	}

	m.Clusters.Set(cluster.Name, Cluster{
		Name:        cluster.Name,
		Status:      cluster.Status,
		Personality: cluster.Personality,
		Events:      cluster.Events,
		Crds:        &crds,
		Nodes:       &nodes,
		Namespace:   &namespaces,
	})
	return nil
}

func (m *RootRouter) UpdateCluster(clusterName string, cluster ClusterContext) error {
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
	log.Print("In AddCustomResource")
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	log.Print("Fetching cluster from map")
	if tmpCluster, ok := m.Clusters.Get(clusterName); ok {
		existingCluster := tmpCluster.(Cluster)
		if existingCluster.Crds.Has(crd.Name) {
			log.Error("crd already exists in cluster")
			return errors.New("crd already exists in cluster")
		}
		log.Info("Adding crd")
		existingCluster.Crds.Set(crd.Name, crd)
		m.Clusters.Set(clusterName, existingCluster)
	} else {
		log.Error("unable to fetch cluster from map for some unknown reason")
		return errors.New("unable to fetch cluster from map for some reason")
	}
	log.Info("Done adding crd")
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

func (m *RootRouter) SetEvents(clusterName, events string) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	if tmpCluster, ok := m.Clusters.Get(clusterName); ok {
		existingCluster := tmpCluster.(Cluster)
		existingCluster.Events = events
		m.Clusters.Set(existingCluster.Name, existingCluster)
	} else {
		return errors.New("unable to fetch cluster from map for some reason")
	}
	return nil
}

func (m *RootRouter) GetAClusterByName(clusterName string) (*Cluster, error) {
	if !m.Clusters.Has(clusterName) {
		return nil, errors.New("cluster does not exist")
	}

	if tmpCluster, ok := m.Clusters.Get(clusterName); ok {
		existingCluster := tmpCluster.(Cluster)
		return &existingCluster, nil
	}
	return nil, errors.New("cluster does not exist")
}

func (m *RootRouter) DeleteAClusterByName(clusterName string) error {
	if !m.Clusters.Has(clusterName) {
		return errors.New("cluster does not exist")
	}

	m.Clusters.Remove(clusterName)
	return nil
}
